package inmemory_test

import (
	"testing"

	"github.com/kudras3r/CommentSystem/internal/storage/inmemory"
	"github.com/kudras3r/CommentSystem/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func setup() *inmemory.IMSt {
	return inmemory.New(logger.New("DEBUG"))
}

func TestCreatePost(t *testing.T) {
	storage := setup()
	post, err := storage.CreatePost("Test Post", "This is a test post", "author1", true)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post", post.Content)
	assert.Equal(t, "author1", post.AuthorID)
	assert.True(t, post.AllowComms)
}

func TestGetPost(t *testing.T) {
	storage := setup()
	post, _ := storage.CreatePost("Test Post", "This is a test post", "author1", true)
	retrievedPost, err := storage.GetPost(post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedPost)
	assert.Equal(t, post.ID, retrievedPost.ID)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
}

func TestCreateComment(t *testing.T) {
	storage := setup()
	post, _ := storage.CreatePost("Test Post", "This is a test post", "author1", true)
	comment, err := storage.CreateComment(post.ID, "This is a test comment", "author2", nil)
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, post.ID, comment.PostID)
	assert.Equal(t, "This is a test comment", comment.Content)
	assert.Equal(t, "author2", comment.AuthorID)
}

func TestGetCommentsByPostID(t *testing.T) {
	storage := setup()
	post, _ := storage.CreatePost("Test Post", "This is a test post", "author1", true)
	storage.CreateComment(post.ID, "First comment", "author2", nil)
	storage.CreateComment(post.ID, "Second comment", "author3", nil)
	comments, err := storage.GetCommentsByPostID(post.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "First comment", comments[0].Content)
	assert.Equal(t, "Second comment", comments[1].Content)
}

func TestGetCommentsWithPagination(t *testing.T) {
	storage := setup()
	post, _ := storage.CreatePost("Test Post", "This is a test post", "author1", true)
	storage.CreateComment(post.ID, "First comment", "author2", nil)
	storage.CreateComment(post.ID, "Second comment", "author3", nil)
	storage.CreateComment(post.ID, "Third comment", "author4", nil)
	storage.CreateComment(post.ID, "Fourth comment", "author5", nil)
	comments, err := storage.GetCommentsByPostID(post.ID, 2, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "First comment", comments[0].Content)
	assert.Equal(t, "Second comment", comments[1].Content)
	comments, err = storage.GetCommentsByPostID(post.ID, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "Third comment", comments[0].Content)
	assert.Equal(t, "Fourth comment", comments[1].Content)
}

func TestCommentsNotAllow(t *testing.T) {
	storage := setup()
	post, _ := storage.CreatePost("Test Post", "This is a test post", "author1", false)
	allow, err := storage.CommentsNotAllow(post.ID)
	assert.NoError(t, err)
	assert.True(t, allow)
}
