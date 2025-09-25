package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListEndpointServiceApi
/* 查看终端节点服务列表
 */type CtvpcListEndpointServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListEndpointServiceApi(client *core.CtyunClient) *CtvpcListEndpointServiceApi {
	return &CtvpcListEndpointServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/list-endpoint-service",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListEndpointServiceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListEndpointServiceRequest) (*CtvpcListEndpointServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Id != nil {
		ctReq.AddParam("id", *req.Id)
	}
	if req.EndpointServiceName != nil {
		ctReq.AddParam("endpointServiceName", *req.EndpointServiceName)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListEndpointServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListEndpointServiceRequest struct {
	RegionID            string  /*  资源池ID  */
	Page                int32   /*  分页参数  */
	PageNo              int32   /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize            int32   /*  每页数据量大小  */
	Id                  *string /*  节点服务id  */
	EndpointServiceName *string /*  终端节点服务名称，该字段为精确匹配  */
	QueryContent        *string /*  支持对终端节点服务名称进行模糊匹配  */
}

type CtvpcListEndpointServiceResponse struct {
	StatusCode   int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListEndpointServiceReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	TotalCount   int32                                        `json:"totalCount"`            /*  总条数  */
	TotalPage    int32                                        `json:"totalPage"`             /*  总页数  */
	CurrentCount int32                                        `json:"currentCount"`          /*  总页数  */
	Error        *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListEndpointServiceReturnObjResponse struct {
	ID             *string                                              `json:"ID,omitempty"`          /*  终端节点服务ID  */
	Name           *string                                              `json:"name,omitempty"`        /*  终端节点服务名称  */
	VpcID          *string                                              `json:"vpcID,omitempty"`       /*  所属的专有网络id  */
	Description    *string                                              `json:"description,omitempty"` /*  描述  */
	RawType        *string                                              `json:"type,omitempty"`        /*  接口还是反向，interface:接口，reverse:反向  */
	AutoConnection *bool                                                `json:"autoConnection"`        /*  是否自动连接  */
	Rules          []*CtvpcListEndpointServiceReturnObjRulesResponse    `json:"rules"`                 /*  接口规则数据  */
	Backends       []*CtvpcListEndpointServiceReturnObjBackendsResponse `json:"backends"`              /*  后端数据  */
	CreatedAt      *string                                              `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt      *string                                              `json:"updatedAt,omitempty"`   /*  更新时间  */
	DnsName        *string                                              `json:"dnsName,omitempty"`     /*  dns 名  */
}

type CtvpcListEndpointServiceReturnObjRulesResponse struct {
	Protocol     *string `json:"protocol,omitempty"` /*  协议，TCP:TCP协议,UDP:UDP协议  */
	ServerPort   int32   `json:"serverPort"`         /*  服务端口(用于创建backend传入)  */
	EndpointPort int32   `json:"endpointPort"`       /*  节点端口(用于创建rule传入)  */
}

type CtvpcListEndpointServiceReturnObjBackendsResponse struct {
	InstanceType *string `json:"instanceType,omitempty"` /*  vm:虚机类型,bm:智能网卡类型,vip:vip类型,lb:负载均衡类型  */
	InstanceID   *string `json:"instanceID,omitempty"`   /*  实例id  */
	ProtocolPort int32   `json:"protocolPort"`           /*  端口  */
	Weight       int32   `json:"weight"`                 /*  权重  */
}
