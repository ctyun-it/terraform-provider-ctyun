package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqInstanceNameV3Api
/* 更改实例名称 V3
 */type RocketmqInstanceNameV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqInstanceNameV3Api(client *core.CtyunClient) *RocketmqInstanceNameV3Api {
	return &RocketmqInstanceNameV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/updateName",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqInstanceNameV3Api) Do(ctx context.Context, credential core.Credential, req *RocketmqInstanceNameV3Request) (*RocketmqInstanceNameV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqInstanceNameV3Request
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
	var resp RocketmqInstanceNameV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqInstanceNameV3Request struct {
	RegionId     string `json:"regionId,omitempty"`     /*  资源池 ID  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例 ID  */
	InstanceName string `json:"instanceName,omitempty"` /*  新的实例名称  */
}

type RocketmqInstanceNameV3Response struct {
	StatusCode int32                                   `json:"statusCode"` // 接口系统层面状态码。成功：800，失败：900
	Message    string                                  `json:"message"`    // 描述状态
	ReturnObj  *RocketmqInstancesInstanceNameReturnObj `json:"returnObj"`  // 返回对象
	Error      string                                  `json:"error"`      // 错误码，只有非成功才有这个字段，方便快速定位问题
}

type RocketmqInstancesInstanceNameReturnObj struct {
	// 空对象，成功时返回空对象 {}
}
