package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
	"github.com/kudras3r/CommentSystem/pkg/config"

	_ "github.com/lib/pq"
)

type pgDB struct {
	DB *sqlx.DB
}

func New(config config.DB) (*pgDB, error) {
	connStr := fmt.Sprintf(
		`host=%s port=%d user=%s 
		password=%s dbname=%s sslmode=disable`,
		config.Host, config.Port, config.User,
		config.Pass, config.Name)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &pgDB{DB: db}, nil
}

func (pg *pgDB) CloseConnection() {
	pg.DB.Close()
}

func (pg *pgDB) GetConnection() *sql.DB {
	return pg.DB.DB
}

func (pg *pgDB) CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error) {
	query := `INSERT INTO posts (title, content, author_id, allow_comms) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, title, content, allow_comms, created_at, rating, author_id`

	var post model.Post
	err := pg.DB.QueryRowx(query, title, content, authorID, allowComment).StructScan(&post)
	if err != nil {
		return nil, storage.FailedToInsert(err)
	}

	return &post, nil
}

func (pg *pgDB) GetPost(id string) (*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, rating, allow_comms 
              FROM posts WHERE id = $1`

	var post model.Post
	err := pg.DB.Get(&post, query, id)
	if err != nil {
		return nil, storage.FailedToGetWithId(storage.POST, id, err)
	}

	return &post, nil
}

func (pg *pgDB) GetAllPosts(limit, offset int) ([]*model.Post, error) {
	query := `SELECT id, title, content, allow_comms, created_at, rating, author_id
		FROM posts 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var posts []*model.Post
	if err := pg.DB.Select(&posts, query, limit, offset); err != nil {
		return nil, storage.FailedToGetComments(err)
	}

	return posts, nil
}

func (pg *pgDB) CreateComment(postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	query := `INSERT INTO comments (content, rating, author_id, post_id, parent_id) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, content, created_at`

	var comment model.Comment
	err := pg.DB.QueryRowx(query, content, 0, authorID, postID, parentID).StructScan(&comment)
	if err != nil {
		return nil, storage.FailedToInsert(err)
	}

	return &comment, nil
}

func (pg *pgDB) GetCommentsByPostID(id string, limit int, offset int) ([]*model.Comment, error) {
	query := `SELECT id, content, created_at, rating, author_id, post_id, parent_id
		FROM comments 
		WHERE post_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 
		OFFSET $3`

	var comments []*model.Comment
	if err := pg.DB.Select(&comments, query, id, limit, offset); err != nil {
		return nil, storage.FailedToGetComments(err)
	}

	return comments, nil
}

func (pg *pgDB) GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error) {
	query := `SELECT id, content, created_at, rating, author_id, post_id, parent_id
		FROM comments 
		WHERE parent_id = $1
		ORDER BY created_at DESC 
		LIMIT $2
		OFFSET $3`

	var comments []*model.Comment
	if err := pg.DB.Select(&comments, query, parent, limit, offset); err != nil {
		return nil, storage.FailedToGetComments(err)
	}

	return comments, nil
}

func (pg *pgDB) CommentsNotAllow(id string) (bool, error) {
	var allowComms bool

	err := pg.DB.Get(&allowComms, "SELECT allow_comms FROM posts WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, storage.NoWithID(id, storage.POST)
		}
		return false, storage.FailedToGetPosts(err)
	}

	return allowComms, nil
}
