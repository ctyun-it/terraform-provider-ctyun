package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcVpcVpcpeerListLabelsApi
/* 获取对等链接绑定的标签
 */type CtvpcVpcVpcpeerListLabelsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVpcVpcpeerListLabelsApi(client *core.CtyunClient) *CtvpcVpcVpcpeerListLabelsApi {
	return &CtvpcVpcVpcpeerListLabelsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/vpcpeer/list-labels",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVpcVpcpeerListLabelsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVpcVpcpeerListLabelsRequest) (*CtvpcVpcVpcpeerListLabelsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("vpcPeerID", req.VpcPeerID)
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVpcVpcpeerListLabelsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVpcVpcpeerListLabelsRequest struct {
	RegionID   string /*  区域ID  */
	VpcPeerID  string /*  对等链接 ID  */
	PageNumber int32  /*  列表的页码，默认值为 1  */
	PageSize   int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10  */
}

type CtvpcVpcVpcpeerListLabelsResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcVpcVpcpeerListLabelsReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcVpcVpcpeerListLabelsReturnObjResponse struct {
	Results      []*CtvpcVpcVpcpeerListLabelsReturnObjResultsResponse `json:"results"`      /*  绑定的标签列表  */
	TotalCount   int32                                                `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                `json:"totalPage"`    /*  总页数  */
}

type CtvpcVpcVpcpeerListLabelsReturnObjResultsResponse struct {
	LabelID    *string `json:"labelID,omitempty"`    /*  标签 id  */
	LabelKey   *string `json:"labelKey,omitempty"`   /*  标签名  */
	LabelValue *string `json:"labelValue,omitempty"` /*  标签值  */
}
