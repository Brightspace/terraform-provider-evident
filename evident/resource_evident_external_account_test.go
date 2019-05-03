package evident

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestEvidentExternalAccountBasic(t *testing.T) {
	updateState("")
	resource.Test(t, resource.TestCase{
		/*
		* we might need to add precheck in order to make sure environment
		* variables are added . skipped for now.
		 */
		PreCheck:     func() {},
		Providers:    testEvidentProviders,
		CheckDestroy: testEvidentExternalAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEvidentExternalAccountBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testEvidentExternalAccountExists("evident_external_account.test_account"),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "arn", fakeArn),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "external_id", fakeExternalId),
				),
			},
		},
	})
}

func TestEvidentExternalAccountUpdate(t *testing.T) {
	updateState("")
	resource.Test(t, resource.TestCase{
		/*
		* we might need to add precheck in order to make sure environment
		* variables are added . skipped for now.
		 */
		PreCheck:     func() {},
		Providers:    testEvidentProviders,
		CheckDestroy: testEvidentExternalAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEvidentExternalAccountBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testEvidentExternalAccountExists("evident_external_account.test_account"),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "arn", fakeArn),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "external_id", fakeExternalId),
				),
			},
			{
				Config: testEvidentExternalAccountUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testEvidentExternalAccountExists("evident_external_account.test_account"),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "arn", updatedFakeArn),
					resource.TestCheckResourceAttr(
						"evident_external_account.test_account", "external_id", updatedFakeExternalId),
				),
			},
		},
	})
}

func testEvidentExternalAccountBasicConfig() string {
	return fmt.Sprintf(`
		resource "evident_external_account" "test_account" {
			arn        = "%s"
			name        = "%s"
			external_id = "%s"
			team_id     = "%s"
		}
	`, fakeArn, fakeName, fakeExternalId, fakeTeamId)
}

func testEvidentExternalAccountUpdateConfig() string {
	return fmt.Sprintf(`
		resource "evident_external_account" "test_account" {
			arn        = "%s"
			name        = "%s"
			external_id = "%s"
			team_id     = "%s"
		}
	`, updatedFakeArn, fakeName, updatedFakeExternalId, updatedFakeTeamId)
}

func testEvidentExternalAccountExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		account := rs.Primary.ID
		config := testEvidentProvider.Meta().(*Config)
		client := config.EvidentClient
		_, err := client.get(account)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testEvidentExternalAccountDestroy(s *terraform.State) error {
	config := testEvidentProvider.Meta().(*Config)
	client := config.EvidentClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "evident_external_account" {
			continue
		}

		_, err := client.get(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("External account still exists")
		}
	}

	return nil
}
