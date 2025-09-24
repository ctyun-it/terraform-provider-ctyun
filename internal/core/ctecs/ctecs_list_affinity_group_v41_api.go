package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListAffinityGroupV41Api
/* 查询云主机组列表或者详情<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />&emsp;&emsp;匹配查找：可以通过部分字段进行匹配筛选数据，无符合条件的为空
 */type CtecsListAffinityGroupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListAffinityGroupV41Api(client *core.CtyunClient) *CtecsListAffinityGroupV41Api {
	return &CtecsListAffinityGroupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/affinity-group/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListAffinityGroupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListAffinityGroupV41Request) (*CtecsListAffinityGroupV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListAffinityGroupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListAffinityGroupV41Request struct {
	RegionID        string `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AffinityGroupID string `json:"affinityGroupID,omitempty"` /*  云主机组ID，获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8324&data=87">查询云主机组列表或者详情</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8316&data=87">创建云主机组</a><br />   */
	QueryContent    string `json:"queryContent,omitempty"`    /*  模糊匹配查询内容（匹配字段：affinityGroupID、affinityGroupName）  */
	PageNo          int32  `json:"pageNo,omitempty"`          /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize        int32  `json:"pageSize,omitempty"`        /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListAffinityGroupV41Response struct {
	StatusCode  int32                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                      `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                      `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                      `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListAffinityGroupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListAffinityGroupV41ReturnObjResponse struct {
	CurrentCount int32                                                `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListAffinityGroupV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListAffinityGroupV41ReturnObjResultsResponse struct {
	AffinityGroupID     string                                                                `json:"affinityGroupID,omitempty"`   /*  云主机组ID  */
	AffinityGroupName   string                                                                `json:"affinityGroupName,omitempty"` /*  云主机组名称  */
	AffinityGroupPolicy *CtecsListAffinityGroupV41ReturnObjResultsAffinityGroupPolicyResponse `json:"affinityGroupPolicy"`         /*  云主机组策略  */
	CreatedTime         string                                                                `json:"createdTime,omitempty"`       /*  创建时间  */
	UpdatedTime         string                                                                `json:"updatedTime,omitempty"`       /*  更新时间  */
	Deleted             *bool                                                                 `json:"deleted"`                     /*  是否删除  */
}

type CtecsListAffinityGroupV41ReturnObjResultsAffinityGroupPolicyResponse struct {
	PolicyType     int32  `json:"policyType,omitempty"`     /*  云主机组策略类型<br />取值范围：<br />0：强制反亲和，<br />1：强制亲和，<br />2：软反亲和，<br />3：软亲和，<br />4：电力反亲和性  */
	PolicyTypeName string `json:"policyTypeName,omitempty"` /*  云主机组策略类型名称<br />取值范围：<br />anti-affinity：强制反亲和性，<br />affinity：强制亲和性，<br />soft-anti-affinity：反亲和性，<br />soft-affinity：亲和性，<br />power-anti-affinity：电力反亲和性  */
}
