package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanListSdwanApi
/* 查询Sdwan列表 */
type SdwanListSdwanApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanListSdwanApi(client *core.CtyunClient) *SdwanListSdwanApi {
	return &SdwanListSdwanApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanListSdwanApi) Do(ctx context.Context, credential core.Credential, req *SdwanListSdwanRequest) (*SdwanListSdwanResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.SdwanID != nil && *req.SdwanID != "" {
		ctReq.AddParam("sdwanID", *req.SdwanID)
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
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanListSdwanResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanListSdwanRequest struct {
	SdwanID   *string `json:"sdwanID,omitempty"`   /*  SD-WAN ID  */
	ProjectID *string `json:"projectID,omitempty"` /*  企业项目ID  */
	PageNo    int32   `json:"pageNo"`              /*  页码  */
	PageSize  int32   `json:"pageSize"`            /*  每页记录数目  */
}

type SdwanListSdwanResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                          `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                          `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                          `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanListSdwanReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                          `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanListSdwanReturnObjResponse struct {
	Message      *string                                  `json:"message"`      /*  message  */
	TotalCount   int32                                    `json:"totalCount"`   /*  总数  */
	CurrentCount int32                                    `json:"currentCount"` /*  页码  */
	Code         *string                                  `json:"code"`         /*  code  */
	Result       []*SdwanListSdwanReturnObjResultResponse `json:"result"`       /*  列表  */
}

type SdwanListSdwanReturnObjResultResponse struct {
	SdwanID       *string `json:"sdwanID"`       /*  SD-WAN ID  */
	SdwanName     *string `json:"sdwanName"`     /*  名称  */
	ProjectID     *string `json:"projectID"`     /*  企业项目ID  */
	Description   *string `json:"description"`   /*  描述  */
	EcName        *string `json:"ecName"`        /*  云间高速name  */
	EcID          *string `json:"ecID"`          /*  云间高速 id  */
	CustomerID    *string `json:"customerID"`    /*  客户id  */
	EdgeNum       int32   `json:"edgeNum"`       /*  设备数量  */
	OnlineEdgeNum int32   `json:"onlineEdgeNum"` /*  在线设备数量  */
	CreateTime    *string `json:"createTime"`    /*  创建时间  */
}
