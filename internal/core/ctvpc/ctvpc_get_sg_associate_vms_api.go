package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcGetSgAssociateVmsApi
/* 获取安全组绑定机器列表
 */type CtvpcGetSgAssociateVmsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcGetSgAssociateVmsApi(client *core.CtyunClient) *CtvpcGetSgAssociateVmsApi {
	return &CtvpcGetSgAssociateVmsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/get-sg-associate-vms",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcGetSgAssociateVmsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcGetSgAssociateVmsRequest) (*CtvpcGetSgAssociateVmsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("securityGroupID", req.SecurityGroupID)
	if req.ProjectID != nil {
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
	var resp CtvpcGetSgAssociateVmsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcGetSgAssociateVmsRequest struct {
	RegionID        string  /*  区域id  */
	SecurityGroupID string  /*  安全组ID  */
	ProjectID       *string /*  企业项目 ID，默认为'0'  */
	PageNo          int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize        int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10  */
}

type CtvpcGetSgAssociateVmsResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcGetSgAssociateVmsReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcGetSgAssociateVmsReturnObjResponse struct {
	Results      []*CtvpcGetSgAssociateVmsReturnObjResultsResponse `json:"results"`      /*  业务数据  */
	TotalCount   int32                                             `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                             `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                             `json:"totalPage"`    /*  总页数  */
}

type CtvpcGetSgAssociateVmsReturnObjResultsResponse struct {
	InstanceID    *string `json:"instanceID,omitempty"`    /*  主机 ID  */
	InstanceName  *string `json:"instanceName,omitempty"`  /*  主机名  */
	InstanceType  *string `json:"instanceType,omitempty"`  /*  主机类型：VM / BM  */
	InstanceState *string `json:"instanceState,omitempty"` /*  主机状态  */
	PrivateIp     *string `json:"privateIp,omitempty"`     /*  私有 ipv4  */
	PrivateIpv6   *string `json:"privateIpv6,omitempty"`   /*  私有 ipv6  */
}
