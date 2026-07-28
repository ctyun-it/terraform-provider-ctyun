package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqSourceBindingV3Api
/* 查交换机绑定v3
 */type RocketmqSourceBindingV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqSourceBindingV3Api(client *core.CtyunClient) *RocketmqSourceBindingV3Api {
	return &RocketmqSourceBindingV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/binding/exchange/source",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *RocketmqSourceBindingV3Api) Do(ctx context.Context, credential core.Credential, req *RocketmqSourceBindingV3Request) (*RocketmqSourceBindingV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	ctReq.AddParam("exchange", req.Exchange)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp RocketmqSourceBindingV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqSourceBindingV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
	Vhost      string `json:"vhost,omitempty"`      /*  虚拟主机  */
	Exchange   string `json:"exchange,omitempty"`   /*  交换器名称  */
}

type RocketmqSourceBindingV3Response struct {
	ReturnObj  *RocketmqSourceBindingV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                                    `json:"message"`    /*  描述状态  */
	StatusCode string                                    `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                                    `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type RocketmqSourceBindingV3ReturnObjResponse struct {
	Data []*RocketmqSourceBindingV3ReturnObjDataResponse `json:"data"` /*  绑定数据  */
}

type RocketmqSourceBindingV3ReturnObjDataResponse struct {
	Source           string `json:"source"`           /*  绑定来源  */
	Vhost            string `json:"vhost"`            /*  虚拟主机  */
	Destination      string `json:"destination"`      /*  绑定目标  */
	Destination_type string `json:"destination_type"` /*  绑定目标类型，只有一种类型queue  */
	Routing_key      string `json:"routing_key"`      /*  绑定键  */
	Arguments        string `json:"arguments"`        /*  绑定参数  */
	Properties_key   string `json:"properties_key"`   /*  properties键  */
}
