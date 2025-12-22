package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatModifyPrivatenatDnatApi
/* 修改DNAT
 */type CtnatModifyPrivatenatDnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatModifyPrivatenatDnatApi(client *core.CtyunClient) *CtnatModifyPrivatenatDnatApi {
	return &CtnatModifyPrivatenatDnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/modify-dnat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatModifyPrivatenatDnatApi) Do(ctx context.Context, credential core.Credential, req *CtnatModifyPrivatenatDnatRequest) (*CtnatModifyPrivatenatDnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatModifyPrivatenatDnatRequest
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
	var resp CtnatModifyPrivatenatDnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatModifyPrivatenatDnatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  私网NAT所在的Region。  */
	NatGatewayID string `json:"natGatewayID,omitempty"` //NAT网关Id
	DnatID       string `json:"dnatID,omitempty"`       /*  要修改的DNAT的ID。  */
	ExternalIP   string `json:"externalIP,omitempty"`   /*  中转IP  */
	ExternalPort int32  `json:"externalPort,omitempty"` /*  外部端口（1-65535）  */
	InternalIP   string `json:"internalIP,omitempty"`   /*  内部IP  */
	PortID       string `json:"portID,omitempty"`       /*  网卡ID  */
	InternalPort int32  `json:"internalPort,omitempty"` /*  内部端口（1-65535）  */
	Protocol     string `json:"protocol,omitempty"`     /*  协议: tcp/udp  */
	Description  string `json:"description,omitempty"`  /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_\-+= <>?:\"{} &#124; ,.\/;'[\]·~！@#￥%……&*（） ——\-+={}&#124;《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtnatModifyPrivatenatDnatResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatModifyPrivatenatDnatReturnObjResponse `json:"returnObj"`   /*  接口业务数据  */
}

type CtnatModifyPrivatenatDnatReturnObjResponse struct{}
