package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtelbListElbLabelsApi
/* 获取负载均衡绑定的标签
 */type CtelbListElbLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListElbLabelsApi(client *core.CtyunClient) *CtelbListElbLabelsApi {
	return &CtelbListElbLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-labels",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListElbLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtelbListElbLabelsRequest) (*CtelbListElbLabelsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("elbID", req.ElbID)
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
	var resp CtelbListElbLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListElbLabelsRequest struct {
	RegionID string /*  区域ID  */
	ElbID    string /*  负载均衡 ID  */
	PageNo   int32  /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtelbListElbLabelsResponse struct {
	StatusCode   int32                                  `json:"statusCode,omitempty"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                 `json:"message,omitempty"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                 `json:"description,omitempty"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                 `json:"errorCode,omitempty"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	TotalCount   int32                                  `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                  `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                  `json:"totalPage,omitempty"`    /*  总页数  */
	ReturnObj    []*CtelbListElbLabelsReturnObjResponse `json:"returnObj"`              /*  返回结果  */
	Error        string                                 `json:"error,omitempty"`        /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListElbLabelsReturnObjResponse struct {
	Results []*CtelbListElbLabelsReturnObjResultsResponse `json:"results"` /*  绑定的标签列表  */
}

type CtelbListElbLabelsReturnObjResultsResponse struct {
	LabelID    string `json:"labelID,omitempty"`    /*  标签 id  */
	LabelKey   string `json:"labelKey,omitempty"`   /*  标签名  */
	LabelValue string `json:"labelValue,omitempty"` /*  标签值  */
}
