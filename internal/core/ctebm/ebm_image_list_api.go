package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbmImageListApi
/* 通过参数查询物理机可支持的镜像
 */type EbmImageListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmImageListApi(client *core.CtyunClient) *EbmImageListApi {
	return &EbmImageListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/image-list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmImageListApi) Do(ctx context.Context, credential core.Credential, req *EbmImageListRequest) (*EbmImageListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("deviceType", req.DeviceType)
	if req.ImageType != nil {
		ctReq.AddParam("imageType", *req.ImageType)
	}
	if req.ImageUUID != nil {
		ctReq.AddParam("imageUUID", *req.ImageUUID)
	}
	if req.OsName != nil {
		ctReq.AddParam("osName", *req.OsName)
	}
	if req.OsVersion != nil {
		ctReq.AddParam("osVersion", *req.OsVersion)
	}
	if req.OsType != nil {
		ctReq.AddParam("osType", *req.OsType)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmImageListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmImageListRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8317&data=87">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */DeviceType string /*  物理机套餐类型<br/><a href="https://www.ctyun.cn/document/10027724/10040107">查询资源池内物理机套餐</a><br /><a href="https://www.ctyun.cn/document/10027724/10040124">查询指定物理机的套餐信息</a>
	 */ImageType *string /*  镜像类型，可选择：private(私有镜像)、standard(标准镜像)、shared(共享镜像)；默认为standard
	 */ImageUUID *string /*  物理机镜像UUID
	 */OsName *string /*  操作系统名词，例如windows、ubuntu、centos等
	 */OsVersion *string /*  操作系统的具体版本信息
	 */OsType *string /*  操作系统类型，取值范围：linux，windows
	 */PageNo int32 /*  页码，默认值:1
	 */PageSize int32 /*  每页记录数目，取值范围:[1~10000]，默认值:10，单页最大记录不超过10000
	 */
}

type EbmImageListResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmImageListReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmImageListReturnObjResponse struct {
	CurrentCount int32 `json:"currentCount"` /*  当前页数量
	 */TotalCount int32 `json:"totalCount"` /*  总记录数
	 */TotalPage int32 `json:"totalPage"` /*  总页数
	 */Results []*EbmImageListReturnObjResultsResponse `json:"results"` /*  分页明细，元素类型是results，定义请参考表results
	 */
}

type EbmImageListReturnObjResultsResponse struct {
	NameZh *string `json:"nameZh"` /*  中文名称
	 */Format *string `json:"format"` /*  规格;包括squashfs,qcow2
	 */ImageType *string `json:"imageType"` /*  镜像类型;包括standard,private,shared，默认为standard
	 */IsShared *bool `json:"isShared"` /*  镜像是否共享
	 */Version *string `json:"version"` /*  版本
	 */ImageUUID *string `json:"imageUUID"` /*  镜像uuid
	 */NameEn *string `json:"nameEn"` /*  英文名称
	 */LayoutType *string `json:"layoutType"` /*  布局类型;包括lvm,direct
	 */Os *EbmImageListReturnObjResultsOsResponse `json:"os"` /*  操作系统
	 */
}

type EbmImageListReturnObjResultsOsResponse struct {
	Uuid *string `json:"uuid"` /*  操作系统uuid
	 */SuperUser *string `json:"superUser"` /*  超级管理员
	 */Platform *string `json:"platform"` /*  平台
	 */Version *string `json:"version"` /*  版本
	 */Architecture *string `json:"architecture"` /*  支持的机器类型
	 */NameEn *string `json:"nameEn"` /*  英文名称
	 */Bits int32 `json:"bits"` /*  比特数
	 */OsType *string `json:"osType"` /*  操作系统类别
	 */NameZh *string `json:"nameZh"` /*  中文名称
	 */
}
