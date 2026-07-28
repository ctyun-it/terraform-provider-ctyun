package opensearch

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OpensearchUnsubscribeInstanceApi
/* 退订 OpenSearch 集群实例<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
<b>注意事项：</b><br />
&emsp;&emsp;退订操作将立即生效，请谨慎操作<br />
&emsp;&emsp;退订后实例将被删除，数据无法恢复<br />
*/type OpensearchUnsubscribeInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchUnsubscribeInstanceApi(client *core.CtyunClient) *OpensearchUnsubscribeInstanceApi {
	return &OpensearchUnsubscribeInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/os/openapi/v1/order/refund",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchUnsubscribeInstanceApi) Do(ctx context.Context, credential core.Credential, req *UnsubscribeInstanceRequest) (*UnsubscribeInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	// GET 请求，参数通过 Query 传递
	ctReq.AddParam("clusterId", req.ClusterID)

	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp UnsubscribeInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type UnsubscribeInstanceRequest struct {
	ClusterID string `json:"clusterId"` /*  实例 id  */
}

type UnsubscribeInstanceResponse struct {
	StatusCode int32                  `json:"statusCode"` /*  状态码，成功："200"，失败："500"  */
	Error      string                 `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string                 `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  []UnsubscribeReturnObj `json:"returnObj"`  /*  返回结果数组  */
}

type UnsubscribeReturnObj struct {
	ErrorMessage string  `json:"errorMessage"` /*  错误信息  */
	Submitted    bool    `json:"submitted"`    /*  提交订单  */
	NewOrderId   string  `json:"newOrderId"`   /*  订单 id  */
	NewOrderNo   string  `json:"newOrderNo"`   /*  订单编号  */
	TotalPrice   float64 `json:"totalPrice"`   /*  价格  */
}
