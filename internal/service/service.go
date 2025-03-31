package service

import (
	"fmt"
	"strings"

	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
	"github.com/kudras3r/CommentSystem/pkg/logger"
)

const (
	MAXCOMMENTLEN = 2000
	DEFAULTLIMIT  = 10
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

func validatePagination(first, offset *int32, log *logger.Logger) (int, int, error) {
	var sB strings.Builder
	sB.WriteString("validate pagination: ")

	var ifirst, ioffset int
	if first == nil {
		ifirst = DEFAULTLIMIT
		sB.WriteString(fmt.Sprintf("first is nil -> bydef set in %d ", DEFAULTLIMIT))
	} else {
		ifirst = int(*first)
		sB.WriteString(fmt.Sprintf("first is %d ", ifirst))
	}
	if offset != nil {
		ioffset = int(*offset)
		sB.WriteString(fmt.Sprintf("offset is %d", ioffset))
	}

	log.Info(sB.String())

	if ifirst < 0 || ioffset < 0 {
		log.Errorf("invalid limit or offset: %d, %d", ifirst, ioffset)
		return -1, -1, InvalidLimitOrOffset(ifirst, ioffset)
	}

	return ifirst, ioffset, nil
}

func (s *Service) ChildrenHandler(parentID string, first, offset *int32) ([]*model.Comment, error) {
	ifirst, ioffset, err := validatePagination(first, offset, s.log)
	if err != nil {
		return nil, err
	}
	s.log.Infof("getting childrens | parentID: %s, first: %d, offset: %d", parentID, ifirst, ioffset)
	return s.storage.GetCommentsByParent(parentID, ifirst, ioffset)
}

func (s *Service) CreatePostHandler(title, content, authorID string, allowComms bool) (*model.Post, error) {
	s.log.Infof("creating post | title: %s, content: %s, authorID: %s, allowComms: %t", content, title, authorID, allowComms)

	return s.storage.CreatePost(title, content, authorID, allowComms)
}

func (s *Service) CreateCommentHandler(postID, content, authorID string, parentID *string) (*model.Comment, error) {
	s.log.Infof("creating comment | postID: %s, authorID: %s", postID, authorID)

	if len(content) > MAXCOMMENTLEN {
		s.log.Warnf("comment is too long: %d", len(content))
		return nil, CommentIsTooLong()
	}

	notAllow, err := s.storage.CommentsNotAllow(postID)
	if notAllow {
		s.log.Warnf("comment not allow at post with id: %s", postID)
		return nil, CommentsNotAllow(postID)
	}
	if err != nil {
		s.log.Errorf("error at %sCreateCommentHandler() : %v", filePath, err)
		return nil, err
	}

	if parentID != nil {
		parentComment, err := s.storage.GetComment(*parentID)
		if err != nil {
			s.log.Errorf("error fetching parent comment: %v", err)
			return nil, err
		}

		if parentComment.PostID != postID {
			s.log.Errorf("parent comment with ID %s does not belong to post %s", *parentID, postID)
			return nil, PostAnCommRelationError()
		}
	}

	return s.storage.CreateComment(postID, content, authorID, parentID)
}

func (s *Service) CommentsHandler(postID string, first, offset *int32) ([]*model.Comment, error) {
	ifirst, ioffset, err := validatePagination(first, offset, s.log)
	if err != nil {
		return nil, err
	}
	s.log.Infof("getting comments | postID: %s, first: %d, offset:%d", postID, ifirst, ioffset)

	return s.storage.GetCommentsByPostID(postID, ifirst, ioffset)
}

func (s *Service) PostsHandler(first, offset *int32) ([]*model.Post, error) {
	ifirst, ioffset, err := validatePagination(first, offset, s.log)
	if err != nil {
		return nil, err
	}
	s.log.Infof("getting posts | first: %d, offset: %d", ifirst, ioffset)

	return s.storage.GetPosts(ifirst, ioffset)
}

func (s *Service) PostHandler(id string) (*model.Post, error) {
	s.log.Infof("getting post with id %s", id)

	return s.storage.GetPost(id)
}
