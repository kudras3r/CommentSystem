package service

import (
	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
	"github.com/kudras3r/CommentSystem/pkg/logger"
)

const (
	MACCOMMENTLEN = 2000
	filePath      = "internal/service/service.go/"
)

type Service struct {
	storage storage.Storage
	log     *logger.Logger
}

func New(storage storage.Storage, log *logger.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
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
	s.log.Infof("geting childrens | parentID: %s, first: %d, offset: %d", parentID, ifirst, ioffset)
	return s.storage.GetCommentsByParent(parentID, ifirst, ioffset)
}

func (s *Service) CreatePostHandler(title, content, authorID string, allowComms bool) (*model.Post, error) {
	s.log.Infof("creating post | title: %s, content: %s, authorID: %s, allowComms: %t", content, title, authorID, allowComms)
	return s.storage.CreatePost(title, content, authorID, allowComms)
}

func (s *Service) CreateCommentHandler(postID, content, authorID string, parentID *string) (*model.Comment, error) {
	s.log.Infof("creating comment | postID: %s, authorID: %s", postID, authorID)

	if len(content) > MACCOMMENTLEN {
		s.log.Warnf("comment is too long: %d", len(content))
		return nil, CommentIsTooLong()
	}

	notAllow, err := s.storage.CommentsNotAllow(postID)
	if notAllow {
		s.log.Warnf("comment not allow comments: %d", len(content))
		return nil, CommentsNotAllow(postID)
	}
	if err != nil {
		s.log.Errorf("error at %sCreateCommentHandler() : %v", filePath, err)
		return nil, err
	}

	return s.storage.CreateComment(postID, content, authorID, parentID)
}

func (s *Service) CommentsHandler(postID string, first, offset *int32) ([]*model.Comment, error) {
	ifirst, ioffset := validatePagination(first, offset)
	s.log.Infof("getting comments | postID: %s, first: %d, offset:%d", postID, ifirst, ioffset)
	return s.storage.GetCommentsByPostID(postID, ifirst, ioffset)
}

func (s *Service) PostsHandler(first, offset *int32) ([]*model.Post, error) {
	ifirst, ioffset := validatePagination(first, offset)
	s.log.Infof("getting posts | first: %d, offset: %d", ifirst, ioffset)
	return s.storage.GetPosts(ifirst, ioffset)
}

func (s *Service) PostHandler(id string) (*model.Post, error) {
	s.log.Infof("getting post with id %s", id)
	return s.storage.GetPost(id)
}
