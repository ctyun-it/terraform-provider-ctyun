package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsUpdateIopsEbsApi
/* 支持修改XSSD类型云硬盘的预配置IOPS。
 */type EbsUpdateIopsEbsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsUpdateIopsEbsApi(client *core.CtyunClient) *EbsUpdateIopsEbsApi {
	return &EbsUpdateIopsEbsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/update-iops-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsUpdateIopsEbsApi) Do(ctx context.Context, credential core.Credential, req *EbsUpdateIopsEbsRequest) (*EbsUpdateIopsEbsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsUpdateIopsEbsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsUpdateIopsEbsRequest struct {
	ProvisionedIops int32 `json:"provisionedIops"` /*  XSSD类型云硬盘的预配置IOPS值，最小值为1，其他类型的盘不支持设置。具体取值范围如下：
	●XSSD-0：（基础IOPS（min{1800+12×容量， 10000}） + 预配置IOPS） ≤ min{500×容量，100000}
	●XSSD-1：（基础IOPS（min{1800+50×容量， 50000}） + 预配置IOPS） ≤ min{500×容量，100000}
	●XSSD-2：（基础IOPS（min{3000+50×容量， 100000}） + 预配置IOPS） ≤ min{500×容量，1000000}  */
	DiskID   string  `json:"diskID,omitempty"`   /*  云硬盘ID。  */
	RegionID *string `json:"regionID,omitempty"` /*  如本地语境支持保存regionID，那么建议传递。  */
}

type EbsUpdateIopsEbsResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     *string `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码.<br/>参考结果码(#通用结果码)。  */
	Error       *string `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码. 请参考结果码(#通用结果码（大驼峰格式）)。  */
}
