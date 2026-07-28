package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpInstanceDeleteApi
/* 注销实例，实例将不可恢复，谨慎操作。
 */type AmqpInstanceDeleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpInstanceDeleteApi(client *core.CtyunClient) *AmqpInstanceDeleteApi {
	return &AmqpInstanceDeleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpInstanceDeleteApi) Do(ctx context.Context, credential core.Credential, req *AmqpInstanceDeleteRequest) (*AmqpInstanceDeleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpInstanceDeleteRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpInstanceDeleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
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
