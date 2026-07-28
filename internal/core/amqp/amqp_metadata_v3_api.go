package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpMetadataV3Api
/* 查询实例交换机队列以及虚拟机数量
 */type AmqpMetadataV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpMetadataV3Api(client *core.CtyunClient) *AmqpMetadataV3Api {
	return &AmqpMetadataV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/metadata",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpMetadataV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpMetadataV3Request) (*AmqpMetadataV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpMetadataV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpMetadataV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
}

type AmqpMetadataV3Response struct {
	ReturnObj  *AmqpMetadataV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                           `json:"message"`    /*  描述状态  */
	StatusCode string                           `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                           `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpMetadataV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
