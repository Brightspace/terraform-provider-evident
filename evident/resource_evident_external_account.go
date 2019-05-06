package evident

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceExternalAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceExternalAccountCreate,
		Read:   resourceExternalAccountRead,
		Update: resourceExternalAccountUpdate,
		Delete: resourceExternalAccountDelete,

		Schema: map[string]*schema.Schema{
			"arn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Amazon Resource Name for the IAM role",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},
			"external_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "External Identifier set on the role",
			},
			"team_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team the external account belongs to",
			},
		},
	}
}

func resourceExternalAccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	name := d.Get("name").(string)
	arn := d.Get("arn").(string)
	externalID := d.Get("external_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[DEBUG] external_account set: (ARN: %q, Name: %q, ExternalID: %q)", arn, name, externalID)
	account, err := client.add(name, arn, externalID, teamID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] external_account added: (Name: %q, ID: %q)", name, account.ID)
	d.SetId(account.GetIdString())

	time.Sleep(5 * time.Second)

	return resourceExternalAccountRead(d, meta)
}

func resourceExternalAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	log.Printf("[DEBUG] external_account get: (ID: %q)", d.Id())
	account, err := client.get(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	log.Printf("[DEBUG] external_account read: (ARN: %q, Name: %q, ExternalID: %q)", account.Attributes.Arn, account.Attributes.Name, account.Attributes.ExternalID)
	d.Set("name", account.Attributes.Name)
	d.Set("arn", account.Attributes.Arn)
	d.Set("external_id", account.Attributes.ExternalID)

	return nil
}

func resourceExternalAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	name := d.Get("name").(string)
	arn := d.Get("arn").(string)
	externalID := d.Get("external_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[DEBUG] external_account set: (ARN: %q, Name: %q, ExternalID: %q)", arn, name, externalID)
	account, err := client.update(d.Id(), name, arn, externalID, teamID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] external_account updated: (Name: %q, ID: %q)", name, account.ID)
	d.SetId(account.GetIdString())

	time.Sleep(5 * time.Second)

	return resourceExternalAccountRead(d, meta)
}

func resourceExternalAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	log.Printf("[DEBUG] external_account delete: (ID: %q)", d.Id())
	_, err := client.delete(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
