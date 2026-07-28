package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageExportImageApi
/* 导出 1 份云主机私有镜像到指定的对象存储的桶<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：<br />1. 在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分<br />2. 接口请求成功后，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=7116&data=105&isNormal=1&vid=99" target="_blank">查询桶内碎片列表</a>接口来查询桶内所有分段上传的碎片，从而获知大致的导出进度。若涉及镜像文件格式转换（如镜像实际为 RAW 格式，导出时指定了 QCOW2 格式），则您需等待一段时间后才能查询到相应的桶内碎片，具体等待时间受镜像大小等因素影响
 */type CtimageExportImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageExportImageApi(client *core.CtyunClient) *CtimageExportImageApi {
	return &CtimageExportImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/image/export",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageExportImageApi) Do(ctx context.Context, credential core.Credential, req *CtimageExportImageRequest) (*CtimageExportImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtimageExportImageRequest
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
	var resp CtimageExportImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageExportImageRequest struct {
	Bucket          string `json:"bucket,omitempty"`          /*  存储桶名。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6247&data=105&isNormal=1&vid=99" target="_blank">查询所有桶</a>接口来查询您拥有的桶的列表，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6271&data=105&isNormal=1&vid=99" target="_blank">查询桶信息</a>接口来查询 1 个桶的详细信息。注意：所指定的桶应是属于您的、存储类型为 STANDARD 的桶  */
	Filename        string `json:"filename,omitempty"`        /*  导出文件名称。注意：长度为 2~100 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾  */
	ImageID         string `json:"imageID,omitempty"`         /*  镜像 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息。注意：<br />1. 所指定的镜像应是镜像状态为 active、镜像类型为空（系统盘镜像）或 data_disk_image 的云主机私有镜像<br />2. 仅在传入的 regionID 参数所指定的资源池具备导出系统盘镜像的功能时支持镜像类型为空（系统盘镜像）的云主机私有镜像。数据盘镜像同理<br />3. 所指定的镜像在导出完成前一般不能再导出  */
	RegionID        string `json:"regionID,omitempty"`        /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	ImageFileFormat string `json:"imageFileFormat,omitempty"` /*  镜像文件格式。取值范围（值：描述）：<br />raw：RAW 格式（默认值），<br />qcow2：QCOW2 格式，<br />vhd：VHD 格式，<br />vmdk：VMDK 格式<br />注意：在此参数值为 vhd 或 vmdk 时，若传入的 regionID 参数所指定的资源池不具备导出成 VHD 或 VMDK 格式的镜像文件的功能，则此参数将使用默认值  */
}

type CtimageExportImageResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                               `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageExportImageReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageExportImageReturnObjResponse struct{}
