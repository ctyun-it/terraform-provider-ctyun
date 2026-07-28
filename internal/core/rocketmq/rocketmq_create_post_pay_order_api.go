package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqCreatePostPayOrderApi
/* 开通 RocketMQ 实例 */
type RocketmqCreatePostPayOrderApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqCreatePostPayOrderApi(client *core.CtyunClient) *RocketmqCreatePostPayOrderApi {
	return &RocketmqCreatePostPayOrderApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/createPostPayOrder",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqCreatePostPayOrderApi) Do(ctx context.Context, credential core.Credential, req *RocketmqCreatePostPayOrderRequest) (*RocketmqCreatePostPayOrderResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqCreatePostPayOrderRequest
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
	var resp RocketmqCreatePostPayOrderResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type RocketmqCreatePostPayOrderRequest struct {
	RegionId            string `json:"regionId"`                      /*  资源池编码  */
	SpecName            string `json:"specName"`                      /*  规格名称，通过规格接口获取支持的规格列表  */
	BillMode            int32  `json:"billMode,omitempty"`            /*  计费模式<br/>取值范围:<br/>1 为预付费，即包年包月，包周期；<br/>2 为后付费，即按需；<br/>默认值为按需。  */
	CycleCnt            int32  `json:"cycleCnt,omitempty"`            /*  billMode 为 1 时生效且为必选值，表示预付费的订购周期时长，单位：月。<br/>取值范围：1，2，3，4，5，6，12，24，36  */
	AutoRenew           bool   `json:"autoRenew,omitempty"`           /*  预付费实例是否自动续订<br/>取值范围：<br/>true：自动续订。<br/>false：不自动续订。<br/>默认不自动续订。  */
	AutoRenewCycleCount int32  `json:"autoRenewCycleCount,omitempty"` /*  autoRenew 为 true 时生效且为必选值，表示自动续订的周期时长，单位：月。<br/>取值范围：1，2，3，4，5，6，12，24，36  */
	NodeNum             int32  `json:"nodeNum"`                       /*  broker 的节点数，值等于代理数*2，取值范围为 [1，32]，单机版传 1  */
	DiskType            string `json:"diskType"`                      /*  存储类型 SAS、SSD、FAST-SSD  */
	DiskSize            int32  `json:"diskSize"`                      /*  节点存储空间大小，单位 G  */
	AzInfo              string `json:"azInfo"`                        /*  az 信息通过规格接口查询，支持单可用区和三可用区  */
	VpcId               string `json:"vpcId"`                         /*  私有云 ID  */
	SecurityGroupId     string `json:"securityGroupId"`               /*  安全组 ID  */
	SubnetId            string `json:"subnetId"`                      /*  子网 ID  */
	ClusterName         string `json:"clusterName"`                   /*  实例名称  */
}

type RocketmqCreatePostPayOrderResponse struct {
	StatusCode int32                         `json:"statusCode"` // 接口系统层面状态码。成功：800，失败：900
	Message    string                        `json:"message"`    // 描述状态
	ReturnObj  *RocketmqCreateOrderReturnObj `json:"returnObj"`  // 返回对象
	Error      string                        `json:"error"`      // 错误码，只有非成功才有这个字段，方便快速定位问题
}

type RocketmqCreateOrderReturnObj struct {
	Data *OrderData `json:"data"` // 订单核心数据。仅接口调用成功时返回
}
