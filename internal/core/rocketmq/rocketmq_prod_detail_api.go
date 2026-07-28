package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqProdDetailApi
/* 查询产品规格。
 */type RocketmqProdDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqProdDetailApi(client *core.CtyunClient) *RocketmqProdDetailApi {
	return &RocketmqProdDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/queryProd",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *RocketmqProdDetailApi) Do(ctx context.Context, credential core.Credential, req *RocketmqProdDetailRequest) (*RocketmqProdDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp RocketmqProdDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqProdDetailRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
}

type RocketmqProdDetailResponse struct {
	StatusCode int32                                `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                               `json:"message"`    /*  描述状态  */
	ReturnObj  *RocketmqProdDetailReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例”里面的注释  */
	Error      string                               `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type RocketmqProdDetailReturnObjResponse struct {
	Data []*RocketmqFlavorDetail `json:"data"` /*  返回数据  */
}
type RocketmqFlavorDetail struct {
	FlavorID      string   `json:"flavorID"`      /*  规格 id  */
	SpecName      string   `json:"specName"`      /*  规格名称  */
	FlavorType    string   `json:"flavorType"`    /*  规格类型  */
	FlavorName    string   `json:"flavorName"`    /*  规格类型名称  */
	CpuNum        int32    `json:"cpuNum"`        /*  cpu 核数  */
	MemSize       int32    `json:"memSize"`       /*  内存大小，单位 G  */
	MultiQueue    int32    `json:"multiQueue"`    /*  多队列数  */
	Pps           int32    `json:"pps"`           /*  网络最大收发包能力 (万 PPS)  */
	BandwidthBase float32  `json:"bandwidthBase"` /*  基准带宽 (Gbps)  */
	BandwidthMax  float32  `json:"bandwidthMax"`  /*  最大带宽 (Gbps)  */
	CpuArch       string   `json:"cpuArch"`       /*  cpu 架构（x86、arm）  */
	Series        string   `json:"series"`        /*  系列  */
	AzList        []string `json:"azList"`        /*  支持的 az 名称列表  */
	SkuProdId     string   `json:"skuProdId"`     /*  产品 id  */
}
