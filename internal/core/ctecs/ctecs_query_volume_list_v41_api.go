package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryVolumeListV41Api
/* 该接口提供查询云主机的云硬盘列表功能<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */type CtecsQueryVolumeListV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryVolumeListV41Api(client *core.CtyunClient) *CtecsQueryVolumeListV41Api {
	return &CtecsQueryVolumeListV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/volume/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryVolumeListV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryVolumeListV41Request) (*CtecsQueryVolumeListV41Response, error) {
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
	var resp CtecsQueryVolumeListV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryVolumeListV41Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	PageNo     int32  `json:"pageNo,omitempty"`     /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsQueryVolumeListV41Response struct {
	StatusCode  int32                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                    `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                    `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryVolumeListV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsQueryVolumeListV41ReturnObjResponse struct {
	CurrentCount int32                                              `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                              `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                              `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsQueryVolumeListV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsQueryVolumeListV41ReturnObjResultsResponse struct {
	DiskMode     string `json:"diskMode,omitempty"`     /*  云硬盘属性，取值范围：<br />FCSAN（光纤通道协议的SAN网络），<br />ISCSI（小型计算机系统接口），<br />VBD（虚拟块存储设备）  */
	DiskID       string `json:"diskID,omitempty"`       /*  云硬盘ID  */
	DiskType     string `json:"diskType,omitempty"`     /*  用途分类，取值范围：<br />数据盘，<br />系统盘  */
	IsEncrypt    *bool  `json:"isEncrypt"`              /*  云硬盘加密标志，取值范围：<br />true：加密，<br />false：未加密  */
	DiskDataType string `json:"diskDataType,omitempty"` /*  云硬盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD  */
	DiskSize     int32  `json:"diskSize,omitempty"`     /*  云硬盘容量大小，单位GB  */
}
