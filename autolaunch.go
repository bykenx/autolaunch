package autolaunch

import "errors"

var (
	ErrAppPathNotExist = errors.New("app path doesn't exist")
	ErrAppPathIllegal  = errors.New("app path is illegal")
)
