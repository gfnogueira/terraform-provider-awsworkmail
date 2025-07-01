// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsworkmail

import (
	"os"
	"testing"

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
