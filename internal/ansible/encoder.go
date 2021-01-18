package ansible

import (
	"errors"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"io/ioutil"
	"os"
	"strings"
)

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
		return errors.New(fmt.Sprintf("failed to save file '%s'", file))
	}
	return nil
}
