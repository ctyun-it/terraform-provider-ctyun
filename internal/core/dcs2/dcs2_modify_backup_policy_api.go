package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyBackupPolicyApi
/* 修改备份策略
 */type Dcs2ModifyBackupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyBackupPolicyApi(client *core.CtyunClient) *Dcs2ModifyBackupPolicyApi {
	return &Dcs2ModifyBackupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisMgr/modifyBackupPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyBackupPolicyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyBackupPolicyRequest) (*Dcs2ModifyBackupPolicyResponse, error) {
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
	var resp Dcs2ModifyBackupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyBackupPolicyRequest struct {
	RegionId              string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId            string `json:"prodInstId,omitempty"`            /*  实例ID  */
	PreferredBackupPeriod string `json:"preferredBackupPeriod,omitempty"` /*  日期范围，1-7表示周一至周日，多个日期使用英文逗号分隔。  */
	PreferredBackupTime   string `json:"preferredBackupTime,omitempty"`   /*  备份时间，0-23点准点  */
}

type Dcs2ModifyBackupPolicyResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	RequestId  string `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string `json:"code,omitempty"`       /*  响应码描述  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
