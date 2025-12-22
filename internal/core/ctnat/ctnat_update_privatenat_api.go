package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatUpdatePrivatenatApi
/* 修改私网NAT
 */type CtnatUpdatePrivatenatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatUpdatePrivatenatApi(client *core.CtyunClient) *CtnatUpdatePrivatenatApi {
	return &CtnatUpdatePrivatenatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatUpdatePrivatenatApi) Do(ctx context.Context, credential core.Credential, req *CtnatUpdatePrivatenatRequest) (*CtnatUpdatePrivatenatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatUpdatePrivatenatRequest
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
	var resp CtnatUpdatePrivatenatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatUpdatePrivatenatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  私网NAT所在的Region。  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  要修改的私网NAT的ID。  */
	Name         string `json:"name,omitempty"`         /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description  string `json:"description,omitempty"`  /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_\-+= <>?:\"{} &#124; ,.\/;'[\]·~！@#￥%……&*（） ——\-+={}&#124;《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtnatUpdatePrivatenatResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatUpdatePrivatenatReturnObjResponse `json:"returnObj"`   /*  接口返回数据  */
}

type CtnatUpdatePrivatenatReturnObjResponse struct{}
