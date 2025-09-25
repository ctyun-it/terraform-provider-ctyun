package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsInfoDataflowApi
/* 根据dataflowID查询指定资源池下数据流动策略信息
 */type HpfsInfoDataflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsInfoDataflowApi(client *core.CtyunClient) *HpfsInfoDataflowApi {
	return &HpfsInfoDataflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/info-dataflow",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsInfoDataflowApi) Do(ctx context.Context, credential core.Credential, req *HpfsInfoDataflowRequest) (*HpfsInfoDataflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("dataflowID", req.DataflowID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsInfoDataflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsInfoDataflowRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池 ID  */
	DataflowID string `json:"dataflowID,omitempty"` /*  数据流动策略ID  */
}

type HpfsInfoDataflowResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述  */
	Description string                             `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsInfoDataflowReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsInfoDataflowReturnObjResponse struct {
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
