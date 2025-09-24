package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtecsListInstanceBackupPolicyBindInstancesApi
/* 通过云主机备份策略ID获取备份信息，主要获取绑定的云主机信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsListInstanceBackupPolicyBindInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceBackupPolicyBindInstancesApi(client *core.CtyunClient) *CtecsListInstanceBackupPolicyBindInstancesApi {
	return &CtecsListInstanceBackupPolicyBindInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup-policy/list-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceBackupPolicyBindInstancesApi) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceBackupPolicyBindInstancesRequest) (*CtecsListInstanceBackupPolicyBindInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("policyID", req.PolicyID)
	if req.InstanceName != "" {
		ctReq.AddParam("instanceName", req.InstanceName)
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
	var resp CtecsListInstanceBackupPolicyBindInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceBackupPolicyBindInstancesRequest struct {
	RegionID     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID     string /*  云主机备份策略ID，32字节<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6924&data=87&isNormal=1&vid=81">查询云主机备份策略列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6914&data=87&isNormal=1&vid=81">创建云主机备份策略</a>  */
	InstanceName string /*  云主机名称，模糊过滤  */
	PageNo       int32  /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListInstanceBackupPolicyBindInstancesResponse struct {
	StatusCode  int32                                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                                       `json:"message,omitempty"`     /*  错误信息的英文描述  */
	Description string                                                       `json:"description,omitempty"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *CtecsListInstanceBackupPolicyBindInstancesReturnObjResponse `json:"returnObj"`             /*  返回对象  */
}

type CtecsListInstanceBackupPolicyBindInstancesReturnObjResponse struct {
	CurrentCount     int32                                                                          `json:"currentCount,omitempty"` /*  当前页记录数目  */
	CurrentPage      int32                                                                          `json:"currentPage,omitempty"`  /*  当前页码  */
	TotalCount       int32                                                                          `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage        int32                                                                          `json:"totalPage,omitempty"`    /*  总页数  */
	InstancePolicies []*CtecsListInstanceBackupPolicyBindInstancesReturnObjInstancePoliciesResponse `json:"instancePolicies"`       /*  分页明细  */
}

type CtecsListInstanceBackupPolicyBindInstancesReturnObjInstancePoliciesResponse struct {
	Status          string   `json:"status,omitempty"`       /*  云主机状态  */
	AttachedVolumes []string `json:"attachedVolumes"`        /*  云主机所关联的云硬盘ID列表  */
	DisplayName     string   `json:"displayName,omitempty"`  /*  云主机显示名称  */
	InstanceID      string   `json:"instanceID,omitempty"`   /*  云主机ID  */
	RegionID        string   `json:"regionID,omitempty"`     /*  资源池ID  */
	InstanceName    string   `json:"instanceName,omitempty"` /*  云主机名称  */
	CreateTime      string   `json:"createTime,omitempty"`   /*  创建时间  */
	UpdateTime      string   `json:"updateTime,omitempty"`   /*  更新时间  */
}
