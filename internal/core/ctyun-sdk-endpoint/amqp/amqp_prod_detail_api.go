package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// AmqpProdDetailApi
/* 查询产品规格。
 */type AmqpProdDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpProdDetailApi(client *ctyunsdk.CtyunClient) *AmqpProdDetailApi {
	return &AmqpProdDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v3/instances/prodDetail",
		},
	}
}

func (this *AmqpProdDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpProdDetailRequest) (res *AmqpProdDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	builder.AddHeader("regionId", req.RegionId)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointName, builder)
	if err != nil {
		return
	}
	res = &AmqpProdDetailResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
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
	ProdId   string                                               `json:"prodId"`   /*  系列产品id  */
	ProdName string                                               `json:"prodName"` /*  产品系列名称  */
	ProdCode string                                               `json:"prodCode"` /*  产品系列编码  */
	ResItem  AmqpProdDetailReturnObjDataSeriesSkuResItemResponse  `json:"resItem"`  /*  主机信息  */
	DiskItem AmqpProdDetailReturnObjDataSeriesSkuDiskItemResponse `json:"diskItem"` /*  磁盘信息  */
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
	Spec     []struct {
		SpecName    string `json:"specName"`
		Description string `json:"description"`
		Cpu         int32  `json:"cpu"`
		Memory      int32  `json:"memory"`
	} `json:"spec"`
}
