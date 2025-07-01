package awsworkmail

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workmail"
)

// Helper for test cleanup: delete organization by alias
func testAccDeleteOrganizationByAlias(alias string) error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	client := workmail.NewFromConfig(cfg)
	// List organizations and delete the one with the matching alias
	orgs, err := client.ListOrganizations(ctx, &workmail.ListOrganizationsInput{})
	if err != nil {
		return err
	}
	for _, org := range orgs.OrganizationSummaries {
		if org.Alias != nil && *org.Alias == alias {
			_, err := client.DeleteOrganization(ctx, &workmail.DeleteOrganizationInput{
				OrganizationId:  org.OrganizationId,
				DeleteDirectory: true,
			})
			if err != nil {
				return fmt.Errorf("failed to delete org %s: %w", alias, err)
			}
		}
	}
	return nil
}
