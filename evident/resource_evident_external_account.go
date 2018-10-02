package evident

import (
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
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
				ForceNew:    true,
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
				ForceNew:    true,
			},
			"evident_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the external account to update an Amazon IAM credential of",
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

	account, err := client.add(name, arn, externalID, teamID)
	if err != nil {
		return err
	}

	d.SetId(account.ID)
	return resourceExternalAccountRead(d, meta)
}

func resourceExternalAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	account, err := client.get(d.Id())
	if err != nil {
		return err
	}

	d.Set("evident_id", account.ID)

	return nil
}

func resourceExternalAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceExternalAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := config.EvidentClient

	_, err := client.delete(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
