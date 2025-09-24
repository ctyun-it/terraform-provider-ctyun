package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeInstanceVersionApi
/* 查询分布式缓存Redis实例版本信息。
 */type Dcs2DescribeInstanceVersionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeInstanceVersionApi(client *core.CtyunClient) *Dcs2DescribeInstanceVersionApi {
	return &Dcs2DescribeInstanceVersionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/instanceManageMgrServant/describeInstanceVersion",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeInstanceVersionApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeInstanceVersionRequest) (*Dcs2DescribeInstanceVersionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeInstanceVersionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeInstanceVersionRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例 ID  */
}

type Dcs2DescribeInstanceVersionResponse struct {
	StatusCode int32                                         `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                        `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeInstanceVersionReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                        `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                        `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                        `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeInstanceVersionReturnObjResponse struct {
	EngineMajorVersionInfo *Dcs2DescribeInstanceVersionReturnObjEngineMajorVersionInfoResponse `json:"engineMajorVersionInfo"` /*  引擎大版本信息  */
	EngineMinorVersionInfo *Dcs2DescribeInstanceVersionReturnObjEngineMinorVersionInfoResponse `json:"engineMinorVersionInfo"` /*  引擎小版本信息  */
	ProxyVersionInfo       *Dcs2DescribeInstanceVersionReturnObjProxyVersionInfoResponse       `json:"proxyVersionInfo"`       /*  代理版本信息  */
}

type Dcs2DescribeInstanceVersionReturnObjEngineMajorVersionInfoResponse struct {
	EngineMajorVersion           string   `json:"engineMajorVersion,omitempty"` /*  引擎大版本  */
	EngineVersionItems           []string `json:"engineVersionItems"`           /*  引擎所有大版本列表  */
	UpgradableEngineVersionItems []string `json:"upgradableEngineVersionItems"` /*  引擎可升级大版本列表  */
}

type Dcs2DescribeInstanceVersionReturnObjEngineMinorVersionInfoResponse struct {
	EngineMinorVersion                string   `json:"engineMinorVersion,omitempty"`      /*  引擎小版本  */
	UpgradableEngineMinorVersionItems []string `json:"upgradableEngineMinorVersionItems"` /*  引擎可升级小版本列表  */
}

type Dcs2DescribeInstanceVersionReturnObjProxyVersionInfoResponse struct {
	ProxyMinorVersion            string   `json:"proxyMinorVersion,omitempty"`  /*  代理版本  */
	UpgradableProxyMinorVersions []string `json:"upgradableProxyMinorVersions"` /*  可升级代理版本列表  */
}
