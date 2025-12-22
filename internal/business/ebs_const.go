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
	EbsDiskTypeXssd0     = "xssd-0"
	EbsDiskTypeXssd1     = "xssd-1"
	EbsDiskTypeXssd2     = "xssd-2"

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
	EbsDiskTypeXssd0,
	EbsDiskTypeXssd1,
	EbsDiskTypeXssd2,
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
		EbsDiskTypeXssd0,
		EbsDiskTypeXssd1,
		EbsDiskTypeXssd2,
	},
	map[utils.Scene][]any{
		EbsDiskTypeMapScene1: {
			"SATA",
			"SAS",
			"SSD",
			"SSD-genric",
			"FAST-SSD",
			"XSSD-0",
			"XSSD-1",
			"XSSD-2",
		},
	},
)
