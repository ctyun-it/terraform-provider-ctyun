package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsUpdateEbsNameApi
/* 云硬盘名称通常用来标识磁盘，云硬盘创建完成后，您可以修改名称。
 */type EbsUpdateEbsNameApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsUpdateEbsNameApi(client *core.CtyunClient) *EbsUpdateEbsNameApi {
	return &EbsUpdateEbsNameApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs/update-attr-ebs",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsUpdateEbsNameApi) Do(ctx context.Context, credential core.Credential, req *EbsUpdateEbsNameRequest) (*EbsUpdateEbsNameResponse, error) {
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
	var resp EbsUpdateEbsNameResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsUpdateEbsNameRequest struct {
	DiskName string  `json:"diskName,omitempty"` /*  磁盘名称。仅允许英文字母、数字及_或者-，长度为2-63字符，不能以特殊字符开头。  */
	DiskID   string  `json:"diskID,omitempty"`   /*  云硬盘ID。  */
	RegionID *string `json:"regionID,omitempty"` /*  如本地语境支持保存regionID，那么建议传递。  */
}

type EbsUpdateEbsNameResponse struct {
	StatusCode  int32   `json:"statusCode"`            /*  返回状态码(800为成功，900为失败)。  */
	Message     *string `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ErrorCode   *string `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}
