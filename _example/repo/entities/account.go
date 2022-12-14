/*
Code generated by Near Gothic. DO NOT EDIT.
Generated on 2022/10/19.
Manual edits may be lost when next regenerated.

Near Gothic is (C) K Cartlidge, 2022.
No rights are asserted over generated code.
*/

package entities

import (
	"net/http"

	"time"

	pgx "github.com/jackc/pgx/v4"

	"kcartlidge/api/repo/support"
)

// Account is for table `account`.
// A table of user accounts.
type Account struct {
	// Id is for column `id`.
	// It's a PRIMARY KEY.
	// It's filterable/sortable.
	//
	// Default: nextval('account_id_seq'::regclass)
	//
	// The unique account ID.
	Id int64 `sql:"id" json:"id" display:"Id" slug:"id"`

	// EmailAddress is for column `email_address` ("Email Address").
	// It's filterable/sortable.
	// It has a maximum size of 250.
	//
	// The account-holder's contact email address.
	EmailAddress string `sql:"email_address" json:"emailAddress" display:"Email Address" slug:"email-address"`

	// DisplayName is for column `display_name` ("Display Name").
	// It has a maximum size of 50.
	//
	// The account-holder's display name.
	DisplayName string `sql:"display_name" json:"displayName" display:"Display Name" slug:"display-name"`

	// CreatedAt is for column `created_at` ("Created At").
	//
	// When the account was created.
	CreatedAt *time.Time `sql:"created_at" json:"createdAt" display:"Created At" slug:"created-at"`

	// UpdatedAt is for column `updated_at` ("Updated At").
	//
	// Default: now()
	//
	// When the account details were last updated.
	UpdatedAt *time.Time `sql:"updated_at" json:"updatedAt" display:"Updated At" slug:"updated-at"`

	// DeletedAt is for column `deleted_at` ("Deleted At").
	//
	// When (if) the account was deleted.
	DeletedAt *time.Time `sql:"deleted_at" json:"deletedAt" display:"Deleted At" slug:"deleted-at"`
}

// NewAccount gets a new Account.
func NewAccount() *Account {
	d := Account{}
	return &d
}

// NewAccountFromRows gets a new 'account' row as a Account.
// You must have already called '.Next()' on the rows.
func NewAccountFromRows(rows pgx.Rows) (*Account, error) {
	d := Account{}
	err := rows.Scan(&d.Id, &d.EmailAddress, &d.DisplayName, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	return &d, err
}

// NewAccountFromPOST extracts a Account from a POST request.
//
// This is intended mainly for CRUD admin screens.
// It can obviously be used for any purpose, but there are over-posting risks.
//
// Fields are expected to be named as per their slug.
// For example, 'CreatedAt' should be 'created-at' in the POST form data.
//
// The returned `[]string` contains any validation issues.
// These are distinct from actual errors.
func NewAccountFromPOST(r *http.Request) (*Account, []string, []error) {
	issues := []string{}
	errs := []error{}
	d := Account{}
	err := r.ParseForm()
	if err != nil {
		errs = append(errs, err)
	} else {
		d.Id = support.Int64FromPOST(r, "id", errs)
		d.EmailAddress = support.StringFromPOST(r, "email-address", errs)
		d.DisplayName = support.StringFromPOST(r, "display-name", errs)
		d.CreatedAt = support.DateTimeFromPOST(r, "created-at", errs)
		d.UpdatedAt = support.DateTimeFromPOST(r, "updated-at", errs)
		d.DeletedAt = support.DateTimeFromPOST(r, "deleted-at", errs)
		issues = d.Validate()
	}
	return &d, issues, errs
}

// Validate performs basic validation on this Account.
func (item *Account) Validate() (issues []string) {
	if len(item.EmailAddress) > 250 {
		issues = append(issues, "Email Address cannot be longer than 250.")
	}
	if len(item.DisplayName) > 50 {
		issues = append(issues, "Display Name cannot be longer than 50.")
	}
	return issues
}
