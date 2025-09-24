package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsInstanceBackupPolicyBindInstancesApi
/* 该接口备份策略绑定云主机<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsInstanceBackupPolicyBindInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsInstanceBackupPolicyBindInstancesApi(client *core.CtyunClient) *CtecsInstanceBackupPolicyBindInstancesApi {
	return &CtecsInstanceBackupPolicyBindInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-policy/bind-instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsInstanceBackupPolicyBindInstancesApi) Do(ctx context.Context, credential core.Credential, req *CtecsInstanceBackupPolicyBindInstancesRequest) (*CtecsInstanceBackupPolicyBindInstancesResponse, error) {
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
	var resp CtecsInstanceBackupPolicyBindInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsInstanceBackupPolicyBindInstancesRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID       string `json:"policyID,omitempty"`       /*  云主机备份策略ID，32字节<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=6924&data=100">查询云主机备份策略列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=6914&data=100">创建云主机备份策略</a>  */
	InstanceIDList string `json:"instanceIDList,omitempty"` /*  云主机ID列表，多台使用英文逗号分割，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
}

type CtecsInstanceBackupPolicyBindInstancesResponse struct {
	StatusCode  int32                                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)。  */
	ErrorCode   string                                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                                   `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                   `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsInstanceBackupPolicyBindInstancesReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsInstanceBackupPolicyBindInstancesReturnObjResponse struct {
	PolicyID       string `json:"policyID,omitempty"`       /*  云主机备份策略ID  */
	InstanceIDList string `json:"instanceIDList,omitempty"` /*  一台或多台云主机ID，多台则使用英文逗号分割  */
}
