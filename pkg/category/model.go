package category

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	errNotFound           = errors.New("category not found")
	errCategoryNotCreated = errors.New("category not created")
	errReportNotCreated   = errors.New("report not created")
)

// Category describes category created by user
type Category struct {
	UID         uuid.UUID
	UserUID     uuid.UUID
	Name        string
	Description string
}

// Report describes report submitted by user
type Report struct {
	UID         uuid.UUID
	CategoryUID uuid.UUID
	PostUID     uuid.UUID
	CommentUID  uuid.UUID
	Reason      string
	CreatedAt   time.Time
}

type datastore interface {
	getAllCategories(int32, int32) ([]*Category, error)
	getCategoryInfo(uuid.UUID) (*Category, error)
	createCategory(string, string, uuid.UUID) (*Category, error)
	getAllReports(uuid.UUID, int32, int32) ([]*Report, error)
	createReport(uuid.UUID, uuid.UUID, uuid.UUID, string) (*Report, error)
	deleteReport(uuid.UUID) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAllCategories(pageSize, pageNumber int32) ([]*Category, error) {
	query := "SELECT uid, user_uid, name, description FROM categories LIMIT $1 OFFSET $2"
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
		err := rows.Scan(&uid, &userUID, &category.Name, &category.Description)
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
	query := "SELECT user_uid, name, description FROM categories WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	result := new(Category)
	var stringUserUID string
	switch err := row.Scan(&stringUserUID, &result.Name, &result.Description); err {
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

func (db *db) createCategory(name, description string, userUID uuid.UUID) (*Category, error) {
	category := new(Category)

	query := "INSERT INTO categories (uid, user_uid, name, description) VALUES ($1, $2, $3, $4)"
	uid := uuid.New()

	category.UID = uid
	category.UserUID = userUID
	category.Name = name
	category.Description = description

	result, err := db.Exec(query, category.UID.String(), userUID.String(), name, description)
	if err != nil {
		return nil, err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errCategoryNotCreated
	}

	return category, nil
}

func (db *db) getAllReports(categoryUID uuid.UUID, pageSize, pageNumber int32) ([]*Report, error) {
	query := `SELECT uid, post_uid, comment_uid, reason, created_at
	          FROM reports
	          WHERE category_uid=$1
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, categoryUID.String(), pageSize, lastRecord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Report, 0)
	for rows.Next() {
		report := new(Report)
		var uid, postUID, commentUID string
		err := rows.Scan(&uid, &postUID, &commentUID, &report.Reason, &report.CreatedAt)
		if err != nil {
			return nil, err
		}

		report.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, err
		}

		report.PostUID, err = uuid.Parse(postUID)
		if err != nil {
			return nil, err
		}

		report.CommentUID, err = uuid.Parse(commentUID)
		if err != nil {
			return nil, err
		}

		report.CategoryUID = categoryUID

		result = append(result, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) createReport(categoryUID, postUID, commentUID uuid.UUID, reason string) (*Report, error) {
	report := new(Report)

	query := "INSERT INTO reports (uid, category_uid, post_uid, comment_uid, reason, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	uid := uuid.New()

	report.UID = uid
	report.CategoryUID = categoryUID
	report.PostUID = postUID
	report.CommentUID = commentUID
	report.Reason = reason
	report.CreatedAt = time.Now()

	result, err := db.Exec(query,
		report.UID.String(), report.CategoryUID.String(), report.PostUID.String(),
		report.CommentUID.String(), report.Reason, report.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errReportNotCreated
	}

	return report, nil
}

func (db *db) deleteReport(uid uuid.UUID) error {
	query := "DELETE FROM reports WHERE uid=$1"
	result, err := db.Exec(query, uid.String())
	if err != nil {
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}
