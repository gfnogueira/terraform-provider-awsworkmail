package awsworkmail

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// domainResourceModel describes the resource data model.
type domainResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Domain         types.String `tfsdk:"domain"`
	MXRecords      types.List   `tfsdk:"mx_records"`
}

type domainResource struct {
	cfg aws.Config
}

func NewDomainResource() resource.Resource {
	return &domainResource{}
}

func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the WorkMail domain (domain name)",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the WorkMail organization",
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Domain name to add to WorkMail organization",
			},
			"mx_records": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of MX records to configure in your DNS for this domain.",
			},
		},
	}
}

func (r *domainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(aws.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected aws.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.cfg = cfg
}

func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data domainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := workmail.NewFromConfig(r.cfg)

	// Register the domain with WorkMail
	input := &workmail.RegisterMailDomainInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		DomainName:     aws.String(data.Domain.ValueString()),
	}

	_, err := client.RegisterMailDomain(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error registering WorkMail domain", err.Error())
		return
	}

	data.ID = data.Domain

	getDomainInput := &workmail.GetMailDomainInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		DomainName:     aws.String(data.Domain.ValueString()),
	}

	domainOutput, err := client.GetMailDomain(ctx, getDomainInput)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving WorkMail domain details", err.Error())
		return
	}

	var mxRecords []string
	for _, record := range domainOutput.Records {
		if record.Type != nil && *record.Type == "MX" && record.Value != nil {
			mxRecords = append(mxRecords, *record.Value)
		}
	}

	mxRecordsList, diags := types.ListValueFrom(ctx, types.StringType, mxRecords)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MXRecords = mxRecordsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data domainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := workmail.NewFromConfig(r.cfg)

	// Get domain information
	input := &workmail.GetMailDomainInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		DomainName:     aws.String(data.Domain.ValueString()),
	}

	output, err := client.GetMailDomain(ctx, input)
	if err != nil {
		// If domain is not found, remove from state
		if strings.Contains(err.Error(), "EntityNotFoundException") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading WorkMail domain", err.Error())
		return
	}

	// Update MX records
	var mxRecords []string
	for _, record := range output.Records {
		if record.Type != nil && *record.Type == "MX" && record.Value != nil {
			mxRecords = append(mxRecords, *record.Value)
		}
	}

	mxRecordsList, diags := types.ListValueFrom(ctx, types.StringType, mxRecords)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MXRecords = mxRecordsList

	// Save the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data domainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For domains, updates are mostly read-only operations
	// The domain name itself cannot be changed, only other attributes

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data domainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := workmail.NewFromConfig(r.cfg)

	// Deregister the domain from WorkMail
	input := &workmail.DeregisterMailDomainInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		DomainName:     aws.String(data.Domain.ValueString()),
	}

	_, err := client.DeregisterMailDomain(ctx, input)
	if err != nil {
		// If domain is already gone, that's fine
		if !strings.Contains(err.Error(), "EntityNotFoundException") {
			resp.Diagnostics.AddError("Error deregistering WorkMail domain", err.Error())
		}
	}
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected import ID format: <organization_id>,<domain>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
