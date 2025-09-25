package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaListResourceV3Api
/* 获取标签绑定的资源列表v3
 */type CtgkafkaListResourceV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaListResourceV3Api(client *core.CtyunClient) *CtgkafkaListResourceV3Api {
	return &CtgkafkaListResourceV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/resourceTag/listResource",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaListResourceV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaListResourceV3Request) (*CtgkafkaListResourceV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("tagId", req.TagId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaListResourceV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaListResourceV3Request struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池编码  */
	TagId    string `json:"tagId,omitempty"`    /*  标签ID  */
}

type CtgkafkaListResourceV3Response struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message,omitempty"`    /*  状态信息  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息  */
}
