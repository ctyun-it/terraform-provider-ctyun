package business

const (
	EcsStatusBackingup     = "backingup"
	EcsStatusCreating      = "creating"
	EcsStatusExpired       = "expired"
	EcsStatusFreezing      = "freezing"
	EcsStatusRebuild       = "rebuild"
	EcsStatusRestarting    = "restarting"
	EcsStatusRunning       = "running"
	EcsStatusStarting      = "starting"
	EcsStatusStopped       = "stopped"
	EcsStatusStopping      = "stopping"
	EcsStatusError         = "error"
	EcsStatusSnapshotting  = "snapshotting"
	EcsStatusUnsubscribed  = "unsubscribed"
	EcsStatusUnsubscribing = "unsubscribing"

	EcsFlavorType_CPU                  = "CPU"
	EcsFlavorType_CPU_C3               = "CPU_C3"
	EcsFlavorType_CPU_C6               = "CPU_C6"
	EcsFlavorType_CPU_C7               = "CPU_C7"
	EcsFlavorType_CPU_c7ne             = "CPU_c7ne"
	EcsFlavorType_CPU_C8               = "CPU_C8"
	EcsFlavorType_CPU_D3               = "CPU_D3"
	EcsFlavorType_CPU_FC1              = "CPU_FC1"
	EcsFlavorType_CPU_FM1              = "CPU_FM1"
	EcsFlavorType_CPU_FS1              = "CPU_FS1"
	EcsFlavorType_CPU_HC1              = "CPU_HC1"
	EcsFlavorType_CPU_HM1              = "CPU_HM1"
	EcsFlavorType_CPU_HS1              = "CPU_HS1"
	EcsFlavorType_CPU_IP3              = "CPU_IP3"
	EcsFlavorType_CPU_IR3              = "CPU_IR3"
	EcsFlavorType_CPU_IP3_2            = "CPU_IP3_2"
	EcsFlavorType_CPU_IR3_2            = "CPU_IR3_2"
	EcsFlavorType_CPU_KC1              = "CPU_KC1"
	EcsFlavorType_CPU_KM1              = "CPU_KM1"
	EcsFlavorType_CPU_KS1              = "CPU_KS1"
	EcsFlavorType_CPU_M2               = "CPU_M2"
	EcsFlavorType_CPU_M3               = "CPU_M3"
	EcsFlavorType_CPU_M6               = "CPU_M6"
	EcsFlavorType_CPU_M7               = "CPU_M7"
	EcsFlavorType_CPU_M8               = "CPU_M8"
	EcsFlavorType_CPU_S2               = "CPU_S2"
	EcsFlavorType_CPU_S3               = "CPU_S3"
	EcsFlavorType_CPU_S6               = "CPU_S6"
	EcsFlavorType_CPU_S7               = "CPU_S7"
	EcsFlavorType_CPU_S8               = "CPU_S8"
	EcsFlavorType_CPU_s8r              = "CPU_s8r"
	EcsFlavorType_GPU_N_V100_V_FMGQ    = "GPU_N_V100_V_FMGQ"
	EcsFlavorType_GPU_N_V100_V         = "GPU_N_V100_V"
	EcsFlavorType_GPU_N_V100S_V        = "GPU_N_V100S_V"
	EcsFlavorType_GPU_N_V100S_V_FMGQ   = "GPU_N_V100S_V_FMGQ"
	EcsFlavorType_GPU_N_T4_V           = "GPU_N_T4_V"
	EcsFlavorType_GPU_N_G7_V           = "GPU_N_G7_V"
	EcsFlavorType_GPU_N_V100           = "GPU_N_V100"
	EcsFlavorType_GPU_N_V100_SHIPINYUN = "GPU_N_V100_SHIPINYUN"
	EcsFlavorType_GPU_N_V100_SUANFA    = "GPU_N_V100_SUANFA"
	EcsFlavorType_GPU_N_P2V_RENMIN     = "GPU_N_P2V_RENMIN"
	EcsFlavorType_GPU_N_V100S          = "GPU_N_V100S"
	EcsFlavorType_GPU_N_T4             = "GPU_N_T4"
	EcsFlavorType_GPU_N_T4_AIJISUAN    = "GPU_N_T4_AIJISUAN"
	EcsFlavorType_GPU_N_T4_ASR         = "GPU_N_T4_ASR"
	EcsFlavorType_GPU_N_T4_JX          = "GPU_N_T4_JX"
	EcsFlavorType_GPU_N_T4_SHIPINYUN   = "GPU_N_T4_SHIPINYUN"
	EcsFlavorType_GPU_N_T4_SUANFA      = "GPU_N_T4_SUANFA"
	EcsFlavorType_GPU_N_T4_YUNYOUXI    = "GPU_N_T4_YUNYOUXI"
	EcsFlavorType_GPU_N_PI7            = "GPU_N_PI7"
	EcsFlavorType_GPU_N_P8A            = "GPU_N_P8A"
	EcsFlavorType_GPU_A_PAK1           = "GPU_A_PAK1"
	EcsFlavorType_GPU_C_PCH1           = "GPU_C_PCH1"

	EcsFlavorSeries_S   = "S"
	EcsFlavorSeries_C   = "C"
	EcsFlavorSeries_M   = "M"
	EcsFlavorSeries_HS  = "HS"
	EcsFlavorSeries_HC  = "HC"
	EcsFlavorSeries_HM  = "HM"
	EcsFlavorSeries_FS  = "FS"
	EcsFlavorSeries_FC  = "FC"
	EcsFlavorSeries_FM  = "FM"
	EcsFlavorSeries_KS  = "KS"
	EcsFlavorSeries_KC  = "KC"
	EcsFlavorSeries_P   = "P"
	EcsFlavorSeries_G   = "G"
	EcsFlavorSeries_IP3 = "IP3"
)

var EcsFlavorTypes = []string{
	EcsFlavorType_CPU,
	EcsFlavorType_CPU_C3,
	EcsFlavorType_CPU_C6,
	EcsFlavorType_CPU_C7,
	EcsFlavorType_CPU_c7ne,
	EcsFlavorType_CPU_C8,
	EcsFlavorType_CPU_D3,
	EcsFlavorType_CPU_FC1,
	EcsFlavorType_CPU_FM1,
	EcsFlavorType_CPU_FS1,
	EcsFlavorType_CPU_HC1,
	EcsFlavorType_CPU_HM1,
	EcsFlavorType_CPU_HS1,
	EcsFlavorType_CPU_IP3,
	EcsFlavorType_CPU_IR3,
	EcsFlavorType_CPU_IP3_2,
	EcsFlavorType_CPU_IR3_2,
	EcsFlavorType_CPU_KC1,
	EcsFlavorType_CPU_KM1,
	EcsFlavorType_CPU_KS1,
	EcsFlavorType_CPU_M2,
	EcsFlavorType_CPU_M3,
	EcsFlavorType_CPU_M6,
	EcsFlavorType_CPU_M7,
	EcsFlavorType_CPU_M8,
	EcsFlavorType_CPU_S2,
	EcsFlavorType_CPU_S3,
	EcsFlavorType_CPU_S6,
	EcsFlavorType_CPU_S7,
	EcsFlavorType_CPU_S8,
	EcsFlavorType_CPU_s8r,
	EcsFlavorType_GPU_N_V100_V_FMGQ,
	EcsFlavorType_GPU_N_V100_V,
	EcsFlavorType_GPU_N_V100S_V,
	EcsFlavorType_GPU_N_V100S_V_FMGQ,
	EcsFlavorType_GPU_N_T4_V,
	EcsFlavorType_GPU_N_G7_V,
	EcsFlavorType_GPU_N_V100,
	EcsFlavorType_GPU_N_V100_SHIPINYUN,
	EcsFlavorType_GPU_N_V100_SUANFA,
	EcsFlavorType_GPU_N_P2V_RENMIN,
	EcsFlavorType_GPU_N_V100S,
	EcsFlavorType_GPU_N_T4,
	EcsFlavorType_GPU_N_T4_AIJISUAN,
	EcsFlavorType_GPU_N_T4_ASR,
	EcsFlavorType_GPU_N_T4_JX,
	EcsFlavorType_GPU_N_T4_SHIPINYUN,
	EcsFlavorType_GPU_N_T4_SUANFA,
	EcsFlavorType_GPU_N_T4_YUNYOUXI,
	EcsFlavorType_GPU_N_PI7,
	EcsFlavorType_GPU_N_P8A,
	EcsFlavorType_GPU_A_PAK1,
	EcsFlavorType_GPU_C_PCH1,
}

var EcsFlavorSeries = []string{
	EcsFlavorSeries_S,
	EcsFlavorSeries_C,
	EcsFlavorSeries_M,
	EcsFlavorSeries_HS,
	EcsFlavorSeries_HC,
	EcsFlavorSeries_HM,
	EcsFlavorSeries_FS,
	EcsFlavorSeries_FC,
	EcsFlavorSeries_FM,
	EcsFlavorSeries_KS,
	EcsFlavorSeries_KC,
	EcsFlavorSeries_P,
	EcsFlavorSeries_G,
	EcsFlavorSeries_IP3,
}

var EcsStatus = []string{
	EcsStatusBackingup,
	EcsStatusCreating,
	EcsStatusExpired,
	EcsStatusFreezing,
	EcsStatusRebuild,
	EcsStatusRestarting,
	EcsStatusRunning,
	EcsStatusStarting,
	EcsStatusStopped,
	EcsStatusStopping,
	EcsStatusError,
	EcsStatusSnapshotting,
	EcsStatusUnsubscribed,
	EcsStatusUnsubscribing,
}
