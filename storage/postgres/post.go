package postgres

import (
	"database/sql"
	"fmt"

	"github.com/MuhammadyusufAdhamov/blog_project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPost(db *sqlx.DB) repo.PostStorageI {
	return &postRepo{
		db: db,
	}
}

func (ur *postRepo) Create(post *repo.Post) (*repo.Post, error) {
	query := `
		insert into posts (
			title,
			description,
			image_url,
			user_id,
			category_id
		) values ($1, $2, $3, $4, $5)
		returning id, created_at
	`

	row := ur.db.QueryRow(
		query,
		post.Title,
		post.Description,
		post.ImageUrl,
		post.UserID,
		post.CategoryID,
	)

	err := row.Scan(
		&post.ID,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ur *postRepo) Get(id int64) (*repo.Post, error) {
	var result repo.Post

	query := `
		select
			id,
			title,
			description,
			image_url,
			user_id,
			category_id,
			created_at,
			views_count
		from posts
		where id=$1
	`

	row := ur.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.ImageUrl,
		&result.UserID,
		&result.CategoryID,
		&result.CreatedAt,
		&result.ViewsCount,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pr *postRepo) GetAll(params *repo.GetAllPostsParams) (*repo.GetAllPostsResult, error) {
	result := repo.GetAllPostsResult {
		Posts: make([]*repo.Post, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" limit %d offset %d", params.Limit, offset)

	filter := "where true"
	if params.Search != "" {
		filter += " and title like '%" + params.Search + "%' "
	}

	if params.UserID != 0 {
		filter += fmt.Sprintf(" and user_id=%d ", params.UserID)
	}

	if params.CategoryID != 0 {
		filter += fmt.Sprintf(" and category_id=%d ", params.CategoryID)
	}

	orderBy := " order by created_at desc "
	if params.SortByData != "" {
		orderBy = fmt.Sprintf(" order by created_at %s ", params.SortByData)
	}

	query := `
		SELECT
			id,
			title,
			description,
			image_url,
			user_id,
			category_id,
			created_at,
			views_count
		FROM posts
		` + filter + orderBy + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p repo.Post

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.ImageUrl,
			&p.UserID,
			&p.CategoryID,
			&p.CreatedAt,
			&p.ViewsCount,
		)
		if err != nil {
			return nil, err
		}

		result.Posts = append(result.Posts, &p)
	}

	queryCount := `select count(1) from posts` + filter
	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


func (ur *postRepo) UpdatePost(post *repo.Post) (*repo.Post, error) {
	query := `update posts set
				title=$1,
				description=$2,
				image_url=$3,
				user_id=$4,
				category_id=$5
			where id=$6
			returning created_at
			`
	
	row := ur.db.QueryRow(
		query,
		post.Title,
		post.Description,
		post.ImageUrl,
		post.UserID,
		post.CategoryID,
		post.ID,
	)

	err := row.Scan(
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func(ur *postRepo) DeletePost(id int64) error {
	query := `delete from posts where id=$1
			returning id`

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