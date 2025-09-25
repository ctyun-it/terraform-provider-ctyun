package business

const (
	ZosAclPrivate         = "private"           // 私有
	ZosAclPublicRead      = "public-read"       // 公共读
	ZosAclPublicReadWrite = "public-read-write" // 公共读写

	ZosStorageTypeStandard   = "STANDARD"    // 标准
	ZosStorageTypeStandardIA = "STANDARD_IA" // 低频
	ZosStorageTypeGlacier    = "GLACIER"     // 归档

	ZosAzPolicySingle = "single-az" // 单az
	ZosAzPolicyMulti  = "multi-az"  //多az
)
