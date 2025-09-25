package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsNewDataflowApi
/* 创建数据流动策略
 */type HpfsNewDataflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsNewDataflowApi(client *core.CtyunClient) *HpfsNewDataflowApi {
	return &HpfsNewDataflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/new-dataflow",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsNewDataflowApi) Do(ctx context.Context, credential core.Credential, req *HpfsNewDataflowRequest) (*HpfsNewDataflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsNewDataflowRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsNewDataflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsNewDataflowRequest struct {
	RegionID            string `json:"regionID,omitempty"`            /*  资源池 ID  */
	SfsUID              string `json:"sfsUID,omitempty"`              /*  并行文件唯一id  */
	SfsDirectory        string `json:"sfsDirectory,omitempty"`        /*  并行文件目录，目录名仅允许数字、字母、下划线、连接符、中文组成，每级目录最大长度为255字节，最大全路径长度为4096字节，如果参数为mydir/、mydir、/mydir或/mydir/，则都视为输入/mydir的目录  */
	BucketName          string `json:"bucketName,omitempty"`          /*  对象存储的桶名称  */
	BucketPrefix        string `json:"bucketPrefix,omitempty"`        /*  对象存储桶的前缀  */
	AutoImport          bool   `json:"autoImport"`                    /*  是否打开自动导入  */
	AutoExport          bool   `json:"autoExport"`                    /*  是否打开自动导出  */
	ImportDataType      string `json:"importDataType,omitempty"`      /*  导入的数据类型，data/metadata，自动导入开关打开时必填  */
	ExportDataType      string `json:"exportDataType,omitempty"`      /*  导出的数据类型，仅支持data，自动导出开关打开时必填  */
	ImportTrigger       string `json:"importTrigger,omitempty"`       /*  导入的触发条件，仅支持new（创建），自动导入开关打开时必填  */
	ExportTrigger       string `json:"exportTrigger,omitempty"`       /*  导出的触发条件，支持new（创建）/changed（新增+修改）自动导出开关打开时必填  */
	DataflowDescription string `json:"dataflowDescription,omitempty"` /*  数据流动策略的描述，最高支持128字符  */
}

type HpfsNewDataflowResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                            `json:"message"`     /*  响应描述  */
	Description string                            `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsNewDataflowReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                            `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                            `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsNewDataflowReturnObjResponse struct {
	RegionID  string                                       `json:"regionID"`  /*  资源所属资源池 ID  */
	Resources []*HpfsNewDataflowReturnObjResourcesResponse `json:"resources"` /*  资源明细  */
}

type HpfsNewDataflowReturnObjResourcesResponse struct {
	DataflowID string `json:"dataflowID"` /*  数据流动策略ID  */
	SfsUID     string `json:"sfsUID"`     /*  数据流动策略所属并行文件ID  */
}
