package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsListSfsBySfsprotocolApi
/* 查询指定协议类型的⽂件系统列表，协议类型有NFS、CIFS两种
 */type SfsListSfsBySfsprotocolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListSfsBySfsprotocolApi(client *core.CtyunClient) *SfsListSfsBySfsprotocolApi {
	return &SfsListSfsBySfsprotocolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-sfs-by-sfsprotocol",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListSfsBySfsprotocolApi) Do(ctx context.Context, credential core.Credential, req *SfsListSfsBySfsprotocolRequest) (*SfsListSfsBySfsprotocolResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsProtocol", req.SfsProtocol)
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
	var resp SfsListSfsBySfsprotocolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListSfsBySfsprotocolRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	SfsProtocol string `json:"sfsProtocol,omitempty"` /*  协议类型，nfs或cifs  */
	PageSize    int32  `json:"pageSize,omitempty"`    /*  每页包含的元素个数，默认10  */
	PageNo      int32  `json:"pageNo,omitempty"`      /*  列表的分页页码，默认1  */
}

type SfsListSfsBySfsprotocolResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                    `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                    `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsListSfsBySfsprotocolReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsListSfsBySfsprotocolReturnObjResponse struct {
	Total    int32 `json:"total"`    /*  文件系统读写权限信息总数  */
	PageSize int32 `json:"pageSize"` /*  每页个数  */
	PageNo   int32 `json:"pageNo"`   /*  当前页数  */
}
