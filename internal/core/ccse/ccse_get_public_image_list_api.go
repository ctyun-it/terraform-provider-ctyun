package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetPublicImageListApi
/* 获取公共镜像及其信息列表
 */type CcseGetPublicImageListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetPublicImageListApi(client *core.CtyunClient) *CcseGetPublicImageListApi {
	return &CcseGetPublicImageListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/osimages",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetPublicImageListApi) Do(ctx context.Context, credential core.Credential, req *CcseGetPublicImageListRequest) (*CcseGetPublicImageListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("flavorName", req.FlavorName)
	ctReq.AddParam("vmType", req.VmType)
	if req.ProjectId != "" {
		ctReq.AddParam("projectId", req.ProjectId)
	}
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetPublicImageListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetPublicImageListRequest struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=5851&amp;data=87">资源池列表查询</a>  */
	FlavorName string `json:"flavorName,omitempty"` /*  云主机规格名称  */
	VmType     string `json:"vmType,omitempty"`     /*  查询镜像类型，默认ecs。
	ecs - 云主机镜像
	ebm - 裸金属镜像  */
	ProjectId string `json:"projectId,omitempty"` /*  企业项目ID  */
	AzName    string `json:"azName,omitempty"`    /*  平台可用区名称，查询裸金属镜像时必填  */
}

type CcseGetPublicImageListResponse struct {
	StatusCode int32                                      `json:"statusCode"` /*  状态码  */
	Message    string                                     `json:"message"`    /*  提示信息  */
	ReturnObj  []*CcseGetPublicImageListReturnObjResponse `json:"returnObj"`  /*  返回数据列表  */
	Error      string                                     `json:"error"`      /*  错误码  */
}

type CcseGetPublicImageListReturnObjResponse struct {
	ImageClass string `json:"imageClass"` /*  镜像类别
	ECS - 云主机
	BMS - 裸金属物理机  */
	ImageID   string `json:"imageID"`   /*  镜像ID  */
	ImageName string `json:"imageName"` /*  镜像名称  */
	OsType    string `json:"osType"`    /*  操作系统类型
	linux - Linux操作系统
	windows - Windows操作系统  */
	OsDistro     string `json:"osDistro"`     /*  操作系统发行版  */
	OsVersion    string `json:"osVersion"`    /*  操作系统版本  */
	Visibility   string `json:"visibility"`   /*  镜像是否公共可见  */
	Architecture string `json:"architecture"` /*  系统架构
	aarch64 - AArch64 架构
	loongarch64 - LoongArch64 架构
	sw_64 - sw_64 架构
	x86_64 - x86_64 架构  */
}
