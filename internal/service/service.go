package service

import (
	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
)

const (
	MACCOMMENTLEN = 2000
)

type Service struct {
	storage *storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: &storage,
	}
}

func validatePagination(first, offset *int32) (int, int) {
	var ifirst, ioffset int
	if first == nil {
		ifirst = 10
	} else {
		ifirst = int(*first)
	}
	if offset != nil {
		ioffset = int(*offset)
	}
	return ifirst, ioffset
}

func (s *Service) ChildrenHandler(parentID string, first, offset *int32) ([]*model.Comment, error) {
	ifirst, ioffset := validatePagination(first, offset)
	return (*s.storage).GetCommentsByParent(parentID, ifirst, ioffset)
}

func (s *Service) CreatePostHandler(title, content, authorID string, allowComms bool) (*model.Post, error) {
	return (*s.storage).CreatePost(title, content, authorID, allowComms)
}

func (s *Service) CreateCommentHandler(postID, content, authorID string, parentID *string) (*model.Comment, error) {
	if len(content) > MACCOMMENTLEN {
		return nil, CommentIsTooLong()
	}

	notAllow, err := (*s.storage).CommentsNotAllow(postID)
	if notAllow {
		return nil, CommentsNotAllow(postID)
	}
	if err != nil {
		return nil, err
	}

	return (*s.storage).CreateComment(postID, content, authorID, parentID)
}

func (s *Service) CommentsHandler(postID string, first, offset *int32) ([]*model.Comment, error) {
	ifirst, ioffset := validatePagination(first, offset)
	return (*s.storage).GetCommentsByPostID(postID, ifirst, ioffset)
}

func (s *Service) PostsHandler(first, offset *int32) ([]*model.Post, error) {
	ifirst, ioffset := validatePagination(first, offset)
	return (*s.storage).GetAllPosts(ifirst, ioffset)
}

func (s *Service) PostHandler(id string) (*model.Post, error) {
	return (*s.storage).GetPost(id)
}
