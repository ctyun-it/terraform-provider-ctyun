package planmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequiresReplaceIfStateNotNull_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		stateValue types.String
		planValue  types.String
		expected   bool
	}{
		"state null, plan has value": {
			stateValue: types.StringNull(),
			planValue:  types.StringValue("new-value"),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has same value": {
			stateValue: types.StringValue("same-value"),
			planValue:  types.StringValue("same-value"),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has different value": {
			stateValue: types.StringValue("old-value"),
			planValue:  types.StringValue("new-value"),
			expected:   true, // 应该触发替换
		},
		"state has value, plan is null": {
			stateValue: types.StringValue("old-value"),
			planValue:  types.StringNull(),
			expected:   true, // 应该触发替换
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				StateValue: testCase.stateValue,
				PlanValue:  testCase.planValue,
			}
			resp := &planmodifier.StringResponse{}

			mod := RequiresReplaceIfStateNotNullModifier()
			mod.PlanModifyString(context.Background(), req, resp)

			if resp.RequiresReplace != testCase.expected {
				t.Errorf("expected RequiresReplace to be %t, got %t", testCase.expected, resp.RequiresReplace)
			}
		})
	}
}

func TestRequiresReplaceIfStateNotNull_Bool(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		stateValue types.Bool
		planValue  types.Bool
		expected   bool
	}{
		"state null, plan has value": {
			stateValue: types.BoolNull(),
			planValue:  types.BoolValue(true),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has same value": {
			stateValue: types.BoolValue(true),
			planValue:  types.BoolValue(true),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has different value": {
			stateValue: types.BoolValue(true),
			planValue:  types.BoolValue(false),
			expected:   true, // 应该触发替换
		},
		"state has value, plan is null": {
			stateValue: types.BoolValue(true),
			planValue:  types.BoolNull(),
			expected:   true, // 应该触发替换
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.BoolRequest{
				StateValue: testCase.stateValue,
				PlanValue:  testCase.planValue,
			}
			resp := &planmodifier.BoolResponse{}

			mod := RequiresReplaceIfStateNotNullModifier()
			mod.PlanModifyBool(context.Background(), req, resp)

			if resp.RequiresReplace != testCase.expected {
				t.Errorf("expected RequiresReplace to be %t, got %t", testCase.expected, resp.RequiresReplace)
			}
		})
	}
}

func TestRequiresReplaceIfStateNotNull_Int64(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		stateValue types.Int64
		planValue  types.Int64
		expected   bool
	}{
		"state null, plan has value": {
			stateValue: types.Int64Null(),
			planValue:  types.Int64Value(123),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has same value": {
			stateValue: types.Int64Value(123),
			planValue:  types.Int64Value(123),
			expected:   false, // 不应该触发替换
		},
		"state has value, plan has different value": {
			stateValue: types.Int64Value(123),
			planValue:  types.Int64Value(456),
			expected:   true, // 应该触发替换
		},
		"state has value, plan is null": {
			stateValue: types.Int64Value(123),
			planValue:  types.Int64Null(),
			expected:   true, // 应该触发替换
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.Int64Request{
				StateValue: testCase.stateValue,
				PlanValue:  testCase.planValue,
			}
			resp := &planmodifier.Int64Response{}

			mod := RequiresReplaceIfStateNotNullModifier()
			mod.PlanModifyInt64(context.Background(), req, resp)

			if resp.RequiresReplace != testCase.expected {
				t.Errorf("expected RequiresReplace to be %t, got %t", testCase.expected, resp.RequiresReplace)
			}
		})
	}
}
