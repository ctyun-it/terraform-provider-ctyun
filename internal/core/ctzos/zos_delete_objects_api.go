package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosDeleteObjectsApi
/* 批量删除指定的对象。
 */type ZosDeleteObjectsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosDeleteObjectsApi(client *core.CtyunClient) *ZosDeleteObjectsApi {
	return &ZosDeleteObjectsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/delete-objects",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosDeleteObjectsApi) Do(ctx context.Context, credential core.Credential, req *ZosDeleteObjectsRequest) (*ZosDeleteObjectsResponse, error) {
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
	var resp ZosDeleteObjectsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosDeleteObjectsRequest struct {
	Bucket   string                         `json:"bucket,omitempty"`   /*  桶名  */
	RegionID string                         `json:"regionID,omitempty"` /*  区域 ID  */
	Delete   *ZosDeleteObjectsDeleteRequest `json:"delete"`             /*  要删除的对象  */
}

type ZosDeleteObjectsDeleteRequest struct {
	Objects []*ZosDeleteObjectsDeleteObjectsRequest `json:"objects"` /*  对象的数组  */
	Quiet   *bool                                   `json:"quiet"`   /*  静默模式，默认 false。若为 true， 则响应不会返回每个对象的删除结果，仅返回失败的结果  */
}

type ZosDeleteObjectsDeleteObjectsRequest struct {
	Key       string `json:"key,omitempty"`       /*  对象名  */
	VersionID string `json:"versionID,omitempty"` /*  版本ID，在开启多版本时可使用  */
}

type ZosDeleteObjectsResponse struct {
	ReturnObj   *ZosDeleteObjectsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	StatusCode  int64                              `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                             `json:"message,omitempty"`     /*  状态描述  */
	Description string                             `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                             `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                             `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosDeleteObjectsReturnObjResponse struct {
	Deleted        []*ZosDeleteObjectsReturnObjDeletedResponse `json:"deleted"`                  /*  已删除的对象  */
	RequestCharged string                                      `json:"requestCharged,omitempty"` /*  如果存在，则表明请求者已成功为请求收费。  */
	Errors         []*ZosDeleteObjectsReturnObjErrorsResponse  `json:"errors"`                   /*  删除失败的信息  */
}

type ZosDeleteObjectsReturnObjDeletedResponse struct {
	Key                   string `json:"key,omitempty"`                   /*  对象名  */
	VersionID             string `json:"versionID,omitempty"`             /*  版本ID，在开启多版本时可使用  */
	DeleteMarker          *bool  `json:"deleteMarker"`                    /*  指定永久删除的版本化对象是（true）还是不是（false）删除标记。在简单的 DELETE 中，此标头指示是否（true）或不（false）创建了删除标记  */
	DeleteMarkerVersionID string `json:"deleteMarkerVersionID,omitempty"` /*  作为DELETE操作的结果而创建的删除标记的版本ID。如果你删除一个特定的对象版本，这个头返回的值是被删除的对象版本的版本ID  */
}

type ZosDeleteObjectsReturnObjErrorsResponse struct {
	Key       string `json:"key,omitempty"`       /*  对象名  */
	VersionID string `json:"versionID,omitempty"` /*  版本ID，在开启多版本时可使用  */
	Code      string `json:"code,omitempty"`      /*  错误码  */
	Message   string `json:"message,omitempty"`   /*  错误消息  */
}
