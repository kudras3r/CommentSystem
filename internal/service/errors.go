package service

import "fmt"

func CommentIsTooLong() error {
	return fmt.Errorf("comment is too long! max comment len: %d", MAXCOMMENTLEN)
}

func CommentsNotAllow(id string) error {
	return fmt.Errorf("comments not allow at post with id %s", id)
}

func InvalidLimitOrOffset(limit, offset int) error {
	return fmt.Errorf("invalid limit %d or offset %d", limit, offset)
}

func PostAnCommRelationError() error {
	return fmt.Errorf("parent comment does not belong to the current post")
}
