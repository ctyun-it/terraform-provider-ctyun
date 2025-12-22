package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanEdgeLanNicInfoListApi
/* 查询edge lan口信息 */
type SdwanSdwanEdgeLanNicInfoListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanEdgeLanNicInfoListApi(client *core.CtyunClient) *SdwanSdwanEdgeLanNicInfoListApi {
	return &SdwanSdwanEdgeLanNicInfoListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-lan-nic-info/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanEdgeLanNicInfoListApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanEdgeLanNicInfoListRequest) (*SdwanSdwanEdgeLanNicInfoListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanSdwanEdgeLanNicInfoListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanEdgeLanNicInfoListRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanSdwanEdgeLanNicInfoListResponse struct {
	StatusCode  int32                                          `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                                        `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                        `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                        `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanEdgeLanNicInfoListReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                        `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanSdwanEdgeLanNicInfoListReturnObjResponse struct {
	Result       *SdwanSdwanEdgeLanNicInfoListReturnObjResultResponse `json:"result"`       /*  端口信息  */
	TotalCount   int32                                                `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                                `json:"currentCount"` /*  页码  */
}

type SdwanSdwanEdgeLanNicInfoListReturnObjResultResponse struct {
	NicList      []*SdwanSdwanEdgeLanNicInfoListReturnObjResultNicListResponse      `json:"nicList"`      /*  端口名称  */
	SlaveNicList []*SdwanSdwanEdgeLanNicInfoListReturnObjResultSlaveNicListResponse `json:"slaveNicList"` /*  备用端口名称  */
}

type SdwanSdwanEdgeLanNicInfoListReturnObjResultNicListResponse struct {
	PortName   *string `json:"portName"`   /*  端口名称  */
	PortType   *string `json:"portType"`   /*  本参数表示端口类型<br/><br/>取值范围：<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	PortStatus *string `json:"portStatus"` /*  本参数表示端口状态<br/><br/>取值范围：<br/>up:开启<br/>down:关闭  */
}

type SdwanSdwanEdgeLanNicInfoListReturnObjResultSlaveNicListResponse struct {
	PortName   *string `json:"portName"`   /*  端口名称  */
	PortType   *string `json:"portType"`   /*  本参数表示端口类型<br/><br/>取值范围：<br/>WAN1:WAN1<br/>WAN2:WAN2<br/>LAN:LAN<br/>LTE:LTE  */
	PortStatus *string `json:"portStatus"` /*  本参数表示端口状态<br/><br/>取值范围：<br/>up:开启<br/>down:关闭  */
}
