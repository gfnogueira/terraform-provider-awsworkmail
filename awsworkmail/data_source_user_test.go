package awsworkmail

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	orgID := os.Getenv("WORKMAIL_ORGANIZATION_ID")
	userID := os.Getenv("WORKMAIL_USER_ID")
	if orgID == "" || userID == "" {
		t.Skip("WORKMAIL_ORGANIZATION_ID and WORKMAIL_USER_ID must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: testAccDataSourceUserConfig(orgID, userID),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.awsworkmail_user.test", "name"),
				resource.TestCheckResourceAttrSet("data.awsworkmail_user.test", "email"),
				resource.TestCheckResourceAttrSet("data.awsworkmail_user.test", "state"),
			),
		}},
	})
}

func testAccDataSourceUserConfig(orgID, userID string) string {
	return `
data "awsworkmail_user" "test" {
  organization_id = "` + orgID + `"
  user_id        = "` + userID + `"
}
`
}
