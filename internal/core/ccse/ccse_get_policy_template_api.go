package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetPolicyTemplateApi
/* 调用该接口查看策略治理规则模板详情
 */type CcseGetPolicyTemplateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetPolicyTemplateApi(client *core.CtyunClient) *CcseGetPolicyTemplateApi {
	return &CcseGetPolicyTemplateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/policies/{policyName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetPolicyTemplateApi) Do(ctx context.Context, credential core.Credential, req *CcseGetPolicyTemplateRequest) (*CcseGetPolicyTemplateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("policyName", req.PolicyName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetPolicyTemplateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetPolicyTemplateRequest struct {
	PolicyName string /*  策略治理规则名称，您可以通过查询策略治理规则列表接口来获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18948&data=128&isNormal=1&vid=121">查询策略治理规则列表</a>  */
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetPolicyTemplateResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                  `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetPolicyTemplateReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetPolicyTemplateReturnObjResponse struct {
	PolicyId         int64  `json:"policyId,omitempty"`         /*  规则ID  */
	Name             string `json:"name,omitempty"`             /*  规则名称  */
	Category         string `json:"category,omitempty"`         /*  规则模板类型  */
	Desc             string `json:"desc,omitempty"`             /*  规则模板描述  */
	Action           string `json:"action,omitempty"`           /*  规则治理动作，取值：<br /> enforce：拦截违规部署<br /> inform：告警  */
	Severity         string `json:"severity,omitempty"`         /*  规则治理等级，取值：<br /> high<br /> medium<br /> low  */
	NoConfig         int32  `json:"noConfig,omitempty"`         /*  是否需要配置策略，取值：<br /> 0：需要<br /> 1：不需要  */
	Template         string `json:"template,omitempty"`         /*  规则模板详情  */
	InstanceTemplate string `json:"instanceTemplate,omitempty"` /*  实例规则模板详情  */
	Status           int32  `json:"status,omitempty"`           /*  状态，取值：<br /> 1：有效<br /> 2：删除  */
	CreatedTime      string `json:"createdTime,omitempty"`      /*  创建时间  */
	CreatedBy        int64  `json:"createdBy,omitempty"`        /*  创建人  */
	ModifiedTime     string `json:"modifiedTime,omitempty"`     /*  修改时间  */
	ModifiedBy       int64  `json:"modifiedBy,omitempty"`       /*  修改人  */
}
