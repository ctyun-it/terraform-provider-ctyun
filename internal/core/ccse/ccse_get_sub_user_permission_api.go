package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseGetSubUserPermissionApi
/* 调用该接口查询子账号集群授权信息。
 */type CcseGetSubUserPermissionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetSubUserPermissionApi(client *core.CtyunClient) *CcseGetSubUserPermissionApi {
	return &CcseGetSubUserPermissionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/binding",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetSubUserPermissionApi) Do(ctx context.Context, credential core.Credential, req *CcseGetSubUserPermissionRequest) (*CcseGetSubUserPermissionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("userId", strconv.FormatInt(int64(req.UserId), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetSubUserPermissionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetSubUserPermissionRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	UserId int64 /*  用户ID，可在云容器引擎控制台 > 安全管理 > 授权页面查看子账号列表及对应的用户ID  */
}

type CcseGetSubUserPermissionResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                       `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  []*CcseGetSubUserPermissionReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetSubUserPermissionReturnObjResponse struct {
	BindingName string `json:"bindingName,omitempty"` /*  绑定标识  */
	BindingKind string `json:"bindingKind,omitempty"` /*  角色绑定类型  */
	Namespace   string `json:"namespace,omitempty"`   /*  角色绑定命名空间，为空表示全部命名空间  */
	ClusterRole string `json:"clusterRole,omitempty"` /*  绑定的角色名  */
}
