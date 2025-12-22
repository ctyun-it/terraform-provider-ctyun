package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCreateEdgeOspfBackupApi
/* 增加智能网关ospf主备关系 */
type SdwanCreateEdgeOspfBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCreateEdgeOspfBackupApi(client *core.CtyunClient) *SdwanCreateEdgeOspfBackupApi {
	return &SdwanCreateEdgeOspfBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-ospf-backup/create",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCreateEdgeOspfBackupApi) Do(ctx context.Context, credential core.Credential, req *SdwanCreateEdgeOspfBackupRequest) (*SdwanCreateEdgeOspfBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCreateEdgeOspfBackupRequest
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
	var resp SdwanCreateEdgeOspfBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCreateEdgeOspfBackupRequest struct {
	MasterEdgeID   string `json:"masterEdgeID"`   /*  主智能网关ID  */
	SlaveEdgeID    string `json:"slaveEdgeID"`    /*  备智能网关ID  */
	MasterPriority string `json:"masterPriority"` /*  主优先级  */
	SlavePriority  string `json:"slavePriority"`  /*  备优先级  */
}

type SdwanCreateEdgeOspfBackupResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
