package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateSubUserPermissionsApi
/* 调用该接口全量更新子账号集群授权信息。
 */type CcseUpdateSubUserPermissionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateSubUserPermissionsApi(client *core.CtyunClient) *CcseUpdateSubUserPermissionsApi {
	return &CcseUpdateSubUserPermissionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/binding",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateSubUserPermissionsApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateSubUserPermissionsRequest) (*CcseUpdateSubUserPermissionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
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
	var resp CcseUpdateSubUserPermissionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateSubUserPermissionsRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	UserIds []int64                                    `json:"userIds"` /*  用户ID，可在云容器引擎控制台 > 安全管理 > 授权页面查看子账号列表及对应的用户ID  */
	Data    []*CcseUpdateSubUserPermissionsDataRequest `json:"data"`    /*  权限设置列表  */
}

type CcseUpdateSubUserPermissionsDataRequest struct {
	Namespace   string `json:"namespace,omitempty"`   /*  命名空间，不填默认全部命名空间  */
	ClusterRole string `json:"clusterRole,omitempty"` /*  需要绑定的clusterRole，创建，更新需要填写；删除不需要填
	目前支持的检查项及含义如下
	ccse:preset:admin	管理员
	ccse:preset:ops	运维人员
	ccse:preset:dev	开发人员
	ccse:preset:view	受限用户  */
	BindingName string `json:"bindingName,omitempty"` /*  clusterRoleBinding或者roleBinding的名称，用于表示权限绑定，创建不需要填写，删除或更新必填；删除或更新通过查询子账号集群授权信息获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18049&data=128&isNormal=1&vid=121">查询子账号集群授权信息</a>  */
	BindingKind string `json:"bindingKind,omitempty"` /*  rolebinding类型，创建不需要填写，删除或更新必填；删除或更新通过查询子账号集群授权信息获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18049&data=128&isNormal=1&vid=121">查询子账号集群授权信息</a>  */
	OperType int32 `json:"operType,omitempty"` /*  操作类型。0-增加，1-删除，2-更新。不填默认为0  */
}

type CcseUpdateSubUserPermissionsResponse struct {
	StatusCode int32                                          `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                         `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseUpdateSubUserPermissionsReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码  */
}

type CcseUpdateSubUserPermissionsReturnObjResponse struct{}
