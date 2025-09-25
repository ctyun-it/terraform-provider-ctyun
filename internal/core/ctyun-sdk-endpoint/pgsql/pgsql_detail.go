package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDetailApi(client *ctyunsdk.CtyunClient) *PgsqlDetailApi {
	return &PgsqlDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/product/get-paas-product",
		},
	}
}

func (this *PgsqlDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDetailRequest, header *PgsqlDetailRequestHeader) (detailResp *PgsqlDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	if req.ProdInstId == "" {
		err = errors.New("missing required field: ProdInstId")
		return
	}
	builder.AddParam("prodInstId", req.ProdInstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	detailResp = &PgsqlDetailResponse{}
	err = resp.Parse(detailResp)
	if err != nil {
		return
	}
	return detailResp, nil
}

type PgsqlDetailRequest struct {
	ProdInstId string `json:"prodInstId"`
}

type PgsqlDetailResponse struct {
	StatusCode int32                         `json:"statusCode"` // 接口状态码，参考下方状态码
	Error      string                        `json:"error"`      // 错误码
	Message    string                        `json:"message"`    // 描述信息
	ReturnObj  *PgsqlDetailResponseReturnObj `json:"returnObj"`  // 返回对象，包含具体的返回数据
}

type PgsqlDetailRequestHeader struct {
	ProjectID *string `json:"projectId"`
	RegionID  string  `json:"regionId"`
}
type PgsqlDetailResponseReturnObj struct {
	Alive               int32  `json:"alive"`               // 实例是否存活,0:存活，-1:异常
	CreateTime          string `json:"createTime"`          // 创建时间，格式：yyyy:MM:dd HH:mm:ss
	DiskRated           int32  `json:"diskRated"`           // 磁盘使用率
	DiskSize            int32  `json:"diskSize"`            // 实例数据最大磁盘空间，单位G
	DiskTotal           string `json:"diskTotal"`           // 磁盘总大小
	DiskType            string `json:"diskType"`            // 磁盘类型，例如：SATA，SSD
	DiskUsed            string `json:"diskUsed"`            // 已使用磁盘
	ExpireTime          string `json:"expireTime"`          // 过期时间，格式：yyyy-MM-dd HH:mm:ss
	MachineSpec         string `json:"machineSpec"`         // 实例机器规格
	NetName             string `json:"netName"`             // VPC名称
	OrderId             int64  `json:"orderId"`             // 订单id
	OuterProdInstId     string `json:"outerProdInstId"`     // 对外的实例ID，对应PaaS平台
	ProdDbEngine        string `json:"prodDbEngine"`        // 数据库实例引擎
	ProdInstFlag        string `json:"prodInstFlag"`        // 实例标识
	ProdInstId          string `json:"prodInstId"`          // 实例id
	ProdInstName        string `json:"prodInstName"`        // 实例名称
	ProdInstSetName     string `json:"prodInstSetName"`     // 实例集群名称
	ProdOrderStatus     int32  `json:"prodOrderStatus"`     // 订单状态，0：正常，1：冻结，2：删除，3：操作中，4：失败,2005:扩容中
	ProdRunningStatus   int32  `json:"prodRunninStatus"`    // 实例状态,0:运行中 1:重启中 2:备份中 3:恢复中 1001:已停止 1006:恢复失败 1007:VIP不可用 1008:GATEWAY不可用 1009:主库不可用 1010:备库不可用 1021:实例维护中 2000:开通中 2002:已退订 2005:扩容中 2011:冻结
	ProdType            int32  `json:"prodType"`            // 实例部署方式 0：单机部署,1：主备部署
	ReadPort            int32  `json:"readPort"`            // 读端口
	SecurityGroup       string `json:"securityGroup"`       // 安全组名称
	Subnet              string `json:"subnet"`              // 子网名称
	UserId              int64  `json:"userId"`              // 用户ID
	Vip                 string `json:"vip"`                 // 虚拟ip地址
	Vip6                string `json:"vip6"`                // vip6
	VpcId               string `json:"vpcId"`               // VPCID
	WritePort           string `json:"writePort"`           // 写端口
	ParameterGroupUsed  string `json:"parameterGroupUsed"`  // 参数配置所使用的名称，对应PaaS平台
	SpuCode             string `json:"spuCode"`             // 规格id
	ToolType            int32  `json:"toolType"`            // 备份工具类型，1：pg_baseback，2：pgbackrest，3：s3
	SubnetId            string `json:"subnetId"`            // 子网id
	SecurityGroupId     string `json:"securityGroupId"`     // 安全组id
	HostSeries          string `json:"hostSeries"`          // 性能类型
	BackupDiskType      string `json:"backupDiskType"`      // 备份空间类型
	BackupDiskSize      string `json:"backupDiskSize"`      // 备份空间大小,当为对象存储时，不返回值
	BackupUsageDiskSize string `json:"backupUsageDiskSize"` // 备份已使用空间
	BackupFreeSpace     string `json:"backupFreeSpace"`     // 备份免费空间
	BillMode            int32  `json:"billMode"`            // 计费模式,1包周期,2按需
}
