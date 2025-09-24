package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutObjectStorageTypeApi
/* 转换对象的存储类型，例如从标准存储转换为归档存储。
 */type ZosPutObjectStorageTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutObjectStorageTypeApi(client *core.CtyunClient) *ZosPutObjectStorageTypeApi {
	return &ZosPutObjectStorageTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-object-storage-type",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutObjectStorageTypeApi) Do(ctx context.Context, credential core.Credential, req *ZosPutObjectStorageTypeRequest) (*ZosPutObjectStorageTypeResponse, error) {
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
	var resp ZosPutObjectStorageTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutObjectStorageTypeRequest struct {
	Bucket       string `json:"bucket,omitempty"`       /*  桶名  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域 ID  */
	Key          string `json:"key,omitempty"`          /*  需要转换存储类型的对象名称  */
	StorageClass string `json:"storageClass,omitempty"` /*  需要转换的存储类，支持标准存储：STANDARD，低频存储： STANDARD_IA，归档存储：GLACIER  */
	VersionID    string `json:"versionID,omitempty"`    /*  对象版本号，当桶开启版本控制时可使用  */
}

type ZosPutObjectStorageTypeResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
