package postgres_test

import (
	"testing"

	"github.com/MuhammadyusufAdhamov/blog_project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)


func createPost(t *testing.T) *repo.Post {
	user := createUser(t)
	category := createCategory(t)
	post, err := strg.Post().Create(&repo.Post{
		Title: faker.Sentence(),
		Description: faker.Sentence(),
		UserID: user.ID,
		CategoryID: category.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, post)

	return post
}

func TestGetPost(t *testing.T) {
	c := createPost(t)

	post, err := strg.Post().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post)
}

func TestCreatePost(t *testing.T) {
	createPost(t)
}

func TestUpdatePost(t *testing.T) {
	p := createPost(t)

	user := createPost(t)
	category := createPost(t)
	post, err := strg.Post().UpdatePost(&repo.Post{
		Title: faker.Sentence(),
		Description: faker.Sentence(),
		UserID: user.ID,
		CategoryID: category.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, post.Title, p.Title)
}

func TestDeletePost(t *testing.T) {
	c := createPost(t)


	err := strg.Post().DeletePost(c.ID)
	require.NoError(t, err)
}

func TestGetAllPosts(t *testing.T) {
	createPost(t)

	result, err := strg.Post().GetAll(&repo.GetAllPostsParams{
		Limit: 3,
		Page: 1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, int(result.Count), 1)
}