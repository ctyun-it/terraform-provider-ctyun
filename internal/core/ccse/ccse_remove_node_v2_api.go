package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseRemoveNodeV2Api
/* 在指定集群下移除节点
 */type CcseRemoveNodeV2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseRemoveNodeV2Api(client *core.CtyunClient) *CcseRemoveNodeV2Api {
	return &CcseRemoveNodeV2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/nodes/remove",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseRemoveNodeV2Api) Do(ctx context.Context, credential core.Credential, req *CcseRemoveNodeV2Request) (*CcseRemoveNodeV2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CcseRemoveNodeV2Request
		RegionId   interface{} `json:"regionId,omitempty"`
		ClusterId  interface{} `json:"clusterId,omitempty"`
		NodePoolId interface{} `json:"nodePoolId,omitempty"`
	}{
		req, nil, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseRemoveNodeV2Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseRemoveNodeV2Request struct {
	ClusterId  string `json:"clusterId,omitempty"`  /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string `json:"nodePoolId,omitempty"` /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string `json:"regionId,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	Nodes     []string `json:"nodes"`               /*  移除节点的名称，可通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18020&data=128&isNormal=1&vid=121" target="_blank">查询节点池详情</a>接口返回的节点池节点信息查询节点名称  */
	LoginType string   `json:"loginType,omitempty"` /*  云主机密码登录类型：目前只支持 password  */
	Password  string   `json:"password,omitempty"`  /*  云主机或物理机，用户登录密码，如果loginType=password，这该项为必填项，满足以下规则：
	物理机：用户密码，满足以下规则：长度在8～30个字符；
	必须包含大小写字母和（至少一个数字或者特殊字符）；
	特殊符号可选：()`~!@#$%&*_-+=\
	云主机：长度在8～30个字符；
	必须包含大写字母、小写字母、数字以及特殊符号中的三项；
	特殊符号可选：()`-!@#$%^&*_-+=｜{}[]:;'<>,.?/且不能以斜线号 / 开头；
	不能包含3个及以上连续字符；
	移除节点时password字段可选择加密，具体加密方法请参见<a href="https://www.ctyun.cn/document/10083472/11002096" target="_blank">password字段加密的方法</a>  */
	ForcedRemoved *bool `json:"forcedRemoved"` /*  是否强制移除，默认是true  */
}

type CcseRemoveNodeV2Response struct {
	StatusCode int32  `json:"statusCode"` /*  响应状态码  */
	RequestId  string `json:"requestId"`  /*  请求ID  */
	Message    string `json:"message"`    /*  响应信息  */
	ReturnObj  *bool  `json:"returnObj"`  /*  响应对象  */
	Error      string `json:"error"`      /*  错误码，参见错误码说明  */
}
