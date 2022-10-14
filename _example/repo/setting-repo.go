package repo

import (
	pgx "github.com/jackc/pgx/v4"

	"kcartlidge/api/repo/connection"
	"kcartlidge/api/repo/entities"
)

// SettingRepo contains data access methods for Setting items.
//
// Specific methods are added for indexed fields.
// General-purpose methods cover unindexed ones.
type SettingRepo struct {
	repo
}

// ---------- Constructor ----------

// NewSettingRepo creates an instance for database access.
func NewSettingRepo(connection *connection.Connection) *SettingRepo {
	r := SettingRepo{}
	r.connection = connection
	r.ResetConditions()
	r.ResetSorting()
	r.ResetLimitAndOffset()
	return &r
}

// ---------- CRUD methods ----------

// List returns all matching Setting items.
func (r *SettingRepo) List() ([]entities.Setting, error) {
	d := make([]entities.Setting, 0)
	cmd := "SELECT id,display_name,details,max_value_length,is_enabled FROM setting "
	err := r.Execute(cmd, func(rows pgx.Rows) error {
		if dd, err := entities.NewSettingFromRows(rows); err != nil {
			return err
		} else {
			d = append(d, *dd)
			return nil
		}
	})
	return d, err
}

// Insert adds a new Setting item.
func (r *SettingRepo) Insert(item entities.Setting) (int64, error) {
	cmd := "INSERT INTO setting (display_name,details,max_value_length,is_enabled) "
	cmd += "VALUES ($1,$2,$3,$4) "
	var p []interface{}
	p = append(p, item.DisplayName)
	p = append(p, item.Details)
	p = append(p, item.MaxValueLength)
	p = append(p, item.IsEnabled)
	return r.ExecuteNonQuery(cmd, p...)
}

// Update modifies a Setting item (all fields except primary keys, which
// are still required anyway in order to know which items to update).
func (r *SettingRepo) Update(id int64, item entities.Setting) (int64, error) {
	cmd := "UPDATE setting "
	cmd += "SET display_name=$1,details=$2,max_value_length=$3,is_enabled=$4 "
	cmd += "WHERE id=$5 "

	// Values to update
	var p []interface{}
	p = append(p, item.DisplayName)
	p = append(p, item.Details)
	p = append(p, item.MaxValueLength)
	p = append(p, item.IsEnabled)

	// Primary key restrictions
	p = append(p, id)
	return r.ExecuteNonQuery(cmd, p...)
}

// Delete removes a Setting item.
func (r *SettingRepo) Delete(id int64) (int64, error) {
	cmd := "DELETE FROM setting "
	cmd += "WHERE id=$1 "
	var p []interface{}
	p = append(p, id)
	return r.ExecuteNonQuery(cmd, p...)
}

// ---------- Paging ----------

// WithLimit adds a restriction on the Setting item(s) returned.
// Overrides the package's MaxRows value (for this instance only).
func (r *SettingRepo) WithLimit(value int) *SettingRepo {
	r.limit = value
	return r
}

// WithOffset skips the given number of Setting item(s) in the result set.
func (r *SettingRepo) WithOffset(value int) *SettingRepo {
	r.offset = value
	return r
}

// ---------- Typed filtering (only indexed fields) ----------

// WhereId adds a filter for Id.
func (r *SettingRepo) WhereId(operator string, value int64) *SettingRepo {
	return r.Where("id", operator, value)
}

// ----------- Typed ordering (only indexed fields) -----------

// SortById adds sorting by Id.
func (r *SettingRepo) SortById() *SettingRepo {
	return r.AddSorting("id", false)
}

// ReverseById adds reverse sorting by Id.
func (r *SettingRepo) ReverseById() *SettingRepo {
	return r.AddSorting("id", true)
}

// ---------- Untyped filtering and ordering (any fields) ----------

// Where adds a clause to the request.
//
// WARNING:
// Prefer the predefined field-specific Where... functions as they use indexed fields.
// Using this method instead is more flexible but may involve unindexed fields.
// Use carefully/sparingly to avoid performance issues in large data sets.
func (r *SettingRepo) Where(thing string, operator string, value interface{}) *SettingRepo {
	r.addCondition(thing, operator, value)
	return r
}

// AddSorting includes an ad-hoc sort by any valid column/thing.
// Indexed fields have their own SortBy... variants.
func (r *SettingRepo) AddSorting(thing string, descending bool) *SettingRepo {
	r.addOrdering(thing, descending)
	return r
}
