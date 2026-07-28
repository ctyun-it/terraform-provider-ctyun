package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EcEcListApi
/* 查询已创建的云间高速 */
type EcEcListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListApi(client *core.CtyunClient) *EcEcListApi {
	return &EcEcListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/express-connect/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListApi) Do(ctx context.Context, credential core.Credential, req *EcEcListRequest) (*EcEcListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if *req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(*req.PageNo), 10))
	}
	if *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.EcID != nil && *req.EcID != "" {
		ctReq.AddParam("ecID", *req.EcID)
	}
	if req.QueryContent != nil && *req.QueryContent != "" {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListRequest struct {
	PageNo       *int32  `json:"pageNo,omitempty"`       /*  页码，从1开始  */
	PageSize     *int32  `json:"pageSize,omitempty"`     /*  每页记录数目  */
	EcID         *string `json:"ecID,omitempty"`         /*  云间高速实例ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊匹配，支持ecID,ecName, ecDescription三个属性  */
}

type EcEcListResponse struct {
	StatusCode  *int32                     `json:"statusCode"`  /*  返回状态码,<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                    `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                    `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                    `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                    `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListReturnObjResponse struct {
	CurrentCount *int32                              `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                              `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                              `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListReturnObjResultsResponse struct {
	EcID          *string `json:"ecID"`          /*  云间高速实例ID  */
	EcName        *string `json:"ecName"`        /*  名称  */
	EcDescription *string `json:"ecDescription"` /*  描述信息  */
	Status        *int32  `json:"status"`        /*  运行状态，<br/>取值范围:<br/>1:不可用<br/>2:可用  */
	CreateDate    *string `json:"createDate"`    /*  创建时间  */
	Vrf           *int32  `json:"vrf"`           /*  云间高速vrf信息  */
	Email         *string `json:"email"`         /*  email  */
	VpcCount      *int32  `json:"vpcCount"`      /*  添加vpc网络实例的数量  */
	CgwCount      *int32  `json:"cgwCount"`      /*  云网关数量  */
	CdaCount      *int32  `json:"cdaCount"`      /*  专线数量  */
	SdwanCount    *int32  `json:"sdwanCount"`    /*  sdwan实例数量  */
	VpnCount      *int32  `json:"vpnCount"`      /*  vpn实例数量  */
	EdsCount      *int32  `json:"edsCount"`      /*  云桌面网络实例数量  */
	Project       *string `json:"project"`       /*  企业项目  */
	PacketStatus  *int32  `json:"packetStatus"`  /*  带宽包状态<br/>本参数表示带宽类型,取值范围:<br/>1:已购买<br/>0:未购买  */
	PacketRate    *int32  `json:"packetRate"`    /*  带宽包总量  */
	CrossReigon   *int32  `json:"crossReigon"`   /*  已配置带宽,已废弃  */
}
