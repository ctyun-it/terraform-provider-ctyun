package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListPolicyRulesApi
/* 调用该接口查看策略治理规则库列表
 */type CcseListPolicyRulesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListPolicyRulesApi(client *core.CtyunClient) *CcseListPolicyRulesApi {
	return &CcseListPolicyRulesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/policies",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListPolicyRulesApi) Do(ctx context.Context, credential core.Credential, req *CcseListPolicyRulesRequest) (*CcseListPolicyRulesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseListPolicyRulesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListPolicyRulesRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseListPolicyRulesResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseListPolicyRulesReturnObjResponse `json:"returnObj"`            /*  Map<String,Array of Strings>类型，返回对象  */
	Error      string                                `json:"error,omitempty"`      /*  错误码  */
}

type CcseListPolicyRulesReturnObjResponse struct{}
