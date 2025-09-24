package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUnTagResourceV3Api
/* 解绑资源v3
 */type CtgkafkaUnTagResourceV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUnTagResourceV3Api(client *core.CtyunClient) *CtgkafkaUnTagResourceV3Api {
	return &CtgkafkaUnTagResourceV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/resourceTag/unTagResource",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUnTagResourceV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUnTagResourceV3Request) (*CtgkafkaUnTagResourceV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUnTagResourceV3Request
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
	var resp CtgkafkaUnTagResourceV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUnTagResourceV3Request struct {
	RegionId       string   `json:"regionId,omitempty"`   /*  资源池编码  */
	TagId          string   `json:"tagId,omitempty"`      /*  标签ID  */
	ProdInstId     string   `json:"prodInstId,omitempty"` /*  实例ID  */
	RawType        string   `json:"type,omitempty"`       /*  资源类型，可选值有INSTANCE,GROUP,TOPIC  */
	ResourceIdList []string `json:"resourceIdList"`       /*  资源id。资源类型是INSTANCE，则资源id是实例id;资源类型是GROUP或TOPIC,则资源id是group名称或topic名称  */
}

type CtgkafkaUnTagResourceV3Response struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message,omitempty"`    /*  状态信息  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息  */
}
