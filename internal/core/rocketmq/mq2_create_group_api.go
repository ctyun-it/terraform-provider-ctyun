package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2CreateGroupApi
/* 创建订阅组 */
type Mq2CreateGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2CreateGroupApi(client *core.CtyunClient) *Mq2CreateGroupApi {
	return &Mq2CreateGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/group/create",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2CreateGroupApi) Do(ctx context.Context, credential core.Credential, req *Mq2CreateGroupRequest) (*Mq2CreateGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2CreateGroupRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2CreateGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2CreateGroupRequest struct {
	ProdInstId              string                                        `json:"prodInstId"`                        /*  实例ID  */
	BrokerNameList          []string                                      `json:"brokerNameList"`                    /*  Broker名字列表  */
	SubscriptionGroupConfig *Mq2CreateGroupSubscriptionGroupConfigRequest `json:"subscriptionGroupConfig,omitempty"` /*  订阅组配置  */
}

type Mq2CreateGroupSubscriptionGroupConfigRequest struct {
	GroupName             string `json:"groupName"`              /*  订阅组名字  */
	ConsumeEnable         bool   `json:"consumeEnable"`          /*  是否允许消费。默认是true  */
	FirstConsumeMechanism string `json:"firstConsumeMechanism"`  /*  首次消费位置 1：客户端指定 2：第一条 3：最新位置  */
	PullMechanism         string `json:"pullMechanism"`          /*  消费机制，目前填“1”  */
	SubscribeMap          *int32 `json:"subscribeMap,omitempty"` /*  订阅TPS阈值，可以不设置  */
}

type Mq2CreateGroupResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
