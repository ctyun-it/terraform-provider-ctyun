package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EcEcListSDWANInstanceApi
/* 查询sdwan网络实例 */
type EcEcListSDWANInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListSDWANInstanceApi(client *core.CtyunClient) *EcEcListSDWANInstanceApi {
	return &EcEcListSDWANInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ec/sdwan-instance/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListSDWANInstanceApi) Do(ctx context.Context, credential core.Credential, req *EcEcListSDWANInstanceRequest) (*EcEcListSDWANInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("ecID", req.EcID)
	if req.SdwanID != nil && *req.SdwanID != "" {
		ctReq.AddParam("sdwanID", *req.SdwanID)
	}
	if req.CgwID != nil && *req.CgwID != "" {
		ctReq.AddParam("cgwID", *req.CgwID)
	}
	if req.InstanceID != nil && *req.InstanceID != "" {
		ctReq.AddParam("instanceID", *req.InstanceID)
	}
	if req.QueryContent != nil && *req.QueryContent != "" {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.IsAuth != nil && *req.IsAuth != 0 {
		ctReq.AddParam("isAuth", strconv.FormatInt(int64(*req.IsAuth), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListSDWANInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListSDWANInstanceRequest struct {
	EcID         string  `json:"ecID"`                   /*  云间高速ID  */
	SdwanID      *string `json:"sdwanID,omitempty"`      /*  sdwan ID  */
	CgwID        *string `json:"cgwID,omitempty"`        /*  云网关ID  */
	InstanceID   *string `json:"instanceID,omitempty"`   /*  网络实例ID  */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊查询，支持sdwanID，网络实例ID，云网关ID三个属性  */
	IsAuth       *int32  `json:"isAuth,omitempty"`       /*  是否是跨账号实例<br/>取值包括：<br/>0:本账号<br/>1:跨账号  */
}

type EcEcListSDWANInstanceResponse struct {
	StatusCode  *int32                                  `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                 `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                 `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                 `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                 `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListSDWANInstanceReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListSDWANInstanceReturnObjResponse struct {
	CurrentCount *int32                                           `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                           `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                           `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListSDWANInstanceReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组  */
}

type EcEcListSDWANInstanceReturnObjResultsResponse struct {
	EcID              *string   `json:"ecID"`              /*  云间高速ID  */
	CgwID             *string   `json:"cgwID"`             /*  云网关ID  */
	CgwName           *string   `json:"cgwName"`           /*  云网关名称  */
	DcID              *string   `json:"dcID"`              /*  云网关资源池ID  */
	SdwanID           *string   `json:"sdwanID"`           /*  sdwan ID  */
	SdwanName         *string   `json:"sdwanName"`         /*  sdwan名称  */
	InstanceID        *string   `json:"instanceID"`        /*  网络实例的ID  */
	DefaultRtbID      *string   `json:"defaultRtbID"`      /*  默认路由表ID  */
	DefaultRtbName    *string   `json:"defaultRtbName"`    /*  默认路由表名称  */
	CIDR              []*string `json:"cidr"`              /*  v4 CIDR列表  */
	V6CIDR            []*string `json:"v6CIDR"`            /*  v6 CIDR列表  */
	CreateDate        *string   `json:"createDate"`        /*  创建时间  */
	Status            *string   `json:"status"`            /*  状态描述，取值包括：creating:加载中, running:已连接, removing:卸载中, flushing:路由待更新, error:失败  */
	Weights           *int32    `json:"weights"`           /*  权重，sdwan默认60  */
	RedundantType     *int32    `json:"redundantType"`     /*  冗余类型，取值包括：1：主备, 2：负载, 0：无  */
	RedundantInstUUID *string   `json:"redundantInstUUID"` /*  冗余负载实例的UUID  */
	RedundantInstName *string   `json:"redundantInstName"` /*  冗余负载实例的名称  */
	RedundantInstType *string   `json:"redundantInstType"` /*  冗余负载实例的类型，取值包括：2: 云专线, 3: sdwan, 4: vpn  */
	RedundantInstID   *string   `json:"redundantInstID"`   /*  冗余负载侧ID，即当前为cdaID, sdwanID，vpnID  */
	IsAuth            *int32    `json:"isAuth"`            /*  是否是跨账号实例，取值包括：0:本账号, 1:跨账号  */
	RouteLearn        *int32    `json:"routeLearn"`        /*  路由学习开关，开启后云网关自动学习网络实例路由，取值范围: 1:学习, 0:不学习, 默认学习  */
	RouteSync         *int32    `json:"routeSync"`         /*  路由同步开关，开启后云网关路由自动同步到网络实例，取值范围: 1:同步, 0:不同步, 默认同步  */
}
