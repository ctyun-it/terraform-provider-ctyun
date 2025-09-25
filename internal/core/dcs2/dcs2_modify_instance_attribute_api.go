package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceAttributeApi
/* 修改实例部分信息，包含：实例描述信息、实例退订保护、实例维护时间。
 */type Dcs2ModifyInstanceAttributeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceAttributeApi(client *core.CtyunClient) *Dcs2ModifyInstanceAttributeApi {
	return &Dcs2ModifyInstanceAttributeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/modifyInstanceAttribute",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceAttributeApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceAttributeRequest) (*Dcs2ModifyInstanceAttributeResponse, error) {
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
	var resp Dcs2ModifyInstanceAttributeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceAttributeRequest struct {
	RegionId         string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId       string `json:"prodInstId,omitempty"`      /*  实例ID  */
	Description      string `json:"description,omitempty"`     /*  实例描述信息，小于1024个字符  */
	ProtectionStatus *bool  `json:"protectionStatus"`          /*  退订保护开关<li>true：开启<li>false：关闭<br>说明：传入空字符串或空格时自动转为 null  */
	MaintenanceTime  string `json:"maintenanceTime,omitempty"` /*  实例维护时间窗<li>格式：HH:mm-HH:mm<li>总时长必须为2小时<li>开始时间范围：00:00-22:00  */
}

type Dcs2ModifyInstanceAttributeResponse struct {
	StatusCode int32                                         `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                        `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifyInstanceAttributeReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                        `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                        `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                        `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyInstanceAttributeReturnObjResponse struct{}
