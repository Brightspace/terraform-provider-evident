package evident

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The working directory where to run.",
				DefaultFunc: schema.EnvDefaultFunc("EVIDENT_ACCESS_KEY", nil),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Create command",
				DefaultFunc: schema.EnvDefaultFunc("EVIDENT_SECRET_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"evident_external_account": resourceExternalAccount(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := Evident{
		Credentials: Credentials{
			AccessKey: []byte(d.Get("access_key").(string)),
			SecretKey: []byte(d.Get("secret_key").(string)),
		},
		RetryMaximum: 5,
	}
	config := Config{
		EvidentClient: client,
	}

	return &config, nil
}
