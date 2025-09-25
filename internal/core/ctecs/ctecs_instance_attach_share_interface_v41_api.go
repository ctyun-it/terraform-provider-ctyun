package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsInstanceAttachShareInterfaceV41Api
/* 给云主机添加共享网卡<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项：</b><br />&emsp;&emsp;该接口给云主机添加共享网卡后，共享网卡不能通过普通网卡方式进行卸载，只能随着退订云主机时进行删除，请谨慎使用
 */type CtecsInstanceAttachShareInterfaceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsInstanceAttachShareInterfaceV41Api(client *core.CtyunClient) *CtecsInstanceAttachShareInterfaceV41Api {
	return &CtecsInstanceAttachShareInterfaceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/share-interface/attach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsInstanceAttachShareInterfaceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsInstanceAttachShareInterfaceV41Request) (*CtecsInstanceAttachShareInterfaceV41Response, error) {
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
	var resp CtecsInstanceAttachShareInterfaceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsInstanceAttachShareInterfaceV41Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SubnetID   string `json:"subnetID,omitempty"`   /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10098380">基本概念</a>来查找子网的相关定义 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a><br />注：在多可用区类型资源池下，subnetID通常以“subnet-”开头，非多可用区类型资源池subnetID为uuid格式  */
}

type CtecsInstanceAttachShareInterfaceV41Response struct {
	StatusCode  int32                                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                                 `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见**错误码说明**  */
	Error       string                                                 `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码，详见**错误码说明**  */
	Message     string                                                 `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                 `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsInstanceAttachShareInterfaceV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsInstanceAttachShareInterfaceV41ReturnObjResponse struct {
	NicID string `json:"nicID,omitempty"` /*  网卡ID  */
}
