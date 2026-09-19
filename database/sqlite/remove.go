package sqlite

import (
	"errors"
	"os"
)

// 测试后需要删除sql.db
func remove(key string) error {
	if key != "Dont Input THIS" {
		return errors.New("WHY YOU WAHT TO REMOVE?")
	}
	if err := db.Close(); err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	return os.Remove(wd + "/sql.db")
}
