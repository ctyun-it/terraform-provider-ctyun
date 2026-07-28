package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanGplotDetailApi
/* 查询拓扑节点详情信息 */
type SdwanGetSdwanGplotDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanGplotDetailApi(client *core.CtyunClient) *SdwanGetSdwanGplotDetailApi {
	return &SdwanGetSdwanGplotDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/gplot-detail/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanGplotDetailApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanGplotDetailRequest) (*SdwanGetSdwanGplotDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanGplotDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanGplotDetailRequest struct {
	SdwanID string `json:"sdwanID"` /*  sdwan的id  */
	EdgeID  string `json:"edgeID"`  /*  sdwan的id  */
}

type SdwanGetSdwanGplotDetailResponse struct {
	StatusCode   int32     `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode    *string   `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string   `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string   `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32     `json:"totalCount"`   /*  总数量  */
	CurrentCount int32     `json:"currentCount"` /*  页码  */
	EdgeName     *string   `json:"edgeName"`     /*  edge名称  */
	Status       *string   `json:"status"`       /*  本参数表示设备状态<br/><br/>取值范围:<br/>online:在线<br/>offline:下线  */
	IpPrefix     []*string `json:"ipPrefix"`     /*  edge子网  ，值类型为string  */
	LinkInfo     []*string `json:"linkInfo"`     /*  edge端口  ，值类型为string  */
	ActiveLink   []*string `json:"activeLink"`   /*  活跃edge端口  ，值类型为string  */
	Error        *string   `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
