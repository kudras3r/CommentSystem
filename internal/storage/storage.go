package storage

import "github.com/kudras3r/CommentSystem/internal/storage/model"

type Storage interface {
	CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error)
	GetPost(id string) (*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	CreateComment(comment model.Comment) (*model.Comment, error)
	GetCommentsByPostID(id string, after string, limit int) ([]*model.Comment, error)
}
