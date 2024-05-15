package business

import "terraform-provider-ctyun/internal/utils"

const (
	EnterpriseProjectStatusEnable  = "enable"
	EnterpriseProjectStatusDisable = "disable"
)

var EnterpriseProjectStatuses = []string{
	EnterpriseProjectStatusEnable,
	EnterpriseProjectStatusDisable,
}

const (
	EnterpriseProjectStatusMapScene1 = iota
)

var EnterpriseProjectStatusMap = utils.Must(
	[]any{
		EnterpriseProjectStatusEnable,
		EnterpriseProjectStatusDisable,
	},
	map[utils.Scene][]any{
		EnterpriseProjectStatusMapScene1: {
			1,
			2,
		},
	},
)
