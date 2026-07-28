package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpCanExpandProdApi
/* 查询产品可扩容规格。
 */type AmqpCanExpandProdApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpCanExpandProdApi(client *core.CtyunClient) *AmqpCanExpandProdApi {
	return &AmqpCanExpandProdApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/canExpandProd",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpCanExpandProdApi) Do(ctx context.Context, credential core.Credential, req *AmqpCanExpandProdRequest) (*AmqpCanExpandProdResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpCanExpandProdResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpCanExpandProdRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID，如果填入，则返回指定实例信息  */
}

type AmqpCanExpandProdResponse struct {
	StatusCode string                              `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                              `json:"message"`    /*  描述状态  */
	ReturnObj  *AmqpCanExpandProdReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例”里面的注释  */
	Error      string                              `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpCanExpandProdReturnObjResponse struct{}
