package autolaunch

import "os"

func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
