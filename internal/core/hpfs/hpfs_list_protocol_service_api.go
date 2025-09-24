package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListProtocolServiceApi
/* 查询协议服务列表
 */type HpfsListProtocolServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListProtocolServiceApi(client *core.CtyunClient) *HpfsListProtocolServiceApi {
	return &HpfsListProtocolServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-protocol-service",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListProtocolServiceApi) Do(ctx context.Context, credential core.Credential, req *HpfsListProtocolServiceRequest) (*HpfsListProtocolServiceResponse, error) {
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
	if req.ProtocolServiceStatus != "" {
		ctReq.AddParam("protocolServiceStatus", req.ProtocolServiceStatus)
	}
	if req.ProtocolSpec != "" {
		ctReq.AddParam("protocolSpec", req.ProtocolSpec)
	}
	if req.ProtocolType != "" {
		ctReq.AddParam("protocolType", req.ProtocolType)
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
	var resp HpfsListProtocolServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListProtocolServiceRequest struct {
	RegionID              string `json:"regionID,omitempty"`              /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	AzName                string `json:"azName,omitempty"`                /*  多可用区下的可用区名字，不传为查询全部  */
	SfsUID                string `json:"sfsUID,omitempty"`                /*  并行文件唯一id  */
	ProtocolServiceStatus string `json:"protocolServiceStatus,omitempty"` /*  协议状态  */
	ProtocolSpec          string `json:"protocolSpec,omitempty"`          /*  协议规格  */
	ProtocolType          string `json:"protocolType,omitempty"`          /*  协议类型  */
	PageSize              int32  `json:"pageSize,omitempty"`              /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo                int32  `json:"pageNo,omitempty"`                /*  列表的分页页码，默认值为1  */
}

type HpfsListProtocolServiceResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800为成功，900为处理中/失败，详见errorCode)  */
	Message     string                                    `json:"message"`     /*  响应描述  */
	Description string                                    `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListProtocolServiceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListProtocolServiceReturnObjResponse struct {
	TotalCount   int32 `json:"totalCount"`   /*  指定条件下协议服务总数  */
	CurrentCount int32 `json:"currentCount"` /*  当前页码下查询回来的协议服务数  */
	PageSize     int32 `json:"pageSize"`     /*  每页包含的元素个数范围(1-50)  */
	PageNo       int32 `json:"pageNo"`       /*  列表的分页页码  */
}
