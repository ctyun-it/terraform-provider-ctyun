package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	BandwidthStatusActive   = "active"
	BandwidthStatusExpired  = "expired"
	BandwidthStatusFreezing = "freezing"
)

const (
	BandwidthStatusMapScene1 = iota
)

var BandwidthStatusMap = utils.Must(
	[]any{
		BandwidthStatusActive,
		BandwidthStatusExpired,
		BandwidthStatusFreezing,
	},
	map[utils.Scene][]any{
		BandwidthStatusMapScene1: {
			"ACTIVE",
			"EXPIRED",
			"FREEZING",
		},
	},
)
