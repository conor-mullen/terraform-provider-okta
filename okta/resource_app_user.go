package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceAppUser() *schema.Resource {
	return &schema.Resource{
		// No point in having an exist function, since only the group has to exist
		Create: resourceAppUserCreate,
		Exists: resourceAppUserExists,
		Read:   resourceAppUserRead,
		Update: resourceAppUserUpdate,
		Delete: resourceAppUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate group with",
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group associated with the application",
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"profile": &schema.Schema{
				Type:      schema.TypeString,
				StateFunc: normalizeDataJSON,
				Computed:  true,
			},
		},
	}
}

func resourceAppUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getOktaClientFromMetadata(m)
	g, _, err := client.Application.GetApplicationUser(
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)

	return g != nil, err
}

func resourceAppUserCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	u, _, err := client.Application.AssignUserToApplication(
		d.Get("app_id").(string),
		*getAppUser(d),
	)

	if err != nil {
		return err
	}

	d.SetId(u.Id)

	return resourceAppUserRead(d, m)
}

func resourceAppUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, _, err := client.Application.UpdateApplicationUser(
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		*getAppUser(d),
	)

	if err != nil {
		return err
	}

	return resourceAppUserRead(d, m)
}

func resourceAppUserRead(d *schema.ResourceData, m interface{}) error {
	u, _, err := getOktaClientFromMetadata(m).Application.GetApplicationUser(
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)

	if err != nil {
		return err
	}

	d.Set("profile", u.Profile)
	d.Set("username", u.Credentials.UserName)

	return nil
}

func resourceAppUserDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationUser(
		d.Get("app_id").(string),
		d.Get("user_id").(string),
		nil,
	)
	return err
}

func getAppUser(d *schema.ResourceData) *okta.AppUser {
	return &okta.AppUser{
		Id: d.Get("user_id").(string),
		Credentials: &okta.AppUserCredentials{
			UserName: d.Get("username").(string),
			Password: &okta.AppUserPasswordCredential{
				Value: d.Get("password").(string),
			},
		},
	}
}
