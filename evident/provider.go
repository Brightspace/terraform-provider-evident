package evident

import (
	"github.com/Brightspace/terraform-provider-evident/evident/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the Evident access key. It must be provided, but it can also be sourced from the `EVIDENT_ACCESS_KEY` environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("EVIDENT_ACCESS_KEY", nil),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the Evident secret key. It must be provided, but it can also be sourced from the `EVIDENT_SECRET_KEY` environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("EVIDENT_SECRET_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"evident_external_account": resourceExternalAccount(),
			"evident_external_account_aws": resourceExternalAccount(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"evident_external_account_aws": dataSourceAwsExternalAccount(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := api.Evident{
		Credentials: api.Credentials{
			AccessKey: []byte(d.Get("access_key").(string)),
			SecretKey: []byte(d.Get("secret_key").(string)),
		},
		RetryMaximum: 10,
	}

	config := Config{
		EvidentClient: client,
	}

	return &config, nil
}
