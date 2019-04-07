package category

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	errNotFound           = errors.New("category not found")
	errCategoryNotCreated = errors.New("category not created")
)

// Category describes category created by user
type Category struct {
	UID     uuid.UUID
	UserUID uuid.UUID
	Name    string
}

type datastore interface {
	getAllCategories(int32, int32) ([]*Category, error)
	getCategoryInfo(uuid.UUID) (*Category, error)
	createCategory(string, uuid.UUID) (*Category, error)
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAllCategories(pageSize, pageNumber int32) ([]*Category, error) {
	query := "SELECT uid, user_uid, name FROM categories LIMIT $1 OFFSET $2"
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, pageSize, lastRecord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Category, 0)
	for rows.Next() {
		category := new(Category)
		var uid, userUID string
		err := rows.Scan(&uid, &userUID, &category.Name)
		if err != nil {
			return nil, err
		}

		category.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, err
		}

		category.UserUID, err = uuid.Parse(userUID)
		if err != nil {
			return nil, err
		}

		result = append(result, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) getCategoryInfo(uid uuid.UUID) (*Category, error) {
	query := "SELECT user_uid, name FROM categories WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	result := new(Category)
	var stringUserUID string
	switch err := row.Scan(&stringUserUID, &result.Name); err {
	case nil:
		result.UID = uid
		userUID, err := uuid.Parse(stringUserUID)
		if err != nil {
			return nil, err
		}

		result.UserUID = userUID
		return result, nil
	case sql.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err
	}
}

func (db *db) createCategory(name string, userUID uuid.UUID) (*Category, error) {
	category := new(Category)

	query := "INSERT INTO categories (uid, user_uid, name) VALUES ($1, $2, $3)"
	uid := uuid.New()

	category.UID = uid
	category.UserUID = userUID
	category.Name = name

	result, err := db.Exec(query, category.UID.String(), userUID.String(), name)
	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errCategoryNotCreated
	}

	return category, nil
}
