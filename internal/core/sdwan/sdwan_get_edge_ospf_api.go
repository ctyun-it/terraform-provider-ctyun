package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanGetEdgeOspfApi
/* 查找智能网关ospf信息 */
type SdwanGetEdgeOspfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeOspfApi(client *core.CtyunClient) *SdwanGetEdgeOspfApi {
	return &SdwanGetEdgeOspfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-ospf/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeOspfApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeOspfRequest) (*SdwanGetEdgeOspfResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeOspfResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeOspfRequest struct {
	EdgeID string `json:"edgeID"` /*  智能网关ID  */
}

type SdwanGetEdgeOspfResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeOspfReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeOspfReturnObjResponse struct {
	Result       []*SdwanGetEdgeOspfReturnObjResultResponse `json:"result"`       /*  查询ospf  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeOspfReturnObjResultResponse struct {
	OspfIntfIP *string `json:"ospfIntfIP"` /*  ospf接口IP  */
	OspfEnable *bool   `json:"ospfEnable"` /*  是否启用ospf  */
	AuthEnable *bool   `json:"authEnable"` /*  是否开启md5认证  */
	HelloTime  int32   `json:"helloTime"`  /*  上传Hello包之间的周期性间隔  */
	DeadTime   int32   `json:"deadTime"`   /*  等待接收Hello数据包的时间  */
	AreaID     int32   `json:"areaID"`     /*  区域ID  */
	RouterID   *string `json:"routerID"`   /*  ospf业务标识  */
	AreaType   *string `json:"areaType"`   /*  本参数表示区域类型<br/><br/>取值范围：<br/>NSSA:NSSA  */
	KeyID      *string `json:"keyID"`      /*  md5认证的key-id  */
	MdsKey     *string `json:"mdsKey"`     /*  md5认证的mds-key  */
}
