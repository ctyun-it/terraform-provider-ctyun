package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListGatewayApi
/* 查询已创建的云网关 */
type EcEcListGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListGatewayApi(client *core.CtyunClient) *EcEcListGatewayApi {
	return &EcEcListGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/cloud-gateway/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListGatewayApi) Do(ctx context.Context, credential core.Credential, req *EcEcListGatewayRequest) (*EcEcListGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.CgwID != nil && *req.CgwID != "" {
		ctReq.AddParam("cgwID", *req.CgwID)
	}
	if req.QueryContent != nil && *req.QueryContent != "" {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.Region != nil && *req.Region != "" {
		ctReq.AddParam("region", *req.Region)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListGatewayRequest struct {
	EcID         string  `json:"ecID"`                   /*  云间高速实例ID  */
	CgwID        *string `json:"cgwID,omitempty"`        /*  云网关实例ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊匹配，支持cgwID,cgwName,cgwDescription三个属性  */
	Region       *string `json:"region,omitempty"`       /*  地域信息，不填默认查询全部<br/>取值如下<br/>1：中国大陆<br/>2:亚太  */
}

type EcEcListGatewayResponse struct {
	StatusCode  *int32                            `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                           `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                           `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                           `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                           `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListGatewayReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListGatewayReturnObjResponse struct {
	CurrentCount *int32                                     `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                     `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                     `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListGatewayReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListGatewayReturnObjResultsResponse struct {
	CgwID          *string `json:"cgwID"`          /*  云网关实例ID  */
	Region         *int64  `json:"region"`         /*  Integer地域信息，不填默认查询全部<br/>取值如下<br/>1：中国大陆<br/>2:亚太  */
	CgwName        *string `json:"cgwName"`        /*  云网关名称  */
	CgwDescription *string `json:"cgwDescription"` /*  云网关描述  */
	EcID           *string `json:"ecID"`           /*  云间高速实例ID  */
	DcID           *string `json:"dcID"`           /*  资源池ID信息  */
	DcType         *string `json:"dcType"`         /*  资源池类型，<br/>取值范围:<br/>"CNP":CNP资源池<br/>"MAZ":MAZ资源池  */
	DcName         *string `json:"dcName"`         /*  资源池名称  */
	RtbCnt         *int32  `json:"rtbCnt"`         /*  路由表数量  */
	RouteCnt       *int32  `json:"routeCnt"`       /*  路由数量  */
	PolicyCount    *int32  `json:"policyCount"`    /*  流量策略数量  */
	InsCnt         *int32  `json:"insCnt"`         /*  连接实例数量  */
	CreateDate     *string `json:"createDate"`     /*  创建时间  */
	DefaultRtbID   *string `json:"defaultRtbID"`   /*  云网关默认路由表ID  */
	HasMonitor     *bool   `json:"hasMonitor"`     /*  是否支持监控 <br/>true:支持，false:不支持  */
}
