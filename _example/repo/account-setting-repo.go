package repo

import (
	pgx "github.com/jackc/pgx/v4"

	"kcartlidge/api/repo/connection"
	"kcartlidge/api/repo/entities"
)

// AccountSettingRepo contains data access methods for Account Setting items.
//
// Specific methods are added for indexed fields.
// General-purpose methods cover unindexed ones.
type AccountSettingRepo struct {
	repo
}

// ---------- Constructor ----------

// NewAccountSettingRepo creates an instance for database access.
func NewAccountSettingRepo(connection *connection.Connection) *AccountSettingRepo {
	r := AccountSettingRepo{}
	r.connection = connection
	r.ResetConditions()
	r.ResetSorting()
	r.ResetLimitAndOffset()
	return &r
}

// ---------- CRUD methods ----------

// List returns all matching Account Setting items.
func (r *AccountSettingRepo) List() ([]entities.AccountSetting, error) {
	d := make([]entities.AccountSetting, 0)
	cmd := "SELECT id,account_id,setting_id,value,updated_at FROM account_setting "
	err := r.Execute(cmd, func(rows pgx.Rows) error {
		if dd, err := entities.NewAccountSettingFromRows(rows); err != nil {
			return err
		} else {
			d = append(d, *dd)
			return nil
		}
	})
	return d, err
}

// Insert adds a new Account Setting item.
func (r *AccountSettingRepo) Insert(item entities.AccountSetting) (int64, error) {
	cmd := "INSERT INTO account_setting (account_id,setting_id,value,updated_at) "
	cmd += "VALUES ($1,$2,$3,$4) "
	var p []interface{}
	p = append(p, item.AccountId)
	p = append(p, item.SettingId)
	p = append(p, item.Value)
	p = append(p, item.UpdatedAt)
	return r.ExecuteNonQuery(cmd, p...)
}

// Update modifies a Account Setting item (all fields except primary keys, which
// are still required anyway in order to know which items to update).
func (r *AccountSettingRepo) Update(id int64, item entities.AccountSetting) (int64, error) {
	cmd := "UPDATE account_setting "
	cmd += "SET account_id=$1,setting_id=$2,value=$3,updated_at=$4 "
	cmd += "WHERE id=$5 "

	// Values to update
	var p []interface{}
	p = append(p, item.AccountId)
	p = append(p, item.SettingId)
	p = append(p, item.Value)
	p = append(p, item.UpdatedAt)

	// Primary key restrictions
	p = append(p, id)
	return r.ExecuteNonQuery(cmd, p...)
}

// Delete removes a Account Setting item.
func (r *AccountSettingRepo) Delete(id int64) (int64, error) {
	cmd := "DELETE FROM account_setting "
	cmd += "WHERE id=$1 "
	var p []interface{}
	p = append(p, id)
	return r.ExecuteNonQuery(cmd, p...)
}

// ---------- Paging ----------

// WithLimit adds a restriction on the Account Setting item(s) returned.
// Overrides the package's MaxRows value (for this instance only).
func (r *AccountSettingRepo) WithLimit(value int) *AccountSettingRepo {
	r.limit = value
	return r
}

// WithOffset skips the given number of Account Setting item(s) in the result set.
func (r *AccountSettingRepo) WithOffset(value int) *AccountSettingRepo {
	r.offset = value
	return r
}

// ---------- Typed filtering (only indexed fields) ----------

// WhereId adds a filter for Id.
func (r *AccountSettingRepo) WhereId(operator string, value int64) *AccountSettingRepo {
	return r.Where("id", operator, value)
}

// ----------- Typed ordering (only indexed fields) -----------

// SortById adds sorting by Id.
func (r *AccountSettingRepo) SortById() *AccountSettingRepo {
	return r.AddSorting("id", false)
}

// ReverseById adds reverse sorting by Id.
func (r *AccountSettingRepo) ReverseById() *AccountSettingRepo {
	return r.AddSorting("id", true)
}

// ---------- Untyped filtering and ordering (any fields) ----------

// Where adds a clause to the request.
//
// WARNING:
// Prefer the predefined field-specific Where... functions as they use indexed fields.
// Using this method instead is more flexible but may involve unindexed fields.
// Use carefully/sparingly to avoid performance issues in large data sets.
func (r *AccountSettingRepo) Where(thing string, operator string, value interface{}) *AccountSettingRepo {
	r.addCondition(thing, operator, value)
	return r
}

// AddSorting includes an ad-hoc sort by any valid column/thing.
// Indexed fields have their own SortBy... variants.
func (r *AccountSettingRepo) AddSorting(thing string, descending bool) *AccountSettingRepo {
	r.addOrdering(thing, descending)
	return r
}
