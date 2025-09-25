package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsInfoByNameSfsApi
/* 根据海量文件名称和资源池ID，查询文件系统详情
 */type OceanfsInfoByNameSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsInfoByNameSfsApi(client *core.CtyunClient) *OceanfsInfoByNameSfsApi {
	return &OceanfsInfoByNameSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oceanfs/info-by-name-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsInfoByNameSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsInfoByNameSfsRequest) (*OceanfsInfoByNameSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sfsName", req.SfsName)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp OceanfsInfoByNameSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsInfoByNameSfsRequest struct {
	SfsName  string `json:"sfsName,omitempty"`  /*  海量文件名称  */
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
}

type OceanfsInfoByNameSfsResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
