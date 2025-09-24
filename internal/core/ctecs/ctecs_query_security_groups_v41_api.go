package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtecsQuerySecurityGroupsV41Api
/* 该接口提供用户查询用户安全组列表的功能<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */type CtecsQuerySecurityGroupsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQuerySecurityGroupsV41Api(client *core.CtyunClient) *CtecsQuerySecurityGroupsV41Api {
	return &CtecsQuerySecurityGroupsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/vpc/query-security-groups",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQuerySecurityGroupsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQuerySecurityGroupsV41Request) (*CtecsQuerySecurityGroupsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != "" {
		ctReq.AddParam("vpcID", req.VpcID)
	}
	if req.QueryContent != "" {
		ctReq.AddParam("queryContent", req.QueryContent)
	}
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.InstanceID != "" {
		ctReq.AddParam("instanceID", req.InstanceID)
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
	var resp CtecsQuerySecurityGroupsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQuerySecurityGroupsV41Request struct {
	RegionID     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	VpcID        string /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常以“vpc-”开头，非多可用区类型资源池vpcID为uuid格式  */
	QueryContent string /*  模糊匹配查询内容（匹配字段：id、name）  */
	ProjectID    string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	InstanceID   string /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	PageNumber   int32  /*  页码，取值范围：正整数（≥1），注：旧字段，后续可能废弃；默认值为1  */
	PageNo       int32  /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsQuerySecurityGroupsV41Response struct {
	TotalCount   int32                                           `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                           `json:"currentCount,omitempty"` /*  分页查询时每页的行数  */
	TotalPage    int32                                           `json:"totalPage,omitempty"`    /*  总页数  */
	StatusCode   int32                                           `json:"statusCode,omitempty"`   /*  返回状态码（800为成功，900为失败）  */
	ReturnObj    []*CtecsQuerySecurityGroupsV41ReturnObjResponse `json:"returnObj"`              /*  详细结果  */
	ErrorCode    string                                          `json:"errorCode,omitempty"`    /*  错误码，为product.module.code三段式码  */
	Error        string                                          `json:"error,omitempty"`        /*  错误码，为product.module.code三段式码  */
	Message      string                                          `json:"message,omitempty"`      /*  英文描述信息  */
	Description  string                                          `json:"description,omitempty"`  /*  中文描述信息  */
}

type CtecsQuerySecurityGroupsV41ReturnObjResponse struct {
	SecurityGroupName     string                                                               `json:"securityGroupName,omitempty"` /*  安全组名称  */
	Id                    string                                                               `json:"id,omitempty"`                /*  安全组id  */
	VmNum                 string                                                               `json:"vmNum,omitempty"`             /*  相关云主机  */
	Origin                string                                                               `json:"origin,omitempty"`            /*  表示是否是默认安全组  */
	VpcName               string                                                               `json:"vpcName,omitempty"`           /*  vpc名称  */
	VpcID                 string                                                               `json:"vpcID,omitempty"`             /*  安全组所属的专有网络  */
	CreationTime          string                                                               `json:"creationTime,omitempty"`      /*  创建时间  */
	Description           string                                                               `json:"description,omitempty"`       /*  安全组描述信息  */
	SecurityGroupRuleList []*CtecsQuerySecurityGroupsV41ReturnObjSecurityGroupRuleListResponse `json:"securityGroupRuleList"`       /*  安全组规则信息  */
}

type CtecsQuerySecurityGroupsV41ReturnObjSecurityGroupRuleListResponse struct {
	Direction             string `json:"direction,omitempty"`             /*  出方向-egress、入方向-ingress  */
	Action                string `json:"action,omitempty"`                /*  拒绝策略:允许-accept 拒绝-drop  */
	Origin                string `json:"origin,omitempty"`                /*  来源  */
	Priority              int32  `json:"priority,omitempty"`              /*  优先级:0~100  */
	Ethertype             string `json:"ethertype,omitempty"`             /*  IP类型:IPv4、IPv6  */
	Protocol              string `json:"protocol,omitempty"`              /*  协议: ANY、TCP、UDP、ICMP、ICMP6  */
	RawRange              string `json:"range,omitempty"`                 /*  接口范围/ICMP类型:1-65535  */
	DestCidrIp            string `json:"destCidrIp,omitempty"`            /*  远端地址:0.0.0.0/0  */
	Description           string `json:"description,omitempty"`           /*  安全组规则描述信息  */
	CreateTime            string `json:"createTime,omitempty"`            /*  创建时间，UTC时间  */
	Id                    string `json:"id,omitempty"`                    /*  唯一标识ID  */
	SecurityGroupID       string `json:"securityGroupID,omitempty"`       /*  安全组ID  */
	RemoteSecurityGroupID string `json:"remoteSecurityGroupID,omitempty"` /*  远端安全组id  */
	PrefixListID          string `json:"prefixListID,omitempty"`          /*  前缀列表id  */
}
