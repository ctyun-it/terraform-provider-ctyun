package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsRegionDetailsSfsApi
/* 用于展示弹性文件下列信息：
 */ /* - 所支持的协议类型
 */ /* - 所支持的存储类型
 */ /* - 资源池下的az列表
 */ /* - 资源池的售罄情况
 */type SfsSfsRegionDetailsSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsRegionDetailsSfsApi(client *core.CtyunClient) *SfsSfsRegionDetailsSfsApi {
	return &SfsSfsRegionDetailsSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/region/storagetype",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsRegionDetailsSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsRegionDetailsSfsRequest) (*SfsSfsRegionDetailsSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.SfsType != "" {
		ctReq.AddParam("sfsType", req.SfsType)
	}
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
	var resp SfsSfsRegionDetailsSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsRegionDetailsSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsType  string `json:"sfsType,omitempty"`  /*  文件系统类型(capacity/performance/hdd_e)  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的数量  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  当前页码  */
}

type SfsSfsRegionDetailsSfsResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                   `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                   `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsRegionDetailsSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsRegionDetailsSfsReturnObjResponse struct {
	TotalCount        int32                                                       `json:"totalCount"`        /*  查询到的支持弹性文件详情的类型数量  */
	CurrentCount      int32                                                       `json:"currentCount"`      /*  当前页码的元素个数  */
	RegionSupportList []*SfsSfsRegionDetailsSfsReturnObjRegionSupportListResponse `json:"regionSupportList"` /*  可用区信息列表，参考 [regionSupportList]  */
	PageSize          int32                                                       `json:"pageSize"`          /*  每页包含的数量  */
	PageNo            int32                                                       `json:"pageNo"`            /*  当前页码  */
}

type SfsSfsRegionDetailsSfsReturnObjRegionSupportListResponse struct {
	ZoneName        string   `json:"zoneName"`        /*  可用区 Name  */
	ZoneID          string   `json:"zoneID"`          /*  可用区ID  */
	StorageType     string   `json:"storageType"`     /*  文件存储类型(capacity/performance/hdd_e)  */
	ProtocolType    []string `json:"protocolType"`    /*  文件传输协议类型。例：['nfs','cifs']  */
	RemainingStatus *bool    `json:"remainingStatus"` /*  该类型是否售罄  */
}
