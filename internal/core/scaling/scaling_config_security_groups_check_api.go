package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingConfigSecurityGroupsCheckApi
/*  用于检查该账户下，哪些安全组被伸缩配置所使用<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingConfigSecurityGroupsCheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingConfigSecurityGroupsCheckApi(client *core.CtyunClient) *ScalingConfigSecurityGroupsCheckApi {
	return &ScalingConfigSecurityGroupsCheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/config/securitygroups-check",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingConfigSecurityGroupsCheckApi) Do(ctx context.Context, credential core.Credential, req *ScalingConfigSecurityGroupsCheckRequest) (*ScalingConfigSecurityGroupsCheckResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingConfigSecurityGroupsCheckRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ScalingConfigSecurityGroupsCheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingConfigSecurityGroupsCheckRequest struct {
	RegionID            string   `json:"regionID,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SecurityGroupIDList []string `json:"securityGroupIDList"` /*  安全组ID列表，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
}

type ScalingConfigSecurityGroupsCheckResponse struct {
	StatusCode  int32                                              `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingConfigSecurityGroupsCheckReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                             `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingConfigSecurityGroupsCheckReturnObjResponse struct {
	SecurityGroupIDList []string `json:"securityGroupIDList"` /*  安全组ID列表  */
}
