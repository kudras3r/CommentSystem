package storage

import "fmt"

// TODO more informative..

func FailedToConnect() error {
	return fmt.Errorf("failed to connect db")
}

func FailedToInsert(err error) error {
	return fmt.Errorf("failed to insert data %v : ", err)
}

func FailedToGetComments(err error) error {
	return fmt.Errorf("failed to get comments %v : ", err)
}

func FailedToGetPosts(err error) error {
	return fmt.Errorf("failed to get posts %v : ", err)
}

func NoWithID(id string, kind string) error {
	return fmt.Errorf("no %s with id %s", kind, id)
}

func FailedToGetWithId(kind, id string, err error) error {
	return fmt.Errorf("failed to get %s with id %s %v", kind, id, err)
}
