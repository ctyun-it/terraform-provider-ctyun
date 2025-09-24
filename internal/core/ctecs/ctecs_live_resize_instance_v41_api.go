package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsLiveResizeInstanceV41Api
/* 该接口提供云主机热变配功能，即开机状态实现变更规格<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;确认当前云主机是否可进行热变配，您可以通过接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=13078&data=87">查询云主机支持的热变配规格信息</a>获取当前云主机是否可以进行热变配，以及可以热变配规格信息<br />&emsp;&emsp;代金券：只支持预付费用户抵扣包周期资源的金额，且不可超过包周期资源的金额<br /><b>热变配当前支持规格和镜像信息为（不同资源池下的镜像和规格支持情况不同，以查询云主机支持热变配规格信息接口的返回值为准）：</b><br />支持的云主机镜像：<br />&emsp;&emsp;CentOS：CentOS 7.6 64位、CentOS 7.8 64位、CentOS 7.9 64位、CentOS 8.0 64位、CentOS 8.1 64位、CentOS 8.2 64位、CentOS 8.4 64位<br />&emsp;&emsp;CTyunOS：CTyunOS 2.0.1-21.06.4 64位、CTyunOS 3-23.01 64位<br />&emsp;&emsp;KylinOS：KylinOS V10 SP1 64位、KylinOS V10 SP2 64位<br />&emsp;&emsp;其他：openEuler 22.03 SP2 64位、UnionTechOS V20 1050u1e 64位<br />支持的云主机规格：<br />&emsp;&emsp;除二代机以外的规格且vcpu≥32
 */type CtecsLiveResizeInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsLiveResizeInstanceV41Api(client *core.CtyunClient) *CtecsLiveResizeInstanceV41Api {
	return &CtecsLiveResizeInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/live-resize",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsLiveResizeInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsLiveResizeInstanceV41Request) (*CtecsLiveResizeInstanceV41Response, error) {
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
	var resp CtecsLiveResizeInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsLiveResizeInstanceV41Request struct {
	RegionID        string  `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID      string  `json:"instanceID,omitempty"`  /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	FlavorID        string  `json:"flavorID,omitempty"`    /*  云主机规格ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10118193">规格说明</a>了解弹性云主机的选型基本信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=13078&data=87">查询云主机支持的热变配规格信息</a>  */
	ClientToken     string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个clientToken值，则代表为同一个请求。保留时间为24小时  */
	PayVoucherPrice float32 `json:"payVoucherPrice"`       /*  代金券，满足以下规则：<br />两位小数，不足两位自动补0，超过两位小数无效；<br />不可为负数；<br />注：字段为0时表示不使用代金券，默认不使用  */
}

type CtecsLiveResizeInstanceV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码部分   */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码，详见错误码部分    */
	Message     string                                       `json:"message,omitempty"`     /*  英文描述信息    */
	Description string                                       `json:"description,omitempty"` /*  中文描述信息    */
	ReturnObj   *CtecsLiveResizeInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsLiveResizeInstanceV41ReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  主订单ID。调用方在拿到masterOrderID之后，可以使用masterOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单号  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
}
