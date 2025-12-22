package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcRuleListAreasApi
/* 查询跨区互通 */
type EcEcRuleListAreasApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcRuleListAreasApi(client *core.CtyunClient) *EcEcRuleListAreasApi {
	return &EcEcRuleListAreasApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cross-region/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcRuleListAreasApi) Do(ctx context.Context, credential core.Credential, req *EcEcRuleListAreasRequest) (*EcEcRuleListAreasResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcRuleListAreasResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcRuleListAreasRequest struct{}

type EcEcRuleListAreasResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcRuleListAreasReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcRuleListAreasReturnObjResponse struct {
	CurrentCount *int32                                       `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                       `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                       `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcRuleListAreasReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcRuleListAreasReturnObjResultsResponse struct {
	AreaID           *string `json:"areaID"`           /*  跨区互通ID  */
	AreaA            *string `json:"areaA"`            /*  area,_A  */
	AreaB            *string `json:"areaB"`            /*  area,_B  */
	EcID             *string `json:"ecID"`             /*  云间高速实例ID  */
	PacketID         *string `json:"packetID"`         /*  带宽包ID  */
	SourceCgwID      *string `json:"sourceCgwID"`      /*  本端网关ID  */
	SourceCgwName    *string `json:"sourceCgwName"`    /*  本端网关名称  */
	SourceDcID       *string `json:"sourceDcID"`       /*  本端资源池ID  */
	SourceDcName     *string `json:"sourceDcName"`     /*  本端资源池名称  */
	DestCgwID        *string `json:"destCgwID"`        /*  对端网关ID  */
	DestCgwName      *string `json:"destCgwName"`      /*  对端网关名称  */
	DestDcID         *string `json:"destDcID"`         /*  对端资源池ID  */
	DestDcName       *string `json:"destDcName"`       /*  对端资源池名称  */
	PacketName       *string `json:"packetName"`       /*  带宽包名称  */
	Rate             *int32  `json:"rate"`             /*  带宽值（MB）  */
	Status           *string `json:"status"`           /*  状态描述<br/>取值包括：<br/>'creating'：加载中<br/>'running'：已连接<br/>'removing'：卸载中<br/>'rate_updating'：路由待更新<br/>'error'：失败  */
	PacketHasExpired *bool   `json:"packetHasExpired"` /*  带宽包是否到期  */
}
