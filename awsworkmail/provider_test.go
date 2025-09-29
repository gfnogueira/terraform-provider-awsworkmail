// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsworkmail

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"awsworkmail": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Skip("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Skip("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func TestProviderAssumeRoleValidation(t *testing.T) {
	p := &AwsWorkMailProvider{}
	
	// Test that assume_role schema is properly configured
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}
	
	p.Schema(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Errorf("Provider schema should not have errors: %v", resp.Diagnostics)
	}
	
	// Check that assume_role attribute exists in schema
	assumeRoleAttr, exists := resp.Schema.Attributes["assume_role"]
	if !exists {
		t.Error("assume_role attribute should exist in provider schema")
	}
	
	if !assumeRoleAttr.IsOptional() {
		t.Error("assume_role should be optional")
	}
}
