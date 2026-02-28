package redmine

import (
	"context"
	"fmt"

	rm "github.com/nixys/nxs-go-redmine/v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

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

	d.StreamListItem(ctx, userRowFromObject(user))

	return nil, nil
}
