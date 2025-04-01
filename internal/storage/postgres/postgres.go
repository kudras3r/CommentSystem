package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
	"github.com/kudras3r/CommentSystem/pkg/config"
	"github.com/kudras3r/CommentSystem/pkg/logger"

	_ "github.com/lib/pq"
)

const (
	filePath = "internal/storage/postgres/postgres.go/"
)

type PgDB struct {
	DB  *sqlx.DB
	Log *logger.Logger
}

func New(config config.DB, log *logger.Logger) *PgDB {
	connStr := fmt.Sprintf(
		`host=%s port=%d user=%s
		password=%s dbname=%s sslmode=disable`,
		config.Host, config.Port, config.User,
		config.Pass, config.Name)

	var attempt int
	var err error
	var db *sqlx.DB
	for ; attempt < 5; attempt++ {
		time.Sleep(time.Millisecond * 200)
		db, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Errorf("%sNew() failed to connect db attempt %d: %v", filePath, attempt, err)
		}
	}
	if err != nil {
		log.Fatalf("failed to connect db %d times", attempt)
	}

	log.Infof("%sNew() successfully connected to db", filePath)
	return &PgDB{
		DB:  db,
		Log: log,
	}
}

func (pg *PgDB) CloseConnection() {
	pg.Log.Warnf("%sCloseConnection() pg connection is closed", filePath)
	pg.DB.Close()
}

func (pg *PgDB) GetConnection() *sql.DB {
	pg.Log.Debugf("%sGetConnection() retrieving db connection", filePath)
	return pg.DB.DB
}

func (pg *PgDB) CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error) {
	query := `INSERT INTO posts (title, content, author_id, allow_comms) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, title, content, allow_comms, created_at, rating, author_id`

	pg.Log.Debugf("%sCreatePost() executing query: %s with params: title=%s, content=%s, authorID=%s, allowComment=%t", filePath, query, title, content, authorID, allowComment)
	var post model.Post
	err := pg.DB.QueryRowx(query, title, content, authorID, allowComment).StructScan(&post)
	if err != nil {
		pg.Log.Errorf("%sCreatePost() error: %v", filePath, err)
		return nil, storage.FailedToInsert(err)
	}

	pg.Log.Infof("%sCreatePost() post created with id: %s", filePath, post.ID)
	return &post, nil
}

func (pg *PgDB) GetPost(id string) (*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, rating, allow_comms 
              FROM posts WHERE id = $1`

	pg.Log.Debugf("%sGetPost() executing query: %s with id: %s", filePath, query, id)
	var post model.Post
	err := pg.DB.Get(&post, query, id)
	if err != nil {
		pg.Log.Errorf("%sGetPost() error: %v", filePath, err)
		return nil, storage.FailedToGetWithId(storage.POST, id, err)
	}

	pg.Log.Infof("%sGetPost() post retrieved with id: %s", filePath, post.ID)
	return &post, nil
}

func (pg *PgDB) GetPosts(limit, offset int) ([]*model.Post, error) {
	query := `SELECT id, title, content, allow_comms, created_at, rating, author_id
		FROM posts 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	pg.Log.Debugf("%sGetPost() executing query: %s with limit=%d, offset=%d", filePath, query, limit, offset)
	var posts []*model.Post = []*model.Post{}
	if err := pg.DB.Select(&posts, query, limit, offset); err != nil {
		pg.Log.Errorf("%sGetPosts() error: %v", filePath, err)
		return nil, storage.FailedToGetComments(err)
	}

	pg.Log.Infof("%sGetPosts() retrieved %d posts", filePath, len(posts))
	return posts, nil
}

func (pg *PgDB) CreateComment(postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	query := `INSERT INTO comments (content, rating, author_id, post_id, parent_id) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, content, created_at`

	pg.Log.Debugf("%sCreateComment() executing query: %s with params: content=%s, authorID=%s, postID=%s, parentID=%v", filePath, query, content, authorID, postID, parentID)
	var comment model.Comment
	err := pg.DB.QueryRowx(query, content, 0, authorID, postID, parentID).StructScan(&comment)
	if err != nil {
		pg.Log.Errorf("%sCreateComment() error: %v", filePath, err)
		return nil, storage.FailedToInsert(err)
	}

	pg.Log.Infof("%sCreateComment() comment created with id: %s", filePath, comment.ID)
	return &comment, nil
}

func (pg *PgDB) GetCommentsByPostID(id string, limit int, offset int) ([]*model.Comment, error) {
	query := `SELECT id, content, created_at, rating, author_id, post_id, parent_id
		FROM comments 
		WHERE post_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 
		OFFSET $3`

	pg.Log.Debugf("%sGetCommentsByPostID() executing query: %s with id=%s, limit=%d, offset=%d", filePath, query, id, limit, offset)
	var comments []*model.Comment = []*model.Comment{}
	if err := pg.DB.Select(&comments, query, id, limit, offset); err != nil {
		pg.Log.Errorf("%sGetCommentsByPostID() error: %v", filePath, err)
		return nil, storage.FailedToGetComments(err)
	}

	pg.Log.Infof("%sGetCommentsByPostID() retrieved %d comments for postID: %s", filePath, len(comments), id)
	return comments, nil
}

func (pg *PgDB) GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error) {
	query := `SELECT id, content, created_at, rating, author_id, post_id, parent_id
		FROM comments 
		WHERE parent_id = $1
		ORDER BY created_at DESC 
		LIMIT $2
		OFFSET $3`

	pg.Log.Debugf("%sGetCommentsByParent() executing query: %s with parentID=%s, limit=%d, offset=%d", filePath, query, parent, limit, offset)
	var comments []*model.Comment = []*model.Comment{}
	if err := pg.DB.Select(&comments, query, parent, limit, offset); err != nil {
		pg.Log.Errorf("%sGetCommentsByParent() error: %v", filePath, err)
		return nil, storage.FailedToGetComments(err)
	}

	pg.Log.Infof("%sGetCommentsByParent() retrieved %d comments for parentID: %s", filePath, len(comments), parent)
	return comments, nil
}

func (pg *PgDB) GetComment(id string) (*model.Comment, error) {
	query := `SELECT id, content, created_at, rating, author_id, post_id, parent_id
              FROM comments WHERE id = $1`

	var comment model.Comment
	err := pg.DB.Get(&comment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			pg.Log.Infof("%sGetComment() comment with id %s not found", filePath, id)
			return nil, storage.NoWithID(id, storage.COMM)
		}
		pg.Log.Errorf("%sGetComment() error : %v", filePath, err)
		return nil, storage.FailedToGetComments(err)
	}

	return &comment, nil
}

func (pg *PgDB) CommentsNotAllow(id string) (bool, error) {
	var allowComms bool
	err := pg.DB.Get(&allowComms, "SELECT allow_comms FROM posts WHERE id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			pg.Log.Infof("%sCommentsNotAllow() no rows found for postID: %s", filePath, id)
			return false, storage.NoWithID(id, storage.POST)
		}
		pg.Log.Errorf("%sCommentsNotAllow() error: %v", filePath, err)
		return false, storage.FailedToGetPosts(err)
	}

	pg.Log.Infof("%sCommentsNotAllow() comments allowed for postID: %s: %t", filePath, id, allowComms)
	return !allowComms, nil
}
