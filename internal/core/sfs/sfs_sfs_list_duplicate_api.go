package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListDuplicateApi
/* 查询复制配置
 */type SfsSfsListDuplicateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListDuplicateApi(client *core.CtyunClient) *SfsSfsListDuplicateApi {
	return &SfsSfsListDuplicateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-duplicate",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListDuplicateApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListDuplicateRequest) (*SfsSfsListDuplicateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
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
	var resp SfsSfsListDuplicateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListDuplicateRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  源/目的弹性文件功能系统唯一 ID  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*   页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type SfsSfsListDuplicateResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListDuplicateReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码，大驼峰形式  */
}

type SfsSfsListDuplicateReturnObjResponse struct {
	List         []*SfsSfsListDuplicateReturnObjListResponse `json:"list"`         /*  查询信息列表  */
	CurrentCount int32                                       `json:"currentCount"` /*  当前查询到的复制关系个数  */
	TotalCount   int32                                       `json:"totalCount"`   /*  资源池下用户复制关系总数  */
	PageSize     int32                                       `json:"pageSize"`     /*  每页包含的元素个数  */
	PageNo       int32                                       `json:"pageNo"`       /*  页号  */
}

type SfsSfsListDuplicateReturnObjListResponse struct {
	Status         string `json:"status"`         /*  跨域复制状态  */
	DuplicateUID   string `json:"duplicateUID"`   /*  跨域复制任务ID  */
	SrcRegionID    string `json:"srcRegionID"`    /*  源文件系统资源池ID  */
	DstRegionID    string `json:"dstRegionID"`    /*  目标文件系统资源池ID  */
	SrcAzName      string `json:"srcAzName"`      /*  源文件AZ  */
	DstAzName      string `json:"dstAzName"`      /*  目标文件AZ  */
	SrcSfsName     string `json:"srcSfsName"`     /*  源文件系统名称  */
	DstSfsName     string `json:"dstSfsName"`     /*  目标文件系统名称  */
	SrcSfsUID      string `json:"srcSfsUID"`      /*  源文件系统ID  */
	DstSfsUID      string `json:"dstSfsUID"`      /*  目标文件系统ID  */
	SrcCephID      string `json:"srcCephID"`      /*  源文件系统ceph底层的id  */
	DstCephID      string `json:"dstCephID"`      /*  目标文件系统ceph底层的id  */
	Last_sync_time int64  `json:"last_sync_time"` /*  最后一次同步时刻，epoch 时戳，精度毫秒  */
	LastSyncTime   int64  `json:"lastSyncTime"`   /*  最后一次同步时刻，epoch 时戳，精度毫秒  */
}
