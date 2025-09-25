package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcUpdateRouteTableRulesAttributeApi
/* 修改路由表规则
 */type CtvpcUpdateRouteTableRulesAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcUpdateRouteTableRulesAttributeApi(client *core.CtyunClient) *CtvpcUpdateRouteTableRulesAttributeApi {
	return &CtvpcUpdateRouteTableRulesAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpc/route-table/modify-rules",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcUpdateRouteTableRulesAttributeApi) Do(ctx context.Context, credential core.Credential, req *CtvpcUpdateRouteTableRulesAttributeRequest) (*CtvpcUpdateRouteTableRulesAttributeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcUpdateRouteTableRulesAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcUpdateRouteTableRulesAttributeRequest struct {
	ClientToken  string                                                  `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID     string                                                  `json:"regionID,omitempty"`     /*  区域id  */
	RouteTableID string                                                  `json:"routeTableID,omitempty"` /*  路由表 id  */
	RouteRules   []*CtvpcUpdateRouteTableRulesAttributeRouteRulesRequest `json:"routeRules"`             /*  路由表规则列表  */
}

type CtvpcUpdateRouteTableRulesAttributeRouteRulesRequest struct {
	NextHopID   string `json:"nextHopID,omitempty"`   /*  下一跳设备 id  */
	NextHopType string `json:"nextHopType,omitempty"` /*  vpcpeering / havip / bm / vm / natgw/ igw /igw6 / dc / ticc / vpngw / enic  */
	Destination string `json:"destination,omitempty"` /*  无类别域间路由  */
	IpVersion   int32  `json:"ipVersion"`             /*  4 标识 ipv4, 6 标识 ipv6  */
	Description string `json:"description,omitempty"` /*  规则描述,支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·~！@#￥%……&*（） —— -+={},《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	RouteRuleID string `json:"routeRuleID,omitempty"` /*  路由规则 id  */
}

type CtvpcUpdateRouteTableRulesAttributeResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	Error       *string `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}
