package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyBigKeyPolicyApi
/* 设置分布式缓存Redis实例大key自动分析策略。
 */type Dcs2ModifyBigKeyPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyBigKeyPolicyApi(client *core.CtyunClient) *Dcs2ModifyBigKeyPolicyApi {
	return &Dcs2ModifyBigKeyPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/keyAnalysisMgrServant/modifyBigKeyPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyBigKeyPolicyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyBigKeyPolicyRequest) (*Dcs2ModifyBigKeyPolicyResponse, error) {
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
	var resp Dcs2ModifyBigKeyPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyBigKeyPolicyRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	ModifyType string `json:"modifyType,omitempty"` /*  操作类型<li> 0：开启自动扫描配置<li>1：关闭自动扫描配置<li>2：修改自动扫描配置  */
	Days       string `json:"days,omitempty"`       /*  日期范围，1-7表示周一至周日，多个日期使用英文逗号分隔。  */
	Hours      string `json:"hours,omitempty"`      /*  整点  */
}

type Dcs2ModifyBigKeyPolicyResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifyBigKeyPolicyReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyBigKeyPolicyReturnObjResponse struct{}
