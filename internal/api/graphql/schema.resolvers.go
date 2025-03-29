package graphql

import (
	"context"
	"fmt"

	"github.com/kudras3r/CommentSystem/internal/storage/model"
)

// TODO service level

func (r *commentResolver) Children(ctx context.Context, obj *model.Comment, first *int32, offset *int32) ([]*model.Comment, error) {
	var ifirst, ioffset int
	if first == nil {
		ifirst = 10
	} else {
		ifirst = int(*first)
	}
	if offset != nil {
		ioffset = int(*offset)
	}

	return r.Storage.GetCommentsByParent(obj.ID, ifirst, ioffset)
}

func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID string, allowComms bool) (*model.Post, error) {
	return r.Storage.CreatePost(title, content, authorID, allowComms)
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	return r.Storage.CreateComment(postID, content, authorID, parentID)
}

func (r *postResolver) Comments(ctx context.Context, obj *model.Post, first *int32, offset *int32) ([]*model.Comment, error) {
	var ifirst, ioffset int
	if first == nil {
		ifirst = 10
	} else {
		ifirst = int(*first)
	}
	if offset != nil {
		ioffset = int(*offset)
	}

	return r.Storage.GetCommentsByPostID(obj.ID, ifirst, ioffset)
}

func (r *queryResolver) Posts(ctx context.Context, first *int32, offset *int32) ([]*model.Post, error) {
	var ifirst, ioffset int
	if first == nil {
		ifirst = 10
	} else {
		ifirst = int(*first)
	}
	if offset != nil {
		ioffset = int(*offset)
	}
	return r.Storage.GetAllPosts(ifirst, ioffset)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.Storage.GetPost(id)
}

func (r *subscriptionResolver) NewComment(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: NewComment - newComment"))
}

func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Post() PostResolver { return &postResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
