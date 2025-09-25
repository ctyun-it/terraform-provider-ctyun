package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceConfigV2Api
/* 修改配置参数v2，支持的配置参数非必填，但是应至少填写一个。
 */type Dcs2ModifyInstanceConfigV2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceConfigV2Api(client *core.CtyunClient) *Dcs2ModifyInstanceConfigV2Api {
	return &Dcs2ModifyInstanceConfigV2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceParam/modifyInstanceConfig",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceConfigV2Api) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceConfigV2Request) (*Dcs2ModifyInstanceConfigV2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ModifyInstanceConfigV2Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceConfigV2Request struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID  */
	Appendfsync string `json:"appendfsync,omitempty"` /*  指定日志更新条件<li>no<li>everysec<li>always。<br>说明：2.8，4.0，5.0内核主备、集群主备版本禁止设置no  */
}

type Dcs2ModifyInstanceConfigV2Response struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                       `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifyInstanceConfigV2ReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                       `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                       `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyInstanceConfigV2ReturnObjResponse struct{}
