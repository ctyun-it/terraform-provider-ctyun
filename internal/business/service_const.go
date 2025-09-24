package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	ServiceTypeRegion = "region"
	ServiceTypeGlobal = "global"
)

var ServiceTypes = []string{
	ServiceTypeRegion,
	ServiceTypeGlobal,
}

const (
	ServiceTypeMapScene1 = iota
)

var ServiceTypeMap = utils.Must(
	[]any{
		ServiceTypeRegion,
		ServiceTypeGlobal,
	},
	map[utils.Scene][]any{
		ServiceTypeMapScene1: {
			1,
			2,
		},
	},
)
