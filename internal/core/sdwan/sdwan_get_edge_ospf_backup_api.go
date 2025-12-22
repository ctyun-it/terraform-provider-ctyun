package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeOspfBackupApi
/* 查找智能网关ospf主备关系 */
type SdwanGetEdgeOspfBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeOspfBackupApi(client *core.CtyunClient) *SdwanGetEdgeOspfBackupApi {
	return &SdwanGetEdgeOspfBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-ospf-backup/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeOspfBackupApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeOspfBackupRequest) (*SdwanGetEdgeOspfBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeOspfBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeOspfBackupRequest struct {
	PageNo   int32 `json:"pageNo"`   /*  页数  */
	PageSize int32 `json:"pageSize"` /*  页大小  */
}

type SdwanGetEdgeOspfBackupResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeOspfBackupReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeOspfBackupReturnObjResponse struct {
	Result       []*SdwanGetEdgeOspfBackupReturnObjResultResponse `json:"result"`       /*  查询ospf备份信息  */
	TotalCount   int32                                            `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页的数量  */
}

type SdwanGetEdgeOspfBackupReturnObjResultResponse struct {
	RelationID     *string `json:"relationID"`     /*  主备关系ID  */
	MasterEdgeName *string `json:"masterEdgeName"` /*  主智能网关名称  */
	MasterPriority int32   `json:"masterPriority"` /*  主优先级  */
	SlaveEdgeName  *string `json:"slaveEdgeName"`  /*  备智能网关名称  */
	SlavePriority  int32   `json:"slavePriority"`  /*  备优先级  */
}
