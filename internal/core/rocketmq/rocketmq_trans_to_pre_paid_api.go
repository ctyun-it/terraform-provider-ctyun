package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqTransToPrePaidApi
/* 包周期转按需。
 */type RocketmqTransToPrePaidApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqTransToPrePaidApi(client *core.CtyunClient) *RocketmqTransToPrePaidApi {
	return &RocketmqTransToPrePaidApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/transToPrePaid",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqTransToPrePaidApi) Do(ctx context.Context, credential core.Credential, req *RocketmqTransToPrePaidRequest) (*RocketmqTransToPrePaidResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqTransToPrePaidRequest
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
	var resp RocketmqTransToPrePaidResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqTransToPrePaidRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type RocketmqTransToPrePaidResponse struct {
	StatusCode int32                                    `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                                   `json:"message"`    /*  描述状态。  */
	ReturnObj  *RocketmqTransToPrePaidReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                   `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type RocketmqTransToPrePaidReturnObjResponse struct {
	Data *RocketmqTransToPrePaidReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type RocketmqTransToPrePaidReturnObjDataResponse struct {
	MasterOrderId   string `json:"masterOrderId"`
	MasterOrderNo   string `json:"masterOrderNo"`
	MasterOrderType string `json:"masterOrderType"`
}
