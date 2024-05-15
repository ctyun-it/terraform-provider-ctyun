package business

import "terraform-provider-ctyun/internal/utils"

const (
	SecurityGroupRuleDirectionEgress  = "egress"
	SecurityGroupRuleDirectionIngress = "ingress"

	SecurityGroupRuleActionAccept = "accept"
	SecurityGroupRuleActionDrop   = "drop"

	SecurityGroupRuleProtocolAny  = "any"
	SecurityGroupRuleProtocolTcp  = "tcp"
	SecurityGroupRuleProtocolUdp  = "udp"
	SecurityGroupRuleProtocolIcmp = "icmp"

	SecurityGroupRuleEtherTypeIpv4 = "ipv4"
	SecurityGroupRuleEtherTypeIpv6 = "ipv6"
)

const (
	SecurityGroupRuleDirectionMapScene1 = iota
)

const (
	SecurityGroupRuleActionMapScene1 = iota
)

const (
	SecurityGroupRuleProtocolMapScene1 = iota
)

const (
	SecurityGroupRuleEtherTypeMapScene1 = iota
)

var SecurityGroupRuleDirections = []string{
	SecurityGroupRuleDirectionEgress,
	SecurityGroupRuleDirectionIngress,
}

var SecurityGroupRuleActions = []string{
	SecurityGroupRuleActionAccept,
	SecurityGroupRuleActionDrop,
}

var SecurityGroupRuleProtocols = []string{
	SecurityGroupRuleProtocolAny,
	SecurityGroupRuleProtocolTcp,
	SecurityGroupRuleProtocolUdp,
	SecurityGroupRuleProtocolIcmp,
}

var SecurityGroupRuleEtherTypes = []string{
	SecurityGroupRuleEtherTypeIpv4,
	SecurityGroupRuleEtherTypeIpv6,
}

var SecurityGroupRuleDirectionMap = utils.Must(
	[]any{
		SecurityGroupRuleDirectionEgress,
		SecurityGroupRuleDirectionIngress,
	},
	map[utils.Scene][]any{
		SecurityGroupRuleDirectionMapScene1: {
			"EGRESS",
			"INGRESS",
		},
	},
)

var SecurityGroupRuleActionMap = utils.Must(
	[]any{
		SecurityGroupRuleActionAccept,
		SecurityGroupRuleActionDrop,
	},
	map[utils.Scene][]any{
		SecurityGroupRuleActionMapScene1: {
			"ACCEPT",
			"DROP",
		},
	},
)

var SecurityGroupRuleProtocolMap = utils.Must(
	[]any{
		SecurityGroupRuleProtocolAny,
		SecurityGroupRuleProtocolTcp,
		SecurityGroupRuleProtocolUdp,
		SecurityGroupRuleProtocolIcmp,
	},
	map[utils.Scene][]any{
		SecurityGroupRuleProtocolMapScene1: {
			"ANY",
			"TCP",
			"UDP",
			"ICMP",
		},
	},
)

var SecurityGroupRuleEtherTypeMap = utils.Must(
	[]any{
		SecurityGroupRuleEtherTypeIpv4,
		SecurityGroupRuleEtherTypeIpv6,
	},
	map[utils.Scene][]any{
		SecurityGroupRuleEtherTypeMapScene1: {
			"IPv4",
			"IPv6",
		},
	},
)

// IsEgress 判断是否为出方向
func IsEgress(direction string) bool {
	return direction == SecurityGroupRuleDirectionEgress
}

// IsIngress 判断是否为入方向
func IsIngress(direction string) bool {
	return direction == SecurityGroupRuleDirectionIngress
}
