package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListEndpointServiceReverseRuleApi
/* 列表终端节点服务中转规则(反向访问规则)
 */type CtvpcListEndpointServiceReverseRuleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListEndpointServiceReverseRuleApi(client *core.CtyunClient) *CtvpcListEndpointServiceReverseRuleApi {
	return &CtvpcListEndpointServiceReverseRuleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/list-endpoint-service-reverse-rule",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListEndpointServiceReverseRuleApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListEndpointServiceReverseRuleRequest) (*CtvpcListEndpointServiceReverseRuleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("endpointServiceID", req.EndpointServiceID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
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
	var resp CtvpcListEndpointServiceReverseRuleResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListEndpointServiceReverseRuleRequest struct {
	RegionID          string /*  资源池ID  */
	EndpointServiceID string /*  终端节点服务id  */
	Page              int32  /*  分页参数  */
	PageNo            int32  /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize          int32  /*  每页数据量大小  */
}

type CtvpcListEndpointServiceReverseRuleResponse struct {
	StatusCode  int32                                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListEndpointServiceReverseRuleReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcListEndpointServiceReverseRuleReturnObjResponse struct {
	ReverseRules []*CtvpcListEndpointServiceReverseRuleReturnObjReverseRulesResponse `json:"reverseRules"` /*  终端节点列表数据  */
	TotalCount   int32                                                               `json:"totalCount"`   /*  总条数  */
	TotalPage    int32                                                               `json:"totalPage"`    /*  总页数  */
	CurrentCount int32                                                               `json:"currentCount"` /*  总页数  */
}

type CtvpcListEndpointServiceReverseRuleReturnObjReverseRulesResponse struct {
	ID                *string `json:"ID,omitempty"`                /*  节点服务中转规则ID  */
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  终端节点关联的终端节点服务  */
	EndpointID        *string `json:"endpointID,omitempty"`        /*  节点id  */
	TransitIPAddress  *string `json:"transitIPAddress,omitempty"`  /*  中转ip地址  */
	TransitPort       int32   `json:"transitPort"`                 /*  中转端口,1到65535  */
	Protocol          *string `json:"protocol,omitempty"`          /*  TCP:TCP协议,UDP:UDP协议  */
	TargetIPAddress   *string `json:"targetIPAddress,omitempty"`   /*  目标ip地址  */
	TargetPort        int32   `json:"targetPort"`                  /*  目标端口,1到65535  */
	CreatedAt         *string `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string `json:"updatedAt,omitempty"`         /*  更新时间  */
}
