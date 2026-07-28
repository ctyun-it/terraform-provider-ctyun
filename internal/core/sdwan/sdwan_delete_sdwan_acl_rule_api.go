package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteSdwanAclRuleApi
/* 删除访问控制规则 */
type SdwanDeleteSdwanAclRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteSdwanAclRuleApi(client *core.CtyunClient) *SdwanDeleteSdwanAclRuleApi {
	return &SdwanDeleteSdwanAclRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/acl-rule/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteSdwanAclRuleApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteSdwanAclRuleRequest) (*SdwanDeleteSdwanAclRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteSdwanAclRuleRequest
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
	var resp SdwanDeleteSdwanAclRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteSdwanAclRuleRequest struct {
	AclID       string    `json:"aclID"`                 /*  ACL ID  */
	DeleteRules []*string `json:"deleteRules,omitempty"` /*  删除规则  ，值类型为  */
}

type SdwanDeleteSdwanAclRuleResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanDeleteSdwanAclRuleReturnObjResponse `json:"returnObj"`   /*  结果列表  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanDeleteSdwanAclRuleReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
