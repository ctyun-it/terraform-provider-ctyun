package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryVmMemoryLatestMetricDataV41Api
/* 查询云主机的内存实时监控数据
 */type CtecsQueryVmMemoryLatestMetricDataV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryVmMemoryLatestMetricDataV41Api(client *core.CtyunClient) *CtecsQueryVmMemoryLatestMetricDataV41Api {
	return &CtecsQueryVmMemoryLatestMetricDataV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/vm-mem-latest-metric-data",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryVmMemoryLatestMetricDataV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryVmMemoryLatestMetricDataV41Request) (*CtecsQueryVmMemoryLatestMetricDataV41Response, error) {
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
	var resp CtecsQueryVmMemoryLatestMetricDataV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryVmMemoryLatestMetricDataV41Request struct {
	RegionID     string   `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DeviceIDList []string `json:"deviceIDList"`       /*  云主机ID列表，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	PageNo       int32    `json:"pageNo,omitempty"`   /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	Page         int32    `json:"page,omitempty"`     /*  页码，取值范围：正整数（≥1），注：默认值为1，建议使用pageNo，该参数后续会下线  */
	PageSize     int32    `json:"pageSize,omitempty"` /*   每页记录数目，取值范围：[1, 50]，注：默认值为10    */
}

type CtecsQueryVmMemoryLatestMetricDataV41Response struct {
	StatusCode  int32                                                   `json:"statusCode,omitempty"`  /*  返回码状态(800为成功，900为失败)，默认值：800  */
	ErrorCode   string                                                  `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Message     string                                                  `json:"message,omitempty"`     /*  英文描述信息  */
	MsgDesc     string                                                  `json:"msgDesc,omitempty"`     /*  中文描述信息  */
	Description string                                                  `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResponse `json:"returnObj"`             /*  返回参数，参考表returnObj  */
	Error       string                                                  `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
}

type CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResponse struct {
	Result       []*CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResultResponse `json:"result"`                 /*  result对象  */
	CurrentCount int32                                                           `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                           `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                           `json:"totalPage,omitempty"`    /*  总页数  */
	PageSize     int32                                                           `json:"pageSize,omitempty"`     /*  每页记录数目  */
	Page         int32                                                           `json:"page,omitempty"`         /*  页码  */
}

type CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResultResponse struct {
	FUID             string                                                                `json:"fUID,omitempty"`             /*  唯一键  */
	FuserLastUpdated string                                                                `json:"fuserLastUpdated,omitempty"` /*  用户最近更新时间  */
	RegionID         string                                                                `json:"regionID,omitempty"`         /*  资源池ID  */
	DeviceUUID       string                                                                `json:"deviceUUID,omitempty"`       /*  云主机ID  */
	ItemList         *CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResultItemListResponse `json:"itemList"`                   /*  监控项值列表  */
}

type CtecsQueryVmMemoryLatestMetricDataV41ReturnObjResultItemListResponse struct {
	SamplingTime        int32  `json:"samplingTime,omitempty"`        /*  监控获取时间  */
	Mem_util            string `json:"mem_util,omitempty"`            /*  内存使用率（%）  */
	Free_memory         string `json:"free_memory,omitempty"`         /*  可用内存（Byte）  */
	Used_memory         string `json:"used_memory,omitempty"`         /*  已用内存量（Byte）  */
	Buffer_memory       string `json:"buffer_memory,omitempty"`       /*  Buffers占用量（Byte）  */
	Cache_memory        string `json:"cache_memory,omitempty"`        /*  Cached占用量（Byte）  */
	Process_memory_used string `json:"process_memory_used,omitempty"` /*  进程内存使用率  */
	Pused_memory        string `json:"pused_memory,omitempty"`        /*  内存使用率（细粒度）（%）  */
}
