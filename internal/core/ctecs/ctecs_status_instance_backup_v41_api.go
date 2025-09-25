package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsStatusInstanceBackupV41Api
/* 查询云主机备份状态<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsStatusInstanceBackupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsStatusInstanceBackupV41Api(client *core.CtyunClient) *CtecsStatusInstanceBackupV41Api {
	return &CtecsStatusInstanceBackupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup/status",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsStatusInstanceBackupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsStatusInstanceBackupV41Request) (*CtecsStatusInstanceBackupV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceBackupID", req.InstanceBackupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsStatusInstanceBackupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsStatusInstanceBackupV41Request struct {
	RegionID         string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceBackupID string /*  云主机备份ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033738">产品定义-云主机备份</a>来了解云主机备份<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8354&data=87&isNormal=1&vid=81">查询云主机备份列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8332&data=87&isNormal=1&vid=81">创建云主机备份</a>  */
}

type CtecsStatusInstanceBackupV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsStatusInstanceBackupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsStatusInstanceBackupV41ReturnObjResponse struct {
	InstanceBackupStatus string `json:"instanceBackupStatus,omitempty"` /*  备份状态，取值范围：<br />CREATING: 备份创建中, <br />ACTIVE: 可用， <br />RESTORING: 备份恢复中，<br />DELETING: 删除中，<br />EXPIRED：到期，<br />ERROR：错误  */
}
