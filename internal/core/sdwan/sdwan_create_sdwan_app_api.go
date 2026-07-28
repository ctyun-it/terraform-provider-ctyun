package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanAppApi
/* 应用创建 */
type SdwanCreateSdwanAppApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanAppApi(client *core.CtyunClient) *SdwanCreateSdwanAppApi {
	return &SdwanCreateSdwanAppApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/app/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanAppApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanAppRequest) (*SdwanCreateSdwanAppResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanAppRequest
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
	var resp SdwanCreateSdwanAppResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanAppRequest struct {
	GroupID  string                                `json:"groupID"`            /*  所属应用组ID  */
	AppName  string                                `json:"appName"`            /*  模糊搜索应用组名称  */
	RuleType string                                `json:"ruleType"`           /*  本参数表示规则类型<br/><br/>取值范围:<br/>hostname:网络域名<br/>ace:五元组  */
	RuleList []*SdwanCreateSdwanAppRuleListRequest `json:"ruleList,omitempty"` /*  规则列表  */
}

type SdwanCreateSdwanAppRuleListRequest struct {
	SrcCidr  string  `json:"srcCidr"`            /*  源网段  */
	SrcPort  *string `json:"srcPort,omitempty"`  /*  源端口范围  */
	DstCidr  *string `json:"dstCidr,omitempty"`  /*  目的网段  */
	DstPort  *string `json:"dstPort,omitempty"`  /*  目的端口范围  */
	Protocol *string `json:"protocol,omitempty"` /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	Hostname *string `json:"hostname,omitempty"` /*  访问地址  */
}

type SdwanCreateSdwanAppResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanCreateSdwanAppReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanCreateSdwanAppReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
