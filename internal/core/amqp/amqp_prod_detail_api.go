package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpProdDetailApi
/* 查询产品规格。
 */type AmqpProdDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpProdDetailApi(client *core.CtyunClient) *AmqpProdDetailApi {
	return &AmqpProdDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/prodDetail",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpProdDetailApi) Do(ctx context.Context, credential core.Credential, req *AmqpProdDetailRequest) (*AmqpProdDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpProdDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpProdDetailRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
}

type AmqpProdDetailResponse struct {
	StatusCode string                           `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                           `json:"message"`    /*  描述状态  */
	ReturnObj  *AmqpProdDetailReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例”里面的注释  */
	Error      string                           `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpProdDetailReturnObjResponse struct {
	Data *AmqpProdDetailReturnObjDataResponse `json:"data"` /*  返回数据  */
}

type AmqpProdDetailReturnObjDataResponse struct {
	Series []*AmqpProdDetailReturnObjDataSeriesResponse `json:"series"` /*  产品系列信息  */
}

type AmqpProdDetailReturnObjDataSeriesResponse struct {
	ProdId   string                                          `json:"prodId"`   /*  系列产品id  */
	ProdName string                                          `json:"prodName"` /*  产品系列名称  */
	ProdCode string                                          `json:"prodCode"` /*  产品系列编码  */
	Sku      []*AmqpProdDetailReturnObjDataSeriesSkuResponse `json:"sku"`      /*  产品系列信息  */
}

type AmqpProdDetailReturnObjDataSeriesSkuResponse struct {
	ProdId   string                                                `json:"prodId"`   /*  系列产品id  */
	ProdName string                                                `json:"prodName"` /*  产品系列名称  */
	ProdCode string                                                `json:"prodCode"` /*  产品系列编码  */
	ResItem  *AmqpProdDetailReturnObjDataSeriesSkuResItemResponse  `json:"resItem"`  /*  主机信息  */
	DiskItem *AmqpProdDetailReturnObjDataSeriesSkuDiskItemResponse `json:"diskItem"` /*  磁盘信息  */
}

type AmqpProdDetailReturnObjDataSeriesSkuResItemResponse struct {
	ResType  string                                                         `json:"resType"`  /*  ecs  */
	ResName  string                                                         `json:"resName"`  /*  云服务器  */
	ResItems []*AmqpProdDetailReturnObjDataSeriesSkuResItemResItemsResponse `json:"resItems"` /*  主机规格信息  */
}

type AmqpProdDetailReturnObjDataSeriesSkuDiskItemResponse struct {
	ResType  string   `json:"resType"`  /*  资源类型  */
	ResName  string   `json:"resName"`  /*  资源名称  */
	ResItems []string `json:"resItems"` /*  磁盘类型  */
}

type AmqpProdDetailReturnObjDataSeriesSkuResItemResItemsResponse struct {
	CpuArch  string `json:"cpuArch"`  /*  cpu架构  */
	HostType string `json:"hostType"` /*  主机类型  */
	HostTag  string `json:"hostTag"`  /*  主机Tag  */
}
