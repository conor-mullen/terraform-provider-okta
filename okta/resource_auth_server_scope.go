package okta

import (
	"net/http"

	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAuthServerScope() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAuthServerScopeCreate,
		Exists:   resourceAuthServerScopeExists,
		Read:     resourceAuthServerScopeRead,
		Update:   resourceAuthServerScopeUpdate,
		Delete:   resourceAuthServerScopeDelete,
		Importer: createNestedResourceImporter([]string{"auth_server_id", "id"}),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server scope name",
			},
			"auth_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"consent": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IMPLICIT",
				Description: "EA Feature and thus it is simply ignored if the feature is off",
			},
			"metadata_publish": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALL_CLIENTS",
				Description:  "Whether to publish metadata or not, matching API type despite the fact it could just be a boolean",
				ValidateFunc: validation.StringInSlice([]string{"ALL_CLIENTS", "NO_CLIENTS"}, false),
			},
		},
	}
}

func buildAuthServerScope(d *schema.ResourceData) *sdk.AuthorizationServerScope {
	return &sdk.AuthorizationServerScope{
		Consent:         d.Get("consent").(string),
		Description:     d.Get("description").(string),
		MetadataPublish: d.Get("metadata_publish").(string),
		Name:            d.Get("name").(string),
	}
}

func resourceAuthServerScopeCreate(d *schema.ResourceData, m interface{}) error {
	authServerScope := buildAuthServerScope(d)
	c := getSupplementFromMetadata(m)
	responseAuthServerScope, _, err := c.CreateAuthorizationServerScope(d.Get("auth_server_id").(string), *authServerScope, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServerScope.Id)

	return resourceAuthServerScopeRead(d, m)
}

func resourceAuthServerScopeExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServerScope(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerScopeRead(d *schema.ResourceData, m interface{}) error {
	authServerScope, err := fetchAuthServerScope(d, m)
	if err != nil {
		return err
	}

	d.Set("name", authServerScope.Name)
	d.Set("description", authServerScope.Description)
	d.Set("metadata_publish", authServerScope.MetadataPublish)

	if authServerScope.Consent != "" {
		d.Set("consent", authServerScope.Consent)
	}

	return nil
}

func resourceAuthServerScopeUpdate(d *schema.ResourceData, m interface{}) error {
	authServerScope := buildAuthServerScope(d)
	c := getSupplementFromMetadata(m)
	_, _, err := c.UpdateAuthorizationServerScope(d.Get("auth_server_id").(string), d.Id(), *authServerScope, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerScopeRead(d, m)
}

func resourceAuthServerScopeDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerScope(d.Get("auth_server_id").(string), d.Id())

	return err
}

func fetchAuthServerScope(d *schema.ResourceData, m interface{}) (*sdk.AuthorizationServerScope, error) {
	c := getSupplementFromMetadata(m)
	auth, resp, err := c.GetAuthorizationServerScope(d.Get("auth_server_id").(string), d.Id(), sdk.AuthorizationServerScope{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}
