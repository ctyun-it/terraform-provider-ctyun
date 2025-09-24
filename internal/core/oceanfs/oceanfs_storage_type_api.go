package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// OceanfsStorageTypeApi
/* 用于展示海量文件下列信息：</br>
 */ /* （1）所支持的协议类型</br>
 */ /* （2）资源池所支持的所有文件类型</br>
 */ /* （3）资源池下的az列表</br>
 */ /* （4）资源池的售罄情况</br>
 */type OceanfsStorageTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsStorageTypeApi(client *core.CtyunClient) *OceanfsStorageTypeApi {
	return &OceanfsStorageTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oceanfs/region/storagetype",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsStorageTypeApi) Do(ctx context.Context, credential core.Credential, req *OceanfsStorageTypeRequest) (*OceanfsStorageTypeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
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
	var resp OceanfsStorageTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsStorageTypeRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的数量  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  当前页码  */
}

type OceanfsStorageTypeResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
