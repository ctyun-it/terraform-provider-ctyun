package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanGplotLinkApi
/* 查询拓扑LINK信息 */
type SdwanGetSdwanGplotLinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanGplotLinkApi(client *core.CtyunClient) *SdwanGetSdwanGplotLinkApi {
	return &SdwanGetSdwanGplotLinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/gplot-link/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanGplotLinkApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanGplotLinkRequest) (*SdwanGetSdwanGplotLinkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanGplotLinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanGplotLinkRequest struct {
	SdwanID string `json:"sdwanID"` /*  sdwan的id  */
}

type SdwanGetSdwanGplotLinkResponse struct {
	StatusCode     int32   `json:"statusCode"`     /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode      *string `json:"errorCode"`      /*  业务细分码，为product.module.code三段式码  */
	Message        *string `json:"message"`        /*  失败时的错误描述，一般为英文描述  */
	Description    *string `json:"description"`    /*  失败时的错误描述，一般为中文描述  */
	TotalCount     int32   `json:"totalCount"`     /*  总数量  */
	CurrentCount   int32   `json:"currentCount"`   /*  当前页总数  */
	BranchID       *string `json:"branchID"`       /*  branch的ID  */
	Source         *string `json:"source"`         /*  源edge id  */
	Target         *string `json:"target"`         /*  目的edge id  */
	Status         *string `json:"status"`         /*  本参数表示连接状态<br/><br/>取值范围:<br/>1:两端盒子都在线<br/>0:两端存在不在线盒子  */
	HasDualArmMode *string `json:"hasDualArmMode"` /*  branch的hasDualArmMode  */
	Error          *string `json:"error"`          /*  业务细分码，为product.module.code三段式码  */
}
