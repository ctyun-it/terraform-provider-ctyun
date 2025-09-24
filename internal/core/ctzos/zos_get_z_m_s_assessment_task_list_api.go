package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosGetZMSAssessmentTaskListApi
/* 查询对象存储评估任务列表
 */type ZosGetZMSAssessmentTaskListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetZMSAssessmentTaskListApi(client *core.CtyunClient) *ZosGetZMSAssessmentTaskListApi {
	return &ZosGetZMSAssessmentTaskListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/zms/list-evaluation",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetZMSAssessmentTaskListApi) Do(ctx context.Context, credential core.Credential, req *ZosGetZMSAssessmentTaskListRequest) (*ZosGetZMSAssessmentTaskListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.SearchName != "" {
		ctReq.AddParam("searchName", req.SearchName)
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
	var resp ZosGetZMSAssessmentTaskListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetZMSAssessmentTaskListRequest struct {
	RegionID   string /*  资源池ID  */
	SearchName string /*  模糊搜索评估任务名称  */
	PageNo     int32  /*  页码，默认值为1  */
	PageSize   int32  /*  页大小，默认值 10，取值范围 1~50  */
}

type ZosGetZMSAssessmentTaskListResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
