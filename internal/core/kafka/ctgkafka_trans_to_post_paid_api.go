package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaTransToPostPaidApi
/* 按需实例转为包周期实例，实例到期后自动转为按需付费。 */
type CtgkafkaTransToPostPaidApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTransToPostPaidApi(client *core.CtyunClient) *CtgkafkaTransToPostPaidApi {
	return &CtgkafkaTransToPostPaidApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/transToPostPaid",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaTransToPostPaidApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTransToPostPaidRequest) (*CtgkafkaTransToPostPaidResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaTransToPostPaidRequest
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
	var resp CtgkafkaTransToPostPaidResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtgkafkaTransToPostPaidRequest struct {
	RegionId   string `json:"regionId"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId"` /*  实例ID。  */
}

type CtgkafkaTransToPostPaidResponse struct {
	StatusCode string                                    `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                                   `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTransToPostPaidReturnObjResponse `json:"returnObj,omitempty"`  /*  返回对象。  */
	Error      *string                                   `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTransToPostPaidReturnObjResponse struct {
	Data *CtgkafkaTransToPostPaidReturnObjDataResponse `json:"data,omitempty"` /*  返回数据。  */
}

type CtgkafkaTransToPostPaidReturnObjDataResponse struct {
	MasterOrderId   *string `json:"masterOrderId,omitempty"`   /*  订单id。  */
	MasterOrderNo   *string `json:"masterOrderNo,omitempty"`   /*  订单编号。  */
	MasterOrderType *string `json:"masterOrderType,omitempty"` /*  订单类型。  */
	Source          *string `json:"source,omitempty"`          /*  资源编码。  */
	StatusDate      *string `json:"statusDate,omitempty"`      /*  更新时间。  */
	TotalPrice      float64 `json:"totalPrice"`                /*  总价，单位元。  */
	FinalPrice      float64 `json:"finalPrice"`                /*  订单最终价格，单位元。  */
}
