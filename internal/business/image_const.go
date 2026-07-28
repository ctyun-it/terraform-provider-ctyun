package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	ImageOsDistroAnolis    = "Anolis"
	ImageOsDistroCentos    = "CentOS"
	ImageOsDistroCtyunos   = "CTyunOS"
	ImageOsDistroDebian    = "Debian"
	ImageOsDistroFedora    = "Fedora"
	ImageOsDistroKylin     = "KylinOS"
	ImageOsDistroOpenEuler = "openEuler"
	ImageOsDistroUbuntu    = "Ubuntu"
	ImageOsDistroUnionTech = "UnionTechOS"
	ImageOsDistroWindows   = "Windows Server"

	ImageArchitectureAarch64 = "aarch64"
	ImageArchitectureX8664   = "x86_64"

	ImageBootModeBios = "bios"
	ImageBootModeUefi = "uefi"

	ImageTypeSystemDiskImage = "system"
	ImageTypeDataDiskImage   = "data"

	ImageStatusAccepted       = "accepted"
	ImageStatusActive         = "active"
	ImageStatusDeactivated    = "deactivated"
	ImageStatusDeactivating   = "deactivating"
	ImageStatusDeleted        = "deleted"
	ImageStatusDeleting       = "deleting"
	ImageStatusError          = "error"
	ImageStatusImporting      = "importing"
	ImageStatusKilled         = "killed"
	ImageStatusPending_delete = "pending_delete"
	ImageStatusQueued         = "queued"
	ImageStatusReactivating   = "reactivating"
	ImageStatusRejected       = "rejected"
	ImageStatusSaving         = "saving"
	ImageStatusSyncing        = "syncing"
	ImageStatusUploading      = "uploading"
	ImageStatusWaiting        = "waiting"

	ImageAssociationUserTypeShare   = "share"
	ImageAssociationUserTypeReceive = "receive"

	ImageVisibilityPrivate   = "private"
	ImageVisibilityPublic    = "public"
	ImageVisibilityShared    = "shared"
	ImageVisibilitySafe      = "safe"
	ImageVisibilityCommunity = "community"
	ImageVisibilityApp       = "app"
	ImageVisibilityMarket    = "market"
)

const (
	ImageTypeMapScene1 = iota
)

const (
	ImageVisibilityMapScene1 = iota
)

var ImageOsDistros = []string{
	ImageOsDistroAnolis,
	ImageOsDistroCentos,
	ImageOsDistroCtyunos,
	ImageOsDistroDebian,
	ImageOsDistroFedora,
	ImageOsDistroKylin,
	ImageOsDistroOpenEuler,
	ImageOsDistroUbuntu,
	ImageOsDistroUnionTech,
	ImageOsDistroWindows,
}

var ImageArchitectures = []string{
	ImageArchitectureAarch64,
	ImageArchitectureX8664,
}

var ImageBootModes = []string{
	ImageBootModeBios,
	ImageBootModeUefi,
}

var ImageTypes = []string{
	ImageTypeSystemDiskImage,
	ImageTypeDataDiskImage,
}

var ImageStatuses = []string{
	ImageStatusQueued,
	ImageStatusActive,
	ImageStatusDeleting,
}

var ImageAssociationUserTypes = []string{
	ImageAssociationUserTypeShare,
	ImageAssociationUserTypeReceive,
}

var ImageVisibilities = []string{
	ImageVisibilityPrivate,
	ImageVisibilityPublic,
	ImageVisibilityShared,
	ImageVisibilitySafe,
	ImageVisibilityCommunity,
	ImageVisibilityApp,
	ImageVisibilityMarket,
}

var ImageTypeMap = utils.Must(
	[]any{
		ImageTypeSystemDiskImage,
		ImageTypeDataDiskImage,
	},
	map[utils.Scene][]any{
		ImageTypeMapScene1: {
			"",
			"data_disk_image",
		},
	},
)

var ImageVisibilityMap = map[string]int{
	ImageVisibilityPrivate:   0,
	ImageVisibilityPublic:    1,
	ImageVisibilityShared:    2,
	ImageVisibilitySafe:      3,
	ImageVisibilityCommunity: 4,
	ImageVisibilityApp:       5,
	ImageVisibilityMarket:    6,
}
