package redmine

import (
	"context"
	"fmt"
	"time"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type myAccountRow struct {
	ID              int64
	Login           string
	Admin           bool
	FirstName       string
	LastName        string
	Mail            string
	CreatedOn       *time.Time
	LastLoginOn     *time.Time
	PasswdChangedOn *time.Time
	TwofaScheme     *string
	APIKey          *string
	Status          int64
	CustomFields    []rm.CustomFieldGetObject
	Groups          *[]rm.IDName
	Memberships     *[]rm.UserMembershipObject
}

func tableRedmineMyAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_my_account",
		Description: "The currently authenticated user (API key owner).",
		List: &plugin.ListConfig{
			Hydrate: getMyAccount,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The user ID."},
			{Name: "admin", Type: proto.ColumnType_BOOL, Description: "Whether the user has admin privileges."},
			{Name: "api_key", Type: proto.ColumnType_STRING, Description: "The user's API key."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user was created."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom field values."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "The user's first name."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups the user belongs to."},
			{Name: "last_login_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user last logged in."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "The user's last name."},
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The user's login name."},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The user's email address."},
			{Name: "memberships", Type: proto.ColumnType_JSON, Description: "Project memberships."},
			{Name: "passwd_changed_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the password was last changed."},
			{Name: "status", Type: proto.ColumnType_INT, Description: "User status (0=anonymous, 1=active, 2=registered, 3=locked)."},
			{Name: "twofa_scheme", Type: proto.ColumnType_STRING, Description: "Two-factor authentication scheme."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Login")},
		},
	}
}

//// HELPER FUNCTIONS

func myAccountRowFromObject(u rm.UserObject) myAccountRow {
	row := myAccountRow{
		ID:              u.ID,
		Login:           u.Login,
		Admin:           u.Admin,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Mail:            u.Mail,
		CreatedOn:       parseRedmineTime(u.CreatedOn),
		LastLoginOn:     parseRedmineTime(u.LastLoginOn),
		PasswdChangedOn: parseRedmineTime(u.PasswdChangedOn),
		TwofaScheme:     u.TwofaScheme,
		APIKey:          u.APIKey,
		CustomFields:    u.CustomFields,
		Groups:          u.Groups,
		Memberships:     u.Memberships,
	}

	if u.Status != nil {
		row.Status = int64(*u.Status)
	}

	return row
}

//// HYDRATE FUNCTIONS

func getMyAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	user, _, err := client.UserCurrentGet(rm.UserCurrentGetRequest{
		Includes: []rm.UserInclude{
			rm.UserIncludeGroups,
			rm.UserIncludeMemberships,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	d.StreamListItem(ctx, myAccountRowFromObject(user))

	return nil, nil
}
