package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListInstanceStatusV41Api
/* 该接口提供用户多台云主机状态信息查询功能，用户可以根据此接口的返回值得到多台云主机的状态信息<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */ /* &emsp;&emsp;匹配查找：可以通过部分字段进行匹配筛选数据，无符合条件的为空，在指定多台云主机ID的情况下，只返回匹配到的云主机信息。推荐每次使用单个条件查找
 */type CtecsListInstanceStatusV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceStatusV41Api(client *core.CtyunClient) *CtecsListInstanceStatusV41Api {
	return &CtecsListInstanceStatusV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/instance-status-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceStatusV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceStatusV41Request) (*CtecsListInstanceStatusV41Response, error) {
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
	var resp CtecsListInstanceStatusV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceStatusV41Request struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName         string `json:"azName,omitempty"`         /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	InstanceIDList string `json:"instanceIDList,omitempty"` /*  云主机ID列表，多台使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	PageNo         int32  `json:"pageNo,omitempty"`         /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	ProjectID      string `json:"projectID,omitempty"`      /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
}

type CtecsListInstanceStatusV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                       `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                       `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListInstanceStatusV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsListInstanceStatusV41ReturnObjResponse struct {
	CurrentCount int32                                                    `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                    `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                    `json:"totalPage,omitempty"`    /*  总页数  */
	StatusList   []*CtecsListInstanceStatusV41ReturnObjStatusListResponse `json:"statusList"`             /*  分页明细  */
}

type CtecsListInstanceStatusV41ReturnObjStatusListResponse struct {
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID  */
	InstanceStatus string `json:"instanceStatus,omitempty"` /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
}
