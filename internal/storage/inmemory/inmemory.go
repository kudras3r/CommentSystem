package inmemory

import (
	"strconv"
	"time"

	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/model"
)

type IMSt struct {
	posts    map[string]*model.Post
	comments map[string][]*model.Comment
	pp       uint64
	cp       uint64
}

func New() *IMSt {
	return &IMSt{
		posts:    make(map[string]*model.Post),
		comments: make(map[string][]*model.Comment),
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
	s.posts[post.ID] = post
	s.pp++

	return post, nil
}

func (s *IMSt) GetPost(id string) (*model.Post, error) {
	post, found := s.posts[id]
	if !found {
		return nil, storage.NoWithID(id, storage.POST)
	}

	return post, nil
}

func (s *IMSt) GetPosts(limit, offset int) ([]*model.Post, error) {
	if limit < 0 || offset < 0 {
		return nil, storage.InvalidLimitOrOffset(limit, offset)
	}

	var posts []*model.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	if offset >= len(posts) {
		return []*model.Post{}, nil
	}
	rightBorder := offset + limit
	if rightBorder > len(posts) {
		rightBorder = len(posts)
	}

	return posts[offset:rightBorder], nil
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

	s.comments[postID] = append(s.comments[postID], comm)
	s.cp++

	return comm, nil
}

func (s *IMSt) GetComment(id string) (*model.Comment, error) {
	for _, comments := range s.comments {
		for _, comm := range comments {
			if comm.ID == id {
				return comm, nil
			}
		}
	}

	return nil, storage.NoWithID(id, storage.COMM)
}

func (s *IMSt) GetCommentsByPostID(postID string, limit int, offset int) ([]*model.Comment, error) {
	comments, found := s.comments[postID]
	if !found {
		return nil, storage.NoWithID(postID, storage.POST)
	}
	if limit < 0 || offset < 0 {
		return nil, storage.InvalidLimitOrOffset(limit, offset)
	}

	var filteredComments []*model.Comment
	for _, comm := range comments {
		if comm.ParentID == nil {
			filteredComments = append(filteredComments, comm)
		}
	}

	if offset >= len(filteredComments) {
		return []*model.Comment{}, nil
	}
	rightBorder := offset + limit
	if rightBorder > len(filteredComments) {
		rightBorder = len(filteredComments)
	}

	return filteredComments[offset:rightBorder], nil
}

func (s *IMSt) GetCommentsByParent(parent string, limit int, offset int) ([]*model.Comment, error) {
	var comments []*model.Comment
	for _, comms := range s.comments {
		for _, comm := range comms {
			if comm.ParentID != nil && *comm.ParentID == parent {
				comments = append(comments, comm)
			}
		}
	}

	if offset >= len(comments) {
		return []*model.Comment{}, nil
	}
	rightBorder := offset + limit
	if rightBorder > len(comments) {
		rightBorder = len(comments)
	}

	return comments[offset:rightBorder], nil
}

func (s *IMSt) CommentsNotAllow(postID string) (bool, error) {
	post, found := s.posts[postID]
	if !found {
		return false, storage.NoWithID(postID, storage.POST)
	}

	return !post.AllowComms, nil
}
