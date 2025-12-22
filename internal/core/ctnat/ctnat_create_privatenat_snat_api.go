package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatCreatePrivatenatSnatApi
/* 创建SNAT规则
 */type CtnatCreatePrivatenatSnatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatCreatePrivatenatSnatApi(client *core.CtyunClient) *CtnatCreatePrivatenatSnatApi {
	return &CtnatCreatePrivatenatSnatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/create-snat",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatCreatePrivatenatSnatApi) Do(ctx context.Context, credential core.Credential, req *CtnatCreatePrivatenatSnatRequest) (*CtnatCreatePrivatenatSnatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatCreatePrivatenatSnatRequest
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
	var resp CtnatCreatePrivatenatSnatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatCreatePrivatenatSnatRequest struct {
	RegionID       string   `json:"regionID,omitempty"`       /*  区域id  */
	NatGatewayID   string   `json:"natGatewayID,omitempty"`   /*  私网NAT ID  */
	SnatIPs        []string `json:"snatIPs"`                  /*  IP地址，必须在中转网段指定的网络范围内  */
	SourceCIDR     string   `json:"sourceCIDR,omitempty"`     /*  源CIDR和sourceSubnetID必须传一个  */
	SourceSubnetID string   `json:"sourceSubnetID,omitempty"` /*  源子网ID和sourceCIDR必须传一个  */
	Description    string   `json:"description,omitempty"`    /*  <支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&\*()\_\-+= <>?:\"{}\&#124;,.\/;'[\]·~！@#￥%……&\*（） ——\-+={} &#124; 《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtnatCreatePrivatenatSnatResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                      `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                      `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatCreatePrivatenatSnatReturnObjResponse `json:"returnObj"`   /*  返回创建结果  */
}

type CtnatCreatePrivatenatSnatReturnObjResponse struct {
	SnatID string `json:"snatID"` /*  SNAT规则的ID。  */
}
