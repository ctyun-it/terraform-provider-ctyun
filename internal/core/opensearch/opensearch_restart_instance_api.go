package opensearch

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OpensearchRestartInstanceApi
/* 重启 OpenSearch 实例<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
*/type OpensearchRestartInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchRestartInstanceApi(client *core.CtyunClient) *OpensearchRestartInstanceApi {
	return &OpensearchRestartInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/os/openapi/v1/cluster/restart",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchRestartInstanceApi) Do(ctx context.Context, credential core.Credential, req *RestartInstanceRequest) (*RestartInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*RestartInstanceRequest
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
	var resp RestartInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RestartInstanceRequest struct {
	ClusterID   string           `json:"clusterId"`             /*  实例 id  */
	RestartType int32            `json:"restartType"`           /*  重启对象：1:实例重启，2:角色重启，3:节点重启  */
	OperateType string           `json:"operateType"`           /*  重启类型：RESTART_COMPONENT（全量重启）/ROLLING_RESTART_INSTANCE（滚动重启）*/
	RestartRole []string         `json:"restartRole,omitempty"` /*  重启角色列表（restartType=2 时必填）：COORDINATE/COLD/MASTER/EXCLUSIVE_MASTER  */
	HostInfo    *RestartHostInfo `json:"hostInfo,omitempty"`    /*  重启节点信息（restartType=3 时必填）*/
}

type RestartHostInfo struct {
	HostIP string `json:"hostIp"` /*  重启节点的 ipv4 地址  */
	Type   string `json:"type"`   /*  重启节点的类型：COORDINATE/COLD/MASTER/EXCLUSIVE_MASTER/ROUTER/ROUTER_CEREBRO  */
}

type RestartInstanceResponse struct {
	StatusCode string           `json:"statusCode"` /*  状态码，成功："200"，失败："500"  */
	Error      string           `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string           `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  RestartReturnObj `json:"returnObj"`  /*  返回结果  */
}

type RestartReturnObj struct {
	Result bool `json:"result"` /*  true:请求提交成功  */
}
