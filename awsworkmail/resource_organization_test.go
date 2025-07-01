package awsworkmail

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganization_Basic(t *testing.T) {
	orgAlias := os.Getenv("TF_AWSWORKMAIL_ORG_ALIAS")
	if orgAlias == "" {
		orgAlias = "tfacc-org-example"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationConfig(orgAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("awsworkmail_organization.test", "alias", orgAlias),
					resource.TestCheckResourceAttrSet("awsworkmail_organization.test", "id"),
				),
			},
		},
	})
}

func testAccOrganizationConfig(alias string) string {
	return fmt.Sprintf(`
resource "awsworkmail_organization" "test" {
  alias = "%s"
}
`, alias)
}
