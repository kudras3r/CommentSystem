package storage

import "fmt"

func FailedToInsert(err error) error {
	return fmt.Errorf("failed to insert data %v", err)
}
