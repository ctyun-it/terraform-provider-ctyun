package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtimageValidateImageFileSrcApi
/* 校验 1 个镜像文件地址<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageValidateImageFileSrcApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageValidateImageFileSrcApi(client *core.CtyunClient) *CtimageValidateImageFileSrcApi {
	return &CtimageValidateImageFileSrcApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/image/validate-file-source",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageValidateImageFileSrcApi) Do(ctx context.Context, credential core.Credential, req *CtimageValidateImageFileSrcRequest) (*CtimageValidateImageFileSrcResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("imageFileSource", req.ImageFileSource)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageValidateImageFileSrcResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageValidateImageFileSrcRequest struct {
	ImageFileSource string `json:"imageFileSource,omitempty"` /*  镜像文件地址，即存储桶内对象的 URL。注意：<br />1. 格式应为 {internetEndpoint}/{bucket}/{key}。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5305&data=89&isNormal=1&vid=83" target="_blank">访问控制 endpoint 查询</a>接口来查询外网 endpoint，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6247&data=105&isNormal=1&vid=99" target="_blank">查询所有桶</a>接口来查询您拥有的桶的列表，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6271&data=105&isNormal=1&vid=99" target="_blank">查询桶信息</a>接口来查询 1 个桶的详细信息，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6273&data=105&isNormal=1&vid=99" target="_blank">查看对象列表</a>接口来查询存储桶内所有对象，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6282&data=105&isNormal=1&vid=99" target="_blank">查询对象是否存在</a>接口来查询 1 个对象的详细信息，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=9&api=6290&data=105&isNormal=1&vid=99" target="_blank">获取对象 ACL </a>接口来查询 1 个对象的访问权限控制列表<br />2. 所指向的桶应是属于您的、存储类型为 STANDARD、未加密的桶<br />3. 所指向的对象应是所指向的桶内的存储类型为 STANDARD、未加密的对象。此对象在非多可用区资源池中还应可公共读。您需自行确保此对象是 QCOW2、RAW、VHD 或 VMDK 格式的镜像文件  */
	RegionID        string `json:"regionID,omitempty"`        /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
}

type CtimageValidateImageFileSrcResponse struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageValidateImageFileSrcReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageValidateImageFileSrcReturnObjResponse struct {
	ArchivedObject         *bool `json:"archivedObject"`         /*  用于表示所指定的对象是否被归档冻结的标识  */
	EncryptedObject        *bool `json:"encryptedObject"`        /*  用于表示所指定的对象是否被加密的标识  */
	ObjectSize             int64 `json:"objectSize,omitempty"`   /*  对象大小。单位为 byte  */
	PubliclyReadableObject *bool `json:"publiclyReadableObject"` /*  用于表示在非多可用区资源池中所指定的对象是否可公共读的标识  */
	ReachableBucket        *bool `json:"reachableBucket"`        /*  用于表示所指定的桶对您而言是否存在的标识  */
	ReachableObject        *bool `json:"reachableObject"`        /*  用于表示所指定的对象对您而言是否存在的标识  */
	ValidEndpoint          *bool `json:"validEndpoint"`          /*  用于表示所指定的 endpoint 是否有效且归属当前资源池的标识  */
	ValidURL               *bool `json:"validURL"`               /*  用于表示所指定的 URL 是否满足基本规范的标识  */
}
