package inmemory

import (
	"strconv"
	"time"

	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
	"github.com/kudras3r/CommentSystem/pkg/logger"
)

const (
	filePath = "internal/storage/inmemory/inmemory.go/"

	maxCommentSliceFill = 70
	initialCommentSliceSize = 256
)

type IMSt struct {
	posts           map[string]*model.Post
	comms           []*model.Comment
	commsByPostID   map[string][]uint64 // indexes
	commsByParentID map[string][]uint64
	pp              uint64
	cp              uint64

	log *logger.Logger
}

func New(log *logger.Logger) *IMSt {
	return &IMSt{
		posts:           make(map[string]*model.Post),
		comms:           make([]*model.Comment, initialCommentSliceSize),
		commsByPostID:   make(map[string][]uint64),
		commsByParentID: make(map[string][]uint64),

		log: log,
	}
}

func (s *IMSt) increaseCommentSliceSize() {
	if len(s.comms) * 100 / cap(s.comms) >= maxCommentSliceFill {
		newSize := cap(s.comms) * 2
		newComms := make([]*model.Comment, len(s.comms), newSize)
		copy(newComms, s.comms)
		s.comms = newComms
	}
}

func (s *IMSt) CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error) {
	s.log.Infof("%sCreatePost() creating post with title: %s, content: %s, authorID: %s, allowComment: %t", filePath, title, content, authorID, allowComment)

	post := &model.Post{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		AllowComms: allowComment,
		CreatedAt:  time.Now().Format(time.RFC3339),
		ID:         strconv.FormatUint(s.pp, 10),
	}

	s.posts[post.ID] = post
	s.pp++

	return post, nil
}

func (s *IMSt) GetPost(id string) (*model.Post, error) {
	s.log.Infof("%sGetPost() fetching post with id: %s", filePath, id)

	post, found := s.posts[id]
	if !found {
		s.log.Errorf("%sGetPost() post not found with id: %s", filePath, id)
		return nil, storage.NoWithID(id, storage.POST)
	}

	return post, nil
}

func (s *IMSt) GetPosts(limit, offset int) ([]*model.Post, error) {
	s.log.Infof("%sGetPosts() fetching posts with limit: %d, offset: %d", filePath, limit, offset)

	var posts []*model.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	if offset >= len(posts) {
		s.log.Warnf("%sGetPosts() offset %d is too big, max : %d", filePath, offset, len(posts))
		return []*model.Post{}, nil
	}

	rightBorder := offset + limit
	if rightBorder > len(posts) {
		rightBorder = len(posts)
	}

	return posts[offset:rightBorder], nil
}

func (s *IMSt) CreateComment(postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	s.log.Infof("%sCreateComment() creating comment for postID: %s, authorID: %s", filePath, postID, authorID)
        s.increaseCommentSliceSize()

	comm := &model.Comment{
		PostID:    postID,
		Content:   content,
		AuthorID:  authorID,
		ParentID:  parentID,
		CreatedAt: time.Now().Format(time.RFC3339),
		ID:        strconv.FormatUint(s.cp, 10),
	}

	if parentID != nil {
		_, err := s.GetComment(*parentID)
		if err != nil {
			s.log.Errorf("%sCreateComment() no parentID %s", filePath, *parentID)
			return nil, storage.NoParentWithID(*parentID)
		}
		s.commsByParentID[*parentID] = append(s.commsByParentID[*parentID], s.cp)
	}
	if _, found := s.posts[postID]; found {
		s.commsByPostID[postID] = append(s.commsByPostID[postID], s.cp)
	} else {
		s.log.Errorf("%sCreateComment() no postID %s", filePath, postID)
		return nil, storage.NoWithID(postID, storage.POST)
	}

	s.comms[s.cp] = comm
	s.cp++

	return comm, nil
}

func (s *IMSt) GetCommentsByPostID(postID string, limit int, offset int) ([]*model.Comment, error) {
	s.log.Infof("%sGetCommentsByPostID() fetching comments for postID: %s with limit: %d, offset: %d", filePath, postID, limit, offset)

	commIDs, found := s.commsByPostID[postID]
	if !found {
		return []*model.Comment{}, nil
	}

	commsByPost := make([]*model.Comment, len(commIDs))
	for i, id := range commIDs {
		commsByPost[i] = s.comms[id]
	}

	var resComms []*model.Comment
	for _, comm := range commsByPost {
		if comm.ParentID == nil {
			resComms = append(resComms, comm)
		}
	}

	ln := len(resComms)
	if offset >= ln {
		s.log.Warnf("%sGetCommentsByPostID() offset %d is too big, max : %d", filePath, offset, ln)
		return []*model.Comment{}, nil
	}
	rightBorder := offset + limit
	if rightBorder > ln {
		rightBorder = ln
	}

	return resComms[offset:rightBorder], nil
}

func (s *IMSt) GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error) {
	s.log.Infof("%sGetCommentsByParent() fetching comments with parentID: %s with limit: %d, offset: %d", filePath, parent, limit, offset)

	commIDs, found := s.commsByParentID[parent]
	if !found {
		s.log.Errorf("%sGetCommentsByParent() not found parentID: %s", filePath, parent)
		return []*model.Comment{}, nil
	}
	commsByParent := make([]*model.Comment, len(commIDs))
	for i, id := range commIDs {
		commsByParent[i] = s.comms[id]
	}

	ln := len(commsByParent)
	if offset >= ln {
		s.log.Warnf("%sGetCommentsByParent() offset %d too big, max : %d", filePath, offset, ln)
		return []*model.Comment{}, nil
	}
	rightBorder := offset + limit
	if rightBorder > ln {
		rightBorder = ln
	}

	return commsByParent[offset:rightBorder], nil
}

func (s *IMSt) GetComment(id string) (*model.Comment, error) {
	s.log.Infof("%sGetComment() fetching comment with id: %s", filePath, id)

	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.Errorf("%sGetComment() failed to parse uint id : %v", filePath, err)
		return nil, storage.FailedToGetComments(err)
	}
	if uid > s.cp {
		s.log.Errorf("%sGetComment() requested id : %d max id in storage : %d", filePath, uid, s.cp)
		return nil, storage.NoWithID(id, storage.COMM)
	}

	return s.comms[uid], nil
}

func (s *IMSt) CommentsNotAllow(postID string) (bool, error) {
	s.log.Infof("%sCommentsNotAllow() checking if comments are allowed for postID: %s", filePath, postID)

	post, found := s.posts[postID]
	if !found {
		s.log.Errorf("%sCommentsNotAllow() no post with id: %s", filePath, postID)
		return false, storage.NoWithID(postID, storage.POST)
	}

	return !post.AllowComms, nil
}
