package graphql

import (
	"github.com/kudras3r/CommentSystem/internal/storage"
)

type Resolver struct {
	Storage storage.Storage
}
