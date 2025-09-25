package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsRenameSfsApi
/* 指定文件系统重命名此请求是异步处理，返回800代表请求下发成功，具体结果请使用【并行文件信息查询】确认
 */type HpfsRenameSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsRenameSfsApi(client *core.CtyunClient) *HpfsRenameSfsApi {
	return &HpfsRenameSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/rename-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsRenameSfsApi) Do(ctx context.Context, credential core.Credential, req *HpfsRenameSfsRequest) (*HpfsRenameSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsRenameSfsRequest
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
	var resp HpfsRenameSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsRenameSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  并行文件唯一ID  */
	SfsName  string `json:"sfsName,omitempty"`  /*  文件系统新名称  */
}

type HpfsRenameSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述  */
	Description string `json:"description"` /*  响应描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}
