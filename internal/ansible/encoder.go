package ansible

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"strings"
)

func Encode(file string, database *database.Database) error {
	f := ini.Empty()
	for _, v := range *database.AllGroups() {
		e := v.GetEntriesAsString()
		if e == nil {
			_, _ = f.NewRawSection(v.GetName(), "")
		} else {
			_, _ = f.NewRawSection(v.GetName(), strings.Join(e, "\n"))
		}
	}
	if err := f.SaveTo(file); err != nil {
		return errors.New(fmt.Sprintf("failed to save file '%s'", file))
	}
	return nil
}
