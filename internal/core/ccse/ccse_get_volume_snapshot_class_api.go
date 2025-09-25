package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetVolumeSnapshotClassApi
/* 查询指定集群下的VolumeSnapshotClass详情
 */type CcseGetVolumeSnapshotClassApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetVolumeSnapshotClassApi(client *core.CtyunClient) *CcseGetVolumeSnapshotClassApi {
	return &CcseGetVolumeSnapshotClassApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/apis/snapshot.storage.k8s.io/v1/volumesnapshotclasses/{volumeSnapshotClassName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetVolumeSnapshotClassApi) Do(ctx context.Context, credential core.Credential, req *CcseGetVolumeSnapshotClassRequest) (*CcseGetVolumeSnapshotClassResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("volumeSnapshotClassName", req.VolumeSnapshotClassName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetVolumeSnapshotClassResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetVolumeSnapshotClassRequest struct {
	ClusterId               string /*  集群id， 获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	VolumeSnapshotClassName string /*  资源名称  */
	RegionId                string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetVolumeSnapshotClassResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  string `json:"returnObj,omitempty"`  /*  返回结果  */
	Error      string `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}
