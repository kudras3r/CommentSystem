package storage

import "github.com/kudras3r/CommentSystem/internal/storage/model"

const (
	POST = "post"
	COMM = "comment"
)

type Storage interface {
	CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error)
	GetPost(id string) (*model.Post, error)
	GetAllPosts(limit, offset int) ([]*model.Post, error)
	CreateComment(postID string, content string, authorID string, parentID *string) (*model.Comment, error)
	GetCommentsByPostID(id string, limit int, offset int) ([]*model.Comment, error)
	GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error)
	CommentsNotAllow(id string) (bool, error)
}
