package amqp

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// AmqpInstanceDeleteApi
/* 注销实例，实例将不可恢复，谨慎操作。
 */type AmqpInstanceDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewAmqpInstanceDeleteApi(client *ctyunsdk.CtyunClient) *AmqpInstanceDeleteApi {
	return &AmqpInstanceDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v3/instances/delete",
		},
	}
}

func (this *AmqpInstanceDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *AmqpInstanceDeleteRequest) (res *AmqpInstanceDeleteResponse, err error) {
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
	res = &AmqpInstanceDeleteResponse{}
	err = resp.Parse(res)
	if err != nil {
		return
	}
	return res, nil
}

type AmqpInstanceDeleteRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type AmqpInstanceDeleteResponse struct {
	StatusCode string                               `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                               `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpInstanceDeleteReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                               `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpInstanceDeleteReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据。  */
}
