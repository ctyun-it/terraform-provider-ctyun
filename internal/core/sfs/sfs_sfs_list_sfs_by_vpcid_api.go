package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListSfsByVpcidApi
/* 查询租户指定vpc下的⽂件系统列表
 */type SfsSfsListSfsByVpcidApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListSfsByVpcidApi(client *core.CtyunClient) *SfsSfsListSfsByVpcidApi {
	return &SfsSfsListSfsByVpcidApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-sfs-by-vpcid",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListSfsByVpcidApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListSfsByVpcidRequest) (*SfsSfsListSfsByVpcidResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("vpcID", req.VpcID)
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsListSfsByVpcidResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListSfsByVpcidRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	VpcID    string `json:"vpcID,omitempty"`    /*  vpc ID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数，默认10  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码，默认1  */
}

type SfsSfsListSfsByVpcidResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                 `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListSfsByVpcidReturnObjResponse `json:"returnObj"`   /*  参考returnObj  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListSfsByVpcidReturnObjResponse struct {
	CurrentCount int32 `json:"currentCount"` /*  当前页码下查询回来的用户弹性文件数  */
	TotalCount   int32 `json:"totalCount"`   /*  绑定该vpc的弹性文件总数  */
	PageSize     int32 `json:"pageSize"`     /*  每页包含的元素个数。默认为1  */
	PageNo       int32 `json:"pageNo"`       /*  当前页码。默认为10  */
}
