package graphql

import (
	"context"
	"fmt"

	"github.com/kudras3r/CommentSystem/internal/storage/model"
)

func (r *commentResolver) Children(ctx context.Context, obj *model.Comment, first *int32, offset *int32) ([]*model.Comment, error) {
	return r.Service.ChildrenHandler(obj.ID, first, offset)
}

func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID string, allowComms bool) (*model.Post, error) {
	return r.Service.CreatePostHandler(title, content, authorID, allowComms)
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, content string, authorID string, parentID *string) (*model.Comment, error) {
	return r.Service.CreateCommentHandler(postID, content, authorID, parentID)
}

func (r *postResolver) Comments(ctx context.Context, obj *model.Post, first *int32, offset *int32) ([]*model.Comment, error) {
	return r.Service.CommentsHandler(obj.ID, first, offset)
}

func (r *queryResolver) Posts(ctx context.Context, first *int32, offset *int32) ([]*model.Post, error) {
	return r.Service.PostsHandler(first, offset)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.Service.PostHandler(id)
}

func (r *subscriptionResolver) NewComment(ctx context.Context, postID string) (<-chan *model.Comment, error) { // TODO
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
