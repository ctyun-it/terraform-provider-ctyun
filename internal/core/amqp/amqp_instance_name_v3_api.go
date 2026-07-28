package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpInstanceNameV3Api
/* 更改实例名称V3
 */type AmqpInstanceNameV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpInstanceNameV3Api(client *core.CtyunClient) *AmqpInstanceNameV3Api {
	return &AmqpInstanceNameV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/instanceName",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpInstanceNameV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpInstanceNameV3Request) (*AmqpInstanceNameV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpInstanceNameV3Request
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
	var resp AmqpInstanceNameV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpInstanceNameV3Request struct {
	RegionId     string `json:"regionId,omitempty"`     /*  资源池id  */
	InstanceName string `json:"instanceName,omitempty"` /*  新的实例名称  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例id  */
}

type AmqpInstanceNameV3Response struct {
	ReturnObj  string `json:"returnObj"`  /*  返回对象  */
	Message    string `json:"message"`    /*  描述状态  */
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}
