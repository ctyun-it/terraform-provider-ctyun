package opensearch

import (
	"context"
	"encoding/json"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OpensearchScaleInstanceApi
/* 扩容 OpenSearch 实例<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
*/type OpensearchScaleInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchScaleInstanceApi(client *core.CtyunClient) *OpensearchScaleInstanceApi {
	return &OpensearchScaleInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/os/openapi/v1/cluster/extendInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchScaleInstanceApi) Do(ctx context.Context, credential core.Credential, req *ScaleInstanceRequest) (*ScaleInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScaleInstanceRequest
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
	var resp ScaleInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScaleInstanceRequest struct {
	ClusterID       string `json:"clusterId"`         /*  集群 ID  */
	NodeGroupName   string `json:"nodeGroupName"`     /*  节点组名称：MASTER（数据节点）/COORDINATE（专属协调节点）/COLD（冷数据节点）*/
	IncreaseHostNum int32  `json:"increaseHostNum"`   /*  扩容节点的数量  */
	AutoPay         bool   `json:"autoPay,omitempty"` /*  true：自动支付，false：不自动支付，默认 false  */
}

type ScaleInstanceResponse struct {
	StatusCode int32           `json:"statusCode"` /*  状态码，成功："200"，失败："500"  */
	Error      string          `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string          `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  json.RawMessage `json:"returnObj"`  /*  返回结果  */
}
