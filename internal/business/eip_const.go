package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	EipAssociationTypeVm  = "vm"
	EipAssociationTypeVip = "vip"
	EipAssociationTypeBm  = "bm"

	EipBandwidthTypeStandalone = "standalone"
	EipBandwidthTypeShare      = "share"

	EipStatusActive              = "active"
	EipStatusDown                = "down"
	EipStatusError               = "error"
	EipStatusUpdating            = "updating"
	EipStatusBandingOrUnbangding = "banding_or_unbangding"
	EipStatusDeleting            = "deleting"
	EipStatusDeleted             = "deleted"
	EipStatusExpired             = "expired"
)

const (
	EipDemandBillingTypeBandwidth = "bandwidth"
	EipDemandBillingTypeUpflowc   = "upflowc"
)

const (
	EipAssociationTypeMapScene1 = iota
	EipAssociationTypeMapScene2
)

const (
	EipStatusMapScene1 = iota
)

var EipAssociationTypes = []string{
	EipAssociationTypeVm,
	EipAssociationTypeVip,
	EipAssociationTypeBm,
}

var EipDemandBillingTypes = []string{
	EipDemandBillingTypeBandwidth,
	EipDemandBillingTypeUpflowc,
}

var EipBandwidthTypes = []string{
	EipBandwidthTypeStandalone,
	EipBandwidthTypeShare,
}

var EipStatus = []string{
	EipStatusActive,
	EipStatusDown,
	EipStatusError,
	EipStatusUpdating,
	EipStatusBandingOrUnbangding,
	EipStatusDeleting,
	EipStatusDeleted,
	EipStatusExpired,
}

var EipAssociationTypeMap = utils.Must(
	[]any{
		EipAssociationTypeVm,
		EipAssociationTypeVip,
		EipAssociationTypeBm,
	},
	map[utils.Scene][]any{
		EipAssociationTypeMapScene1: {
			1,
			2,
			3,
		},
		EipAssociationTypeMapScene2: {
			"INSTANCE",
			"VIP",
			"PHYSICALINSTANCE",
		},
	},
)

var EipStatusMap = utils.Must(
	[]any{
		EipStatusActive,
		EipStatusDown,
		EipStatusError,
		EipStatusUpdating,
		EipStatusBandingOrUnbangding,
		EipStatusDeleting,
		EipStatusDeleted,
		EipStatusExpired,
	},
	map[utils.Scene][]any{
		EipStatusMapScene1: {
			"ACTIVE",
			"DOWN",
			"ERROR",
			"UPDATING",
			"BANDING_OR_UNBANGDING",
			"DELETING",
			"DELETED",
			"EXPIRED",
		},
	},
)
