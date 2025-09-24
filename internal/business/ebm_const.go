package business

const (
	EbmExtIpNotUse   = "0"
	EbmExtIpUseExist = "2"

	EbmOrderOnCycle  = "ORDER_ON_CYCLE"
	EbmOrderOnDemand = "ORDER_ON_DEMAND"

	EbmStatusCreating          = "creating"
	EbmStatusStarting          = "starting"
	EbmStatusRunning           = "running"
	EbmStatusStopping          = "stopping"
	EbmStatusStopped           = "stopped"
	EbmStatusRestarting        = "restarting"
	EbmStatusError             = "error"
	EbmStatusResettingPassword = "resetting_password"
	EbmStatusResettingHostname = "resetting_hostname"
)

var EbmDiskTypes = []string{
	EbsDiskTypeSata,
	EbsDiskTypeSas,
	EbsDiskTypeSsd,
}
