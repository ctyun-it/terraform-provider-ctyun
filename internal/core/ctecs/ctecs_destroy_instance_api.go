package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDestroyInstanceApi
/* 销毁一台包周期已退订云主机<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;1. 对于成功销毁，并重复使用clientToken再次请求的情况下，只保证返回第一次使用该clientToken时请求参数对应的主订单ID（masterOrderID）<br />&emsp;&emsp;2. 包周期已退订云主机已不再计费，但占用户资源相关配额，确认该云主机需要销毁的情况下执行该请求
 */type CtecsDestroyInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDestroyInstanceApi(client *core.CtyunClient) *CtecsDestroyInstanceApi {
	return &CtecsDestroyInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/destroy-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDestroyInstanceApi) Do(ctx context.Context, credential core.Credential, req *CtecsDestroyInstanceRequest) (*CtecsDestroyInstanceResponse, error) {
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
	var resp CtecsDestroyInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDestroyInstanceRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。保留时间为24小时  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID  string `json:"instanceID,omitempty"`  /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
}

type CtecsDestroyInstanceResponse struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                 `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                 `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                 `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDestroyInstanceReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDestroyInstanceReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  主订单ID。调用方在拿到masterOrderID之后，可以使用masterOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单号  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
}
