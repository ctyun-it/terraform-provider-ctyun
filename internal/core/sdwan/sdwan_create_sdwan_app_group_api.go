package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateSdwanAppGroupApi
/* 应用组创建 */
type SdwanCreateSdwanAppGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateSdwanAppGroupApi(client *core.CtyunClient) *SdwanCreateSdwanAppGroupApi {
	return &SdwanCreateSdwanAppGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/app-group/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateSdwanAppGroupApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateSdwanAppGroupRequest) (*SdwanCreateSdwanAppGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateSdwanAppGroupRequest
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
	var resp SdwanCreateSdwanAppGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateSdwanAppGroupRequest struct {
	GroupName string  `json:"groupName"`           /*  应用组名称  */
	GroupType *string `json:"groupType,omitempty"` /*  本参数表示应用组类型<br/><br/>取值范围:<br/>custom:自定义<br/>system:系统  */
}

type SdwanCreateSdwanAppGroupResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                      `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                      `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                      `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanCreateSdwanAppGroupReturnObjResponse `json:"returnObj"`   /*  结果列表  */
	Error       *string                                      `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanCreateSdwanAppGroupReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
