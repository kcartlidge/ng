package entities

import (
	"net/http"

	pgx "github.com/jackc/pgx/v4"

	"kcartlidge/app/data/support"
)

// Setting is for table `setting`.
// The settings available for an account.
type Setting struct {
	// Id is for column `id`.
	// It's a PRIMARY KEY.
	// It's filterable/sortable.
	//
	// Default: nextval('setting_id_seq'::regclass)
	//
	// The unique setting ID.
	Id int64 `sql:"id" json:"id" display:"Id" slug:"id"`

	// DisplayName is for column `display_name` ("Display Name").
	// It has a maximum size of 50.
	//
	// The displayable brief name for this setting.
	DisplayName string `sql:"display_name" json:"displayName" display:"Display Name" slug:"display-name"`

	// Details is for column `details`.
	// It has a maximum size of 500.
	//
	// Descriptive details for this setting.
	Details string `sql:"details" json:"details" display:"Details" slug:"details"`

	// MaxValueLength is for column `max_value_length` ("Max Value Length").
	//
	// Default: 30
	//
	// The longest a value for this setting is allowed to be.
	MaxValueLength int64 `sql:"max_value_length" json:"maxValueLength" display:"Max Value Length" slug:"max-value-length"`

	// IsEnabled is for column `is_enabled` ("Is Enabled").
	IsEnabled bool `sql:"is_enabled" json:"isEnabled" display:"Is Enabled" slug:"is-enabled"`
}

// NewSetting gets a new Setting.
func NewSetting() *Setting {
	d := Setting{}
	return &d
}

// NewSettingFromRows gets a new 'setting' row as a Setting.
// You must have already called '.Next()' on the rows.
func NewSettingFromRows(rows pgx.Rows) (*Setting, error) {
	d := Setting{}
	err := rows.Scan(&d.Id, &d.DisplayName, &d.Details, &d.MaxValueLength, &d.IsEnabled)
	return &d, err
}

// NewSettingFromPOST extracts a Setting from a POST request.
//
// This is intended mainly for CRUD admin screens.
// It can obviously be used for any purpose, but there are over-posting risks.
//
// Fields are expected to be named as per their slug.
// For example, 'CreatedAt' should be 'created-at' in the POST form data.
//
// The returned `[]string` contains any validation issues.
// These are distinct from actual errors.
func NewSettingFromPOST(r *http.Request) (*Setting, []string, []error) {
	issues := []string{}
	errs := []error{}
	d := Setting{}
	err := r.ParseForm()
	if err != nil {
		errs = append(errs, err)
	} else {
		d.Id = support.Int64FromPOST(r, "id", errs)
		d.DisplayName = support.StringFromPOST(r, "display-name", errs)
		d.Details = support.StringFromPOST(r, "details", errs)
		d.MaxValueLength = support.Int64FromPOST(r, "max-value-length", errs)
		d.IsEnabled = support.BoolFromPOST(r, "is-enabled", errs)
		issues = d.Validate()
	}
	return &d, issues, errs
}

// Validate performs basic validation on this Setting.
func (item *Setting) Validate() (issues []string) {
	if len(item.DisplayName) > 50 {
		issues = append(issues, "Display Name cannot be longer than 50.")
	}
	if len(item.Details) > 500 {
		issues = append(issues, "Details cannot be longer than 500.")
	}
	return issues
}
