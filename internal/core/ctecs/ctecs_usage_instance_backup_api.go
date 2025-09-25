package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsUsageInstanceBackupApi
/* 查看云主机备份空间占用大小<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsUsageInstanceBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsUsageInstanceBackupApi(client *core.CtyunClient) *CtecsUsageInstanceBackupApi {
	return &CtecsUsageInstanceBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-usage",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsUsageInstanceBackupApi) Do(ctx context.Context, credential core.Credential, req *CtecsUsageInstanceBackupRequest) (*CtecsUsageInstanceBackupResponse, error) {
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
	var resp CtecsUsageInstanceBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsUsageInstanceBackupRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceBackupID string `json:"instanceBackupID,omitempty"` /*  云主机备份ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033738">产品定义-云主机备份</a>来了解云主机备份<br /> 获取：<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8332&data=87&isNormal=1&vid=81">创建云主机备份</a>  */
}

type CtecsUsageInstanceBackupResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  默认值：800  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                     `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                     `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                     `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsUsageInstanceBackupReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsUsageInstanceBackupReturnObjResponse struct {
	Usage int32 `json:"usage,omitempty"` /*  备份占用空间大小，单位为B  */
}
