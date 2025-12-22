package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanListEdgeBindSdwanListApi
/* 查询绑定Sdwan智能网关信息 */
type SdwanListEdgeBindSdwanListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanListEdgeBindSdwanListApi(client *core.CtyunClient) *SdwanListEdgeBindSdwanListApi {
	return &SdwanListEdgeBindSdwanListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-bind-sdwan/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanListEdgeBindSdwanListApi) Do(ctx context.Context, credential core.Credential, req *SdwanListEdgeBindSdwanListRequest) (*SdwanListEdgeBindSdwanListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	if req.Status != nil && *req.Status != "" {
		ctReq.AddParam("status", *req.Status)
	}
	if req.ProjectID != nil && *req.ProjectID != "" {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanListEdgeBindSdwanListResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanListEdgeBindSdwanListRequest struct {
	SdwanID   string  `json:"sdwanID"`             /*  sdwan的id  */
	Status    *string `json:"status,omitempty"`    /*  本参数表示设备状态<br/><br/>取值范围:<br/>online:上线<br/>offline:下线  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目ID  */
	PageNo    int32   `json:"pageNo"`              /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页数目  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询  */
}

type SdwanListEdgeBindSdwanListResponse struct {
	StatusCode     int32   `json:"statusCode"`     /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode      *string `json:"errorCode"`      /*  业务细分码，为product.module.code三段式码  */
	Message        *string `json:"message"`        /*  失败时的错误描述，一般为英文描述  */
	Description    *string `json:"description"`    /*  失败时的错误描述，一般为中文描述  */
	TotalCount     int32   `json:"totalCount"`     /*  总数  */
	CurrentCount   int32   `json:"currentCount"`   /*  当前页数量  */
	EdgeID         *string `json:"edgeID"`         /*  站点ID  */
	EdgeName       *string `json:"edgeName"`       /*  站点名称  */
	Status         *string `json:"status"`         /*  本参数表示设备状态<br/><br/>取值范围:<br/>online:上线<br/>offline:下线  */
	IsAuth         *bool   `json:"isAuth"`         /*  本参数表示设备状态<br/><br/>取值范围:<br/> false: 同账号<br/> true: 跨账号  */
	CustomerEmail  *string `json:"customerEmail"`  /*  所属账号邮箱  */
	CustomerID     *string `json:"customerID"`     /*  所属账号ID  */
	Size           int32   `json:"size"`           /*  基础带宽大小,单位mbps  */
	TotalSize      int32   `json:"totalSize"`      /*  总带宽,单位mbps  */
	DeviceModel    *string `json:"deviceModel"`    /*  本参数表示设备型号<br/><br/>取值范围：<br/>economic:经济型<br/>standard:标准版<br/>enterprise:企业版<br/>enhance:企业增强版<br/>vcpe:虚拟智能网关  */
	CurrentVersion *string `json:"currentVersion"` /*  软件版本号  */
	Error          *string `json:"error"`          /*  业务细分码，为product.module.code三段式码  */
}
