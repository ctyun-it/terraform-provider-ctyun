package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// SecBool *bool转bool
func SecBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// SecBoolValue  *bool转types.Bool
func SecBoolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolValue(false)
	}
	return types.BoolValue(*b)
}
