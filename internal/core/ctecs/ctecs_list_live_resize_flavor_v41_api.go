package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListLiveResizeFlavorV41Api
/* 该接口提供查询当前云主机支持热变配的规格信息，用户可以根据此接口的返回值了解自己可使用的云主机规格有哪些<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;确认云主机是否存在于当前资源池<br />&emsp;&emsp;如无法进行热变配，则返回规格结果为空
 */type CtecsListLiveResizeFlavorV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListLiveResizeFlavorV41Api(client *core.CtyunClient) *CtecsListLiveResizeFlavorV41Api {
	return &CtecsListLiveResizeFlavorV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/flavor/live-resize-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListLiveResizeFlavorV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListLiveResizeFlavorV41Request) (*CtecsListLiveResizeFlavorV41Response, error) {
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
	var resp CtecsListLiveResizeFlavorV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListLiveResizeFlavorV41Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsListLiveResizeFlavorV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                                         `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                         `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListLiveResizeFlavorV41ReturnObjResponse `json:"returnObj"`             /*  返回内容  */
}

type CtecsListLiveResizeFlavorV41ReturnObjResponse struct {
	FlavorList []*CtecsListLiveResizeFlavorV41ReturnObjFlavorListResponse `json:"flavorList"` /*  规格列表  */
}

type CtecsListLiveResizeFlavorV41ReturnObjFlavorListResponse struct {
	CpuInfo       string  `json:"cpuInfo,omitempty"`       /*  cpu架构  */
	BaseBandwidth float32 `json:"baseBandwidth"`           /*  基准带宽  */
	FlavorName    string  `json:"flavorName,omitempty"`    /*  云主机规格名称  */
	FlavorType    string  `json:"flavorType,omitempty"`    /*  规格类型，取值范围：[CPU、CPU_S6、CPU_C6、CPU_M6、CPU_S3、CPU_C3、CPU_M3、CPU_IP3、GPU_N_T4_V、GPU_N_V100、GPU_N_V100_V、GPU_N_P2V_RENMIN、GPU_N_PI7、GPU_N_G7_V、GPU_N_V100、GPU_N_T4_JX]，支持类型会随着功能升级增加  */
	FlavorSeries  string  `json:"flavorSeries,omitempty"`  /*  云主机规格系列，详见枚举值表  */
	NicMultiQueue int32   `json:"nicMultiQueue,omitempty"` /*  网卡多队列数目  */
	Pps           int32   `json:"pps,omitempty"`           /*  最大收发包限制  */
	FlavorCPU     int32   `json:"flavorCPU,omitempty"`     /*  VCPU个数  */
	FlavorRAM     int32   `json:"flavorRAM,omitempty"`     /*  内存  */
	Bandwidth     float32 `json:"bandwidth"`               /*  带宽  */
	FlavorID      string  `json:"flavorID,omitempty"`      /*  云主机规格ID  */
}
