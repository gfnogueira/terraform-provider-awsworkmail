package awsworkmail

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDomainResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "awsworkmail_organization" "test" {
						alias = "terraform-test-org"
					}

					resource "awsworkmail_domain" "test" {
						organization_id = awsworkmail_organization.test.id
						domain         = "test-domain.example.com"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("awsworkmail_domain.test", "domain", "test-domain.example.com"),
					resource.TestCheckResourceAttrSet("awsworkmail_domain.test", "id"),
					resource.TestCheckResourceAttrSet("awsworkmail_domain.test", "organization_id"),
				),
			},
		},
	})
}
