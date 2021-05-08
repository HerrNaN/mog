package mog

import "path"

func lockFilePathOf(filePath string) string {
	filename := path.Base(filePath)
	dir := path.Dir(filePath)
	lockFileName := ".#" + filename
	lockFilePath := path.Join(dir, lockFileName)
	return lockFilePath
}
