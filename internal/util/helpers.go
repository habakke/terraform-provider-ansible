package util

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceToStringArray(d *schema.ResourceData, name string) []string {
	if d.Get(name) == nil {
		return nil
	}
	in, ok := d.Get(name).([]interface{})
	if !ok || len(in) == 0 {
		return nil
	}

	out := make([]string, len(in))
	for i, v := range in {
		out[i] = fmt.Sprint(v)
	}
	return out
}

func ResourceToStringMap(d *schema.ResourceData, name string) map[string]string {
	in := d.Get(name)
	out := make(map[string]string)
	if in == nil {
		return out
	}

	// Have to first convert to a interface map which we can iterate over
	inMap, ok := in.(map[string]interface{})
	if !ok {
		return out
	}
	for k, v := range inMap {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func ResourceToInterfaceMap(d *schema.ResourceData, name string) map[string]interface{} {
	in, ok := d.Get(name).(map[string]interface{})
	if !ok {
		return nil
	}
	return in
}

func ResourceToFloat(d *schema.ResourceData, name string) float64 {
	in, ok := d.Get(name).(string)
	if !ok {
		return 0
	}
	return stringToFloat(in)
}

func stringToFloat(in string) float64 {
	if in == "" {
		return 0
	}

	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0
	}

	return f
}

func ResourceToString(d *schema.ResourceData, name string) string {
	in := d.Get(name)
	if in == nil {
		return ""
	}
	out, ok := in.(string)
	if !ok {
		return ""
	}
	return out
}

func ResourceToBool(d *schema.ResourceData, name string) bool {
	in := d.Get(name)
	if in == nil {
		return false
	}
	out, ok := in.(bool)
	if !ok {
		return false
	}
	return out
}

func ResourceToInt(d *schema.ResourceData, name string) int {
	in := d.Get(name)
	if in == nil {
		return 0
	}
	out, ok := in.(int)
	if !ok {
		return 0
	}
	return out
}
