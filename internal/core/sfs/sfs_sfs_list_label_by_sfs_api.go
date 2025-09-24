package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsListLabelBySfsApi
/* 查询指定文件系统绑定的标签。
 */type SfsSfsListLabelBySfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListLabelBySfsApi(client *core.CtyunClient) *SfsSfsListLabelBySfsApi {
	return &SfsSfsListLabelBySfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-label-by-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListLabelBySfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListLabelBySfsRequest) (*SfsSfsListLabelBySfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sfsUID", req.SfsUID)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsListLabelBySfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListLabelBySfsRequest struct {
	SfsUID   string `json:"sfsUID,omitempty"`   /*  文件系统实例ID  */
	RegionID string `json:"regionID,omitempty"` /*  资源池（区域）ID  */
}

type SfsSfsListLabelBySfsResponse struct {
	StatusCode int32                                  `json:"statusCode"` /*  返回状态码(800为成功，900为失败)  */
	Message    string                                 `json:"message"`    /*  参考[响应示例]  */
	ReturnObj  *SfsSfsListLabelBySfsReturnObjResponse `json:"returnObj"`  /*  参考[returnObj]  */
	ErrorCode  string                                 `json:"errorCode"`  /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error      string                                 `json:"error"`      /*  业务细分码，为Product.Module.Code三段式码大驼峰形式  */
}

type SfsSfsListLabelBySfsReturnObjResponse struct {
	LabelList []*SfsSfsListLabelBySfsReturnObjLabelListResponse `json:"labelList"` /*  文件系统绑定的标签集合  */
}

type SfsSfsListLabelBySfsReturnObjLabelListResponse struct {
	Key   string `json:"key"`   /*  标签键  */
	Value string `json:"value"` /*  标签值  */
}
