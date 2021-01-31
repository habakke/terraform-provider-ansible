package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"io/ioutil"
	"os"
	"strings"
)

// The Encode function encodes the database to an Ansible compatible hosts.ini file
func Encode(file string, database *database.Database) error {
	var s string
	for _, v := range *database.AllGroups() {
		e := v.GetEntriesAsString()
		if e == nil {
			s = s + fmt.Sprintf("[%s]\n", v.GetName())
		} else {
			s = s + fmt.Sprintf("[%s]\n", v.GetName())
			s = s + fmt.Sprintf("%s\n", strings.Join(e, "\n"))
		}
	}
	if err := ioutil.WriteFile(file, []byte(s), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file '%s'", file)
	}
	return nil
}
