package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsLiteEcsCreateVolumeApi
/* 该接口提供用户轻量型云主机新建云硬盘的能力
 */ /* &emsp;&emsp;1. 云硬盘购买完成后自动挂载至轻量型云主机
 */ /* &emsp;&emsp;2. 数据盘到期日与所挂载的实例一致
 */ /* &emsp;&emsp;3. 实例和数据盘必须一起续费，无法单独为实例或数据盘续费
 */ /* &emsp;&emsp;4. 若数据盘所属轻量型云主机已开启自动续费，则数据盘默认开启自动续费
 */ /* &emsp;&emsp;5. 云硬盘计费方式：包年包月
 */ /*
 */ /* <b>准备工作：</b>
 */ /* &emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a>
 */ /* &emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a>
 */ /* &emsp;&emsp;计费模式：确认开通云主机的计费模式及计费项，详细查看<a href="https://www.ctyun.cn/document/10114925/10115664">计费方式及计费项</a>
 */ /* &emsp;&emsp;产品选型：购买轻量型云主机前，请先阅读<a href="https://www.ctyun.cn/document/10114925/10268652">实例套餐</a>了解轻量型云主机的规格套餐，并通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11998&data=87">查询轻量型云主机的规格套餐资源</a>接口，获取当前资源池可用轻量型云主机规格信息
 */ /* <b>注意事项：</b>
 */ /* &emsp;&emsp;成本估算：了解云主机的计费项，详细查看<a href="https://www.ctyun.cn/document/10114925/10115703">价格总览</a>进行成本估算
 */ /* &emsp;&emsp;用户配额：确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询
 */ /* &emsp;&emsp;异步接口：该接口为异步接口，下单成功不代表业务成功
 */type CtecsLiteEcsCreateVolumeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsLiteEcsCreateVolumeApi(client *core.CtyunClient) *CtecsLiteEcsCreateVolumeApi {
	return &CtecsLiteEcsCreateVolumeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/lite/create-volume",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsLiteEcsCreateVolumeApi) Do(ctx context.Context, credential core.Credential, req *CtecsLiteEcsCreateVolumeRequest) (*CtecsLiteEcsCreateVolumeResponse, error) {
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
	var resp CtecsLiteEcsCreateVolumeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsLiteEcsCreateVolumeRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性，保留时间为24小时  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID  string `json:"instanceID,omitempty"`  /*  轻量型云主机ID，您可以查看<a href="https://www.ctyun.cn/products/lite-ecs">轻量型云主机</a>了解轻量型云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11981&data=87">查询轻量型云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11980&data=87">创建轻量型云主机</a>  */
	DiskName    string `json:"diskName,omitempty"`    /*  磁盘名称，名称规则：长度2~63，不支持中文。注：当创建多块云硬盘时（即参数diskCount的值大于1时），第二块盘起名称会追加序号  */
	DiskType    string `json:"diskType,omitempty"`    /*  磁盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	DiskSize    int32  `json:"diskSize,omitempty"`    /*  磁盘容量，单位为GB，取值范围：[10, 32768]  */
	DiskCount   int32  `json:"diskCount,omitempty"`   /*  本地订购磁盘数量。注：不填写默认为1块盘，创建多块盘时填写多个；该参数值受单台轻量型云主机挂载硬盘上限值限制  */
}

type CtecsLiteEcsCreateVolumeResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                     `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                     `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                     `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsLiteEcsCreateVolumeReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsLiteEcsCreateVolumeReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单号  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源所属资源池ID  */
}
