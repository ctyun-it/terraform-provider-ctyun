package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetNodeDetailApi
/* 查询节点详情
 */type CcseGetNodeDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetNodeDetailApi(client *core.CtyunClient) *CcseGetNodeDetailApi {
	return &CcseGetNodeDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodes/{nodeName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetNodeDetailApi) Do(ctx context.Context, credential core.Credential, req *CcseGetNodeDetailRequest) (*CcseGetNodeDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodeName", req.NodeName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetNodeDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetNodeDetailRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	NodeName  string /*  节点名称  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=5851&amp;data=87">资源池列表查询</a>  */
}

type CcseGetNodeDetailResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                              `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetNodeDetailReturnObjResponse `json:"returnObj"`            /*  返回参数  */
	Error      string                              `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetNodeDetailReturnObjResponse struct {
	NodeName         string `json:"nodeName,omitempty"`         /*  集群节点名称  */
	NodeType         int32  `json:"nodeType,omitempty"`         /*  节点类型，取值：<br/>1：master <br/>2：slave  */
	NodeStatus       string `json:"nodeStatus,omitempty"`       /*  节点状态，取值：<br/>normal：健康<br/>abnormal：异常<br/>expulsion：驱逐中  */
	IsSchedule       int32  `json:"isSchedule,omitempty"`       /*  是否调度，取值： 1：是 <br />0：否  */
	IsEvict          int32  `json:"isEvict,omitempty"`          /*  是否驱逐，取值： 1：是 <br />0：否  */
	CreatedTime      string `json:"createdTime,omitempty"`      /*  创建时间  */
	HostIp           string `json:"hostIp,omitempty"`           /*  主机管理ip  */
	HostExtraIpv6    string `json:"hostExtraIpv6,omitempty"`    /*  主机业务ipv6  */
	KubeletVersion   string `json:"kubeletVersion,omitempty"`   /*  Kubelet 版本  */
	PodCidr          string `json:"podCidr,omitempty"`          /*  Pod CIDR  */
	KernelVersion    string `json:"kernelVersion,omitempty"`    /*  内核版本  */
	OsImageVersion   string `json:"osImageVersion,omitempty"`   /*  OS 镜像  */
	KubeProxyVersion string `json:"kubeProxyVersion,omitempty"` /*  KubeProxy 版本  */
	EcsId            string `json:"ecsId,omitempty"`            /*  云主机ID  */
	HostType         string `json:"hostType,omitempty"`         /*  host类型  */
	ZoneName         string `json:"zoneName,omitempty"`         /*  可用区名称  */
	ZoneCode         string `json:"zoneCode,omitempty"`         /*  可用区编码  */
}
