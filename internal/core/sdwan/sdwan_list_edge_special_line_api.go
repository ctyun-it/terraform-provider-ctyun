package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanListEdgeSpecialLineApi
/* 查询智能网关专线配置 */
type SdwanListEdgeSpecialLineApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanListEdgeSpecialLineApi(client *core.CtyunClient) *SdwanListEdgeSpecialLineApi {
	return &SdwanListEdgeSpecialLineApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-special-line/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanListEdgeSpecialLineApi) Do(ctx context.Context, credential core.Credential, req *SdwanListEdgeSpecialLineRequest) (*SdwanListEdgeSpecialLineResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanListEdgeSpecialLineResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanListEdgeSpecialLineRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanListEdgeSpecialLineResponse struct {
	StatusCode   int32   `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode    *string `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32   `json:"totalCount"`   /*  总数  */
	CurrentCount int32   `json:"currentCount"` /*  当前页数量  */
	UserWay      *string `json:"userWay"`      /*  设备分类  */
	SerialNumber *string `json:"serialNumber"` /*  盒子SN  */
	LineType     *string `json:"lineType"`     /*  本参数表示专线类型<br/><br/>取值范围：<br/>dhcp:DHCP<br/>static:静态  */
	LineStatus   *string `json:"lineStatus"`   /*  本参数表示链路状态<br/><br/>取值范围：<br/>processing:使用中<br/>completed:已完成  */
	IPAddress    *string `json:"IPAddress"`    /*  ip地址  */
	SubMask      *string `json:"subMask"`      /*  子网掩码  */
	Error        *string `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
