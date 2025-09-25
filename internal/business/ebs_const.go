package business

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
)

const (
	EbsDiskModeVbd   = "vbd"
	EbsDiskModeIscsi = "iscsi"
	EbsDiskModeFcsan = "fcsan"

	EbsDiskTypeSata      = "sata"
	EbsDiskTypeSas       = "sas"
	EbsDiskTypeSsd       = "ssd"
	EbsDiskTypeSsdGenric = "ssd-genric"
	EbsDiskTypeFastSsd   = "fast-ssd"

	EbsSnapshotStatusAvailable = "available"
)

const (
	EbsDiskModeMapScene1 = iota
)

const (
	EbsDiskTypeMapScene1 = iota
)

var EbsDiskModes = []string{
	EbsDiskModeVbd,
	EbsDiskModeIscsi,
	EbsDiskModeFcsan,
}

var EbsDiskTypes = []string{
	EbsDiskTypeSata,
	EbsDiskTypeSas,
	EbsDiskTypeSsd,
	EbsDiskTypeSsdGenric,
	EbsDiskTypeFastSsd,
}

var EbsDiskModeMap = utils.Must(
	[]any{
		EbsDiskModeVbd,
		EbsDiskModeIscsi,
		EbsDiskModeFcsan,
	},
	map[utils.Scene][]any{
		EbsDiskModeMapScene1: {
			"VBD",
			"ISCSI",
			"FCSAN",
		},
	},
)

var EbsDiskTypeMap = utils.Must(
	[]any{
		EbsDiskTypeSata,
		EbsDiskTypeSas,
		EbsDiskTypeSsd,
		EbsDiskTypeSsdGenric,
		EbsDiskTypeFastSsd,
	},
	map[utils.Scene][]any{
		EbsDiskTypeMapScene1: {
			"SATA",
			"SAS",
			"SSD",
			"SSD-genric",
			"FAST-SSD",
		},
	},
)
