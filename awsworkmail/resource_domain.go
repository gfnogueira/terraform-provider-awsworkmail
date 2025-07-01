package awsworkmail

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type domainResource struct{}

func NewDomainResource() resource.Resource {
	return &domainResource{}
}

type domainResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	OrganizationID types.String   `tfsdk:"organization_id"`
	Domain         types.String   `tfsdk:"domain"`
	MXRecords      []types.String `tfsdk:"mx_records"`
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

func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning("Not Implemented", "AWS SDK v2 does not support WorkMail domain registration. This resource is a stub for documentation only.")
}

func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.AddWarning("Not Implemented", "AWS SDK v2 does not support WorkMail domain registration. This resource is a stub for documentation only.")
}

func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op: not implemented
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Not Implemented", "AWS SDK v2 does not support WorkMail domain registration. This resource is a stub for documentation only.")
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
