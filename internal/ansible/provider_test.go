package ansible

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAnsibleProviders map[string]terraform.ResourceProvider
var testAnsibleProvider *schema.Provider

func init() {
	testAnsibleProvider = Provider().(*schema.Provider)
	testAnsibleProviders = map[string]terraform.ResourceProvider{
		"ansible": testAnsibleProvider,
	}
}

func testAnsiblePreCheck(t *testing.T) {
	// TODO Add pre checks here.
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		assert.Fail(t, fmt.Sprintf("provider internal validation failed: %e", err))
	}
}
