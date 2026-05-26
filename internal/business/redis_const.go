package business

const (
	RedisVersionBasic = "BASIC"
	RedisVersionPlus  = "PLUS"

	RedisDiskTypeSas = "SAS"
	RedisDiskTypeSsd = "SSD"

	RedisHostTypeS  = "S"
	RedisHostTypeC  = "C"
	RedisHostTypeM  = "M"
	RedisHostTypeHS = "HS"
	RedisHostTypeHC = "HC"
	RedisHostTypeKS = "KS"
	RedisHostTypeKC = "KC"
	RedisHostTypeFS = "FS"
	RedisHostTypeFC = "FC"

	RedisEditionStandardSingle          = "StandardSingle"          // 单机
	RedisEditionStandardDual            = "StandardDual"            // 主备
	RedisEditionDirectClusterSingle     = "DirectClusterSingle"     // Cluster单机
	RedisEditionDirectCluster           = "DirectCluster"           // Cluster主备
	RedisEditionClusterOriginalProxy    = "ClusterOriginalProxy"    // Proxy集群
	RedisEditionOriginalMultipleReadLvs = "OriginalMultipleReadLvs" // 读写分离

	RedisStatusRunning          = 0 // 运行中
	RedisStatusActivationFailed = 4 // 运行中
	RedisStatusUnsubscribed     = 8 // 已退订
)

var RedisEngineVersion = []string{"5.0", "6.0", "7.0"}

var RedisHostType = []string{
	RedisHostTypeS,
	RedisHostTypeC,
	RedisHostTypeM,
	RedisHostTypeHS,
	RedisHostTypeHC,
	RedisHostTypeKS,
	RedisHostTypeKC,
	RedisHostTypeFS,
	RedisHostTypeFC,
}

var RedisEdition = []string{
	RedisEditionStandardSingle,
	RedisEditionStandardDual,
	RedisEditionDirectClusterSingle,
	RedisEditionDirectCluster,
	RedisEditionClusterOriginalProxy,
	RedisEditionOriginalMultipleReadLvs,
}

var RedisHostTypeMap = map[string]string{
	"通用型":         RedisHostTypeS,
	"计算增强型":     RedisHostTypeC,
	"内存型":         RedisHostTypeM,
	"海光通用型":     RedisHostTypeHS,
	"海光计算增强型": RedisHostTypeHC,
	"鲲鹏通用型":     RedisHostTypeKS,
	"鲲鹏计算增强型": RedisHostTypeKC,
	"飞腾通用型":     RedisHostTypeFS,
	"飞腾计算增强型": RedisHostTypeFC,
}

var RedisTypeToApiEdition = map[int]string{
	5:  RedisEditionDirectCluster,
	6:  RedisEditionDirectClusterSingle,
	8:  RedisEditionDirectCluster,
	9:  RedisEditionDirectClusterSingle,
	10: RedisEditionStandardDual,
	11: RedisEditionStandardSingle,
	12: RedisEditionStandardDual,
	13: RedisEditionStandardSingle,
	14: RedisEditionOriginalMultipleReadLvs,
	16: RedisEditionOriginalMultipleReadLvs,
	17: RedisEditionClusterOriginalProxy,
	18: RedisEditionClusterOriginalProxy,
}

var RedisTypeToApiVersion = map[int]string{
	5:  "BASIC",
	6:  "BASIC",
	8:  "PLUS",
	9:  "PLUS",
	10: "BASIC",
	11: "BASIC",
	12: "PLUS",
	13: "PLUS",
	14: "BASIC",
	16: "PLUS",
	17: "BASIC",
	18: "PLUS",
}
