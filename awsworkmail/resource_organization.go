package awsworkmail

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type organizationResource struct{}

// NewOrganizationResource returns a new WorkMail organization resource.
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

type organizationResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Alias types.String `tfsdk:"alias"`
}

// Metadata sets the resource type name.
func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the WorkMail organization (alias used as ID)",
			},
			"alias": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Alias for the WorkMail organization",
			},
		},
	}
}

// Create creates a new WorkMail organization.
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data organizationResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}

	client := workmail.NewFromConfig(cfg)

	alias := data.Alias.ValueString()

	out, err := client.CreateOrganization(ctx, &workmail.CreateOrganizationInput{
		Alias: &alias,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating WorkMail organization", err.Error())
		return
	}

	orgID := *out.OrganizationId

	// Wait for organization to become Active
	for i := 0; i < 30; i++ { // up to ~5 minutes
		listOut, err := client.ListOrganizations(ctx, &workmail.ListOrganizationsInput{})
		if err == nil {
			for _, org := range listOut.OrganizationSummaries {
				if org.OrganizationId != nil && *org.OrganizationId == orgID {
					if org.State != nil && *org.State == "Active" {
						goto active
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	resp.Diagnostics.AddWarning("Organization not Active", "Organization did not become Active after 5 minutes. Dependent resources may fail.")
active:
	// Set state so Terraform can track this resource
	data.ID = types.StringValue(orgID)
	data.Alias = types.StringValue(alias)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read is a no-op because WorkMail API lacks a proper describe call by alias.
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data organizationResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update is a no-op because WorkMail does not support updating alias.
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data organizationResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the WorkMail organization.
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data organizationResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}

	client := workmail.NewFromConfig(cfg)

	orgID := data.ID.ValueString()

	_, err = client.DeleteOrganization(ctx, &workmail.DeleteOrganizationInput{
		OrganizationId: &orgID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting WorkMail organization", err.Error())
	}
}

// ImportState imports a WorkMail organization by ID (alias).
func (r *organizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
