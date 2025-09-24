package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsListMountpointApi
/* 根据文件系统resourceID ，查询文件系统挂载点列表
 */type SfsListMountpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListMountpointApi(client *core.CtyunClient) *SfsListMountpointApi {
	return &SfsListMountpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-mountpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListMountpointApi) Do(ctx context.Context, credential core.Credential, req *SfsListMountpointRequest) (*SfsListMountpointResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("resourceID", req.ResourceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsListMountpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListMountpointRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
	ResourceID string `json:"resourceID,omitempty"` /*  文件系统资源ID  */
}

type SfsListMountpointResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                              `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                              `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsListMountpointReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsListMountpointReturnObjResponse struct {
	TotalCount     int32                                               `json:"totalCount"`     /*  挂载点个数  */
	CurrentCount   int32                                               `json:"currentCount"`   /*  当前页码的元素个数	  */
	MountPointList []*SfsListMountpointReturnObjMountPointListResponse `json:"mountPointList"` /*  挂载点信息列表  */
	PageSize       int32                                               `json:"pageSize"`       /*  每页个数  */
	PageNo         int32                                               `json:"pageNo"`         /*  页数  */
}

type SfsListMountpointReturnObjMountPointListResponse struct {
	MountPointID string `json:"mountPointID"` /*  文件系统挂载点ID  */
	SubnetID     string `json:"subnetID"`     /*  子网ID  */
	VpcID        string `json:"vpcID"`        /*  vpc ID  */
	SubnetName   string `json:"subnetName"`   /*  子网名字  */
	VpcName      string `json:"vpcName"`      /*  vpc名字  */
}
