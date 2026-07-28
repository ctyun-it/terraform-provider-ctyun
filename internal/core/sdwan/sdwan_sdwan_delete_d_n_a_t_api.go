package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanDeleteDNATApi
/* 删除DNAT */
type SdwanSdwanDeleteDNATApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanDeleteDNATApi(client *core.CtyunClient) *SdwanSdwanDeleteDNATApi {
	return &SdwanSdwanDeleteDNATApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/dnat/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanDeleteDNATApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanDeleteDNATRequest) (*SdwanSdwanDeleteDNATResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanDeleteDNATRequest
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
	var resp SdwanSdwanDeleteDNATResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanDeleteDNATRequest struct {
	SiteID     string   `json:"siteID"`     /*  站点ID  */
	DnatIDList []string `json:"dnatIDList"` /*  业务ID列表，值类型为string  */
}

type SdwanSdwanDeleteDNATResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanDeleteDNATReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanDeleteDNATReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
