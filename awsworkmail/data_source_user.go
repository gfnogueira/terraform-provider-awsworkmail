package awsworkmail

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces.
var _ datasource.DataSource = &userDataSource{}

// userDataSource is the data source implementation.
type userDataSource struct {
	cfg aws.Config
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for querying AWS WorkMail users.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "The WorkMail Organization ID.",
				Required:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The WorkMail User ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The primary email of the user.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the user.",
				Computed:    true,
			},
		},
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(aws.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected aws.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.cfg = cfg
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data struct {
		OrganizationId types.String `tfsdk:"organization_id"`
		UserId         types.String `tfsdk:"user_id"`
		Name           types.String `tfsdk:"name"`
		Email          types.String `tfsdk:"email"`
		State          types.String `tfsdk:"state"`
	}

	diag := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := workmail.NewFromConfig(d.cfg)

	input := &workmail.DescribeUserInput{
		OrganizationId: aws.String(data.OrganizationId.ValueString()),
		UserId:         aws.String(data.UserId.ValueString()),
	}

	output, err := client.DescribeUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to describe WorkMail user",
			err.Error(),
		)
		return
	}
	if output != nil && output.Name != nil {
		data.Name = types.StringValue(*output.Name)
	}
	if output != nil && output.Email != nil {
		data.Email = types.StringValue(*output.Email)
	}
	if output != nil && output.State != "" {
		data.State = types.StringValue(string(output.State))
	}

	diag = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diag...)
}
