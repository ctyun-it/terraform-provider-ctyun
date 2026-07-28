package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanQosRuleApi
/* 查询qos规则 */
type SdwanGetSdwanQosRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanQosRuleApi(client *core.CtyunClient) *SdwanGetSdwanQosRuleApi {
	return &SdwanGetSdwanQosRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/qos-rule/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanQosRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanQosRuleRequest) (*SdwanGetSdwanQosRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("qosID", req.QosID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanQosRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanQosRuleRequest struct {
	QosID    string `json:"qosID"`    /*  qos策略ID  */
	PageNo   int32  `json:"pageNo"`   /*  页数  */
	PageSize int32  `json:"pageSize"` /*  页大小  */
}

type SdwanGetSdwanQosRuleResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanQosRuleReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanQosRuleReturnObjResponse struct {
	Result       []*SdwanGetSdwanQosRuleReturnObjResultResponse `json:"result"`       /*  查询qos规则  */
	TotalCount   int32                                          `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                          `json:"currentCount"` /*  页码  */
}

type SdwanGetSdwanQosRuleReturnObjResultResponse struct {
	RuleID   *string                                                `json:"ruleID"`   /*  规则id  */
	Priority *string                                                `json:"priority"` /*  优先级  */
	Rate     *string                                                `json:"rate"`     /*  速率  */
	RuleType *string                                                `json:"ruleType"` /*  本参数表示规则类型<br/><br/>取值范围:<br/>ace:五元组<br/>app:应用组  */
	ItemList []*SdwanGetSdwanQosRuleReturnObjResultItemListResponse `json:"itemList"` /*  qos规则元素列表  */
}

type SdwanGetSdwanQosRuleReturnObjResultItemListResponse struct {
	ItemID       *string                                                    `json:"itemID"`       /*  qos策略应用实例ID  */
	ItemName     *string                                                    `json:"itemName"`     /*  qos规则五元组/应用组名称  */
	IpVersion    *string                                                    `json:"ipVersion"`    /*  本参数表示ip版本<br/><br/>取值范围:<br/>ipv4:ipv4<br/>ipv6:ipv6  */
	DstCidr      *string                                                    `json:"dstCidr"`      /*  目的子网  */
	DstPortRange *string                                                    `json:"dstPortRange"` /*  目的端口范围（例如1-200， -1/-1为默认值，表示1-65535)  */
	SrcCidr      *string                                                    `json:"srcCidr"`      /*  源子网  */
	SrcPortRange *string                                                    `json:"srcPortRange"` /*  源端口范围（例如1-200， -1/-1为默认值，表示1-65535)  */
	RuleType     *string                                                    `json:"ruleType"`     /*  本参数表示规则类型<br/><br/>取值范围:<br/>ace:五元组<br/>app:应用组  */
	GroupType    *string                                                    `json:"groupType"`    /*  本参数表示应用组类型<br/><br/>取值范围:<br/>system:系统<br/>custom:自定义  */
	Apps         []*SdwanGetSdwanQosRuleReturnObjResultItemListAppsResponse `json:"apps"`         /*  app列表  */
}

type SdwanGetSdwanQosRuleReturnObjResultItemListAppsResponse struct {
	AppName *string `json:"appName"` /*  app名称  */
	AppID   *string `json:"appID"`   /*  app id  */
}
