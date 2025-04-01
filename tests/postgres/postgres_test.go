package postgres_test

import (
	"testing"

	"github.com/kudras3r/CommentSystem/internal/storage/postgres"
	"github.com/kudras3r/CommentSystem/pkg/config"
	"github.com/kudras3r/CommentSystem/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	log := logger.New("DEBUG")
	dbConfig := config.DB{
		Host: "localhost",
		Port: 5432,
		User: "ozon_keker",
		Pass: "1234",
		Name: "comm_sys_db",
	}

	pgDB := postgres.New(dbConfig, log)
	defer pgDB.CloseConnection()

	post, err := pgDB.CreatePost("Test Post", "Test Content", "1", true)
	assert.Nil(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Post", post.Title)
}

func TestGetPost(t *testing.T) {
	log := logger.New("DEBUG")
	dbConfig := config.DB{
		Host: "localhost",
		Port: 5432,
		User: "ozon_keker",
		Pass: "1234",
		Name: "comm_sys_db",
	}

	pgDB := postgres.New(dbConfig, log)
	defer pgDB.CloseConnection()

	post, err := pgDB.CreatePost("Test Post", "Test Content", "1", true)
	assert.Nil(t, err)
	assert.NotNil(t, post)

	retrievedPost, err := pgDB.GetPost(post.ID)
	assert.Nil(t, err)
	assert.Equal(t, post.ID, retrievedPost.ID)
}

func TestCreateComment(t *testing.T) {
	log := logger.New("DEBUG")
	dbConfig := config.DB{
		Host: "localhost",
		Port: 5432,
		User: "ozon_keker",
		Pass: "1234",
		Name: "comm_sys_db",
	}

	pgDB := postgres.New(dbConfig, log)
	defer pgDB.CloseConnection()

	post, err := pgDB.CreatePost("Test Post", "Test Content", "1", true)
	assert.Nil(t, err)
	assert.NotNil(t, post)

	comment, err := pgDB.CreateComment(post.ID, "Test Comment", "2", nil)
	assert.Nil(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Test Comment", comment.Content)
}

func TestGetCommentsByPostID(t *testing.T) {
	log := logger.New("DEBUG")
	dbConfig := config.DB{
		Host: "localhost",
		Port: 5432,
		User: "ozon_keker",
		Pass: "1234",
		Name: "comm_sys_db",
	}

	pgDB := postgres.New(dbConfig, log)
	defer pgDB.CloseConnection()

	post, err := pgDB.CreatePost("Test Post", "Test Content", "1", true)
	assert.Nil(t, err)
	assert.NotNil(t, post)

	_, err = pgDB.CreateComment(post.ID, "Test Comment 1", "2", nil)
	assert.Nil(t, err)
	_, err = pgDB.CreateComment(post.ID, "Test Comment 2", "3", nil)
	assert.Nil(t, err)

	comments, err := pgDB.GetCommentsByPostID(post.ID, 2, 0)
	assert.Nil(t, err)
	assert.Len(t, comments, 2)
}

func TestCommentsNotAllow(t *testing.T) {
	log := logger.New("DEBUG")
	dbConfig := config.DB{
		Host: "localhost",
		Port: 5432,
		User: "ozon_keker",
		Pass: "1234",
		Name: "comm_sys_db",
	}

	pgDB := postgres.New(dbConfig, log)
	defer pgDB.CloseConnection()

	post, err := pgDB.CreatePost("Test Post", "Test Content", "1", false)
	assert.Nil(t, err)
	assert.NotNil(t, post)

	allow, err := pgDB.CommentsNotAllow(post.ID)
	assert.Nil(t, err)
	assert.True(t, allow)
}
