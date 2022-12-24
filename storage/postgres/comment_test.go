package postgres_test

import (
	"testing"

	"blog_project/storage/repo"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createComment(t *testing.T) *repo.Comment {
	user := createUser(t)
	post := createPost(t)
	comment, err := strg.Comment().Create(&repo.Comment{
		Description: faker.Sentence(),
		UserID:      user.ID,
		PostID:      post.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, comment)

	return comment
}

func TestGetComment(t *testing.T) {
	c := createComment(t)

	comment, err := strg.Comment().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
}

func TestCreateComment(t *testing.T) {
	createComment(t)
}

func TestUpdateComment(t *testing.T) {
	c := createComment(t)

	c.Description = faker.Sentence()

	comment, err := strg.Comment().UpdateComment(c)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.Equal(t, comment.Description, c.Description)
}

func TestDeleteComment(t *testing.T) {
	c := createComment(t)

	err := strg.Comment().DeleteComment(c.ID)
	require.NoError(t, err)
}

func TestGetAllComments(t *testing.T) {
	createComment(t)

	result, err := strg.Comment().GetAll(&repo.GetAllCommentsParams{
		Limit: 3,
		Page:  1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, int(result.Count), 1)
}
