package entities

import (
	"net/http"

	"time"

	pgx "github.com/jackc/pgx/v4"

	"kcartlidge/api/repo/support"
)

// AccountSetting is for table `account_setting` ("Account Setting").
// Settings for a particular account.
type AccountSetting struct {
	// Id is for column `id`.
	// It's a PRIMARY KEY.
	// It's filterable/sortable.
	//
	// Default: nextval('account_setting_id_seq'::regclass)
	//
	// The unique ID for this setting for this account.
	Id int64 `sql:"id" json:"id" display:"Id" slug:"id"`

	// AccountId is for column `account_id` ("Account Id").
	//
	// The account this setting's value applies to.
	AccountId int64 `sql:"account_id" json:"accountId" display:"Account Id" slug:"account-id"`

	// SettingId is for column `setting_id` ("Setting Id").
	//
	// The ID of the setting which has this value.
	SettingId int64 `sql:"setting_id" json:"settingId" display:"Setting Id" slug:"setting-id"`

	// Value is for column `value`.
	// It has a maximum size of 250.
	//
	// Default: ''::character varying
	//
	// The current value for this account setting.
	Value string `sql:"value" json:"value" display:"Value" slug:"value"`

	// UpdatedAt is for column `updated_at` ("Updated At").
	//
	// Default: now()
	//
	// When the value was last updated.
	UpdatedAt *time.Time `sql:"updated_at" json:"updatedAt" display:"Updated At" slug:"updated-at"`
}

// NewAccountSetting gets a new AccountSetting.
func NewAccountSetting() *AccountSetting {
	d := AccountSetting{}
	return &d
}

// NewAccountSettingFromRows gets a new 'account_setting' row as a AccountSetting.
// You must have already called '.Next()' on the rows.
func NewAccountSettingFromRows(rows pgx.Rows) (*AccountSetting, error) {
	d := AccountSetting{}
	err := rows.Scan(&d.Id, &d.AccountId, &d.SettingId, &d.Value, &d.UpdatedAt)
	return &d, err
}

// NewAccountSettingFromPOST extracts a Account Setting from a POST request.
//
// This is intended mainly for CRUD admin screens.
// It can obviously be used for any purpose, but there are over-posting risks.
//
// Fields are expected to be named as per their slug.
// For example, 'CreatedAt' should be 'created-at' in the POST form data.
//
// The returned `[]string` contains any validation issues.
// These are distinct from actual errors.
func NewAccountSettingFromPOST(r *http.Request) (*AccountSetting, []string, []error) {
	issues := []string{}
	errs := []error{}
	d := AccountSetting{}
	err := r.ParseForm()
	if err != nil {
		errs = append(errs, err)
	} else {
		d.Id = support.Int64FromPOST(r, "id", errs)
		d.AccountId = support.Int64FromPOST(r, "account-id", errs)
		d.SettingId = support.Int64FromPOST(r, "setting-id", errs)
		d.Value = support.StringFromPOST(r, "value", errs)
		d.UpdatedAt = support.DateTimeFromPOST(r, "updated-at", errs)
		issues = d.Validate()
	}
	return &d, issues, errs
}

// Validate performs basic validation on this AccountSetting.
func (item *AccountSetting) Validate() (issues []string) {
	if len(item.Value) > 250 {
		issues = append(issues, "Value cannot be longer than 250.")
	}
	return issues
}
