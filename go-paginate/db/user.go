package db

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
)

var (
	ErrCouldNotOpen   = errors.New("Could not open db file")
	ErrCouldNotCreate = errors.New("Could not create table")
	ErrInsertFailed   = errors.New("Could not insert user")
	ErrRetrieveFailed = errors.New("Could not retrieve users")
	ErrInvalidOrder   = errors.New("Invalid order direction")
)

type OrderBy struct {
	Id   string `json:"id,omitempty"`
	Age  string `json:"age,omitempty"`
	Name string `json:"name,omitempty"`
}

type Filters struct {
	Status    string   `json:"status,omitempty"`
	Countries []string `json:"countries,omitempty"`
	Age       int      `json:"age,omitempty"`
	Degree    string   `json:"degree,omitempty"`
}

type User struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	Country string `json:"country,omitempty"`
	Degree  string `json:"degree,omitempty"`
	Status  string `json:"status,omitempty"`
	Site    string `json:"site,omitempty"`
}

type UserTable struct {
	db *sql.DB
}

func Create() (*UserTable, error) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return nil, ErrCouldNotOpen
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		country TEXT NOT NULL,
		degree TEXT,
		status TEXT,
		site TEXT
	);
	`

	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, ErrCouldNotCreate
	}

	return &UserTable{db: db}, nil
}

func (ut *UserTable) Insert(user User) error {
	insertQuery := `
	INSERT INTO users (name, age, country, degree, status, site)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err := ut.db.Exec(
		insertQuery,
		user.Name,
		user.Age,
		user.Country,
		user.Degree,
		user.Status,
		user.Site,
	)
	if err != nil {
		return ErrInsertFailed
	}
	return nil
}

func (ut *UserTable) Show(users []User) bytes.Buffer {
	buf := bytes.Buffer{}
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"ID", "Name", "Age", "Country", "Degree", "Status", "Site"})

	for _, user := range users {
		row := []string{
			fmt.Sprintf("%d", user.Id),
			user.Name,
			fmt.Sprintf("%d", user.Age),
			user.Country,
			user.Degree,
			user.Status,
			user.Site,
		}
		table.Append(row)
	}

	table.Render()

	return buf
}

func (ut *UserTable) Retrieve(opts ...Option) ([]User, error) {
	q := &queryOptions{}
	for _, opt := range opts {
		opt(q)
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(
		"SELECT id, name, age, country, degree, status, site FROM users WHERE 1=1",
	)

	args := []any{}

	if len(q.Status) > 0 {
		statusPlaceholders := strings.Repeat("?,", len(q.Status))
		queryBuilder.WriteString(
			fmt.Sprintf(" AND status IN (%s)", statusPlaceholders[:len(statusPlaceholders)-1]),
		)
		for _, status := range q.Status {
			args = append(args, status)
		}
	}

	if len(q.Country) > 0 {
		countryPlaceholders := strings.Repeat("?,", len(q.Country))
		queryBuilder.WriteString(
			fmt.Sprintf(" AND country IN (%s)", countryPlaceholders[:len(countryPlaceholders)-1]),
		)
		for _, country := range q.Country {
			args = append(args, country)
		}
	}

	if q.OrderBy != "" && q.Order != "" {
		queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", q.OrderBy, q.Order))
	}

	if q.Limit > 0 {
		queryBuilder.WriteString(" LIMIT ?")
		args = append(args, q.Limit)
	}

	query := queryBuilder.String()

	rows, err := ut.db.Query(query, args...)
	if err != nil {
		return nil, ErrRetrieveFailed
	}
	defer rows.Close()

	return collectUsers(rows)
}

func (ut *UserTable) RetrieveAll() ([]User, error) {
	query := "SELECT id, name, age, country, degree, status, site FROM users"

	rows, err := ut.db.Query(query)
	if err != nil {
		return nil, ErrRetrieveFailed
	}
	defer rows.Close()

	return collectUsers(rows)
}

func (ut *UserTable) Paginate(
	pageSize, page int,
	orderBy OrderBy,
	filters Filters,
) ([]User, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(
		"SELECT id, name, age, country, degree, status, site FROM users WHERE 1=1",
	)

	applyFilters(&queryBuilder, filters)

	applyOrderBy(&queryBuilder, orderBy)

	offset := (page - 1) * pageSize
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset))

	args := []any{}
	if filters.Status != "" {
		args = append(args, filters.Status)
	}
	for _, country := range filters.Countries {
		args = append(args, country)
	}
	if filters.Age > 0 {
		args = append(args, filters.Age)
	}
	if filters.Degree != "" {
		args = append(args, filters.Degree)
	}

	rows, err := ut.db.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectUsers(rows)
}

func (ut *UserTable) Close() error {
	return ut.db.Close()
}

type queryOptions struct {
	Limit   int
	Status  []string
	Country []string
	OrderBy string
	Order   string
}

type Option func(*queryOptions)

func WithLimit(limit int) Option {
	return func(q *queryOptions) {
		q.Limit = limit
	}
}

func WithStatus(status []string) Option {
	return func(q *queryOptions) {
		q.Status = status
	}
}

func WithCountry(countries []string) Option {
	return func(q *queryOptions) {
		q.Country = countries
	}
}

func WithOrderBy(column, order string) Option {
	return func(q *queryOptions) {
		if order != "asc" && order != "desc" {
			panic(ErrInvalidOrder)
		}
		q.OrderBy = column
		q.Order = order
	}
}

func collectUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Age,
			&user.Country,
			&user.Degree,
			&user.Status,
			&user.Site,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func applyFilters(queryBuilder *strings.Builder, filters Filters) {
	if filters.Status != "" {
		queryBuilder.WriteString(" AND LOWER(status) = LOWER(?)")
	}
	if len(filters.Countries) > 0 {
		queryBuilder.WriteString(" AND LOWER(country) IN (")
		for i := range filters.Countries {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString("LOWER(?)")
		}
		queryBuilder.WriteString(")")
	}
	if filters.Age > 0 {
		queryBuilder.WriteString(" AND age = ?")
	}
	if filters.Degree != "" {
		queryBuilder.WriteString(" AND LOWER(degree) = LOWER(?)")
	}
}

func applyOrderBy(queryBuilder *strings.Builder, orderBy OrderBy) {
	orderClauses := []string{}
	if orderBy.Id != "" {
		orderClauses = append(orderClauses, "id "+orderBy.Id)
	}
	if orderBy.Age != "" {
		orderClauses = append(orderClauses, "age "+orderBy.Age)
	}
	if orderBy.Name != "" {
		orderClauses = append(orderClauses, "name "+orderBy.Name)
	}

	if len(orderClauses) > 0 {
		queryBuilder.WriteString(" ORDER BY " + strings.Join(orderClauses, ", "))
	}
}
