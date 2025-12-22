package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatCreatePrivatenatDnatApi
/* 创建DNAT规则
 */type CtnatCreatePrivatenatDnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatCreatePrivatenatDnatApi(client *core.CtyunClient) *CtnatCreatePrivatenatDnatApi {
	return &CtnatCreatePrivatenatDnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/create-dnat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatCreatePrivatenatDnatApi) Do(ctx context.Context, credential core.Credential, req *CtnatCreatePrivatenatDnatRequest) (*CtnatCreatePrivatenatDnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatCreatePrivatenatDnatRequest
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
	var resp CtnatCreatePrivatenatDnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatCreatePrivatenatDnatRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  区域id  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  私网NAT的ID  */
	ExternalIP   string `json:"externalIP,omitempty"`   /*  中转IP  */
	ExternalPort int32  `json:"externalPort,omitempty"` /*  对外的端口（1-65535）  */
	InternalPort int32  `json:"internalPort,omitempty"` /*  对应的内部端口（1-65535）  */
	InternalIP   string `json:"internalIP,omitempty"`   /*  对应的内部IP(和portID二选一)  */
	PortID       string `json:"portID,omitempty"`       /*  对应的网卡ID(和internalIP二选一)  */
	Protocol     string `json:"protocol,omitempty"`     /*  协议: tcp, udp  */
	Description  string `json:"description,omitempty"`  /*  <支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&\*()\_\-+= <>?:\"{}\&#124;,.\/;'[\]·~！@#￥%……&\*（） ——\-+={} &#124; 《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtnatCreatePrivatenatDnatResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatCreatePrivatenatDnatReturnObjResponse `json:"returnObj"`   /*  object  */
}

type CtnatCreatePrivatenatDnatReturnObjResponse struct {
	DnatID string `json:"dnatID"` /*  DNAT规则的ID。  */
}
