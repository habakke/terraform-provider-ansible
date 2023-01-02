package ansible

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"ansible": func() (*schema.Provider, error) {
		return New(), nil
	},
}

func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		os.Exit(m.Run())
	}
	resource.TestMain(m)
}

func init() {
}

//nolint:deadcode,unused
func testAnsiblePreCheck(t *testing.T, resourceName string) {
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %v", err)
	}
}
