package business

import "terraform-provider-ctyun/internal/utils"

const (
	ImageOsDistroAnolis    = "anolis"
	ImageOsDistroCentos    = "centos"
	ImageOsDistroCtyunos   = "ctyunos"
	ImageOsDistroDebian    = "debian"
	ImageOsDistroFedora    = "fedora"
	ImageOsDistroKylin     = "kylin"
	ImageOsDistroOpenEuler = "openEuler"
	ImageOsDistroUbuntu    = "ubuntu"
	ImageOsDistroUnionTech = "UnionTech"
	ImageOsDistroWindows   = "windows"

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

	ImageVisibilityPrivate = "private"
	ImageVisibilityPublic  = "public"
	ImageVisibilityShared  = "shared"
	ImageVisibilitySafe    = "safe"
	ImageVisibilityApp     = "app"
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
	ImageVisibilityApp,
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

var ImageVisibilityMap = utils.Must(
	[]any{
		ImageVisibilityPrivate,
		ImageVisibilityPublic,
		ImageVisibilityShared,
		ImageVisibilitySafe,
		ImageVisibilityApp,
	},
	map[utils.Scene][]any{
		ImageVisibilityMapScene1: {
			0,
			1,
			2,
			3,
			4,
		},
	},
)
