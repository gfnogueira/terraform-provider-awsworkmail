package awsworkmail

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// awsworkmail_group resource: manages a group in a WorkMail organization

type groupResource struct{}

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	OrganizationID types.String   `tfsdk:"organization_id"`
	Name           types.String   `tfsdk:"name"`
	Email          types.String   `tfsdk:"email"`
	Members        []types.String `tfsdk:"members"`
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the WorkMail group",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the WorkMail organization",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Group name (login name)",
			},
			"email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Primary email address for the group (optional, can be set after creation)",
			},
			"members": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of user IDs to be members of the group.",
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data groupResourceModel
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

	input := &workmail.CreateGroupInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		Name:           aws.String(data.Name.ValueString()),
	}
	out, err := client.CreateGroup(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating WorkMail group", err.Error())
		return
	}
	data.ID = types.StringValue(*out.GroupId)

	// Add members if provided
	for _, member := range data.Members {
		// Wait for user to be ENABLED
		userEnabled := false
		for i := 0; i < 30; i++ { // up to ~5 minutes
			desc, err := client.DescribeUser(ctx, &workmail.DescribeUserInput{
				OrganizationId: aws.String(data.OrganizationID.ValueString()),
				UserId:         aws.String(member.ValueString()),
			})
			if err == nil && string(desc.State) == "ENABLED" {
				userEnabled = true
				break
			}
			time.Sleep(10 * time.Second)
		}
		if !userEnabled {
			resp.Diagnostics.AddWarning("User not enabled", "User "+member.ValueString()+" was not enabled after 5 minutes. Group membership may fail.")
		}
		_, err := client.AssociateMemberToGroup(ctx, &workmail.AssociateMemberToGroupInput{
			OrganizationId: aws.String(data.OrganizationID.ValueString()),
			GroupId:        out.GroupId,
			MemberId:       aws.String(member.ValueString()),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error adding member to group", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data groupResourceModel
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

	input := &workmail.DescribeGroupInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		GroupId:        aws.String(data.ID.ValueString()),
	}
	out, err := client.DescribeGroup(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error reading WorkMail group", err.Error())
		return
	}
	if out != nil && out.GroupId != nil {
		data.ID = types.StringValue(*out.GroupId)
	}
	if out != nil && out.Name != nil {
		data.Name = types.StringValue(*out.Name)
	}
	if out != nil && out.Email != nil {
		data.Email = types.StringValue(*out.Email)
	}

	// Read group members
	members := []types.String{}
	listInput := &workmail.ListGroupMembersInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		GroupId:        aws.String(data.ID.ValueString()),
	}
	listOut, err := client.ListGroupMembers(ctx, listInput)
	if err == nil {
		for _, m := range listOut.Members {
			if m.Id != nil {
				members = append(members, types.StringValue(*m.Id))
			}
		}
	}
	data.Members = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data groupResourceModel
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

	// Update group name is not supported by AWS, so only manage members
	// Get current members
	currentMembers := map[string]struct{}{}
	listOut, err := client.ListGroupMembers(ctx, &workmail.ListGroupMembersInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		GroupId:        aws.String(data.ID.ValueString()),
	})
	if err == nil {
		for _, m := range listOut.Members {
			if m.Id != nil {
				currentMembers[*m.Id] = struct{}{}
			}
		}
	}
	// Build desired members set
	desiredMembers := map[string]struct{}{}
	for _, m := range data.Members {
		desiredMembers[m.ValueString()] = struct{}{}
	}
	// Add new members
	for m := range desiredMembers {
		if _, exists := currentMembers[m]; !exists {
			// Wait for user to be ENABLED
			userEnabled := false
			for i := 0; i < 30; i++ {
				desc, err := client.DescribeUser(ctx, &workmail.DescribeUserInput{
					OrganizationId: aws.String(data.OrganizationID.ValueString()),
					UserId:         aws.String(m),
				})
				if err == nil && string(desc.State) == "ENABLED" {
					userEnabled = true
					break
				}
				time.Sleep(10 * time.Second)
			}
			if !userEnabled {
				resp.Diagnostics.AddWarning("User not enabled", "User "+m+" was not enabled after 5 minutes. Group membership may fail.")
			}
			_, err := client.AssociateMemberToGroup(ctx, &workmail.AssociateMemberToGroupInput{
				OrganizationId: aws.String(data.OrganizationID.ValueString()),
				GroupId:        aws.String(data.ID.ValueString()),
				MemberId:       aws.String(m),
			})
			if err != nil {
				resp.Diagnostics.AddError("Error adding member to group", err.Error())
				return
			}
		}
	}
	// Remove old members
	for m := range currentMembers {
		if _, exists := desiredMembers[m]; !exists {
			_, err := client.DisassociateMemberFromGroup(ctx, &workmail.DisassociateMemberFromGroupInput{
				OrganizationId: aws.String(data.OrganizationID.ValueString()),
				GroupId:        aws.String(data.ID.ValueString()),
				MemberId:       aws.String(m),
			})
			if err != nil {
				resp.Diagnostics.AddError("Error removing member from group", err.Error())
				return
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data groupResourceModel
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
	_, err = client.DeleteGroup(ctx, &workmail.DeleteGroupInput{
		OrganizationId: aws.String(data.OrganizationID.ValueString()),
		GroupId:        aws.String(data.ID.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting WorkMail group", err.Error())
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected import ID format: <organization_id>,<group_id>",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
