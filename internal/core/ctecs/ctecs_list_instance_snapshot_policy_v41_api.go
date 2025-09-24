package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListInstanceSnapshotPolicyV41Api
/* 该接口提供用户查询快照策略绑定云主机列表的功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */type CtecsListInstanceSnapshotPolicyV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceSnapshotPolicyV41Api(client *core.CtyunClient) *CtecsListInstanceSnapshotPolicyV41Api {
	return &CtecsListInstanceSnapshotPolicyV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot-policy/instance-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceSnapshotPolicyV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceSnapshotPolicyV41Request) (*CtecsListInstanceSnapshotPolicyV41Response, error) {
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
	var resp CtecsListInstanceSnapshotPolicyV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceSnapshotPolicyV41Request struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty"` /*  云主机快照策略ID，32字节<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9600&data=87">查询云主机快照策略列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9588&data=87">创建云主机快照策略</a>  */
	PageNo           int32  `json:"pageNo,omitempty"`           /*  页码，取值范围：正整数（≥1）注：默认值为1  */
	PageSize         int32  `json:"pageSize,omitempty"`         /*  每页记录数目，取值范围：[1, 50]注：默认值为10  */
}

type CtecsListInstanceSnapshotPolicyV41Response struct {
	StatusCode  int32                                                `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                               `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                               `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListInstanceSnapshotPolicyV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListInstanceSnapshotPolicyV41ReturnObjResponse struct {
	CurrentCount int32                                                              `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                              `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                              `json:"totalPage,omitempty"`    /*  总页数  */
	InstanceList []*CtecsListInstanceSnapshotPolicyV41ReturnObjInstanceListResponse `json:"instanceList"`           /*  分页明细  */
}

type CtecsListInstanceSnapshotPolicyV41ReturnObjInstanceListResponse struct {
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	InstanceName   string `json:"instanceName,omitempty"`   /*  云主机名称，不同操作系统下，云主机名称规则有差异<br />Windows：长度为2~15个字符，允许使用大小写字母、数字或连字符（-）。不能以连字符（-）开头或结尾，不能连续使用连字符（-），也不能仅使用数字；<br />其他操作系统：长度为2-64字符，允许使用点（.）分隔字符成多段，每段允许使用大小写字母、数字或连字符（-），但不能连续使用点号（.）或连字符（-），不能以点号（.）或连字符（-）开头或结尾  */
	DisplayName    string `json:"displayName,omitempty"`    /*  云主机显示名称，长度为2-63字符  */
	InstanceStatus string `json:"instanceStatus,omitempty"` /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	VolumeCount    int32  `json:"volumeCount,omitempty"`    /*  磁盘数量  */
}
