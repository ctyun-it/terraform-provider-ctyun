package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListMultipartUploadsApi
/* 查询正在进行中的分段上传。
 */type ZosListMultipartUploadsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListMultipartUploadsApi(client *core.CtyunClient) *ZosListMultipartUploadsApi {
	return &ZosListMultipartUploadsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-multipart-uploads",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListMultipartUploadsApi) Do(ctx context.Context, credential core.Credential, req *ZosListMultipartUploadsRequest) (*ZosListMultipartUploadsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	if req.KeyMarker != "" {
		ctReq.AddParam("keyMarker", req.KeyMarker)
	}
	if req.UploadIDMarker != "" {
		ctReq.AddParam("uploadIDMarker", req.UploadIDMarker)
	}
	if req.MaxUploads != 0 {
		ctReq.AddParam("maxUploads", strconv.FormatInt(int64(req.MaxUploads), 10))
	}
	if req.Prefix != "" {
		ctReq.AddParam("prefix", req.Prefix)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListMultipartUploadsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListMultipartUploadsRequest struct {
	Bucket         string /*  桶名  */
	RegionID       string /*  区域 ID  */
	KeyMarker      string /*  和uploadIdMarker参数一起用于指定返回哪部分分段上传的信息。若未设置uploadIdMarker参数，则返回对象key按字典序大于等于keyMarker的片段信息。若设置了uploadIdMarker参数，则返回对象key大于等于keyMarker且uploadId大于uploadIdMarker的片段信息  */
	UploadIDMarker string /*  和keyMarker参数一起用于指定返回哪部分分段上传的信息，仅当设置了keyMarker参数的时候有效。设置后将返回对象key大于等于keyMarker且uploadId大于uploadIdMarker的片段信息  */
	MaxUploads     int64  /*  单次最多返回的分段上传数据，大小是1-1000，超过1000的数据会被视为1000  */
	Prefix         string /*  Key的前缀，只有以Prefix为开头的Key才会被返回  */
}

type ZosListMultipartUploadsResponse struct {
	StatusCode  int64                                     `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                    `json:"message,omitempty"`     /*  状态描述  */
	Description string                                    `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListMultipartUploadsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                    `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                    `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListMultipartUploadsReturnObjResponse struct {
	Bucket             string                                             `json:"bucket,omitempty"`             /*  启动分段上传的桶的名称  */
	KeyMarker          string                                             `json:"keyMarker,omitempty"`          /*  与请求中设置的 keyMarker 相同  */
	UploadIDMarker     string                                             `json:"uploadIDMarker,omitempty"`     /*  响应结果列表中的起始 uploadId  */
	NextKeyMarker      string                                             `json:"nextKeyMarker,omitempty"`      /*  当列表被截断时，此元素指定应用于后续请求中的 key 标记请求参数的值  */
	NextUploadIDMarker string                                             `json:"nextUploadIDMarker,omitempty"` /*  当列表被截断时，此元素指定应用于后续请求中的 upload-id-marker 请求参数的值  */
	MaxUploads         int64                                              `json:"maxUploads,omitempty"`         /*  可以包含在响应中的分段上传的最大数量  */
	IsTruncated        *bool                                              `json:"isTruncated"`                  /*  指示返回的分段列表是否被截断。值为 true 表示列表已被截断。如果分段上传的数量超过最大上传允许或指定的限制，则可以截断该列表  */
	Uploads            []*ZosListMultipartUploadsReturnObjUploadsResponse `json:"uploads"`                      /*  与特定分段上传相关的元素的容器。响应可以包含零个或多个 Upload 元素  */
}

type ZosListMultipartUploadsReturnObjUploadsResponse struct {
	UploadID     string                                                    `json:"uploadID,omitempty"`     /*  uploadID  */
	Key          string                                                    `json:"key,omitempty"`          /*  对象名  */
	Initiated    string                                                    `json:"initiated,omitempty"`    /*  分段上传初始化时间，ISO-8601 格式的日期字符串  */
	StorageClass string                                                    `json:"storageClass,omitempty"` /*  存储类，可选的值由 STANDARD，STANDARD_IA，GLACIER  */
	Owner        *ZosListMultipartUploadsReturnObjUploadsOwnerResponse     `json:"owner"`                  /*  所有者  */
	Initiator    *ZosListMultipartUploadsReturnObjUploadsInitiatorResponse `json:"initiator"`              /*  初始化分段上传事件的用户  */
}

type ZosListMultipartUploadsReturnObjUploadsOwnerResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名称  */
	ID          string `json:"ID,omitempty"`          /*  用户ID  */
}

type ZosListMultipartUploadsReturnObjUploadsInitiatorResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名称  */
	ID          string `json:"ID,omitempty"`          /*  用户ID  */
}
