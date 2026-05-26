package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2InstanceQueryProdApi
/* 查询产品规格 */
type Mq2InstanceQueryProdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2InstanceQueryProdApi(client *core.CtyunClient) *Mq2InstanceQueryProdApi {
	return &Mq2InstanceQueryProdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/queryProd",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2InstanceQueryProdApi) Do(ctx context.Context, credential core.Credential, req *Mq2InstanceQueryProdRequest) (*Mq2InstanceQueryProdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2InstanceQueryProdResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2InstanceQueryProdRequest struct {
	RegionId string `json:"regionId"` /*  资源池编码  */
}

type Mq2InstanceQueryProdResponse struct {
	StatusCode *string                                `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2InstanceQueryProdReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2InstanceQueryProdReturnObjResponse struct {
	Data []*Mq2InstanceQueryProdReturnObjDataResponse `json:"data"` /*  规格详情列表  */
}

type Mq2InstanceQueryProdReturnObjDataResponse struct {
	FlavorID      *string   `json:"flavorID"`      /*  规格id  */
	SpecName      *string   `json:"specName"`      /*  规格名称  */
	FlavorType    *string   `json:"flavorType"`    /*  规格类型  */
	FlavorName    *string   `json:"flavorName"`    /*  规格类型名称  */
	CpuNum        *int32    `json:"cpuNum"`        /*  cpu核数  */
	MemSize       *int32    `json:"memSize"`       /*  内存大小，单位G  */
	MultiQueue    *int32    `json:"multiQueue"`    /*  多队列数  */
	Pps           *int32    `json:"pps"`           /*  网络最大收发包能力 (万PPS)  */
	BandwidthBase *float32  `json:"bandwidthBase"` /*  基准带宽 (Gbps)  */
	BandwidthMax  *int32    `json:"bandwidthMax"`  /*  最大带宽 (Gbps)  */
	CpuArch       *string   `json:"cpuArch"`       /*  cpu架构（x86、arm）  */
	Series        *string   `json:"series"`        /*  系列  */
	AzList        []*string `json:"azList"`        /*  支持的az名称列表  */
	SkuProdId     *string   `json:"skuProdId"`     /*  产品id  */
}
