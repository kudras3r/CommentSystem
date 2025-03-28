package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
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
		return nil, err
	}

	return &post, nil
}

func (pg *pgDB) GetPost(id string) (*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, rating, allow_comms 
              FROM posts WHERE id = $1`

	var post model.Post
	err := pg.DB.Get(&post, query, id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
func (pg *pgDB) GetAllPosts() ([]*model.Post, error) { return nil, nil }

func (pg *pgDB) CreateComment(comment model.Comment) (*model.Comment, error) { return nil, nil }

func (pg *pgDB) GetCommentsByPostID(id string, after string, limit int) ([]*model.Comment, error) {
	return nil, nil
}
