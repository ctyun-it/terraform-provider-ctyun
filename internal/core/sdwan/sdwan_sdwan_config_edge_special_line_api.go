package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanConfigEdgeSpecialLineApi
/* 增加专线配置 */
type SdwanSdwanConfigEdgeSpecialLineApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanConfigEdgeSpecialLineApi(client *core.CtyunClient) *SdwanSdwanConfigEdgeSpecialLineApi {
	return &SdwanSdwanConfigEdgeSpecialLineApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/config-edge-special-line",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanConfigEdgeSpecialLineApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanConfigEdgeSpecialLineRequest) (*SdwanSdwanConfigEdgeSpecialLineResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanConfigEdgeSpecialLineRequest
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
	var resp SdwanSdwanConfigEdgeSpecialLineResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanConfigEdgeSpecialLineRequest struct {
	EdgeID                string                                                         `json:"edgeID"`                          /*  edge的ID  */
	SpecialLineConfigList []*SdwanSdwanConfigEdgeSpecialLineSpecialLineConfigListRequest `json:"specialLineConfigList,omitempty"` /*  专线配置  */
}

type SdwanSdwanConfigEdgeSpecialLineSpecialLineConfigListRequest struct {
	EdgePort     string `json:"edgePort"`     /*  端口(端口必须是mpls的)  */
	SerialNumber string `json:"serialNumber"` /*  盒子SN  */
	LineType     string `json:"lineType"`     /*  本参数表示专线类型<br/><br/>取值范围:<br/>static:静态<br/>dhcp:dhcp  */
	LineStatus   string `json:"lineStatus"`   /*  本参数表示链路状态<br/><br/>取值范围:<br/>processing:使用中  */
	IPAddress    string `json:"IPAddress"`    /*  ip地址  */
	SubMask      string `json:"subMask"`      /*  子网掩码  */
}

type SdwanSdwanConfigEdgeSpecialLineResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	OperationID *string `json:"operationID"` /*  操作日志Id  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
