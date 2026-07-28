package business

const (
	// 资源类型常量定义
	ResourceTypeVpc                 = "vpc"
	ResourceTypeSubnet              = "subnet"
	ResourceTypeAcl                 = "acl"
	ResourceTypeSecurityGroup       = "security_group"
	ResourceTypeRouteTable          = "route_table"
	ResourceTypeHavip               = "havip"
	ResourceTypePort                = "port"
	ResourceTypeMulticastDomain     = "multicast_domain"
	ResourceTypeVpcPeer             = "vpc_peer"
	ResourceTypeVpceEndpoint        = "vpce_endpoint"
	ResourceTypeVpceEndpointService = "vpce_endpoint_service"
	ResourceTypeIPv6Gateway         = "ipv6_gateway"
	ResourceTypeElb                 = "elb"
	ResourceTypePrivateNat          = "private_nat"
	ResourceTypeNat                 = "nat"
	ResourceTypeEip                 = "eip"
	ResourceTypeBandwidth           = "bandwidth"
	ResourceTypeIPv6Bandwidth       = "ipv6_bandwidth"
	ResourceTypeListener            = "listener"
	ResourceTypePrefixList          = "prefix_list"
)

var NetResourceTypes = []string{
	ResourceTypeVpc,
	ResourceTypeSubnet,
	ResourceTypeAcl,
	ResourceTypeSecurityGroup,
	ResourceTypeRouteTable,
	ResourceTypeHavip,
	ResourceTypePort,
	ResourceTypeMulticastDomain,
	ResourceTypeVpcPeer,
	ResourceTypeVpceEndpoint,
	ResourceTypeVpceEndpointService,
	ResourceTypeIPv6Gateway,
	ResourceTypeElb,
	ResourceTypePrivateNat,
	ResourceTypeNat,
	ResourceTypeEip,
	ResourceTypeBandwidth,
	ResourceTypeIPv6Bandwidth,
	ResourceTypeListener,
	ResourceTypePrefixList,
}
