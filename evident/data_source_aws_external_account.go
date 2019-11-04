package evident

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAwsExternalAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsExternalAccountRead,

		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the external account",
			},
			"arn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Amazon Resource Name for the IAM role",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name",
			},
			"external_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External Identifier set on the role",
			},
		},
	}
}

func dataSourceAwsExternalAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient
	id := d.Get("account").(string)

	log.Printf("[DEBUG] external_account_aws get: (ID: %q)", id)
	account, err := client.Get(id)
	if err != nil {
		log.Printf("[DEBUG] external_account_aws not found: (ID: %q)", id)
		return err
	}

	log.Printf("[DEBUG] external_account_aws read: (ARN: %q, Name: %q, ExternalID: %q)", account.Attributes.Arn, account.Attributes.Name, account.Attributes.ExternalID)
	d.SetId(id)
	d.Set("name", account.Attributes.Name)
	d.Set("arn", account.Attributes.Arn)
	d.Set("external_id", account.Attributes.ExternalID)

	return nil
}
