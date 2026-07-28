package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbQueryInstNodesApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbQueryInstNodesApi(client *ctyunsdk.CtyunClient) *MongodbQueryInstNodesApi {
	return &MongodbQueryInstNodesApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v1/openApi/queryInstNodes",
		},
	}
}

func (this *MongodbQueryInstNodesApi) Do(ctx context.Context, credentials ctyunsdk.Credential, req *MongodbQueryInstNodesRequest, headers *MongodbQueryInstNodesRequestHeaders) (detailResp *MongodbQueryInstNodesResponse, err error) {
	builder := this.WithCredential(&credentials)

	// 验证必填参数
	if req == nil {
		err = errors.New("request is nil")
		return
	}
	if req.ProdInstId == "" {
		err = errors.New("ProdInstId is empty")
		return
	}

	if headers == nil {
		err = errors.New("headers is nil")
		return
	}
	if headers.RegionID == "" {
		err = errors.New("regionId is empty")
		return
	}

	// 添加查询参数
	builder.AddParam("prodInstId", req.ProdInstId)

	// 添加请求头
	builder.AddHeader("regionId", headers.RegionID)
	builder.AddHeader("prodInstId", req.ProdInstId)

	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}

	detailResp = &MongodbQueryInstNodesResponse{}
	err = resp.Parse(detailResp)
	if err != nil {
		return
	}

	return detailResp, nil
}

// 请求参数
type MongodbQueryInstNodesRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
}

// 请求头参数
type MongodbQueryInstNodesRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 区域ID，必填
}

// 响应结构
type MongodbQueryInstNodesResponse struct {
	StatusCode int32                        `json:"statusCode"`
	Message    *string                      `json:"message"`
	ReturnObj  *QueryInstNodesRespReturnObj `json:"returnObj"`
	Error      string                       `json:"error"`
}

// 节点信息
type QueryInstNodesRespReturnObjVOSInfo struct {
	ResId            int     `json:"resId"`                     // 机器id
	ProdInstId       string  `json:"prodInstId"`                // 实例ID
	Node             string  `json:"node"`                      // ID/名称
	Role             string  `json:"role"`                      // 角色
	Alive            int32   `json:"alive"`                     // 存活状态
	Host             string  `json:"host"`                      // 主机ip
	VpcIpv6          string  `json:"vpcIpv6"`                   // vpc内部的ipv6地址
	Port             string  `json:"port"`                      // 端口
	Type             *string `json:"type,omitempty"`            // 类型
	Memory           string  `json:"memory"`                    // 内存规格
	DiskSize         string  `json:"diskSize"`                  // 储存空间
	UsedSpace        string  `json:"usedSpace"`                 // 已用空间
	ProdInstSetName  *string `json:"prodInstSetName,omitempty"` // dds数据库实例对应的SET名
	AllowBeMaster    bool    `json:"allowBeMaster"`             // 允许切换成为备用节点
	AzDisplayName    string  `json:"azDisplayName"`             // 可用区
	RealIp           *string `json:"realIp,omitempty"`          // 实际IP
	ElasticIp        string  `json:"elasticIp"`                 // 弹性IP地址
	OuterElasticIpId string  `json:"outerElasticIpId"`          // 弹性IPid
	CpuCount         int32   `json:"cpuCount"`                  // cpu规格
}

// 返回对象
type QueryInstNodesRespReturnObj struct {
	List         []QueryInstNodesRespReturnObjVOSInfo `json:"list"`         // 节点信息列表
	ReadOnlyList []QueryInstNodesRespReturnObjVOSInfo `json:"readOnlyList"` // 只读节点信息列表
}
