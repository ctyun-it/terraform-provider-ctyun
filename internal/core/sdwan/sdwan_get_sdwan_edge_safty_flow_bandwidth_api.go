package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanEdgeSaftyFlowBandwidthApi
/* 安全引流带宽查询 */
type SdwanGetSdwanEdgeSaftyFlowBandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanEdgeSaftyFlowBandwidthApi(client *core.CtyunClient) *SdwanGetSdwanEdgeSaftyFlowBandwidthApi {
	return &SdwanGetSdwanEdgeSaftyFlowBandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-safty-flow-bandwidth/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanEdgeSaftyFlowBandwidthApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanEdgeSaftyFlowBandwidthRequest) (*SdwanGetSdwanEdgeSaftyFlowBandwidthResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanEdgeSaftyFlowBandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanEdgeSaftyFlowBandwidthRequest struct {
	PageNo   int32 `json:"pageNo"`   /*  页码  */
	PageSize int32 `json:"pageSize"` /*  每页记录数目  */
}

type SdwanGetSdwanEdgeSaftyFlowBandwidthResponse struct {
	StatusCode  int32                                                 `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                               `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                               `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                               `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanEdgeSaftyFlowBandwidthReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type SdwanGetSdwanEdgeSaftyFlowBandwidthReturnObjResponse struct {
	TotalCount   int32                                                         `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                                         `json:"currentCount"` /*  当前页面数量  */
	Code         *string                                                       `json:"code"`         /*  code  */
	Result       []*SdwanGetSdwanEdgeSaftyFlowBandwidthReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanEdgeSaftyFlowBandwidthReturnObjResultResponse struct {
	InstanceID *string `json:"instanceID"` /*  实例id  */
	Email      *string `json:"email"`      /*  账号邮箱  */
	CustomerID *string `json:"customerID"` /*  用户ID  */
	Bandwidth  *string `json:"bandwidth"`  /*  带宽  */
	StartTime  *string `json:"startTime"`  /*  起始时间  */
	ExpireTime *string `json:"expireTime"` /*  到期时间  */
}
