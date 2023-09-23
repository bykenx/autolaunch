package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/bykenx/autolaunch"
)

func setAutoStart() {
	cwd, _ := os.Getwd()

	autoLaunch := autolaunch.New("hello-world", path.Join(cwd, "hello-world"))
	autoLaunch.Args = append(autoLaunch.Args, "--autostart=1")

	enabled := autoLaunch.Enabled()

	var err error
	var mode string

	if enabled {
		mode = "disable"
		err = autoLaunch.Disable()
	} else {
		mode = "enable"
		err = autoLaunch.Enable()
	}
	if err != nil {
		fmt.Printf("%s auto launch error: %v\n", mode, err)
	} else {
		fmt.Printf("%s auto launch succeed.\n", mode)
	}
}

func normalStart() {
	fmt.Println("hello world from auto launch")
	time.Sleep(time.Second * 100)
}

func getBool(key string) bool {

	var val bool = false

	matchParam := fmt.Sprintf("--%s=1", key)

	for _, arg := range os.Args {
		if arg == matchParam {
			val = true
			break
		}
	}
	return val
}

func main() {

	isAutoStart := getBool("autostart")

	if isAutoStart {
		normalStart()
	} else {
		setAutoStart()
	}
}
