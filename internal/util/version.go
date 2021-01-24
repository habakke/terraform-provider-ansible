package util

import "fmt"

var (
	version   string // build version number
	commit    string // sha1 revision used to build the program
	buildTime string // when the executable was built
	buildBy   string
)

func GetVersionString() string {
	return fmt.Sprintf("terraform-ansible-provider %s-%s %s %s", version, commit, buildTime, buildBy)
}
