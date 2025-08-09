// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsworkmail

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Endpoint types.String `tfsdk:"endpoint"`
	Region   types.String `tfsdk:"region"`
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
		},
	}
}

func (p *AwsWorkMailProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AwsWorkMailProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create AWS config with optional region override
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

	// Pass AWS config to resources and data sources
	resp.DataSourceData = cfg
	resp.ResourceData = cfg
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
