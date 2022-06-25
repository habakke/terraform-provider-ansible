package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"os"
)

func commitAndExport(db *database.Database, path string) error {
	if err := db.Commit(); err != nil {
		return fmt.Errorf("failed to commit database to disk: %s", err.Error())
	}

	if err := Encode(fmt.Sprintf("%s%shosts.ini", path, string(os.PathSeparator)), db); err != nil {
		return fmt.Errorf("failed to export to ansible: %s", err.Error())
	}

	return nil
}
