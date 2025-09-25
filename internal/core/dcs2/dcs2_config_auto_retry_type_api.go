package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ConfigAutoRetryTypeApi
/* 设置迁移任务自动重连。
 */type Dcs2ConfigAutoRetryTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ConfigAutoRetryTypeApi(client *core.CtyunClient) *Dcs2ConfigAutoRetryTypeApi {
	return &Dcs2ConfigAutoRetryTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/transfer/configAutoRetryType",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ConfigAutoRetryTypeApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ConfigAutoRetryTypeRequest) (*Dcs2ConfigAutoRetryTypeResponse, error) {
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
	var resp Dcs2ConfigAutoRetryTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ConfigAutoRetryTypeRequest struct {
	RegionId      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	TaskId        string `json:"taskId,omitempty"`        /*  任务Id<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span><a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15103&data=81&isNormal=1&vid=270">查询数据迁移任务列表</a>  */
	AutoRetryType int32  `json:"autoRetryType,omitempty"` /*  自动重试类型<li>1：失败时，自动重试<li>2：失败时，不进行自动重试  */
}

type Dcs2ConfigAutoRetryTypeResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ConfigAutoRetryTypeReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ConfigAutoRetryTypeReturnObjResponse struct{}
