package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanAddDNATApi
/* 增加DNAT */
type SdwanSdwanAddDNATApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanAddDNATApi(client *core.CtyunClient) *SdwanSdwanAddDNATApi {
	return &SdwanSdwanAddDNATApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/dnat/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanAddDNATApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanAddDNATRequest) (*SdwanSdwanAddDNATResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanAddDNATRequest
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
	var resp SdwanSdwanAddDNATResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanAddDNATRequest struct {
	SiteID       string `json:"siteID"`       /*  站点ID  */
	InternalIP   string `json:"internalIP"`   /*  本端私网IP  */
	Protocol     string `json:"protocol"`     /*  本参数表示协议类型<br/><br/>取值范围:<br/>UDP:UDP</br>TCP:TCP  */
	InternalPort string `json:"internalPort"` /*  外服务端口  */
	ExternalPort string `json:"externalPort"` /*  内网端口  */
}

type SdwanSdwanAddDNATResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
