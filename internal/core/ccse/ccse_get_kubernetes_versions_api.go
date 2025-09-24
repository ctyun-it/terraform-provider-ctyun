package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetKubernetesVersionsApi
/* 调用该接口查询支持Kubernetes版本。
 */type CcseGetKubernetesVersionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetKubernetesVersionsApi(client *core.CtyunClient) *CcseGetKubernetesVersionsApi {
	return &CcseGetKubernetesVersionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/metadata/versions",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetKubernetesVersionsApi) Do(ctx context.Context, credential core.Credential, req *CcseGetKubernetesVersionsRequest) (*CcseGetKubernetesVersionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetKubernetesVersionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetKubernetesVersionsRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
}

type CcseGetKubernetesVersionsResponse struct {
	StatusCode int32                                         `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                        `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  []*CcseGetKubernetesVersionsReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                        `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetKubernetesVersionsReturnObjResponse struct {
	Runtimes []*CcseGetKubernetesVersionsReturnObjRuntimesResponse `json:"runtimes"`           /*  容器运行时  */
	Multi_az string                                                `json:"multi_az,omitempty"` /*  是否支持多AZ  */
	Version  string                                                `json:"version,omitempty"`  /*  版本  */
}

type CcseGetKubernetesVersionsReturnObjRuntimesResponse struct {
	Name    string `json:"name,omitempty"`    /*  容器运行时名称，包括docker、containerd  */
	Version string `json:"version,omitempty"` /*  容器运行时版本  */
}
