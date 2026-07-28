package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetSdwanApp1Api
/* 应用查询 */
type SdwanGetSdwanApp1Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanApp1Api(client *core.CtyunClient) *SdwanGetSdwanApp1Api {
	return &SdwanGetSdwanApp1Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/app/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanApp1Api) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanApp1Request) (*SdwanGetSdwanApp1Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.GroupID != nil && *req.GroupID != "" {
		ctReq.AddParam("groupID", *req.GroupID)
	}
	if req.AppType != nil && *req.AppType != "" {
		ctReq.AddParam("appType", *req.AppType)
	}
	if req.AppID != nil && *req.AppID != "" {
		ctReq.AddParam("appID", *req.AppID)
	}
	if req.AppName != nil && *req.AppName != "" {
		ctReq.AddParam("appName", *req.AppName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanApp1Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanApp1Request struct {
	PageNo   int32   `json:"pageNo"`            /*  页码  */
	PageSize int32   `json:"pageSize"`          /*  每页记录数目  */
	Search   *string `json:"search,omitempty"`  /*  模糊查询  */
	GroupID  *string `json:"groupID,omitempty"` /*  所属应用组ID  */
	AppType  *string `json:"appType,omitempty"` /*  本参数表示应用类型<br/><br/>取值范围:<br/>custom:自定义<br/>system:系统  */
	AppID    *string `json:"appID,omitempty"`   /*  应用ID  */
	AppName  *string `json:"appName,omitempty"` /*  模糊搜索应用组名称  */
}

type SdwanGetSdwanApp1Response struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetSdwanApp1ReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetSdwanApp1ReturnObjResponse struct {
	Message      *string                                     `json:"message"`      /*  message  */
	TotalCount   int32                                       `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                       `json:"currentCount"` /*  当前页数量  */
	Code         *string                                     `json:"code"`         /*  状态码  */
	Result       []*SdwanGetSdwanApp1ReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanGetSdwanApp1ReturnObjResultResponse struct {
	AppName     *string                                        `json:"appName"`     /*  模糊搜索应用组名称  */
	AppID       *string                                        `json:"appID"`       /*  应用ID  */
	AppType     *string                                        `json:"appType"`     /*  本参数表示应用类型<br/><br/>取值范围:<br/>custom:自定义<br/>system:系统  */
	CustomerID  *string                                        `json:"customerID"`  /*  用户ID  */
	CreatedTime *string                                        `json:"createdTime"` /*  创建时间  */
	Protocol    *string                                        `json:"protocol"`    /*  本参数表示协议<br/><br/>取值范围:<br/>tcp:tcp<br/>udp:udp<br/>icmp:icmp  */
	SrcCidr     *string                                        `json:"srcCidr"`     /*  源网段  */
	SrcPort     *string                                        `json:"srcPort"`     /*  源端口范围  */
	DstCidr     *string                                        `json:"dstCidr"`     /*  目的网段  */
	Group       *SdwanGetSdwanApp1ReturnObjResultGroupResponse `json:"group"`       /*  所属应用组  */
}

type SdwanGetSdwanApp1ReturnObjResultGroupResponse struct {
	GroupName   *string `json:"groupName"`   /*  group_name  */
	GroupID     *string `json:"groupID"`     /*  group_id  */
	GroupType   *string `json:"groupType"`   /*  本参数表示应用组类型<br/><br/>取值范围:<br/>custom:自定义<br/>system:系统  */
	CustomerID  *string `json:"customerID"`  /*  用户ID  */
	CreatedTime *string `json:"createdTime"` /*  创建时间  */
}
