package opensearch

// EnableIPv6 constants
const (
	EnableIPv6Open  = "OPEN"
	EnableIPv6Close = "CLOSE"
)

// ClusterType constants
const (
	ClusterTypeOpenSearch    = 1
	ClusterTypeElasticsearch = 2
)

// OSType constants
const (
	OSTypeCTyun = "CTyun"
	OSTypeKylin = "Kylin"
)

// EnableHTTPS constants
const (
	EnableHTTPSOpen  = "OPEN"
	EnableHTTPSClose = "CLOSE"
)

// PayType constants
const (
	PayTypePrepaid  = 1
	PayTypePostpaid = 2
)

// CycleType constants
const (
	CycleTypeMonthly = 2
	CycleTypeYearly  = 3
)

// CycleTypeString constants (for Terraform string attribute)
const (
	CycleTypeMonthlyStr  = "month"
	CycleTypeYearlyStr   = "year"
	CycleTypeOnDemandStr = "on_demand"
)

// IOType constants
const (
	IOTypeSSDGeneric = "SSD-genric"
	IOTypeSAS        = "SAS"
	IOTypeSSD        = "SSD"
	IOTypeXSSD0      = "XSSD-0"
	IOTypeXSSD1      = "XSSD-1"
)

// NodeGroupType constants
const (
	NodeGroupTypeMaster          = "MASTER"
	NodeGroupTypeExclusiveMaster = "EXCLUSIVE_MASTER"
	NodeGroupTypeCoordinate      = "COORDINATE"
	NodeGroupTypeCold            = "COLD"
)
