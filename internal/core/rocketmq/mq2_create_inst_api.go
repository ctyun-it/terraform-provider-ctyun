package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2CreateInstApi
/* 开通消息队列RocketMQ版服务 */
type Mq2CreateInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2CreateInstApi(client *core.CtyunClient) *Mq2CreateInstApi {
	return &Mq2CreateInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/spuInst/createInst",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2CreateInstApi) Do(ctx context.Context, credential core.Credential, req *Mq2CreateInstRequest) (*Mq2CreateInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2CreateInstRequest
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
	var resp Mq2CreateInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2CreateInstRequest struct {
	RegionId string                     `json:"regionId"`        /*  资源池编码  */
	Order    *Mq2CreateInstOrderRequest `json:"order,omitempty"` /*  订单对象  */
}

type Mq2CreateInstOrderRequest struct {
	InstanceCnt string                         `json:"instanceCnt"`    /*  订购数量  */
	CycleCnt    string                         `json:"cycleCnt"`       /*  订购周期，取值范围：值需大于零，不超过384个月  */
	CycleType   string                         `json:"cycleType"`      /*  订购周期类型，取值范围：3表示按月订购，5表示按一年订购，6表示按两年订购、7表示按三年订购，cycleCnt属性为1，cycleType取值为3 表示订购1个月  */
	Item        *Mq2CreateInstOrderItemRequest `json:"item,omitempty"` /*  订单项对象  */
}

type Mq2CreateInstOrderItemRequest struct {
	ItemValue  string `json:"itemValue"`  /*  订购数量  */
	ItemConfig string `json:"itemConfig"` /*  订单项配置对象  */
}

type Mq2CreateInstResponse struct {
	StatusCode *int32 `json:"statusCode"` /*  返回码
	取值范围：800 成功 , 900 失败  */
	ReturnObj *Mq2CreateInstReturnObjResponse `json:"returnObj"` /*  返回对象  */
}

type Mq2CreateInstReturnObjResponse struct {
	Submitted  *bool   `json:"submitted"`  /*  1  */
	NewOrderId *string `json:"newOrderId"` /*  订单Id  */
	NewOrderNo *string `json:"newOrderNo"` /*  订单号  */
	TotalPrice *string `json:"totalPrice"` /*  总价  */
}
