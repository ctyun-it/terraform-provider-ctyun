package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsExtendVolumeV41Api
/* 当云硬盘空间不足时，您可以扩大云硬盘的容量，也就是云硬盘扩容。此接口也支持随云主机订购的系统盘及数据盘的扩容<br />当您扩容成功后，需要将扩容部分的容量划分至原有分区内，或者对扩容部分的云硬盘分配新的分区。详见：<a herf="https://www.ctyun.cn/document/10027696/10029076"></a><br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />&emsp;&emsp;计费模式：确认开通云硬盘的计费模式，详细查看<a href="https://www.ctyun.cn/document/10027696/10028345">计费模式</a><br /><b>注意事项：</b><br />&emsp;&emsp;成本估算：了解云硬盘的<a href="https://www.ctyun.cn/document/10027696/10059594">计费说明</a><br />
 */type CtecsExtendVolumeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsExtendVolumeV41Api(client *core.CtyunClient) *CtecsExtendVolumeV41Api {
	return &CtecsExtendVolumeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/volume/extend",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsExtendVolumeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsExtendVolumeV41Request) (*CtecsExtendVolumeV41Response, error) {
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
	var resp CtecsExtendVolumeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsExtendVolumeV41Request struct {
	DiskSize    int32  `json:"diskSize,omitempty"`    /*  变配后的云硬盘大小，数据盘的取值范围为：<br />●超高IO/高IO/极速型SSD/普通IO：10GB~32768GB<br />●XSSD-0：10GB-65536GB<br />●XSSD-1：20GB-65536GB<br />●XSSD-2：512GB-65536GB<br />系统盘的取值范围为：<br />●超高IO/高IO/极速型SSD/普通IO：40GB~2048GB<br />●XSSD-0：40GB-2048GB<br />●XSSD-1：40GB-2048GB<br />●XSSD-2：512GB-2048GB  */
	DiskID      string `json:"diskID,omitempty"`      /*  磁盘ID。您可以查看<a href="https://www.ctyun.cn/document/10027696/10027930">产品定义-云硬盘</a>来了解云硬盘 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">云硬盘列表查询</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7332&data=48&isNormal=1&vid=45">创建云硬盘</a>  */
	RegionID    string `json:"regionID,omitempty"`    /*  如本地语境支持保存regionID，那么建议传递。资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，则代表为同一个请求。保留时间为24小时  */
}

type CtecsExtendVolumeV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中或失败)  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ErrorDetail *CtecsExtendVolumeV41ErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。<br />其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
	ReturnObj   *CtecsExtendVolumeV41ReturnObjResponse   `json:"returnObj"`             /*  返回参数  */
}

type CtecsExtendVolumeV41ErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMsg        string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息  */
}

type CtecsExtendVolumeV41ReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用masterOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO        string `json:"masterOrderNO,omitempty"`        /*  订单号  */
	MasterResourceID     string `json:"masterResourceID,omitempty"`     /*  主资源ID。云硬盘场景下，无需关心  */
	MasterResourceStatus string `json:"masterResourceStatus,omitempty"` /*  主资源状态。只有主订单资源会返回  */
}
