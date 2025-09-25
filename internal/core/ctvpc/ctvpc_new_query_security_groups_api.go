package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewQuerySecurityGroupsApi
/* 查询用户安全组列表。
 */type CtvpcNewQuerySecurityGroupsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewQuerySecurityGroupsApi(client *core.CtyunClient) *CtvpcNewQuerySecurityGroupsApi {
	return &CtvpcNewQuerySecurityGroupsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/new-query-security-groups",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewQuerySecurityGroupsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewQuerySecurityGroupsRequest) (*CtvpcNewQuerySecurityGroupsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.InstanceID != nil {
		ctReq.AddParam("instanceID", *req.InstanceID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
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
	var resp CtvpcNewQuerySecurityGroupsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewQuerySecurityGroupsRequest struct {
	RegionID     string  /*  区域id  */
	VpcID        *string /*  安全组所在的专有网络ID。  */
	QueryContent *string /*  【模糊查询】  安全组ID或名称  */
	InstanceID   *string /*  实例 ID  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcNewQuerySecurityGroupsResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewQuerySecurityGroupsReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcNewQuerySecurityGroupsReturnObjResponse struct {
	SecurityGroups []*CtvpcNewQuerySecurityGroupsReturnObjSecurityGroupsResponse `json:"securityGroups"` /*  安全组列表  */
	TotalCount     int32                                                         `json:"totalCount"`     /*  列表条目数  */
	CurrentCount   int32                                                         `json:"currentCount"`   /*  分页查询时每页的行数。  */
	TotalPage      int32                                                         `json:"totalPage"`      /*  总页数  */
}

type CtvpcNewQuerySecurityGroupsReturnObjSecurityGroupsResponse struct {
	SecurityGroupName     *string                                                                            `json:"securityGroupName,omitempty"` /*  安全组名称  */
	ProjectID             *string                                                                            `json:"projectID,omitempty"`         /*  项目ID  */
	Id                    *string                                                                            `json:"id,omitempty"`                /*  安全组id  */
	VmNum                 int32                                                                              `json:"vmNum"`                       /*  相关云主机  */
	Origin                *string                                                                            `json:"origin,omitempty"`            /*  表示是否是默认安全组  */
	VpcName               *string                                                                            `json:"vpcName,omitempty"`           /*  vpc名称  */
	VpcID                 *string                                                                            `json:"vpcID,omitempty"`             /*  安全组所属的专有网络。  */
	CreationTime          *string                                                                            `json:"creationTime,omitempty"`      /*  创建时间  */
	Description           *string                                                                            `json:"description,omitempty"`       /*  安全组描述信息。  */
	SecurityGroupRuleList []*CtvpcNewQuerySecurityGroupsReturnObjSecurityGroupsSecurityGroupRuleListResponse `json:"securityGroupRuleList"`       /*  安全组规则信息  */
}

type CtvpcNewQuerySecurityGroupsReturnObjSecurityGroupsSecurityGroupRuleListResponse struct {
	Direction             *string `json:"direction,omitempty"`             /*  出方向-egress、入方向-ingress  */
	Priority              int32   `json:"priority"`                        /*  优先级:0~100  */
	Ethertype             *string `json:"ethertype,omitempty"`             /*  IP类型:IPv4、IPv6  */
	Protocol              *string `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP、ICMP6  */
	RawRange              *string `json:"range,omitempty"`                 /*  接口范围/ICMP类型:1-65535  */
	DestCidrIp            *string `json:"destCidrIp,omitempty"`            /*  远端地址:0.0.0.0/0  */
	Description           *string `json:"description,omitempty"`           /*  安全组规则描述信息。  */
	CreateTime            *string `json:"createTime,omitempty"`            /*  创建时间，UTC时间。  */
	Id                    *string `json:"id,omitempty"`                    /*  唯一标识ID  */
	SecurityGroupID       *string `json:"securityGroupID,omitempty"`       /*  安全组ID  */
	Action                *string `json:"action,omitempty"`                /*  否  */
	Origin                *string `json:"origin,omitempty"`                /*  类型  */
	RemoteSecurityGroupID *string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组id  */
	PrefixListID          *string `json:"prefixListID,omitempty"`          /*  前缀列表id  */
}
