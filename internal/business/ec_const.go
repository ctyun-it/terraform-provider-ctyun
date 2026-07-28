package business

const (
	ECRouteTypeAuto   = "auto"
	ECRouteTypeCustom = "custom"

	EcNextHopTypeVPC    = "vpc"
	EcNextHopCDA        = "cda"
	EcNextHopSDWAN      = "sdwan"
	EcNextHopVPN        = "vpn"
	EcNextHopEDS        = "eds"
	EcNextHopBlackHole  = "black_hole" // 黑洞路由
	EcNextHopBlackCross = "cross"      // 跨域连接

	EcIpVersionIpv4 = "ipv4"
	EcIpVersionIpv6 = "ipv6"

	EcRouteStatusRunning = "running"
	EcRouteStatusStop    = "stop"
)

var EcRouteTypeMap = map[string]string{
	ECRouteTypeAuto:   "1",
	ECRouteTypeCustom: "2",
}

var EcRouteTypeRevMap = map[string]string{
	"1": ECRouteTypeAuto,
	"2": ECRouteTypeCustom,
}
var EcNextHopTypeMap = map[string]string{
	EcNextHopTypeVPC:    "1",
	EcNextHopCDA:        "2",
	EcNextHopSDWAN:      "3",
	EcNextHopVPN:        "4",
	EcNextHopEDS:        "5",
	EcNextHopBlackHole:  "20",
	EcNextHopBlackCross: "30",
}

var EcNextHopTypeRevMap = map[string]string{
	"1":  EcNextHopTypeVPC,
	"2":  EcNextHopCDA,
	"3":  EcNextHopSDWAN,
	"4":  EcNextHopVPN,
	"5":  EcNextHopEDS,
	"20": EcNextHopBlackHole,
	"30": EcNextHopBlackCross,
	"":   "",
}

var EcIpVersionMap = map[string]string{
	EcIpVersionIpv4: "1",
	EcIpVersionIpv6: "2",
}

var EcIpVersionRevMap = map[string]string{
	"1": EcIpVersionIpv4,
	"2": EcIpVersionIpv6,
}

var EcRouteStatusMap = map[string]string{
	EcRouteStatusRunning: "1",
	EcRouteStatusStop:    "2",
}

var EcRouteStatusRevMap = map[string]string{
	"1": EcRouteStatusRunning,
	"2": EcRouteStatusStop,
}
