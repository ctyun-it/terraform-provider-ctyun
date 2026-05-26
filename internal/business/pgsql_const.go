package business

const (
	BillModeCycle          = "1"
	BillModelOnDemand      = "2"
	CaseSensitiveTrue      = "0"
	CaseSensitiveFalse     = "1"
	CaseSensitiveUnCertain = "2"
	OsTypePure             = "0"  // 裸机
	OsTypeWindows          = "1"  // Windows
	OsTypeCentos           = "2"  // Centos
	OsTypeUbuntu           = "3"  // Ubuntu
	OsTypeAndroid          = "4"  // Android
	OsTypeRedHat           = "5"  // RedHat
	OsTypeKylin            = "6"  // kylin
	OsTypeUos              = "7"  // Uos
	OsTypeSuse             = "8"  // Suse
	OsTypeAsianux          = "9"  // Asianus
	OsTypeOpenEuler        = "10" // OpenEuler
	OsTypeCtyunOS          = "11" // CtyunOS
	OsTypeEuler            = "12" // Euler

	PgsqlProdRunningStatusStarted             = 0
	pgsqlProdRunningStatusRestarting          = 1
	PgsqlProdRunningStatusBackup              = 2
	PgsqlProdRunningStatusRecovering          = 3
	PgsqlProdRunningStatusStopped             = 1001
	PgsqlProdRunningStatusRecoveryFailed      = 1006
	PgsqlProdRunningStatusVipUnavailable      = 1007
	PgsqlProdRunningStatusGatewayUnavailable  = 1008
	PgsqlProdRunningStatusMasterUnavailable   = 1009
	PgsqlProdRunningStatusSlaveUnavailable    = 1010
	PgsqlProdRunningStatusInstanceMaintenance = 1021
	PgsqlProdRunningStatusActivating          = 2000
	PgsqlProdRunningStatusUnsubscribed        = 2002
	PgsqlProdRunningStatusExpanding           = 2005
	PgsqlProdRunningStatusFreeze              = 2011

	PgsqlProdOrderStatusRunning    = 0
	PgsqlProdOrderStatusFreeze     = 1
	PgsqlProdOrderStatusDelete     = 2
	PgsqlProdOrderStatusProcessing = 3
	PgsqlProdOrderStatusFailure    = 4
	PgsqlProdOrderStatusExpanding  = 5

	PgsqlBindEipStatusACTIVE              = "ACTIVE"                //已使用
	PgsqlBindEipStatusDOWN                = "DOWN"                  //未使用
	PgsqlBindEipStatusERROR               = "ERROR"                 //中间状态-异常
	PgsqlBindEipStatusUPDATING            = "UPDATING"              //中间状态-更新中
	PgsqlBindEipStatusBANDINGORUNBANGDING = "BANDING_OR_UNBANGDING" //中间状态-绑定或解绑中
	PgsqlBindEipStatusDELETING            = "DELETING"              //中间状态-删除中
	PgsqlBindEipStatusDELETED             = "DELETED"

	PgsqlProdIDS1222    = 10003011
	PgsqlProdIDMS1222   = 10003012
	PgsqlProdIDS1419    = 10003013
	PgsqlProdIDMS1419   = 10003014
	PgsqlProdIDS1322    = 10003015
	PgsqlProdIDMS1322   = 10003016
	PgsqlProdIDRead1222 = 10003017
	PgsqlProdIDRead1322 = 10003018
	PgsqlProdIDRead1419 = 10003019
	PgsqlProdIDS1514    = 10003021
	PgsqlProdIDMS1514   = 10003022
	PgsqlProdIDRead1514 = 10003023
	PgsqlProdIDM2S1222  = 10003024
	PgsqlProdIDM2S1419  = 10003025
	PgsqlProdIDM2S1322  = 10003026
	PgsqlProdIDM2S1514  = 10003027
	PgsqlProdIDS1610    = 10003028
	PgsqlProdIDMS1610   = 10003029
	PgsqlProdIDM2S1610  = 10003031
	PgsqlProdIDRead1610 = 10003030

	PgsqlStorageTypeBackUp = "backup"
	PgsqlStorageTypeMaster = "master"

	PgsqlNodeTypeMaster = "master"

	PgsqlAccountTypeNormal   = "normal"
	PgsqlAccountTypeAdvanced = "advanced"

	PgsqlBackupTypeAuto     = 1
	PgsqlBackupTypeManual   = 2
	PgsqlBackupTypeRecovery = 3

	PgsqlBackupTypeAutoStr     = "auto"
	PgsqlBackupTypeManualStr   = "manual"
	PgsqlBackupTypeRecoveryStr = "recovery"

	PgsqlBackupResultING     = 0
	PgsqlBackupResultING1    = 2
	PgsqlBackupResultSuccess = 3
	PgsqlBackupResultFail    = 5

	PgsqlBackupResultINGStr     = "ing"
	PgsqlBackupResultSuccessStr = "success"
	PgsqlBackupResultFailStr    = "fail"

	PgsqlNodeTypeReadNode = "readNode"

	PgsqlProdSpecNameRead   = "只读实例"
	PgsqlProdSpecNameSingle = "单机版"
	PgsqlProdSpecNameMS     = "一主一备"
	PgsqlProdSpecNameM2S    = "一主两备"
)

var PgsqlBackupTypeMap = map[string]int32{
	PgsqlBackupTypeAutoStr:     PgsqlBackupTypeAuto,
	PgsqlBackupTypeManualStr:   PgsqlBackupTypeManual,
	PgsqlBackupTypeRecoveryStr: PgsqlBackupTypeRecovery,
}

var PgsqlBackupTypeMapConv = map[int32]string{
	PgsqlBackupTypeAuto:     PgsqlBackupTypeAutoStr,
	PgsqlBackupTypeManual:   PgsqlBackupTypeManualStr,
	PgsqlBackupTypeRecovery: PgsqlBackupTypeRecoveryStr,
}

var PgsqlBackupResultMap = map[string]int32{
	PgsqlBackupResultINGStr:     PgsqlBackupResultING,
	PgsqlBackupResultSuccessStr: PgsqlBackupResultSuccess,
	PgsqlBackupResultFailStr:    PgsqlBackupResultFail,
}
var PgsqlBackupResultMapConv = map[int32]string{
	PgsqlBackupResultING:     PgsqlBackupResultINGStr,
	PgsqlBackupResultING1:    PgsqlBackupResultINGStr,
	PgsqlBackupResultFail:    PgsqlBackupResultFailStr,
	PgsqlBackupResultSuccess: PgsqlBackupResultSuccessStr,
}

var PgsqlBillModes = []string{
	BillModeCycle,
	BillModelOnDemand,
}

var PgsqlCaseSensitive = []string{
	CaseSensitiveTrue,
	CaseSensitiveFalse,
	CaseSensitiveUnCertain,
}

var PgsqlOsType = []string{
	OsTypePure,
	OsTypeWindows,
	OsTypeCentos,
	OsTypeUbuntu,
	OsTypeAndroid,
	OsTypeRedHat,
	OsTypeKylin,
	OsTypeUos,
	OsTypeSuse,
	OsTypeAsianux,
	OsTypeOpenEuler,
	OsTypeCtyunOS,
	OsTypeEuler,
}

var PgsqlProdOrderStatus = []int32{
	PgsqlProdOrderStatusRunning,
	PgsqlProdOrderStatusFreeze,
	PgsqlProdOrderStatusDelete,
	PgsqlProdOrderStatusProcessing,
	PgsqlProdOrderStatusFailure,
	PgsqlProdOrderStatusExpanding,
}

var PgsqlProdRunningStatus = []int32{
	PgsqlProdRunningStatusStarted,
	pgsqlProdRunningStatusRestarting,
	PgsqlProdRunningStatusBackup,
	PgsqlProdRunningStatusRecovering,
	PgsqlProdRunningStatusStopped,
	PgsqlProdRunningStatusRecoveryFailed,
	PgsqlProdRunningStatusVipUnavailable,
	PgsqlProdRunningStatusGatewayUnavailable,
	PgsqlProdRunningStatusMasterUnavailable,
	PgsqlProdRunningStatusSlaveUnavailable,
	PgsqlProdRunningStatusInstanceMaintenance,
	PgsqlProdRunningStatusActivating,
	PgsqlProdRunningStatusUnsubscribed,
	PgsqlProdRunningStatusExpanding,
	PgsqlProdRunningStatusFreeze,
}

var PgsqlBindEipStatus = []string{
	MysqlBindEipStatusACTIVE,
	MysqlBindEipStatusDOWN,
	MysqlBindEipStatusERROR,
	MysqlBindEipStatusUPDATING,
	MysqlBindEipStatusBANDINGORUNBANGDING,
	MysqlBindEipStatusDELETING,
	MysqlBindEipStatusDELETED,
}

var PgsqlProdID = []int64{
	PgsqlProdIDS1222,
	PgsqlProdIDMS1222,
	PgsqlProdIDS1419,
	PgsqlProdIDMS1419,
	PgsqlProdIDS1322,
	PgsqlProdIDMS1322,
	PgsqlProdIDRead1222,
	PgsqlProdIDRead1322,
	PgsqlProdIDRead1419,
	PgsqlProdIDS1514,
	PgsqlProdIDMS1514,
	PgsqlProdIDRead1514,
	PgsqlProdIDM2S1222,
	PgsqlProdIDM2S1419,
	PgsqlProdIDM2S1322,
	PgsqlProdIDM2S1514,
	PgsqlProdIDS1610,
	PgsqlProdIDMS1610,
	PgsqlProdIDM2S1610,
	PgsqlProdIDRead1610,
}

var PgsqlProdIDDict = map[string]int64{
	"Single1222":       PgsqlProdIDS1222,
	"MasterSlave1222":  PgsqlProdIDMS1222,
	"Single1419":       PgsqlProdIDS1419,
	"MasterSlave1419":  PgsqlProdIDMS1419,
	"Single1322":       PgsqlProdIDS1322,
	"MasterSlave1322":  PgsqlProdIDMS1322,
	"ReadOnly1222":     PgsqlProdIDRead1222,
	"ReadOnly1322":     PgsqlProdIDRead1322,
	"ReadOnly1419":     PgsqlProdIDRead1419,
	"Single1514":       PgsqlProdIDS1514,
	"MasterSlave1514":  PgsqlProdIDMS1514,
	"ReadOnly1514":     PgsqlProdIDRead1514,
	"Master2Slave1222": PgsqlProdIDM2S1222,
	"Master2Slave1419": PgsqlProdIDM2S1419,
	"Master2Slave1322": PgsqlProdIDM2S1322,
	"Master2Slave1514": PgsqlProdIDM2S1514,
	"Single1610":       PgsqlProdIDS1610,
	"MasterSlave1610":  PgsqlProdIDMS1610,
	"Master2Slave1610": PgsqlProdIDM2S1610,
	"ReadOnly1610":     PgsqlProdIDRead1610,
}

var PgsqlReadNodeVersionProdIdDict = map[string]string{
	"12.22": "ReadOnly1222",
	"14.19": "ReadOnly1419",
	"15.14": "ReadOnly1514",
	"13.22": "ReadOnly1322",
	"16.10": "ReadOnly1610",
}

var PgsqlProdIDRevDict = map[int64]string{
	PgsqlProdIDS1222:    "Single1222",
	PgsqlProdIDMS1222:   "MasterSlave1222",
	PgsqlProdIDS1419:    "Single1419",
	PgsqlProdIDMS1419:   "MasterSlave1419",
	PgsqlProdIDS1322:    "Single1322",
	PgsqlProdIDMS1322:   "MasterSlave1322",
	PgsqlProdIDRead1222: "ReadOnly1222",
	PgsqlProdIDRead1322: "ReadOnly1322",
	PgsqlProdIDRead1419: "ReadOnly1419",
	PgsqlProdIDS1514:    "Single1514",
	PgsqlProdIDMS1514:   "MasterSlave1514",
	PgsqlProdIDRead1514: "ReadOnly1514",
	PgsqlProdIDM2S1222:  "Master2Slave1222",
	PgsqlProdIDM2S1419:  "Master2Slave1419",
	PgsqlProdIDM2S1322:  "Master2Slave1322",
	PgsqlProdIDM2S1514:  "Master2Slave1514",
	PgsqlProdIDS1610:    "Single1610",
	PgsqlProdIDMS1610:   "MasterSlave1610",
	PgsqlProdIDM2S1610:  "Master2Slave1610",
	PgsqlProdIDRead1610: "ReadOnly1610",
}

var PgsqlProdIds = []string{
	"Single1222",
	"MasterSlave1222",
	"Single1419",
	"MasterSlave1419",
	"Single1322",
	"MasterSlave1322",
	"ReadOnly1222",
	"ReadOnly1322",
	"ReadOnly1419",
	"Single1514",
	"MasterSlave1514",
	"ReadOnly1514",
	"Master2Slave1222",
	"Master2Slave1419",
	"Master2Slave1322",
	"Master2Slave1514",
	"Single1610",
	"MasterSlave1610",
	"Master2Slave1610",
	"ReadOnly1610",
}

var PgsqlNodeTypeDict = map[string]string{
	"Single1222":       "master",
	"MasterSlave1222":  "master",
	"Single1419":       "master",
	"MasterSlave1419":  "master",
	"Single1322":       "master",
	"MasterSlave1322":  "master",
	"ReadOnly1222":     "readNode",
	"ReadOnly1322":     "readNode",
	"ReadOnly1419":     "readNode",
	"Single1514":       "master",
	"MasterSlave1514":  "master",
	"ReadOnly1514":     "readNode",
	"Master2Slave1222": "master",
	"Master2Slave1419": "master",
	"Master2Slave1322": "master",
	"Master2Slave1514": "master",
	"Single1610":       "readNode",
	"MasterSlave1610":  "master",
	"Master2Slave1610": "master",
	"ReadOnly1610":     "readNode",
}

var PgsqlNodeNumDict = map[string]int32{
	"Single1222":       1,
	"MasterSlave1222":  2,
	"Single1419":       1,
	"MasterSlave1419":  2,
	"Single1322":       1,
	"MasterSlave1322":  2,
	"ReadOnly1222":     -1,
	"ReadOnly1322":     -1,
	"ReadOnly1419":     -1,
	"Single1514":       1,
	"MasterSlave1514":  2,
	"ReadOnly1514":     -1,
	"Master2Slave1222": 3,
	"Master2Slave1419": 3,
	"Master2Slave1322": 3,
	"Master2Slave1514": 3,
	"Single1610":       1,
	"MasterSlave1610":  2,
	"Master2Slave1610": 3,
	"ReadOnly1610":     -1,
}

var PgsqlProdVersionDict = map[string]string{
	"Single1222":       "12.22",
	"MasterSlave1222":  "12.22",
	"Single1419":       "14.19",
	"MasterSlave1419":  "14.19",
	"Single1322":       "13.22",
	"MasterSlave1322":  "13.22",
	"ReadOnly1222":     "12.22",
	"ReadOnly1322":     "13.22",
	"ReadOnly1419":     "14.19",
	"Single1514":       "15.14",
	"MasterSlave1514":  "15.14",
	"ReadOnly1514":     "15.14",
	"Master2Slave1222": "12.22",
	"Master2Slave1419": "14.19",
	"Master2Slave1322": "13.22",
	"Master2Slave1514": "15.14",
	"Single1610":       "16.10",
	"MasterSlave1610":  "16.10",
	"Master2Slave1610": "16.10",
	"ReadOnly1610":     "16.10",
}
var PgsqlInstanceSeriesDict = map[string]string{
	"S": "1",
	"C": "2",
	"M": "3",
}

var PgsqlAccountTypes = []string{
	PgsqlAccountTypeAdvanced,
	PgsqlAccountTypeNormal,
}

var PgsqlProdSpecNodeNumDict = map[string]int32{
	PgsqlProdSpecNameRead:   1,
	PgsqlProdSpecNameSingle: 1,
	PgsqlProdSpecNameMS:     2,
	PgsqlProdSpecNameM2S:    3,
}
