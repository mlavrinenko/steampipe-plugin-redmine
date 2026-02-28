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

type userRow struct {
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
	Status          int64
	CustomFields    []rm.CustomFieldGetObject
	Groups          *[]rm.IDName
	Memberships     *[]rm.UserMembershipObject
}

func tableRedmineUser() *plugin.Table {
	return &plugin.Table{
		Name:        "redmine_user",
		Description: "Users in the Redmine instance. Listing users requires admin privileges.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		Columns: []*plugin.Column{
			// Key columns first
			{Name: "id", Type: proto.ColumnType_INT, Description: "The user ID."},
			{Name: "status", Type: proto.ColumnType_INT, Description: "User status (0=anonymous, 1=active, 2=registered, 3=locked)."},
			// Remaining columns alphabetically
			{Name: "admin", Type: proto.ColumnType_BOOL, Description: "Whether the user has admin privileges."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user was created."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom field values."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "The user's first name."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups the user belongs to (populated only on single-user get)."},
			{Name: "last_login_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user last logged in."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "The user's last name."},
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The user's login name."},
			{Name: "mail", Type: proto.ColumnType_STRING, Description: "The user's email address."},
			{Name: "memberships", Type: proto.ColumnType_JSON, Description: "Project memberships (populated only on single-user get)."},
			{Name: "passwd_changed_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the password was last changed."},
			{Name: "twofa_scheme", Type: proto.ColumnType_STRING, Description: "Two-factor authentication scheme."},
			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The display name for this resource.", Transform: transform.FromField("Login")},
		},
	}
}

//// HELPER FUNCTIONS

func userRowFromObject(u rm.UserObject) userRow {
	row := userRow{
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

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	userID := d.EqualsQuals["id"].GetInt64Value()

	user, _, err := client.UserSingleGet(userID, rm.UserSingleGetRequest{
		Includes: []rm.UserInclude{
			rm.UserIncludeGroups,
			rm.UserIncludeMemberships,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	return userRowFromObject(user), nil
}

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var filters *rm.UserGetRequestFilters
	if d.EqualsQuals["status"] != nil {
		filters = rm.UserGetRequestFiltersInit().
			StatusSet(rm.UserStatus(d.EqualsQuals["status"].GetInt64Value()))
	}

	var offset int64
	var pageSize int64 = 100
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < pageSize {
		pageSize = *d.QueryContext.Limit
	}

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		result, _, err := client.UserMultiGet(rm.UserMultiGetRequest{
			Filters: filters,
			Offset:  offset,
			Limit:   pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %w", err)
		}

		for _, user := range result.Users {
			d.StreamListItem(ctx, userRowFromObject(user))

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if int64(len(result.Users)) < pageSize {
			break
		}
		offset += pageSize
	}

	return nil, nil
}
