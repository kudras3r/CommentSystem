package storage

import "fmt"

// TODO more informative..

func FailedToInsert(err error) error {
	return fmt.Errorf("failed to insert data %v : ", err)
}

func FailedToGetComments(err error) error {
	return fmt.Errorf("failed to get comments %v : ", err)
}

func FailedToGetPosts(err error) error {
	return fmt.Errorf("failed to get posts %v : ", err)
}

func NoWithID(id uint64, kind string) error {
	return fmt.Errorf("no %s with id %d", kind, id)
}

func InvalidLimitOrOffset(limit, offset int) error {
	return fmt.Errorf("invalid limit %d or offset %d", limit, offset)
}

func FailedToGetWithId(kind, id string, err error) error {
	return fmt.Errorf("failed to get %s with id %s %v", kind, id, err)
}
