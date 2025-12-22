package business

const (
	VpceServiceTypeInterface = "interface"
	VpceServiceTypeReverse   = "reverse"

	VpceServiceConnectionUp   = "up"
	VpceServiceConnectionDown = "down"

	IPVersionIpv4 = "ipv4"
	IPVersionIpv6 = "ipv6"
)

var IPVersionDict = map[string]int32{
	IPVersionIpv4: 4,
	IPVersionIpv6: 6,
}
