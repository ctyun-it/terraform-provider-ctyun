package crs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CrsCreateInstanceVpceLinkedVpcsV2Api
/* VPC接入 */
type CrsCreateInstanceVpceLinkedVpcsV2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCrsCreateInstanceVpceLinkedVpcsV2Api(client *core.CtyunClient) *CrsCreateInstanceVpceLinkedVpcsV2Api {
	return &CrsCreateInstanceVpceLinkedVpcsV2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/createInstanceVpceLinkedVpcs",
			ContentType:  "application/json",
		},
	}
}

func (a *CrsCreateInstanceVpceLinkedVpcsV2Api) Do(ctx context.Context, credential core.Credential, req *CrsCreateInstanceVpceLinkedVpcsV2Request) (*CrsCreateInstanceVpceLinkedVpcsV2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("Content-Type", req.ContentType)
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CrsCreateInstanceVpceLinkedVpcsV2Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CrsCreateInstanceVpceLinkedVpcsV2Request struct {
	ContentType string                                             `json:"Content-Type"`      /*  类型  */
	RegionId    string                                             `json:"regionId"`          /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	VpcList     []*CrsCreateInstanceVpceLinkedVpcsV2VpcListRequest `json:"vpcList,omitempty"` /*  接入VPC列表  */
}

type CrsCreateInstanceVpceLinkedVpcsV2VpcListRequest struct {
	VpcId    string  `json:"vpcId"`              /*  需要接入的VPC id  */
	SubnetId *string `json:"subnetId,omitempty"` /*  指定VPCE创建的子网id（如果不传该参数，则随机选择vpc下任一子网创建VPCE）  */
}

type CrsCreateInstanceVpceLinkedVpcsV2Response struct {
	StatusCode int32                                                 `json:"statusCode"` /*  响应码 （800为请求成功，900为失败 ）  */
	Message    *string                                               `json:"message"`    /*  返回信息  */
	Error      *string                                               `json:"error"`      /*  错误码  */
	ReturnObj  []*CrsCreateInstanceVpceLinkedVpcsV2ReturnObjResponse `json:"returnObj"`  /*  VPC接入结果列表  */
}

type CrsCreateInstanceVpceLinkedVpcsV2ReturnObjResponse struct {
	VpcLinkResult []*CrsCreateInstanceVpceLinkedVpcsV2ReturnObjVpcLinkResultResponse `json:"vpcLinkResult"` /*  vpc接入结果  */
}

type CrsCreateInstanceVpceLinkedVpcsV2ReturnObjVpcLinkResultResponse struct {
	Vpc       *string `json:"vpc"`       /*  VPC Id  */
	State     *string `json:"state"`     /*  VPC接入状态（ACTIVE-接入，INACTIVE-未接入，ERROR-错误）  */
	VpceState *string `json:"vpceState"` /*  VPCE状态（ACTIVE-接入，INACTIVE-未接入，CREATING-创建中，ERROR-错误）  */
	DnsState  *string `json:"dnsState"`  /*  DNS状态（ACTIVE-接入，INACTIVE-未接入，CREATING-创建中，ERROR-错误）  */
	Msg       *string `json:"msg"`       /*  错误的信息  */
}
