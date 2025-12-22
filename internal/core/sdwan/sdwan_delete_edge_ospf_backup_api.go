package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanDeleteEdgeOspfBackupApi
/* 删除智能网关ospf主备关系 */
type SdwanDeleteEdgeOspfBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanDeleteEdgeOspfBackupApi(client *core.CtyunClient) *SdwanDeleteEdgeOspfBackupApi {
	return &SdwanDeleteEdgeOspfBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-ospf-backup/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanDeleteEdgeOspfBackupApi) Do(ctx context.Context, credential core.Credential, req *SdwanDeleteEdgeOspfBackupRequest) (*SdwanDeleteEdgeOspfBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanDeleteEdgeOspfBackupRequest
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
	var resp SdwanDeleteEdgeOspfBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanDeleteEdgeOspfBackupRequest struct {
	RelationIDList []string `json:"RelationIDList"` /*  主备关系ID列表  ，值类型为string  */
}

type SdwanDeleteEdgeOspfBackupResponse struct {
	StatusCode  int32                                         `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                       `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                       `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                       `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*SdwanDeleteEdgeOspfBackupReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                       `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanDeleteEdgeOspfBackupReturnObjResponse struct {
	OperationID *string `json:"operationID"` /*  操作日志Id  */
}
