package business

const (
	LbResourceTypeInternal = "internal" //内网负载均衡
	LbResourceTypeExternal = "external" //公网负载均衡

	AdminStatusDown   = "down"   //管理状态: DOWN
	AdminStatusActive = "active" //管理状态: ACTIVE

	// elb rule 状态
	ElbRuleStatusACTIVE = "ACTIVE"
	ElbRuleStatusDOWN   = "DOWN"

	// condition 类型
	ElbRuleConditionTypeServerName = "server_name"
	ElbRuleConditionTypeUrlPath    = "url_path"
	// 匹配类型
	ElbRuleMatchTypeABSOLUTE = "ABSOLUTE"
	ElbRuleMatchTypePREFIX   = "PREFIX"
	ElbRuleMatchTypeREG      = "REG"

	ElbTargetIPStatusOffline = "offline"
	ElbTargetIPStatusOnline  = "online"
	ElbTargetIPStatusUnknown = "unknown"

	ElbTargetTypeVM = "VM"
	ElbTargetTypeBM = "BM"

	ElbRuleActionTypeForward  = "forward"
	ElbRuleActionTypeRedirect = "redirect"
	ElbRuleActionTypeDeny     = "deny"

	ElbTargetInstanceTypeVM  = "VM"
	ElbTargetInstanceTypeBM  = "BM"
	ElbTargetInstanceTypeECI = "ECI"
	ElbTargetInstanceTypeIP  = "IP"

	LbSlaNameSmall   = "elb.s1.small"
	LbSlaNameDefault = "elb.default"

	PgElbSlaNameS2Small = "elb.s2.small"
	PgElbSlaNameS3Small = "elb.s3.small"
	PgElbSlaNameS4Small = "elb.s4.small"
	PgElbSlaNameS5Small = "elb.s5.small"
	PgElbSlaNameS2Large = "elb.s2.large"
	PgElbSlaNameS3Large = "elb.s3.large"
	PgElbSlaNameS4Large = "elb.s4.large"
	PgElbSlaNameS5Large = "elb.s5.large"

	ELbCycleTypeMonth = "month"
	ELbCycleTypeYear  = "year"

	//资源状态常量
	ElbStatusStarted    = "started"     //启用
	ElbStatusRenewed    = "renewed"     //续订
	ElbStatusRefunded   = "refunded"    //退订
	ElbStatusDestroyed  = "destroyed"   //销毁
	ElbStatusFailed     = "failed"      //失败
	ElbStatusStarting   = "starting"    //正在启动
	ElbStatusChanged    = "changed"     //变配
	ElbStatusExpired    = "expired"     //过期
	ElbStatusUnknown    = "unknown"     //未知
	ElbStatusInProgress = "in_progress" //

	ListenerProtocolTCP   = "TCP"
	ListenerProtocolUDP   = "UDP"
	ListenerProtocolHTTP  = "HTTP"
	ListenerProtocolHTTPS = "HTTPS"

	//访问控制类型
	ListenerAccessControlTypeClose = "Close" //未启用
	ListenerAccessControlTypeWhite = "White" //白名单
	ListenerAccessControlTypeBlack = "Black" //黑名单

	ListenerDefaultActionTypeForward  = "forward"
	ListenerDefaultActionTypeRedirect = "redirect"

	// 后端主机组的调度算法
	TargetGroupAlgorithmRR  = "rr"  //轮询
	TargetGroupAlgorithmWRR = "wrr" //带权重轮询
	TargetGroupAlgorithmLC  = "lc"  //最少连接
	TargetGroupAlgorithmSH  = "sh"  //源IP哈希

	TargetGroupSessionStickyModeCLOSE    = "CLOSE"   //关闭
	TargetGroupSessionStickyModeINSERT   = "INSERT"  //插入
	TargetGroupSessionStickyModeREWRITE  = "REWRITE" //重写
	TargetGroupSessionStickyModeSourceIP = "SOURCE_IP"

	// 健康检查协议
	HealthCheckProtocolHTTP = "HTTP"
	HealthCheckProtocolTCP  = "TCP"
	HealthCheckProtocolUDP  = "UDP"

	//HTTP请求的方法
	HTTPMethodGET     = "GET"
	HTTPMethodPOST    = "POST"
	HTTPMethodPUT     = "PUT"
	HTTPMethodHEAD    = "HEAD"
	HTTPMethodDELETE  = "DELETE"
	HTTPMethodTRACE   = "TRACE"
	HTTPMethodOPTIONS = "OPTIONS"
	HTTPMethodPATCH   = "PATCH"
	HTTPMethodCONNECT = "CONNECT"

	CertificateTypeServer = "Server"
	CertificateTypeCA     = "Ca"
)

var LbResourceType = []string{LbResourceTypeInternal, LbResourceTypeExternal}
var AdminStatusName = []string{AdminStatusActive, AdminStatusDown}
var ElbRuleStatus = []string{ElbRuleStatusACTIVE, ElbRuleStatusDOWN}
var ElbRuleConditionTypes = []string{ElbRuleConditionTypeServerName, ElbRuleConditionTypeUrlPath}
var ElbRuleMatchTypes = []string{ElbRuleMatchTypeABSOLUTE, ElbRuleMatchTypePREFIX, ElbRuleMatchTypeREG}
var ElbTargetIpStatus = []string{ElbTargetIPStatusOffline, ElbTargetIPStatusOnline, ElbTargetIPStatusUnknown}
var ElbTargetType = []string{ElbTargetTypeVM, ElbTargetTypeBM}
var ElbRuleActionType = []string{ElbRuleActionTypeForward, ElbRuleActionTypeRedirect}
var ElbTargetInstanceType = []string{ElbTargetInstanceTypeVM, ElbTargetInstanceTypeBM, ElbTargetInstanceTypeECI}
var ElbSlaNames = []string{LbSlaNameSmall, LbSlaNameDefault}
var ElbCycleTypes = []string{ELbCycleTypeMonth, ELbCycleTypeYear}
var PgElbSlaNames = []string{PgElbSlaNameS2Small, PgElbSlaNameS3Small, PgElbSlaNameS4Small, PgElbSlaNameS5Small, PgElbSlaNameS2Large, PgElbSlaNameS3Large, PgElbSlaNameS4Large, PgElbSlaNameS5Large}
var ElbMasterResourceStatus = []string{
	ElbStatusStarted,
	ElbStatusRenewed,
	ElbStatusRefunded,
	ElbStatusDestroyed,
	ElbStatusFailed,
	ElbStatusStarting,
	ElbStatusChanged,
	ElbStatusExpired,
	ElbStatusUnknown,
	ElbStatusInProgress,
}
var ListenerProtocols = []string{ListenerProtocolTCP, ListenerProtocolUDP, ListenerProtocolHTTP, ListenerProtocolHTTPS}
var ListenerAccessControlTypes = []string{ListenerAccessControlTypeClose, ListenerAccessControlTypeWhite, ListenerAccessControlTypeBlack}
var ListenerDefaultActionTypes = []string{ListenerDefaultActionTypeForward, ListenerDefaultActionTypeRedirect}
var TargetGroupAlgorithms = []string{TargetGroupAlgorithmWRR, TargetGroupAlgorithmLC, TargetGroupAlgorithmSH}
var TargetGroupSessionStickyModes = []string{TargetGroupSessionStickyModeCLOSE, TargetGroupSessionStickyModeINSERT, TargetGroupSessionStickyModeREWRITE, TargetGroupSessionStickyModeSourceIP}
var HealthCheckProtocols = []string{HealthCheckProtocolUDP, HealthCheckProtocolTCP, HealthCheckProtocolHTTP}
var HttpMethods = []string{
	HTTPMethodGET,
	HTTPMethodPOST,
	HTTPMethodPUT,
	HTTPMethodHEAD,
	HTTPMethodDELETE,
	HTTPMethodTRACE,
	HTTPMethodOPTIONS,
	HTTPMethodPATCH,
	HTTPMethodCONNECT,
}
var CertificateTypes = []string{CertificateTypeServer, CertificateTypeCA}
