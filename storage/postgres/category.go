package postgres

import (
	"database/sql"
	"fmt"

	"github.com/MuhammadyusufAdhamov/blog_project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategory(db *sqlx.DB) repo.CategoryStorageI {
	return &categoryRepo{
		db: db,
	}
}

func (cr *categoryRepo) Create(category *repo.Category) (*repo.Category, error) {
	query := `
		insert into categories(title) values($1)
		returning id, created_at
	`

	row := cr.db.QueryRow(
		query,
		category.Title,
	)

	err := row.Scan(
		&category.ID,
		&category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func(cr *categoryRepo) Get(id int64) (*repo.Category, error) {
	var result repo.Category

	quey := `
		select
			id,
			title,
			created_at
		from categories
		where id = $1
	`

	row := cr.db.QueryRow(quey, id)
	err := row.Scan(
		&result.ID,
		&result.Title,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (cr *categoryRepo) GetAll(params *repo.GetAllCategoriesParams) (*repo.GetAllCategoriesResult, error) {
	result := repo.GetAllCategoriesResult{
		Categories: make([]*repo.Category, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d ", params.Limit, offset)

	filter := ""
	if params.Search != "" {
		filter += " where title like '%" + params.Search + "%' "
	}

	query := `
		select
			id,
			title,
			created_at
		from categories
	`	+ filter + `
	order by created_at desc
	`	+ limit

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c repo.Category

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Categories = append(result.Categories, &c)
	}

	queyCount := `select count(1) from categories ` + filter
	err = cr.db.QueryRow(queyCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


func (ur *categoryRepo) UpdateCategory(category *repo.Category) (*repo.Category, error) {
	query := `update categories set title=$1 where id=$2
			returning created_at
			`
	
	err := ur.db.QueryRow(query,category.Title,category.ID,).Scan(&category.CreatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func(ur *categoryRepo) DeleteCategory(id int64) error {
	query := `delete from categories where id=$1`

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}