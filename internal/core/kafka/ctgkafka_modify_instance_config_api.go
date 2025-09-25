package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaModifyInstanceConfigApi
/* 修改实例配置
 */type CtgkafkaModifyInstanceConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaModifyInstanceConfigApi(client *core.CtyunClient) *CtgkafkaModifyInstanceConfigApi {
	return &CtgkafkaModifyInstanceConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/modifyInstanceConfig",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaModifyInstanceConfigApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaModifyInstanceConfigRequest) (*CtgkafkaModifyInstanceConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaModifyInstanceConfigRequest
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
	var resp CtgkafkaModifyInstanceConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaModifyInstanceConfigRequest struct {
	RegionId        string `json:"regionId,omitempty"`        /*  资源池编码  */
	ProdInstId      string `json:"prodInstId,omitempty"`      /*  实例ID  */
	LogRetentionMs  string `json:"logRetentionMs,omitempty"`  /*  消息保留时长  */
	MessageMaxBytes string `json:"messageMaxBytes,omitempty"` /*  最大消息大小  */
}

type CtgkafkaModifyInstanceConfigResponse struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message,omitempty"`    /*  描述状态  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息  */
}
