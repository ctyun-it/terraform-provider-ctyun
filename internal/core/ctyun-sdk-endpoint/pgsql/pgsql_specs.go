package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlSpecsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlSpecsApi(client *ctyunsdk.CtyunClient) *PgsqlSpecsApi {
	return &PgsqlSpecsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/template",
		},
	}
}

func (this *PgsqlSpecsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlSpecsRequest, header *PgsqlSpecsRequestHeader) (specsResponse *PgsqlSpecsResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.ProdType == "" {
		err = errors.New("missing required field: ProdType")
		return
	}
	if req.ProdCode == "" {
		err = errors.New("missing required field: ProdCode")
		return
	}
	if req.RegionId == "" {
		err = errors.New("missing required field: RegionId")
		return
	}
	if req.InstanceType == "" {
		err = errors.New("missing required field: InstanceType")
		return
	}
	builder.AddParam("prodType", req.ProdType)
	builder.AddParam("prodCode", req.ProdCode)
	builder.AddParam("regionId", req.RegionId)
	builder.AddParam("instanceType", req.InstanceType)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	specsResponse = &PgsqlSpecsResponse{}
	err = resp.Parse(specsResponse)
	if err != nil {
		return
	}
	return specsResponse, nil
}

type PgsqlSpecsRequest struct {
	ProdType     string `json:"prodType"`     // 产品类型，0=UNKNOWN, 1=RDS, 2=NoSql, 3=TOOL, 4=MemDB
	ProdCode     string `json:"prodCode"`     // 产品编码 :POSTGRESQL
	RegionId     string `json:"regionId"`     // 区域id
	InstanceType string `json:"instanceType"` // 实例类型，1=通用型，2=计算增强型，3=内存优化型，4=直通（未用到）
}

type PgsqlSpecsRequestHeader struct {
	ProjectID string `json:"projectId"`
}

type PgsqlSpecsResponse struct {
	StatusCode int32               `json:"statusCode"`
	Message    string              `json:"message"`
	Error      string              `json:"error"`
	ReturnObj  PgsqlSpecsReturnObj `json:"returnObj"`
}

type PgsqlSpecsInstSpec struct {
	SpecId              string   `json:"specId"`              // 废弃
	ProdPerformanceSpec string   `json:"prodPerformanceSpec"` // 规格名称
	AzList              []string `json:"azList"`              // 该规格支持的AZ列表
	SpecName            string   `json:"specName"`            // 主机世代完整名称
	CpuType             string   `json:"cpuType"`             // CPU类型
	Generation          string   `json:"generation"`          // 主机世代缩写
	MinRate             string   `json:"minRate"`             // 带宽下限
	MaxRate             string   `json:"maxRate"`             // 带宽上限
}

type PgsqlSpecsHostInst struct {
	HostTypeName          string   `json:"hostTypeName"`          // 类型名称
	HostType              string   `json:"hostType"`              // 节点类型, 取值范围 master（主节点）、readnode（只读节点）
	ProdPerformanceSpeces []string `json:"prodPerformanceSpeces"` // 支持的性能指标规格列表
	HostDefaultNum        int      `json:"hostDefaultNum"`        // 节点默认数量
}

type PgsqlSpecsHostConfig struct {
	HostInsts PgsqlSpecsHostInst `json:"hostInsts"`
}

type PgsqlSpecsReturnObj struct {
	ProdId           int64                `json:"prodId"`           // 产品id
	ProdCode         string               `json:"prodCode"`         // 产品编码
	ProdSpecName     string               `json:"prodSpecName"`     // 产品名称
	ProdSpecDesc     string               `json:"prodSpecDesc"`     // 产品描述
	InstanceDesc     string               `json:"instanceDesc"`     // 实例描述
	ProdVersion      string               `json:"prodVersion"`      // 产品版本
	HostSpec         string               `json:"hostSpec"`         // 主机规格
	LvsSpec          string               `json:"lvsSpec"`          // lvs规格
	InstSpecInfoList []PgsqlSpecsInstSpec `json:"instSpecInfoList"` // AZ支持的产品规格信息，以及规格代S6/S7
	ProdHostConfig   PgsqlSpecsHostConfig `json:"prodHostConfig"`   // 主机配置
}
