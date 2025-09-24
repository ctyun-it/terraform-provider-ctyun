package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaProdDetailApi
/* 查询产品规格。
 */type CtgkafkaProdDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaProdDetailApi(client *core.CtyunClient) *CtgkafkaProdDetailApi {
	return &CtgkafkaProdDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/prodDetail",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaProdDetailApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaProdDetailRequest) (*CtgkafkaProdDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaProdDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaProdDetailRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
}

type CtgkafkaProdDetailResponse struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"  */
	Message    string                               `json:"message,omitempty"`    /*  描述状态  */
	ReturnObj  *CtgkafkaProdDetailReturnObjResponse `json:"returnObj"`            /*  返回对象。此参数所包含的参数请见“响应示例”里面的注释  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息  */
	RequestId  string                               `json:"requestId,omitempty"`
}

type CtgkafkaProdDetailReturnObjResponse struct {
	Data struct {
		Series []CtgkafkaProdDetailReturnObjResponseSeries `json:"series"`
		OsList []string                                    `json:"osList"`
	} `json:"data"`
}

type CtgkafkaProdDetailReturnObjResponseSeries struct {
	ProdId   string                                   `json:"prodId"`
	ProdName string                                   `json:"prodName"`
	ProdCode string                                   `json:"prodCode"`
	Sku      []CtgkafkaProdDetailReturnObjResponseSku `json:"sku"`
}

type CtgkafkaProdDetailReturnObjResponseSku struct {
	ProdId   string                                         `json:"prodId"`
	ProdName string                                         `json:"prodName"`
	ProdCode string                                         `json:"prodCode"`
	ResItem  CtgkafkaProdDetailReturnObjResponseSkuResItem  `json:"resItem"`
	DiskItem CtgkafkaProdDetailReturnObjResponseSkuDiskItem `json:"diskItem"`
}

type CtgkafkaProdDetailReturnObjResponseSkuResItem struct {
	ResType  string `json:"resType"`
	ResName  string `json:"resName"`
	ResItems []struct {
		CpuArch  string `json:"cpuArch"`
		HostType string `json:"hostType"`
		Spec     []struct {
			SpecName     string `json:"specName"`
			Description  string `json:"description"`
			Tps          int32  `json:"tps"`
			MaxPartition int32  `json:"maxPartition"`
			Flow         int32  `json:"flow"`
			Cpu          int32  `json:"cpu"`
			Memory       int32  `json:"memory"`
		} `json:"spec"`
	} `json:"resItems"`
}

type CtgkafkaProdDetailReturnObjResponseSkuDiskItem struct {
	ResType  string   `json:"resType"`
	ResName  string   `json:"resName"`
	ResItems []string `json:"resItems"`
}
