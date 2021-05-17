package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"io/ioutil"
	"log"
	"os"
)

// Encode function encodes the database to an Ansible compatible hosts.ini file
func Encode(file string, database *database.Database) error {
	var s string
	for _, v := range *database.AllGroups() {
		ek := v.GetEntities()
		if len(ek) == 0 {
			s = s + fmt.Sprintf("[%s]\n", v.GetName())
		} else {
			s = s + fmt.Sprintf("[%s]\n", v.GetName())
			for _, k := range ek {
				e, err := v.GetEntity(k)
				if err != nil {
					log.Fatalf("failed to lookup entity '%s'", k)
				}

				es, err := encodeEntity(e)
				if err != nil {
					log.Fatalf("failed to encode entity %e", err)
				}
				s = s + fmt.Sprintf("%s\n", es)
			}
			s = s + "\n"
		}
	}
	if err := ioutil.WriteFile(file, []byte(s), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file '%s'", file)
	}
	return nil
}

func encodeEntity(e interface{}) (string, error) {
	switch t := e.(type) {
	case *database.Host:
		return encodeHost(e.(*database.Host)), nil
	case *database.Group:
		return encodeGroup(e.(*database.Group)), nil
	default:
		return "", fmt.Errorf("unknown entity type %s", t)
	}
}

func encodeGroup(g *database.Group) string {
	return g.GetName()
}

func encodeHost(h *database.Host) string {
	s := h.GetName()
	for _, vk := range h.GetVariableNames() {
		v, err := h.GetVariable(vk)
		if err != nil {
			log.Fatalf("unalbe to find expected host variable '%s'", vk)
		}
		s = s + fmt.Sprintf(" %s=%s", vk, v)
	}
	return s
}
