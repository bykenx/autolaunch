package autolaunch

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"text/template"
)

// @see https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html
const launchAgentsTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>{{.AppName}}</string>
	<key>ProgramArguments</key>
	<array>
	{{.ArgsTemplate}}
	</array>
	<key>RunAtLoad</key>
	<true/>
	{{.OtherTemplate}}
</dict>
</plist>
`

type autoLaunch struct {
	AppName           string
	AppPath           string
	Args              []string
	StartInterval     *int
	StandardErrorPath string
	StandardOutPath   string
}

// get auto launch item
//
// appName: app name, also the plist name use to create LaunchAgents,
// appPath: path to execute file which you want to auto launch
func New(appName, appPath string) *autoLaunch {
	return &autoLaunch{
		AppName: appName,
		AppPath: appPath,
	}
}

// enable auto launch item
func (m *autoLaunch) Enable() error {
	if !pathExists(m.AppPath) {
		return ErrAppPathNotExist
	}
	if !path.IsAbs(m.AppPath) {
		return ErrAppPathIllegal
	}
	tpl := template.Must(template.New("").Parse(launchAgentsTemplate))

	f, err := os.OpenFile(m.getPlistFilePath(), os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}
	defer f.Close()

	var argsTemplateList = []string{
		fmt.Sprintf("<string>%s</string>", m.AppPath),
	}

	for _, arg := range m.Args {
		argsTemplateList = append(argsTemplateList, fmt.Sprintf("<string>%s</string>", arg))
	}

	var otherTemplateList = []string{}

	if m.StartInterval != nil {
		otherTemplateList = append(otherTemplateList, fmt.Sprintf("<key>StartInterval</key>\n<integer>%d</integer>", *m.StartInterval))
	}

	if m.StandardOutPath == "" {
		otherTemplateList = append(otherTemplateList, "<key>StandardOutPath</key>\n<integer>/dev/null</integer>")
	} else {
		otherTemplateList = append(otherTemplateList, fmt.Sprintf("<key>StandardOutPath</key>\n<integer>%s</integer>", m.StandardOutPath))
	}

	if m.StandardErrorPath == "" {
		otherTemplateList = append(otherTemplateList, "<key>StandardErrorPath</key>\n<integer>/dev/null</integer>")
	} else {
		otherTemplateList = append(otherTemplateList, fmt.Sprintf("<key>StandardErrorPath</key>\n<integer>%s</integer>", m.StandardErrorPath))
	}

	return tpl.Execute(f, map[string]interface{}{
		"AppName":       m.AppName,
		"ArgsTemplate":  strings.Join(argsTemplateList, "\n"),
		"OtherTemplate": strings.Join(otherTemplateList, "\n"),
	})
}

// disable auto launch item
func (m autoLaunch) Disable() error {
	launchAgentsFilePath := m.getPlistFilePath()

	if !pathExists(launchAgentsFilePath) {
		return nil
	}

	return os.Remove(launchAgentsFilePath)
}

func (m autoLaunch) getPlistFilePath() string {
	usr, _ := user.Current()

	return path.Join(usr.HomeDir, "Library/LaunchAgents/", fmt.Sprintf("%s.plist", m.AppName))
}

// return auto launch item has enabled
func (m autoLaunch) Enabled() bool {
	return pathExists(m.getPlistFilePath())
}
