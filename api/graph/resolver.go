package graph

import (
	"context"

	"github.com/kudras3r/CommentSystem/api/graph/model"
)

type Resolver struct {
	posts []*model.Post
}

func (r *Resolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	return r.posts, nil
}
