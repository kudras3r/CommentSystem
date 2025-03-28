package graph

import (
	"github.com/kudras3r/CommentSystem/internal/storage"
)

type Resolver struct {
	Storage storage.Storage
}

// func (r *Resolver) GetPost(ctx context.Context) ([]*model.Post, error) {
// 	return r.posts, nil
// }
