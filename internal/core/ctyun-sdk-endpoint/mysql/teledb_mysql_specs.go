package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbMysqlSpecsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbMysqlSpecsApi(client *ctyunsdk.CtyunClient) *TeledbMysqlSpecsApi {
	return &TeledbMysqlSpecsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/template",
		},
	}
}

func (this *TeledbMysqlSpecsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbMysqlSpecsRequest, header *TeledbMysqlSpecsRequestHeader) (specsResponse *TeledbMysqlSpecsResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.ProdType == "" {
		err = errors.New("prodType 为空")
		return
	}
	builder.AddParam("prodType", req.ProdType)
	if req.ProdCode == "" {
		err = errors.New("prodCode 为空")
	}
	builder.AddParam("prodCode", req.ProdCode)
	if req.RegionID == "" {
		err = errors.New("regionId 为空")
	}
	builder.AddParam("regionId", req.RegionID)
	if req.InstanceType == "" {
		err = errors.New("instanceType 为空")
	}
	builder.AddParam("instanceType", req.InstanceType)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	specsResponse = &TeledbMysqlSpecsResponse{}
	err = resp.Parse(specsResponse)
	if err != nil {
		return
	}
	return specsResponse, nil
}

type TeledbMysqlSpecsRequest struct {
	ProdType     string `json:"prod_type"`
	ProdCode     string `json:"prod_code"`
	RegionID     string `json:"region_id"`
	InstanceType string `json:"instance_type"`
}

type TeledbMysqlSpecsRequestHeader struct {
	ProjectID *string `json:"project_id"`
}

type InstSpecInfo struct {
	SpecId              string   `json:"specId"`              // 废弃
	ProdPerformanceSpec string   `json:"prodPerformanceSpec"` // 规格名称
	AzList              []string `json:"azList"`              // 该规格支持的AZ列表
	SpecName            string   `json:"specName"`            // 主机世代完整名称
	CpuType             string   `json:"cpuType"`             // CPU类型
	Generation          string   `json:"generation"`          // 主机世代缩写
	MinRate             string   `json:"minRate"`             // 带宽下限
	MaxRate             string   `json:"maxRate"`             // 带宽上限
}

type HostInst struct {
	HostTypeName          string   `json:"hostTypeName"`
	HostType              string   `json:"hostType"`
	ProdPerformanceSpeces []string `json:"prodPerformanceSpeces"`
	HostDefaultNum        int32    `json:"hostDefaultNum"`
}

type ProdHostConfig struct {
	HostInsts []HostInst `json:"hostInsts"`
}

type TeledbMysqlSpecsResponseReturnObjData struct {
	ProdId           int64          `json:"prodId"`
	ProdCode         string         `json:"prodCode"`
	ProdSpecName     string         `json:"prodSpecName"`
	ProdSpecDesc     string         `json:"prodSpecDesc"`
	InstanceDesc     string         `json:"instanceDesc"`
	ProdVersion      string         `json:"prodVersion"`
	HostSpec         string         `json:"hostSpec"`
	LvsSpec          string         `json:"lvsSpec"`
	InstSpecInfoList []InstSpecInfo `json:"instSpecInfoList"`
	ProdHostConfig   ProdHostConfig `json:"prodHostConfig"`
}

type TeledbMysqlSpecsResponseReturnObj struct {
	Data []TeledbMysqlSpecsResponseReturnObjData `json:"data"`
}

type TeledbMysqlSpecsResponse struct {
	StatusCode int32                              `json:"statusCode"` // 接口状态码
	Error      string                             `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                             `json:"message"`    // 描述信息
	ReturnObj  *TeledbMysqlSpecsResponseReturnObj `json:"returnObj"`
}
