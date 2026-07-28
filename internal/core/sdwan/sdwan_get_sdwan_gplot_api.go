package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetSdwanGplotApi
/* 查询拓扑Node信息 */
type SdwanGetSdwanGplotApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetSdwanGplotApi(client *core.CtyunClient) *SdwanGetSdwanGplotApi {
	return &SdwanGetSdwanGplotApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/gplot/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetSdwanGplotApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetSdwanGplotRequest) (*SdwanGetSdwanGplotResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("sdwanID", req.SdwanID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetSdwanGplotResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetSdwanGplotRequest struct {
	SdwanID string `json:"sdwanID"` /*  sdwan的id  */
}

type SdwanGetSdwanGplotResponse struct {
	StatusCode   int32   `json:"statusCode"`   /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode    *string `json:"errorCode"`    /*  业务细分码，为product.module.code三段式码  */
	Message      *string `json:"message"`      /*  失败时的错误描述，一般为英文描述  */
	Description  *string `json:"description"`  /*  失败时的错误描述，一般为中文描述  */
	TotalCount   int32   `json:"totalCount"`   /*  总数量  */
	CurrentCount int32   `json:"currentCount"` /*  当前页数量  */
	EdgeID       *string `json:"edgeID"`       /*  edge的ID  */
	EdgeName     *string `json:"edgeName"`     /*  edge名称  */
	Status       *string `json:"status"`       /*  本参数表示设备状态<br/><br/>取值范围:<br/>online:在线<br/>offline:下线  */
	UseType      *string `json:"useType"`      /*  本参数表示使用方式<br/><br/>取值范围:<br/>singleNode:单机<br/>activeStandby:双机  */
	x            *string `json:"x"`            /*  edge的x坐标  */
	y            *string `json:"y"`            /*  edge的y坐标  */
	Master       *string `json:"master"`       /*  edge的主SN  */
	Slave        *string `json:"slave"`        /*  edge的备SN  */
	Error        *string `json:"error"`        /*  业务细分码，为product.module.code三段式码  */
}
