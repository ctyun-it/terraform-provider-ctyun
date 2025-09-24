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
}

var RedisEdition = []string{
	RedisEditionStandardSingle,
	RedisEditionStandardDual,
	RedisEditionDirectClusterSingle,
	RedisEditionDirectCluster,
	RedisEditionClusterOriginalProxy,
	RedisEditionOriginalMultipleReadLvs,
}
