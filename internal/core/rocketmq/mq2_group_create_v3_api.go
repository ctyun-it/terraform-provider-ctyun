package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2GroupCreateV3Api
/* 创建消费组 */
type Mq2GroupCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2GroupCreateV3Api(client *core.CtyunClient) *Mq2GroupCreateV3Api {
	return &Mq2GroupCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/group/create",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2GroupCreateV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2GroupCreateV3Request) (*Mq2GroupCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2GroupCreateV3Request
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
	var resp Mq2GroupCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2GroupCreateV3Request struct {
	RegionId                string                                          `json:"regionId"`                          /*  资源池编码  */
	ProdInstId              string                                          `json:"prodInstId"`                        /*  实例id  */
	Remark                  string                                          `json:"remark"`                            /*  实例备注  */
	BrokerNameList          []string                                        `json:"brokerNameList"`                    /*  Broker名称列表  */
	SubscriptionGroupConfig *Mq2GroupCreateV3SubscriptionGroupConfigRequest `json:"subscriptionGroupConfig,omitempty"` /*  订阅组配置信息  */
}

type Mq2GroupCreateV3SubscriptionGroupConfigRequest struct {
	ConsumeEnable          bool   `json:"consumeEnable"`          /*  是否允许消费  */
	ConsumeBroadcastEnable bool   `json:"consumeBroadcastEnable"` /*  是否开启广播消费  */
	FirstConsumeMechanism  int32  `json:"firstConsumeMechanism"`  /*  首次消费机制  */
	GroupName              string `json:"groupName"`              /*  订阅组名称  */
	PullMechanism          int32  `json:"pullMechanism"`          /*  拉取机制  */
	RetryMaxTimes          int32  `json:"retryMaxTimes"`          /*  最大重试次数  */
}

type Mq2GroupCreateV3Response struct {
	StatusCode int32                              `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                            `json:"message"`    /*  描述状态。  */
	ReturnObj  *Mq2GroupCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      *string                            `json:"error"`      /*  错误码，描述错误信息。  */
}

type Mq2GroupCreateV3ReturnObjResponse struct{}
