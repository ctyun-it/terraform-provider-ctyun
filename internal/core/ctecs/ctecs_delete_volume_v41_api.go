package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDeleteVolumeV41Api
/* 您可以退订/删除一块包周期/按需的云硬盘，以释放存储空间资源。退订/删除云硬盘后，将停止对云硬盘收费<br />包周期云硬盘退订时，按照原始订单实付价格折算退订金额并进行返还<br />当云硬盘被退订/删除后，云硬盘的数据将无法被访问<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsDeleteVolumeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDeleteVolumeV41Api(client *core.CtyunClient) *CtecsDeleteVolumeV41Api {
	return &CtecsDeleteVolumeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/volume/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDeleteVolumeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDeleteVolumeV41Request) (*CtecsDeleteVolumeV41Response, error) {
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
	var resp CtecsDeleteVolumeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDeleteVolumeV41Request struct {
	ClientToken       string `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性。保留时间为24小时，使用同一个clientToken值，则代表为同一个请求  */
	DiskID            string `json:"diskID,omitempty"`            /*  磁盘ID。您可以查看<a href="https://www.ctyun.cn/document/10027696/10027930">产品定义-云硬盘</a>来了解云硬盘 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">云硬盘列表查询</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7332&data=48&isNormal=1&vid=45">创建云硬盘</a>  */
	RegionID          string `json:"regionID,omitempty"`          /*  如本地语境支持保存regionID，那么建议传递。资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DeleteSnapWithEbs string `json:"deleteSnapWithEbs,omitempty"` /*  设置快照是否随盘删除，只能设置为true  */
}

type CtecsDeleteVolumeV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中或失败)  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDeleteVolumeV41ReturnObjResponse   `json:"returnObj"`             /*  返回参数  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式码  */
	ErrorDetail *CtecsDeleteVolumeV41ErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode<br />其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
}

type CtecsDeleteVolumeV41ReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  退订订单号，可以使用该订单号确认资源的最终退订状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  退订订单号  */
}

type CtecsDeleteVolumeV41ErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMsg        string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息  */
}
