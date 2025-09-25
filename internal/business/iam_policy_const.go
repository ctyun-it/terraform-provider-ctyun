package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	PolicyRangeRegion = "region"
	PolicyRangeGlobal = "global"

	PolicyEffectAllow = "allow"
	PolicyEffectDeny  = "deny"
)

const (
	PolicyRangeMapScene1 = iota
	PolicyRangeMapScene2
)

const (
	PolicyEffectMapScene1 = iota
)

var PolicyRanges = []string{
	PolicyRangeRegion,
	PolicyRangeGlobal,
}

var PolicyEffects = []string{
	PolicyEffectAllow,
	PolicyEffectDeny,
}

var PolicyRangeMap = utils.Must(
	[]any{
		PolicyRangeRegion,
		PolicyRangeGlobal,
	},
	map[utils.Scene][]any{
		PolicyRangeMapScene1: {
			1,
			2,
		},
		PolicyRangeMapScene2: {
			"PROJECT_SERVICE",
			"GLOBAL_SERVICE",
		},
	},
)

var PolicyEffectMap = utils.Must(
	[]any{
		PolicyEffectAllow,
		PolicyEffectDeny,
	},
	map[utils.Scene][]any{
		PolicyEffectMapScene1: {
			"Allow",
			"Deny",
		},
	},
)
