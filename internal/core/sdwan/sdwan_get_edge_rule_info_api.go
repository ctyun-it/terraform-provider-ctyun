package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeRuleInfoApi
/* 查找智能网关规则策略 */
type SdwanGetEdgeRuleInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeRuleInfoApi(client *core.CtyunClient) *SdwanGetEdgeRuleInfoApi {
	return &SdwanGetEdgeRuleInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-rule-info/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeRuleInfoApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeRuleInfoRequest) (*SdwanGetEdgeRuleInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeRuleInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeRuleInfoRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanGetEdgeRuleInfoResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeRuleInfoReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeRuleInfoReturnObjResponse struct {
	Result *SdwanGetEdgeRuleInfoReturnObjResultResponse `json:"result"` /*  acl和qos返回参数  */
}

type SdwanGetEdgeRuleInfoReturnObjResultResponse struct {
	Acl *SdwanGetEdgeRuleInfoReturnObjResultAclResponse   `json:"acl"` /*  acl信息  */
	Qos []*SdwanGetEdgeRuleInfoReturnObjResultQosResponse `json:"qos"` /*  qos信息  */
}

type SdwanGetEdgeRuleInfoReturnObjResultAclResponse struct {
	AclID   *string `json:"aclID"`   /*  访问控制id  */
	AclName *string `json:"aclName"` /*  访问控制名称  */
}

type SdwanGetEdgeRuleInfoReturnObjResultQosResponse struct {
	QOSID   *string `json:"QOSID"`   /*  qos id  */
	QOSName *string `json:"QOSName"` /*  qos名称  */
}
