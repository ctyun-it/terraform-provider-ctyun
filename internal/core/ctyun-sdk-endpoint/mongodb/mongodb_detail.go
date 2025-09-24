package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbQueryDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbQueryDetailApi(client *ctyunsdk.CtyunClient) *MongodbQueryDetailApi {
	return &MongodbQueryDetailApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v1/openApi/describeDBInstanceAttribute",
		},
	}
}

func (this *MongodbQueryDetailApi) Do(ctx context.Context, credentials ctyunsdk.Credential, req *MongodbQueryDetailRequest, headers *MongodbQueryDetailRequestHeaders) (detailResp *MongodbQueryDetailResponse, err error) {
	builder := this.WithCredential(&credentials)
	_, err = builder.WriteJson(req)
	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	if headers.RegionID == "" {
		err = errors.New("regionId is empty")
		return
	}
	builder.AddHeader("regionId", headers.RegionID)
	if req.ProdInstId == "" {
		err = errors.New("ProdInstId is empty")
		return
	}
	builder.AddParam("prodInstId", req.ProdInstId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	detailResp = &MongodbQueryDetailResponse{}
	err = resp.Parse(detailResp)
	if err != nil {
		return
	}
	return detailResp, nil
}

type MongodbQueryDetailRequest struct {
	ProdInstId string `json:"prodInstId"` //实例ID，必填
}
type MongodbQueryDetailRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"`
	RegionID  string  `json:"regionId"`
}

type MongodbQueryDetailResponse struct {
	StatusCode int32                `json:"statusCode"`
	Message    *string              `json:"message"`
	ReturnObj  *DetailRespReturnObj `json:"returnObj"`
	Error      string               `json:"error"`
}

type DetailRespReturnObjBackupInfo struct {
	UsePercent string `json:"usepercent"` // 使用百分比
	Size       string `json:"size"`       // 存储总大小
	Used       string `json:"used"`       // 已使用大小
}

type DetailRespReturnObjVOSInfo struct {
	RealIp           *string `json:"realIp,omitempty"`          // 实际IP
	AllowBeMaster    bool    `json:"allowBeMaster"`             // 是否允许成为主节点
	OuterElasticIpId string  `json:"outerElasticIpId"`          // 弹性IP ID
	Memory           string  `json:"memory"`                    // 内存
	Role             string  `json:"role"`                      // 节点角色
	Alive            int32   `json:"alive"`                     // 存活状态
	UsedSpace        *string `json:"usedSpace,omitempty"`       // 已使用空间
	VpcIpv6          string  `json:"vpcIpv6"`                   // VPC内IPv6地址
	Type             *string `json:"type,omitempty"`            // 节点类型
	ResId            int     `json:"resId"`                     // 资源ID
	ElasticIp        string  `json:"elasticIp"`                 // 弹性IP
	Node             string  `json:"node"`                      // 节点名称
	DiskSize         *string `json:"diskSize,omitempty"`        // 磁盘大小
	AzDisplayName    string  `json:"azDisplayName"`             // 可用区显示名
	Host             string  `json:"host"`                      // 主机地址
	ProdInstSetName  *string `json:"prodInstSetName,omitempty"` // 产品实例集名
	ProdInstId       string  `json:"prodInstId"`                // 产品实例ID
	CpuCount         int32   `json:"cpuCount"`                  // CPU数量
	AzId             *string `json:"azId,omitempty"`            // az id
}

type DetailRespReturnObj struct {
	Backup            *DetailRespReturnObjBackupInfo `json:"backup"`
	ProdInstName      string                         `json:"prodInstName"`      // 产品实例名称
	ProdType          int32                          `json:"prodType"`          // 产品类型
	EnableSSL         int32                          `json:"enableSSL"`         // 是否启用SSL
	DiskSize          int32                          `json:"diskSize"`          // 磁盘大小
	CreateTime        int64                          `json:"createTime"`        // 创建时间 (毫秒)
	Port              string                         `json:"port"`              // 端口号
	ProdRunningStatus int32                          `json:"prodRunningStatus"` // 产品运行状态
	DiskRate          float64                        `json:"diskRate"`          // 磁盘占用率
	NodeInfoVOS       []DetailRespReturnObjVOSInfo   `json:"nodeInfoVOS"`
	Host              string                         `json:"host"`        // 主机地址
	TenantId          int64                          `json:"tenantId"`    // 租户ID
	ProdInstId        string                         `json:"prodInstId"`  // 产品实例ID
	DiskType          string                         `json:"diskType"`    // 磁盘类型
	MachineSpec       string                         `json:"machineSpec"` // 机器规格
}

type AzInfo struct {
	ResId  int64  `json:"resId"`
	Role   string `json:"role"`
	AzName string `json:"azName"`
	AzId   string `json:"azId"`
}
