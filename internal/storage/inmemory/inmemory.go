package inmemory

import (
	"strconv"
	"time"

	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
)

const (
	defaultPostsCap = 128
	defaultCommCap  = 512
	MAXCAPCOEF      = 70
)

type IMSt struct {
	posts    []*model.Post
	comments []*model.Comment
	pp       uint64
	cp       uint64
}

func New() *IMSt {
	return &IMSt{
		posts:    make([]*model.Post, defaultPostsCap),
		comments: make([]*model.Comment, defaultCommCap),
	}
}

func (s *IMSt) increseCapacity(kind string) {
	switch kind {
	case storage.POST:
		newStorage := make([]*model.Post, len(s.posts)*2)
		copy(newStorage, s.posts)
		s.posts = newStorage

	case storage.COMM:
		newStorage := make([]*model.Comment, len(s.comments)*2)
		copy(newStorage, s.comments)
		s.comments = newStorage
	}
}

func (s *IMSt) CreatePost(title, content, authorID string, allowComment bool) (*model.Post, error) {
	post := &model.Post{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		AllowComms: allowComment,
		CreatedAt:  time.Now().Format(time.RFC3339),
		ID:         strconv.FormatUint(s.pp, 10),
	}
	s.posts[s.pp] = post
	s.pp++

	if s.pp/uint64(len(s.posts))*100 > MAXCAPCOEF {
		s.increseCapacity(storage.POST)
	}

	return post, nil
}

func (s *IMSt) GetPost(id string) (*model.Post, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, storage.FailedToGetPosts(err)
	}
	if uid >= s.pp {
		return nil, storage.NoWithID(id, storage.POST)
	}

	return s.posts[uid], nil
}

func (s *IMSt) GetAllPosts(limit, offset int) ([]*model.Post, error) {
	if limit < 0 || offset < 0 {
		return nil, storage.InvalidLimitOrOffset(limit, offset)
	}
	if s.pp == 0 {
		return []*model.Post{}, nil
	}

	return s.posts[offset : limit+offset], nil
}

func (s *IMSt) CreateComment(postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	comm := &model.Comment{
		PostID:    postID,
		Content:   content,
		AuthorID:  authorID,
		ParentID:  parentID,
		CreatedAt: time.Now().Format(time.RFC3339),
		ID:        strconv.FormatUint(s.cp, 10),
	}
	s.comments[s.pp] = comm
	s.cp++

	if s.cp/uint64(len(s.comments))*100 > MAXCAPCOEF {
		s.increseCapacity(storage.COMM)
	}

	return comm, nil
}

func (s *IMSt) GetComment(id string) (*model.Comment, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, storage.FailedToGetComments(err)
	}
	if uid >= s.cp {
		return nil, storage.NoWithID(id, storage.COMM)
	}

	return s.comments[uid], nil
}

func (s *IMSt) GetCommentsByPostID(id string, limit int, offset int) ([]*model.Comment, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, storage.FailedToGetWithId(storage.POST, id, err)
	}
	if uid >= s.pp {
		return nil, storage.NoWithID(id, storage.POST)
	}
	if limit < 0 || offset < 0 {
		return nil, storage.InvalidLimitOrOffset(limit, offset)
	}

	var comments []*model.Comment
	for i := uint64(0); i < s.cp; i++ {
		comm := s.comments[i]
		if comm.PostID == id {
			comments = append(comments, comm)
		}
	}

	return comments[offset : limit+offset], nil
}

func (s *IMSt) GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error) {
	uid, err := strconv.ParseUint(parent, 10, 64)
	if err != nil {
		return nil, storage.FailedToGetWithId(storage.COMM, parent, err)
	}
	if uid >= s.cp {
		return nil, storage.NoWithID(parent, storage.POST)
	}
	if limit < 0 || offset < 0 {
		return nil, storage.InvalidLimitOrOffset(limit, offset)
	}

	var comments []*model.Comment
	for i := uint64(0); i < s.cp; i++ {
		comm := s.comments[i]
		if comm.ParentID == &parent {
			comments = append(comments, comm)
		}
	}

	return comments[offset : limit+offset], nil
}

func (s *IMSt) CommentsNotAllow(id string) (bool, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return false, storage.FailedToGetWithId(storage.COMM, id, err)
	}
	if uid >= s.cp {
		return false, storage.NoWithID(id, storage.POST)
	}

	return !s.posts[uid].AllowComms, nil
}
