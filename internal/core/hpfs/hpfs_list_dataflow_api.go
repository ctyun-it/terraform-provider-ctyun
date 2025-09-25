package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListDataflowApi
/* 查询资源池 ID 下，所有的数据流动策略详情
 */type HpfsListDataflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListDataflowApi(client *core.CtyunClient) *HpfsListDataflowApi {
	return &HpfsListDataflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-dataflow",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListDataflowApi) Do(ctx context.Context, credential core.Credential, req *HpfsListDataflowRequest) (*HpfsListDataflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.SfsUID != "" {
		ctReq.AddParam("sfsUID", req.SfsUID)
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
	var resp HpfsListDataflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListDataflowRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	AzName   string `json:"azName,omitempty"`   /*  多可用区下的可用区名字，不传为查询全部  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  数据流动策略所属并行文件UID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码，默认值为1  */
}

type HpfsListDataflowResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述  */
	Description string                             `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListDataflowReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListDataflowReturnObjResponse struct {
	DataflowList []*HpfsListDataflowReturnObjDataflowListResponse `json:"dataflowList"` /*  返回的数据流动策略列表  */
	TotalCount   int32                                            `json:"totalCount"`   /*  资源池下用户数据流动策略总数  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页码下查询回来的用户数据流动策略数  */
	PageSize     int32                                            `json:"pageSize"`     /*  每页包含的元素个数范围(1-50)  */
	PageNo       int32                                            `json:"pageNo"`       /*  列表的分页页码  */
}

type HpfsListDataflowReturnObjDataflowListResponse struct {
	RegionID            string `json:"regionID"`            /*  资源池id  */
	DataflowID          string `json:"dataflowID"`          /*  数据流动策略id  */
	SfsUID              string `json:"sfsUID"`              /*  并行文件唯一id  */
	SfsDirectory        string `json:"sfsDirectory"`        /*  并行文件下目录  */
	BucketName          string `json:"bucketName"`          /*  对象存储的桶名称  */
	BucketPrefix        string `json:"bucketPrefix"`        /*  对象存储桶的前缀  */
	AutoImport          *bool  `json:"autoImport"`          /*  是否打开自动导入  */
	AutoExport          *bool  `json:"autoExport"`          /*  是否打开自动导出  */
	ImportDataType      string `json:"importDataType"`      /*  导入的数据类型  */
	ExportDataType      string `json:"exportDataType"`      /*  导出的数据类型  */
	ImportTrigger       string `json:"importTrigger"`       /*  导入的触发条件，多个条件用英文逗号隔开  */
	ExportTrigger       string `json:"exportTrigger"`       /*  导出的触发条件，多个条件用英文逗号隔开  */
	DataflowDescription string `json:"dataflowDescription"` /*  数据流动策略的描述  */
	CreateTime          int64  `json:"createTime"`          /*  数据流动策略创建时间  */
	DataflowStatus      string `json:"dataflowStatus"`      /*  数据流动策略的状态，creating/updating/available/syncing/deleting/fail/error。creating：策略创建中；updating：策略更新中；available：策略可用（未打开自动导入导出开）；syncing：策略同步中（打开了自动导入或自动导出，数据持续流动中，即使没有数据正在流动也是同步中）；deleting：策略删除中；fail：策略异常（异常原因可见dataflowFailMsg，该状态的策略可更新、可恢复）；error：策略创建失败（异常原因可见dataflowFailMsg，该状态的策略只能删除，无法恢复，不占用配额）  */
	DataflowFailTime    int64  `json:"dataflowFailTime"`    /*  数据流动策略异常发生时间，当数据流动策略状态dataflowStatus为fail或error才有此字段  */
	DataflowFailMsg     string `json:"dataflowFailMsg"`     /*  数据流动策略异常原因，当数据流动策略状态dataflowStatus为fail或error才有此字段  */
	AzName              string `json:"azName"`              /*  多可用区下的可用区名字  */
}
