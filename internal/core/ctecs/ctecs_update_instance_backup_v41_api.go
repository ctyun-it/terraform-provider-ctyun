package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsUpdateInstanceBackupV41Api
/* 更改云主机备份名称和描述<br/><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsUpdateInstanceBackupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsUpdateInstanceBackupV41Api(client *core.CtyunClient) *CtecsUpdateInstanceBackupV41Api {
	return &CtecsUpdateInstanceBackupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsUpdateInstanceBackupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsUpdateInstanceBackupV41Request) (*CtecsUpdateInstanceBackupV41Response, error) {
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
	var resp CtecsUpdateInstanceBackupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsUpdateInstanceBackupV41Request struct {
	RegionID                  string `json:"regionID,omitempty"`                  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceBackupID          string `json:"instanceBackupID,omitempty"`          /*  云主机备份ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033738">产品定义-云主机备份</a>来了解云主机备份<br /> 获取：<span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8332&data=87&isNormal=1&vid=81">创建云主机备份</a>   */
	InstanceBackupName        string `json:"instanceBackupName,omitempty"`        /*  云主机备份名称。满足以下规则：长度为2-63字符，头尾不支持输入空格  */
	InstanceBackupDescription string `json:"instanceBackupDescription,omitempty"` /*  云主机备份描述，字符长度不超过256字符  */
}

type CtecsUpdateInstanceBackupV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsUpdateInstanceBackupV41ReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtecsUpdateInstanceBackupV41ReturnObjResponse struct {
	InstanceBackupID string `json:"instanceBackupID,omitempty"` /*  云主机备份ID  */
}
