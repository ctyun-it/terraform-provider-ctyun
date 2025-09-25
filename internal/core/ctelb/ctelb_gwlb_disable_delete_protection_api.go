package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbGwlbDisableDeleteProtectionApi
/* 网关负载均衡关闭删除保护
 */type CtelbGwlbDisableDeleteProtectionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbDisableDeleteProtectionApi(client *core.CtyunClient) *CtelbGwlbDisableDeleteProtectionApi {
	return &CtelbGwlbDisableDeleteProtectionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/gwlb/disable-delete-protection",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbDisableDeleteProtectionApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbDisableDeleteProtectionRequest) (*CtelbGwlbDisableDeleteProtectionResponse, error) {
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
	var resp CtelbGwlbDisableDeleteProtectionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbDisableDeleteProtectionRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	GwLbID   string `json:"gwLbID,omitempty"`   /*  网关负载均衡ID  */
}

type CtelbGwlbDisableDeleteProtectionResponse struct {
	StatusCode  int32                                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbDisableDeleteProtectionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbDisableDeleteProtectionReturnObjResponse struct {
	GwLbID string `json:"gwLbID,omitempty"` /*  网关负载均衡 ID  */
}
