package awsworkmail

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// awsworkmail_user resource: manages a user in a WorkMail organization

type userResource struct{}

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	DisplayName    types.String `tfsdk:"display_name"`
	FirstName      types.String `tfsdk:"first_name"`
	LastName       types.String `tfsdk:"last_name"`
	Password       types.String `tfsdk:"password"`
	Email          types.String `tfsdk:"email"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the WorkMail user",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the WorkMail organization",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User name (login name)",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display name for the user",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for the user (must meet AWS WorkMail requirements)",
			},
			"email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Primary email address for the user (optional, can be set after creation)",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the user is enabled in WorkMail.",
			},
			"first_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "First name of the user (optional)",
			},
			"last_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Last name of the user (optional)",
			},
		},
	}
}

func isEntityStateException(err error) bool {
	return err != nil && strings.Contains(err.Error(), "EntityStateException")
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}
	client := workmail.NewFromConfig(cfg)

	input := &workmail.CreateUserInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		Name:           aws.String(data.Name.ValueString()),
		DisplayName:    aws.String(data.DisplayName.ValueString()),
		Password:       aws.String(data.Password.ValueString()),
	}
	if !data.FirstName.IsNull() && data.FirstName.ValueString() != "" {
		input.FirstName = aws.String(data.FirstName.ValueString())
	}
	if !data.LastName.IsNull() && data.LastName.ValueString() != "" {
		input.LastName = aws.String(data.LastName.ValueString())
	}
	out, err := client.CreateUser(ctx, input)
	if err != nil {
		if isEntityStateException(err) {
			resp.Diagnostics.AddError(
				"WorkMail user cannot be created due to entity state",
				"AWS returned EntityStateException. This usually means the user is in a state that does not allow creation or update. Please check if the user already exists, is enabled, or is registered to WorkMail. If the problem persists, try deleting and recreating the user in the AWS Console.\n\nOriginal error: "+err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError("Error creating WorkMail user", err.Error())
		return
	}
	data.ID = types.StringValue(*out.UserId)

	// Enable or disable user if needed
	enabled := true
	if !data.Enabled.IsNull() {
		enabled = data.Enabled.ValueBool()
	}
	if enabled {
		_, _ = client.RegisterToWorkMail(ctx, &workmail.RegisterToWorkMailInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			EntityId:       aws.String(data.ID.ValueString()),
			Email:          aws.String(data.Email.ValueString()),
		})
	} else {
		_, _ = client.DeregisterFromWorkMail(ctx, &workmail.DeregisterFromWorkMailInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			EntityId:       aws.String(data.ID.ValueString()),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}
	client := workmail.NewFromConfig(cfg)

	input := &workmail.DescribeUserInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		UserId:         aws.String(data.ID.ValueString()),
	}
	out, err := client.DescribeUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error reading WorkMail user", err.Error())
		return
	}
	if out != nil && out.UserId != nil && *out.UserId != "" {
		data.ID = types.StringValue(*out.UserId)
	}
	if out != nil && out.Name != nil {
		data.Name = types.StringValue(*out.Name)
	}
	if out != nil && out.DisplayName != nil {
		data.DisplayName = types.StringValue(*out.DisplayName)
	}
	if out != nil && out.Email != nil {
		data.Email = types.StringValue(*out.Email)
	}
	if out != nil && out.State != "" {
		data.Enabled = types.BoolValue(out.State == "ENABLED")
	}
	if out != nil && out.FirstName != nil {
		data.FirstName = types.StringValue(*out.FirstName)
	}
	if out != nil && out.LastName != nil {
		data.LastName = types.StringValue(*out.LastName)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		var stateData userResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		data.ID = stateData.ID
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// Only display name and password can be updated
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}
	client := workmail.NewFromConfig(cfg)

	if !data.DisplayName.IsNull() || !data.FirstName.IsNull() || !data.LastName.IsNull() {
		updateInput := &workmail.UpdateUserInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			UserId:         aws.String(data.ID.ValueString()),
		}
		if !data.DisplayName.IsNull() {
			updateInput.DisplayName = aws.String(data.DisplayName.ValueString())
		}
		if !data.FirstName.IsNull() {
			updateInput.FirstName = aws.String(data.FirstName.ValueString())
		}
		if !data.LastName.IsNull() {
			updateInput.LastName = aws.String(data.LastName.ValueString())
		}
		_, err := client.UpdateUser(ctx, updateInput)
		if err != nil {
			if isEntityStateException(err) {
				resp.Diagnostics.AddError(
					"WorkMail user cannot be updated due to entity state",
					"AWS returned EntityStateException. This usually means the user is in a state that does not allow updates (e.g., not enabled, not registered, or disabled). Please check the user's state in the AWS Console. If the user is stuck, try deleting and recreating the user.\n\nOriginal error: "+err.Error(),
				)
				return
			}
			resp.Diagnostics.AddError("Error updating WorkMail user attributes", err.Error())
			return
		}
	}
	if !data.Password.IsNull() {
		_, err := client.ResetPassword(ctx, &workmail.ResetPasswordInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			UserId:         aws.String(data.ID.ValueString()),
			Password:       aws.String(data.Password.ValueString()),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error resetting WorkMail user password", err.Error())
			return
		}
	}
	// Enable or disable user if needed
	enabled := true
	if !data.Enabled.IsNull() {
		enabled = data.Enabled.ValueBool()
	}
	if enabled {
		_, _ = client.RegisterToWorkMail(ctx, &workmail.RegisterToWorkMailInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			EntityId:       aws.String(data.ID.ValueString()),
			Email:          aws.String(data.Email.ValueString()),
		})
	} else {
		_, _ = client.DeregisterFromWorkMail(ctx, &workmail.DeregisterFromWorkMailInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			EntityId:       aws.String(data.ID.ValueString()),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("AWS Config Error", err.Error())
		return
	}
	client := workmail.NewFromConfig(cfg)
	_, err = client.DeleteUser(ctx, &workmail.DeleteUserInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		UserId:         aws.String(data.ID.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting WorkMail user", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected import ID format: <organization_id>,<user_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
