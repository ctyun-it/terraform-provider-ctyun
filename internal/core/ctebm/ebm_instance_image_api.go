package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbmInstanceImageApi
/* 根据实例ID查询所使用的镜像信息
 */type EbmInstanceImageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmInstanceImageApi(client *core.CtyunClient) *EbmInstanceImageApi {
	return &EbmInstanceImageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/instance-image",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmInstanceImageApi) Do(ctx context.Context, credential core.Credential, req *EbmInstanceImageRequest) (*EbmInstanceImageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("instanceUUID", req.InstanceUUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmInstanceImageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmInstanceImageRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */
}

type EbmInstanceImageResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmInstanceImageReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmInstanceImageReturnObjResponse struct {
	Results []*EbmInstanceImageReturnObjResultsResponse `json:"results"` /*  分页明细,元素类型是results,定义请参考表results
	 */
}

type EbmInstanceImageReturnObjResultsResponse struct {
	NameZh *string `json:"nameZh"` /*  名称
	 */Format *string `json:"format"` /*  规格;包括squashfs,qcow2
	 */ImageType *string `json:"imageType"` /*  镜像类型;包括standard,private,shared
	 */IsShared *bool `json:"isShared"` /*  镜像是否共享
	 */Version *string `json:"version"` /*  版本
	 */ImageUUID *string `json:"imageUUID"` /*  镜像uuid
	 */NameEn *string `json:"nameEn"` /*  名称
	 */LayoutType *string `json:"layoutType"` /*  布局类型;包括lvm,direct
	 */Os *EbmInstanceImageReturnObjResultsOsResponse `json:"os"` /*  操作系统
	 */
}

type EbmInstanceImageReturnObjResultsOsResponse struct {
	Uuid *string `json:"uuid"` /*  操作系统uuid
	 */SuperUser *string `json:"superUser"` /*  超级管理员
	 */Platform *string `json:"platform"` /*  平台
	 */Version *string `json:"version"` /*  版本
	 */Architecture *string `json:"architecture"` /*  支持的机器类型
	 */NameEn *string `json:"nameEn"` /*  名称
	 */Bits int32 `json:"bits"` /*  比特数
	 */OsType *string `json:"osType"` /*  操作系统类别
	 */NameZh *string `json:"nameZh"` /*  名称
	 */
}
