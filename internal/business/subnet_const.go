package business

import "terraform-provider-ctyun/internal/utils"

const (
	SubnetTypeCommon = "common" // 普通子网
	SubnetTypeCbm    = "cbm"    // 裸金属子网
)

const (
	SubnetTypeMapScene1 = iota
)

var SubnetTypes = []string{
	SubnetTypeCommon,
	SubnetTypeCbm,
}

var SubnetTypeMap = utils.Must(
	[]any{
		SubnetTypeCommon,
		SubnetTypeCbm,
	},
	map[utils.Scene][]any{
		SubnetTypeMapScene1: {
			0,
			1,
		},
	},
)
