package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqSpecExtendApi
/* 规格扩容。
 */type RocketmqSpecExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqSpecExtendApi(client *core.CtyunClient) *RocketmqSpecExtendApi {
	return &RocketmqSpecExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/specExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqSpecExtendApi) Do(ctx context.Context, credential core.Credential, req *RocketmqSpecExtendRequest) (*RocketmqSpecExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqSpecExtendRequest
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
	var resp RocketmqSpecExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqSpecExtendRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池 ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例 ID。  */
	CpuNum     int32  `json:"cpuNum,omitempty"`     /*  扩容后的 cpu 核数  */
	MemSize    int32  `json:"memSize,omitempty"`    /*  扩容后的内存大小，单位 G  */
	AutoPay    bool   `json:"autoPay"`              /*  是否自动支付，当实例为按需计费模式不生效。true：自动付费，默认值。false：手动付费。  */
}

type RocketmqSpecExtendResponse struct {
	StatusCode int32                        `json:"statusCode"` // 接口系统层面状态码。成功：800，失败：900
	Message    string                       `json:"message"`    // 描述状态
	ReturnObj  *RocketmqSpecExtendReturnObj `json:"returnObj"`  // 返回对象
	Error      string                       `json:"error"`      // 错误码，只有非成功才有这个字段，方便快速定位问题
}

type RocketmqSpecExtendReturnObj struct {
	Data *OrderData `json:"data"` // 返回的订单相关数据
}
