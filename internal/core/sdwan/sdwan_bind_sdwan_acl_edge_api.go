package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanBindSdwanAclEdgeApi
/* 访问控制与智能接入网关绑定 */
type SdwanBindSdwanAclEdgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanBindSdwanAclEdgeApi(client *core.CtyunClient) *SdwanBindSdwanAclEdgeApi {
	return &SdwanBindSdwanAclEdgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/acl-edge/bind",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanBindSdwanAclEdgeApi) Do(ctx context.Context, credential core.Credential, req *SdwanBindSdwanAclEdgeRequest) (*SdwanBindSdwanAclEdgeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanBindSdwanAclEdgeRequest
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
	var resp SdwanBindSdwanAclEdgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanBindSdwanAclEdgeRequest struct {
	AclID   string   `json:"aclID"`   /*  ACL ID  */
	Action  string   `json:"action"`  /*  动作create/retry  */
	EdgeIds []string `json:"edgeIds"` /*  盒子ID  ，值类型为string  */
}

type SdwanBindSdwanAclEdgeResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanBindSdwanAclEdgeReturnObjResponse `json:"returnObj"`   /*  结果列表  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanBindSdwanAclEdgeReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作id  */
}
