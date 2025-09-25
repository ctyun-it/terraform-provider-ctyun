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
	PgsqlProdIDS1417    = 10003013
	PgsqlProdIDMS1417   = 10003014
	PgsqlProdIDS1320    = 10003015
	PgsqlProdIDMS1320   = 10003016
	PgsqlProdIDRead1222 = 10003017
	PgsqlProdIDRead1320 = 10003018
	PgsqlProdIDRead1417 = 10003019
	PgsqlProdIDS1512    = 10003021
	PgsqlProdIDMS1512   = 10003022
	PgsqlProdIDRead1512 = 10003023
	PgsqlProdIDM2S1222  = 10003024
	PgsqlProdIDM2S1417  = 10003025
	PgsqlProdIDM2S1320  = 10003026
	PgsqlProdIDM2S1512  = 10003027
	PgsqlProdIDS168     = 10003028
	PgsqlProdIDMS168    = 10003029
	PgsqlProdIDM2S168   = 10003031
	PgsqlProdIDRead168  = 10003030

	PgsqlStorageTypeBackUp = "backup"
	PgsqlStorageTypeMaster = "master"
)

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
	PgsqlProdIDS1417,
	PgsqlProdIDMS1417,
	PgsqlProdIDS1320,
	PgsqlProdIDMS1320,
	PgsqlProdIDRead1222,
	PgsqlProdIDRead1320,
	PgsqlProdIDRead1417,
	PgsqlProdIDS1512,
	PgsqlProdIDMS1512,
	PgsqlProdIDRead1512,
	PgsqlProdIDM2S1222,
	PgsqlProdIDM2S1417,
	PgsqlProdIDM2S1320,
	PgsqlProdIDM2S1512,
	PgsqlProdIDS168,
	PgsqlProdIDMS168,
	PgsqlProdIDM2S168,
	PgsqlProdIDRead168,
}

var PgsqlProdIDDict = map[string]int64{
	"Single1222":       PgsqlProdIDS1222,
	"MasterSlave1222":  PgsqlProdIDMS1222,
	"Single1417":       PgsqlProdIDS1417,
	"MasterSlave1417":  PgsqlProdIDMS1417,
	"Single1320":       PgsqlProdIDS1320,
	"MasterSlave1320":  PgsqlProdIDMS1320,
	"ReadOnly1222":     PgsqlProdIDRead1222,
	"ReadOnly1320":     PgsqlProdIDRead1320,
	"ReadOnly1417":     PgsqlProdIDRead1417,
	"Single1512":       PgsqlProdIDS1512,
	"MasterSlave1512":  PgsqlProdIDMS1512,
	"ReadOnly1512":     PgsqlProdIDRead1512,
	"Master2Slave1222": PgsqlProdIDM2S1222,
	"Master2Slave1417": PgsqlProdIDM2S1417,
	"Master2Slave1320": PgsqlProdIDM2S1320,
	"Master2Slave1512": PgsqlProdIDM2S1512,
	"Single168":        PgsqlProdIDS168,
	"MasterSlave168":   PgsqlProdIDMS168,
	"Master2Slave168":  PgsqlProdIDM2S168,
	"ReadOnly168":      PgsqlProdIDRead168,
}

var PgsqlProdIDRevDict = map[int64]string{
	PgsqlProdIDS1222:    "Single1222",
	PgsqlProdIDMS1222:   "MasterSlave1222",
	PgsqlProdIDS1417:    "Single1417",
	PgsqlProdIDMS1417:   "MasterSlave1417",
	PgsqlProdIDS1320:    "Single1320",
	PgsqlProdIDMS1320:   "MasterSlave1320",
	PgsqlProdIDRead1222: "ReadOnly1222",
	PgsqlProdIDRead1320: "ReadOnly1320",
	PgsqlProdIDRead1417: "ReadOnly1417",
	PgsqlProdIDS1512:    "Single1512",
	PgsqlProdIDMS1512:   "MasterSlave1512",
	PgsqlProdIDRead1512: "ReadOnly1512",
	PgsqlProdIDM2S1222:  "Master2Slave1222",
	PgsqlProdIDM2S1417:  "Master2Slave1417",
	PgsqlProdIDM2S1320:  "Master2Slave1320",
	PgsqlProdIDM2S1512:  "Master2Slave1512",
	PgsqlProdIDS168:     "Single168",
	PgsqlProdIDMS168:    "MasterSlave168",
	PgsqlProdIDM2S168:   "Master2Slave168",
	PgsqlProdIDRead168:  "ReadOnly168",
}

var PgsqlProdIds = []string{
	"Single1222",
	"MasterSlave1222",
	"Single1417",
	"MasterSlave1417",
	"Single1320",
	"MasterSlave1320",
	"ReadOnly1222",
	"ReadOnly1320",
	"ReadOnly1417",
	"Single1512",
	"MasterSlave1512",
	"ReadOnly1512",
	"Master2Slave1222",
	"Master2Slave1417",
	"Master2Slave1320",
	"Master2Slave1512",
	"Single168",
	"MasterSlave168",
	"Master2Slave168",
	"ReadOnly168",
}

var PgsqlNodeTypeDict = map[string]string{
	"Single1222":       "master",
	"MasterSlave1222":  "master",
	"Single1417":       "master",
	"MasterSlave1417":  "master",
	"Single1320":       "master",
	"MasterSlave1320":  "master",
	"ReadOnly1222":     "readNode",
	"ReadOnly1320":     "readNode",
	"ReadOnly1417":     "readNode",
	"Single1512":       "master",
	"MasterSlave1512":  "master",
	"ReadOnly1512":     "readNode",
	"Master2Slave1222": "master",
	"Master2Slave1417": "master",
	"Master2Slave1320": "master",
	"Master2Slave1512": "master",
	"Single168":        "readNode",
	"MasterSlave168":   "master",
	"Master2Slave168":  "master",
	"ReadOnly168":      "readNode",
}

var PgsqlNodeNumDict = map[string]int32{
	"Single1222":       1,
	"MasterSlave1222":  2,
	"Single1417":       1,
	"MasterSlave1417":  2,
	"Single1320":       1,
	"MasterSlave1320":  2,
	"ReadOnly1222":     -1,
	"ReadOnly1320":     -1,
	"ReadOnly1417":     -1,
	"Single1512":       1,
	"MasterSlave1512":  2,
	"ReadOnly1512":     -1,
	"Master2Slave1222": 3,
	"Master2Slave1417": 3,
	"Master2Slave1320": 3,
	"Master2Slave1512": 3,
	"Single168":        1,
	"MasterSlave168":   2,
	"Master2Slave168":  3,
	"ReadOnly168":      -1,
}

var PgsqlProdVersionDict = map[string]string{
	"Single1222":       "12.22",
	"MasterSlave1222":  "12.22",
	"Single1417":       "14.17",
	"MasterSlave1417":  "14.17",
	"Single1320":       "13.20",
	"MasterSlave1320":  "13.20",
	"ReadOnly1222":     "12.22",
	"ReadOnly1320":     "13.20",
	"ReadOnly1417":     "14.17",
	"Single1512":       "15.12",
	"MasterSlave1512":  "15.12",
	"ReadOnly1512":     "15.12",
	"Master2Slave1222": "12.22",
	"Master2Slave1417": "14.17",
	"Master2Slave1320": "13.20",
	"Master2Slave1512": "15.12",
	"Single168":        "16.8",
	"MasterSlave168":   "16.8",
	"Master2Slave168":  "16.8",
	"ReadOnly168":      "16.8",
}
var PgsqlInstanceSeriesDict = map[string]string{
	"S": "1",
	"C": "2",
	"M": "3",
}
