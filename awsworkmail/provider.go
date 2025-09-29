// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsworkmail

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure AwsWorkMailProvider satisfies various provider interfaces.
var _ provider.Provider = &AwsWorkMailProvider{}

// AwsWorkMailProvider defines the provider implementation.
type AwsWorkMailProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AwsWorkMailProviderModel describes the provider data model.
type AwsWorkMailProviderModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	Region     types.String `tfsdk:"region"`
	AssumeRole types.Object `tfsdk:"assume_role"`
}

// AssumeRoleModel describes the assume_role configuration.
type AssumeRoleModel struct {
	RoleArn         types.String `tfsdk:"role_arn"`
	SessionName     types.String `tfsdk:"session_name"`
	ExternalId      types.String `tfsdk:"external_id"`
	DurationSeconds types.Int64  `tfsdk:"duration_seconds"`
}

func (p *AwsWorkMailProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awsworkmail"
	resp.Version = p.version
}

func (p *AwsWorkMailProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = pschema.Schema{
		Attributes: map[string]pschema.Attribute{
			"endpoint": pschema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"region": pschema.StringAttribute{
				MarkdownDescription: "AWS region for WorkMail operations. If not specified, uses the standard AWS SDK configuration (environment variables, ~/.aws/config, etc.). WorkMail is only available in select regions.",
				Optional:            true,
			},
			"assume_role": pschema.SingleNestedAttribute{
				MarkdownDescription: "Configuration block for assuming an IAM role. Useful for cross-account access or when Terraform state is stored in a different account than WorkMail resources.",
				Optional:            true,
				Attributes: map[string]pschema.Attribute{
					"role_arn": pschema.StringAttribute{
						MarkdownDescription: "ARN of the IAM role to assume.",
						Required:            true,
					},
					"session_name": pschema.StringAttribute{
						MarkdownDescription: "Session name for the assumed role session. If not specified, generates a default name.",
						Optional:            true,
					},
					"external_id": pschema.StringAttribute{
						MarkdownDescription: "External ID to use when assuming the role.",
						Optional:            true,
					},
					"duration_seconds": pschema.Int64Attribute{
						MarkdownDescription: "Duration of the assumed role session in seconds. Must be between 900 (15 minutes) and 43200 (12 hours). If not specified, uses the role's default.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (p *AwsWorkMailProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AwsWorkMailProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create initial AWS config with optional region override
	var cfg aws.Config
	var err error

	if !data.Region.IsNull() && data.Region.ValueString() != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(data.Region.ValueString()))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError("AWS Configuration Error", "Failed to load AWS configuration: "+err.Error())
		return
	}

	// Handle assume_role if configured
	if !data.AssumeRole.IsNull() {
		var assumeRoleConfig AssumeRoleModel
		resp.Diagnostics.Append(data.AssumeRole.As(ctx, &assumeRoleConfig, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		cfg, err = p.assumeRole(ctx, cfg, assumeRoleConfig)
		if err != nil {
			resp.Diagnostics.AddError("AWS Assume Role Error", "Failed to assume role: "+err.Error())
			return
		}
	}

	// Pass AWS config to resources and data sources
	resp.DataSourceData = cfg
	resp.ResourceData = cfg
}

// assumeRole handles the STS AssumeRole operation
func (p *AwsWorkMailProvider) assumeRole(ctx context.Context, baseCfg aws.Config, assumeRoleConfig AssumeRoleModel) (aws.Config, error) {
	stsClient := sts.NewFromConfig(baseCfg)

	sessionName := "terraform-awsworkmail"
	if !assumeRoleConfig.SessionName.IsNull() && assumeRoleConfig.SessionName.ValueString() != "" {
		sessionName = assumeRoleConfig.SessionName.ValueString()
	}

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(assumeRoleConfig.RoleArn.ValueString()),
		RoleSessionName: aws.String(sessionName),
	}

	if !assumeRoleConfig.ExternalId.IsNull() && assumeRoleConfig.ExternalId.ValueString() != "" {
		input.ExternalId = aws.String(assumeRoleConfig.ExternalId.ValueString())
	}

	if !assumeRoleConfig.DurationSeconds.IsNull() && assumeRoleConfig.DurationSeconds.ValueInt64() > 0 {
		duration := assumeRoleConfig.DurationSeconds.ValueInt64()
		if duration < 900 || duration > 43200 {
			return aws.Config{}, fmt.Errorf("duration_seconds must be between 900 (15 minutes) and 43200 (12 hours), got %d", duration)
		}
		input.DurationSeconds = aws.Int32(int32(duration))
	}

	result, err := stsClient.AssumeRole(ctx, input)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to assume role %s: %w", assumeRoleConfig.RoleArn.ValueString(), err)
	}

	assumedCfg := baseCfg.Copy()
	assumedCfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		*result.Credentials.AccessKeyId,
		*result.Credentials.SecretAccessKey,
		*result.Credentials.SessionToken,
	))

	return assumedCfg, nil
}

func (p *AwsWorkMailProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewUserResource,
		NewGroupResource,
		NewDomainResource,
	}
}

func (p *AwsWorkMailProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource, // Register the new user data source
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AwsWorkMailProvider{
			version: version,
		}
	}
}
