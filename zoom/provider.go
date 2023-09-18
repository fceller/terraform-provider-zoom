package zoom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-zoom/client"
	"terraform-provider-zoom/token"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"zoom_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_ACCOUNT_ID", ""),
			},
			"zoom_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_CLIENT_ID", ""),
			},
			"zoom_client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_CLIENT_SECRET", ""),
			},
			"zoom_timeout_minutes": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_TIMEOUT_MINUTES", 2),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zoom_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zoom_user": dataSourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var accountId string
	if v, ok := d.GetOk("zoom_client_id"); ok {
		accountId = v.(string)
	}

	var clientId string
	if v, ok := d.GetOk("zoom_client_id"); ok {
		clientId = v.(string)
	}

	var clientSecret string
	if v, ok := d.GetOk("zoom_client_secret"); ok {
		clientSecret = v.(string)
	}

	var timeoutMinutes int
	if v, ok := d.GetOk("zoom_timeout_minutes"); ok {
		timeoutMinutes = v.(int)
	}

	c := client.NewClient("", timeoutMinutes)
	err := token.GenerateToken(c, accountId, clientId, clientSecret)

	return c, err
}
