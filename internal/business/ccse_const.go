package business

const (
	CcseClusterSeriesStandard = "cce.standard" // 专有版
	CcseClusterSeriesManaged  = "cce.managed"  // 托管版
	CcseClusterSeriesIcce     = "cce.icce"     // 智算版

	CcseRefundedBizState  = 4
	CcseRefundingBizState = 18

	CcsePluginCalico  = "calico"
	CcsePluginCubecni = "cubecni"

	CcseContainerRuntimeContainerd = "containerd"
	CcseContainerRuntimeDocker     = "docker"

	CcseDeployTypeSingle = "single"
	CcseDeployTypeMulti  = "multi"

	CcseKubeProxyIptables = "iptables"
	CcseKubeProxyIpvs     = "ipvs"

	CcseSeriesTypeManagedbase = "managedbase"
	CcseSeriesTypeManagedpro  = "managedpro"

	CcseSlaveInstanceTypeEcs = "ecs"
	CcseSlaveInstanceTypeEbm = "ebm"
)

var CcseApiServerElbSpecs = []string{"standardI", "standardII", "enhancedI", "enhancedII", "higherI"}

var CcseClusterVersions = []string{"1.31.6", "1.27.8", "1.29.3"}

var CcseDiskTypes = []string{"SATA", "SSD", "SAS"}
