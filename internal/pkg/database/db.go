package database

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func RunSQLScripts(db *gorm.DB, dirPath string) error {

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = db.Exec(string(data)).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
