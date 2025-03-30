package service

import "fmt"

func CommentIsTooLong() error {
	return fmt.Errorf("comment is too long! max comment len: %d", MACCOMMENTLEN)
}

func CommentsNotAllow(id string) error {
	return fmt.Errorf("comments not allow at post with id %s", id)
}
