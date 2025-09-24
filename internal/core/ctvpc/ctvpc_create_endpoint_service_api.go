package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreateEndpointServiceApi
/* 创建终端节点服务
 */type CtvpcCreateEndpointServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreateEndpointServiceApi(client *core.CtyunClient) *CtvpcCreateEndpointServiceApi {
	return &CtvpcCreateEndpointServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/create-endpoint-service",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreateEndpointServiceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreateEndpointServiceRequest) (*CtvpcCreateEndpointServiceResponse, error) {
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
	var resp CtvpcCreateEndpointServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreateEndpointServiceRequest struct {
	ClientToken       string                                    `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string                                    `json:"regionID,omitempty"`     /*  资源池ID  */
	VpcID             string                                    `json:"vpcID,omitempty"`        /*  所属的专有网络id  */
	IpVersion         int32                                     `json:"ipVersion"`              /*  0:ipv4, 1:ipv6（暂不支持）, 2:双栈，默认0  */
	RawType           *string                                   `json:"type,omitempty"`         /*  接口、反向和网关负载均衡，interface:接口，reverse:反向，gwlb:网关负载均衡,默认为 interface  */
	Name              string                                    `json:"name,omitempty"`         /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	InstanceType      *string                                   `json:"instanceType,omitempty"` /*  服务后端实例类型，vm:虚机类型,bm:物理机,vip:vip类型,lb:负载均衡类型,underlay:天翼云内网资源,gwlb:网关负载均衡 当 type 为 interface 时，必填  */
	InstanceID        *string                                   `json:"instanceID,omitempty"`   /*  服务后端实例id, 当 type 为 interface 时，必填  */
	InstanceID6       *string                                   `json:"instanceID6,omitempty"`  /*  后端服务为havip时，v6的havip id  */
	UnderlayIP        *string                                   `json:"underlayIP,omitempty"`   /*  instance_type为天翼云内网资源时，v4的underlay ip  */
	UnderlayIP6       *string                                   `json:"underlayIP6,omitempty"`  /*  instance_type为天翼云内网资源时，v6的underlay ip  */
	SubnetID          *string                                   `json:"subnetID,omitempty"`     /*  服务后端子网id，当type是reverse，且reverseIsUnderlay为false时，必填  */
	AutoConnection    bool                                      `json:"autoConnection"`         /*  是否自动连接，true 表示自动链接，false 表示非自动链接  */
	Rules             []*CtvpcCreateEndpointServiceRulesRequest `json:"rules"`                  /*  节点服务规则, 当 type 为 interface 时，必填  */
	OaType            *string                                   `json:"oaType,omitempty"`       /*  oa 类型，支持: tcp_option / proxy_protocol / close  */
	ServiceCharge     *bool                                     `json:"serviceCharge"`          /*  是否开启服务计费，一旦开启服务计费，不可修改  */
	ForceEnableDns    *bool                                     `json:"forceEnableDns"`         /*  是否强制开启dns  */
	DnsName           *string                                   `json:"dnsName,omitempty"`      /*  dns名称  */
	ReverseIsUnderlay *bool                                     `json:"reverseIsUnderlay"`      /*  反向终端节点服务是否是underlay类型  */
	TransitIP         *string                                   `json:"transitIP,omitempty"`    /*  中转ipv4 ip，当reverseIsUnderlay是True时，必填  */
	TransitIP6        *string                                   `json:"transitIP6,omitempty"`   /*  中转ipv6 ip，当reverseIsUnderlay是True且服务是双栈时，必填  */
}

type CtvpcCreateEndpointServiceRulesRequest struct {
	Protocol     string `json:"protocol,omitempty"` /*  协议，TCP:TCP协议,UDP:UDP协议  */
	ServerPort   int32  `json:"serverPort"`         /*  服务端口(用于创建backend传入)(1-65535)  */
	EndpointPort int32  `json:"endpointPort"`       /*  节点端口(用于创建rule传入)(1-65535)  */
}

type CtvpcCreateEndpointServiceResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreateEndpointServiceReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcCreateEndpointServiceReturnObjResponse struct {
	EndpointService *CtvpcCreateEndpointServiceReturnObjEndpointServiceResponse `json:"endpointService"` /*  创建的终端节点服务信息  */
}

type CtvpcCreateEndpointServiceReturnObjEndpointServiceResponse struct {
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  创建的终端节点 ID  */
}
