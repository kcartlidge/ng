package main

import "encoding/json"

type Schema struct {
	SchemaName  string `json:"schemaName"`
	CodeName    string `json:"codeName"`
	DisplayName string `json:"displayName"`
	JsonName    string `json:"jsonName"`
	SlugName    string `json:"slugName"`

	Owner  string  `json:"owner"`
	Tables []Table `json:"tables"`
}

type Table struct {
	SchemaName        string `json:"schemaName"`
	TableName         string `json:"tableName"`
	CodeName          string `json:"codeName"`
	DisplayName       string `json:"displayName"`
	DisplayNamePlural string `json:"displayNamePlural"`
	JsonName          string `json:"jsonName"`
	SlugName          string `json:"slugName"`
	SlugNamePlural    string `json:"slugNamePlural"`

	Owner       string       `json:"owner"`
	Comment     string       `json:"comment"`
	TableType   string       `json:"tableType"`
	IsUpdatable bool         `json:"isUpdatable"`
	Columns     []Column     `json:"columns"`
	Constraints []Constraint `json:"constraints"`
	Indexes     []Index      `json:"indexes"`

	CodeImports []string `json:"codeImports"`
}

type Column struct {
	Position    int    `json:"position"`
	ColumnName  string `json:"columnName"`
	CodeName    string `json:"codeName"`
	DisplayName string `json:"displayName"`
	JsonName    string `json:"jsonName"`
	SlugName    string `json:"slugName"`

	Comment          string  `json:"comment"`
	IsPrimaryKey     bool    `json:"isPrimaryKey"`
	IsNullable       bool    `json:"isNullable"`
	IsCardinal       bool    `json:"isCardinal"`
	HasMaxLen        bool    `json:"hasMaxLen"`
	HasDefault       bool    `json:"hasDefault"`
	HasPrecision     bool    `json:"hasPrecision"`
	CanFilter        bool    `json:"canFilter"`
	SqlType          string  `json:"sqlType"`
	DataType         string  `json:"dataType"`
	MaxLen           *int    `json:"maxLen,omitempty"`
	ColumnDefault    *string `json:"columnDefault,omitempty"`
	NumericPrecision *int    `json:"numericPrecision,omitempty"`
}

type Constraint struct {
	ConstraintName string `json:"constraintName"`
	CodeName       string `json:"codeName"`
	DisplayName    string `json:"displayName"`
	JsonName       string `json:"jsonName"`
	SlugName       string `json:"slugName"`

	IsPrimaryKey   bool     `json:"isPrimaryKey"`
	IsForeignKey   bool     `json:"isForeignKey"`
	IsUniqueKey    bool     `json:"isUniqueKey"`
	ColumnNames    []string `json:"columnNames"`
	ConstraintType string   `json:"constraintType"`
	ForeignTable   *string  `json:"foreignTable,omitempty"`
	ForeignColumn  *string  `json:"foreignColumn,omitempty"`
}

type Index struct {
	IndexName   string `json:"indexName"`
	CodeName    string `json:"codeName"`
	DisplayName string `json:"displayName"`
	JsonName    string `json:"jsonName"`
	SlugName    string `json:"slugName"`

	ColumnNames  []string `json:"columnNames"`
	IsPrimaryKey bool     `json:"isPrimaryKey"`
	IsUnique     bool     `json:"isUnique"`
}

func (schema Schema) ToJSON() []byte {
	b, err := json.MarshalIndent(schema, "", "\t")
	check(err)
	return b
}
