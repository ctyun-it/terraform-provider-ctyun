package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListUpdateSpecFlavorV41Api
/* 该接口提供查询当前云主机支持冷变配的规格信息，用户可以根据此接口的返回值了解自己可使用的云主机规格有哪些<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;确认云主机是否存在于当前资源池<br />&emsp;&emsp;该接口查询的规格适用于冷变配操作，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8293&data=87&isNormal=1">云主机修改规格</a>来变更
 */type CtecsListUpdateSpecFlavorV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListUpdateSpecFlavorV41Api(client *core.CtyunClient) *CtecsListUpdateSpecFlavorV41Api {
	return &CtecsListUpdateSpecFlavorV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/flavor/update-spec-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListUpdateSpecFlavorV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListUpdateSpecFlavorV41Request) (*CtecsListUpdateSpecFlavorV41Response, error) {
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
	var resp CtecsListUpdateSpecFlavorV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListUpdateSpecFlavorV41Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsListUpdateSpecFlavorV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                         `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                         `json:"description,omitempty"` /*  中文描述信息    */
	ReturnObj   *CtecsListUpdateSpecFlavorV41ReturnObjResponse `json:"returnObj"`             /*  返回内容  */
}

type CtecsListUpdateSpecFlavorV41ReturnObjResponse struct {
	FlavorList []*CtecsListUpdateSpecFlavorV41ReturnObjFlavorListResponse `json:"flavorList"` /*  规格列表  */
}

type CtecsListUpdateSpecFlavorV41ReturnObjFlavorListResponse struct {
	CpuInfo               string   `json:"cpuInfo,omitempty"`               /*  cpu架构  */
	BaseBandwidth         float32  `json:"baseBandwidth"`                   /*  基准带宽  */
	FlavorName            string   `json:"flavorName,omitempty"`            /*  云主机规格名称  */
	FlavorType            string   `json:"flavorType,omitempty"`            /*  规格类型，取值范围：[CPU、CPU_C3、CPU_C6、CPU_C7、CPU_c7ne、CPU_C8、CPU_D3、CPU_FC1、CPU_FM1、CPU_FS1、CPU_HC1、CPU_HM1、CPU_HS1、CPU_IP3、CPU_IR3、CPU_IP3_2、CPU_IR3_2、CPU_KC1、CPU_KM1、CPU_KS1、CPU_M2、CPU_M3、CPU_M6、CPU_M7、CPU_M8、CPU_S2、CPU_S3、CPU_S6、CPU_S7、CPU_S8、CPU_s8r、GPU_N_V100_V_FMGQ、GPU_N_V100_V、GPU_N_V100S_V、GPU_N_V100S_V_FMGQ、GPU_N_T4_V、GPU_N_G7_V、GPU_N_V100、GPU_N_V100_SHIPINYUN、GPU_N_V100_SUANFA、GPU_N_P2V_RENMIN、GPU_N_V100S、GPU_N_T4、GPU_N_T4_AIJISUAN、GPU_N_T4_ASR、GPU_N_T4_JX、GPU_N_T4_SHIPINYUN、GPU_N_T4_SUANFA、GPU_N_T4_YUNYOUXI、GPU_N_PI7、GPU_N_P8A、GPU_A_PAK1、GPU_C_PCH1]，支持类型会随着功能升级增加  */
	FlavorSeries          string   `json:"flavorSeries,omitempty"`          /*  云主机规格系列，详见枚举值表  */
	NicMultiQueue         int32    `json:"nicMultiQueue,omitempty"`         /*  网卡多队列数目  */
	Pps                   int32    `json:"pps,omitempty"`                   /*  最大收发包限制  */
	FlavorCPU             int32    `json:"flavorCPU,omitempty"`             /*  VCPU个数  */
	FlavorRAM             int32    `json:"flavorRAM,omitempty"`             /*  内存  */
	Bandwidth             float32  `json:"bandwidth"`                       /*  带宽  */
	FlavorID              string   `json:"flavorID,omitempty"`              /*  云主机规格ID  */
	GpuVendor             string   `json:"gpuVendor,omitempty"`             /*  GPU厂商  */
	VideoMemSize          int32    `json:"videoMemSize,omitempty"`          /*  GPU显存大小  */
	GpuType               string   `json:"gpuType,omitempty"`               /*  GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount              int32    `json:"gpuCount,omitempty"`              /*  GPU设备数量  */
	Available             *bool    `json:"available"`                       /*  是否可用<br />true：可用，<br />false：不可用，已售罄  */
	AzList                []string `json:"azList"`                          /*  多az名称列表（非多可用区为\["default"\]）  */
	FlavorSeriesName      string   `json:"flavorSeriesName,omitempty"`      /*  规格系列名称，参照参数flavorSeries说明  */
	DedicatedAreaZone     string   `json:"dedicatedAreaZone,omitempty"`     /*  可用分区ID  */
	DedicatedAreaZoneName string   `json:"dedicatedAreaZoneName,omitempty"` /*  可用分区名称  */
}
