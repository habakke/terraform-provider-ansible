package ansible

import (
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAnsibleProviders map[string]*schema.Provider
var testAnsibleProvider *schema.Provider

func TestMain(m *testing.M) {
	util.ConfigureLogging(util.GetEnv("LOGLEVEL", "info"), true)
	code := m.Run()
	os.Exit(code)
}

func init() {

	testAnsibleProvider = Provider()
	testAnsibleProviders = map[string]*schema.Provider{
		"ansible": testAnsibleProvider,
	}
}

//nolint:deadcode,unused
func testAnsiblePreCheck(t *testing.T, resourceName string) {
	r := testAnsibleProvider.ResourcesMap[strings.Split(resourceName, ".")[0]]
	d := testAnsibleProvider.DataSourcesMap[strings.Split(resourceName, ".")[0]]

	if r != nil && d != nil {
		t.Fatalf("missing resource '%s'", resourceName)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %v", err)
	}
}
