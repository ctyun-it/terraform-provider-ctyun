package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaResourceTagAddtagV3Api
/* 创建标签v3
 */type CtgkafkaResourceTagAddtagV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaResourceTagAddtagV3Api(client *core.CtyunClient) *CtgkafkaResourceTagAddtagV3Api {
	return &CtgkafkaResourceTagAddtagV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/resourceTag/addTag",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaResourceTagAddtagV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaResourceTagAddtagV3Request) (*CtgkafkaResourceTagAddtagV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaResourceTagAddtagV3Request
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
	var resp CtgkafkaResourceTagAddtagV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaResourceTagAddtagV3Request struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池编码  */
	TagName  string `json:"tagName,omitempty"`  /*  标签名称  */
}

type CtgkafkaResourceTagAddtagV3Response struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message,omitempty"`    /*  状态信息  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息  */
}
