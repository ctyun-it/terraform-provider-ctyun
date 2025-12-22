package business

const (
	AclEnable             = "enable"
	AclDisable            = "disable"
	AclDirectionIngress   = "ingress"
	AclDirectionEgress    = "egress"
	PrefixAddressTypeIpv4 = "ipv4"
	PrefixAddressTypeIpv6 = "ipv6"

	AclRuleProtocolTCP  = "tcp"
	AclRuleProtocolUDP  = "udp"
	AclRuleProtocolICMP = "icmp"
	AclRuleProtocolALL  = "all"
	AclRuleActionAccept = "accept"
	AclRuleActionDrop   = "drop"

	AclRuleEnable  = "enable"
	AclRuleDisable = "disable"
)

var PrefixAddressTypeMap = map[string]int32{
	PrefixAddressTypeIpv4: 4,
	PrefixAddressTypeIpv6: 6,
}

var PrefixAddressTyperRevMap = map[int32]string{
	4: PrefixAddressTypeIpv4,
	6: PrefixAddressTypeIpv6,
}

var AclRuleProtocols = []string{AclRuleProtocolTCP, AclRuleProtocolUDP, AclRuleProtocolICMP, AclRuleProtocolALL}
