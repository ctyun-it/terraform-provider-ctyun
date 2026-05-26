package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteCustomDomainApi
/* 删除自定义域名 */
type CfDeleteCustomDomainApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteCustomDomainApi(client *core.CtyunClient) *CfDeleteCustomDomainApi {
	return &CfDeleteCustomDomainApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/domains/customdomains/{domainName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteCustomDomainApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteCustomDomainRequest) (*CfDeleteCustomDomainResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("domainName", req.DomainName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteCustomDomainResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteCustomDomainRequest struct {
	DomainName string `json:"domainName"` /*  域名  */
	RegionId   string `json:"regionId"`   /*  资源池id  */
}

type CfDeleteCustomDomainResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码  */
	Message    *string `json:"message"`    /*  信息  */
	ReturnObj  *bool   `json:"returnObj"`  /*  是否成功  */
}
