package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsPortsAttachInstanceV41Api
/* 网卡绑定云主机
 */type CtecsPortsAttachInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsPortsAttachInstanceV41Api(client *core.CtyunClient) *CtecsPortsAttachInstanceV41Api {
	return &CtecsPortsAttachInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/ports/attach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsPortsAttachInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsPortsAttachInstanceV41Request) (*CtecsPortsAttachInstanceV41Response, error) {
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
	var resp CtecsPortsAttachInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsPortsAttachInstanceV41Request struct {
	ClientToken        string `json:"clientToken,omitempty"`        /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个clientToken值，则代表为同一个请求。保留时间为24小时  */
	RegionID           string `json:"regionID,omitempty"`           /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName             string `json:"azName,omitempty"`             /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	ProjectID          string `json:"projectID,omitempty"`          /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	NetworkInterfaceID string `json:"networkInterfaceID,omitempty"` /*  弹性网卡ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10197673">弹性网卡-弹性网卡基本知识</a>来了解弹性网卡<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5802&data=94">查询弹性网卡列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5789&data=94">创建弹性网卡</a><br />注：多可用区时，为弹性网卡ID；非多可用区时，为子网ID。  */
	InstanceID         string `json:"instanceID,omitempty"`         /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	InstanceType       int32  `json:"instanceType,omitempty"`       /*  实例类型：3-虚拟机  */
}

type CtecsPortsAttachInstanceV41Response struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string `json:"message,omitempty"`     /*  英文描述信息  */
	Description string `json:"description,omitempty"` /*  中文描述信息  */
}
