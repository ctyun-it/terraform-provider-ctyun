package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	SubnetTypeCommon    = "common" // 普通子网
	SubnetTypeEbm       = "cbm"    // 裸金属子网
	SubnetTypeCommonInt = 0
	SubnetTypeEbmInt    = 1
)

const (
	SubnetTypeMapScene1 = iota
)

var SubnetTypes = []string{
	SubnetTypeCommon,
	SubnetTypeEbm,
}

var SubnetTypeMap = utils.Must(
	[]any{
		SubnetTypeCommon,
		SubnetTypeEbm,
	},
	map[utils.Scene][]any{
		SubnetTypeMapScene1: {
			SubnetTypeCommonInt,
			SubnetTypeEbmInt,
		},
	},
)
