package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeDeviceInfoApi
/* 查找智能网关设备信息 */
type SdwanGetEdgeDeviceInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeDeviceInfoApi(client *core.CtyunClient) *SdwanGetEdgeDeviceInfoApi {
	return &SdwanGetEdgeDeviceInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-device/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeDeviceInfoApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeDeviceInfoRequest) (*SdwanGetEdgeDeviceInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("edgeID", req.EdgeID)
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
	var resp SdwanGetEdgeDeviceInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeDeviceInfoRequest struct {
	EdgeID   string `json:"edgeID"`   /*  智能网关ID  */
	PageNo   int32  `json:"pageNo"`   /*  页数  */
	PageSize int32  `json:"pageSize"` /*  页大小  */
}

type SdwanGetEdgeDeviceInfoResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeDeviceInfoReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeDeviceInfoReturnObjResponse struct {
	Result       []*SdwanGetEdgeDeviceInfoReturnObjResultResponse `json:"result"`       /*  查询edge 设备  */
	TotalCount   int32                                            `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                            `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeDeviceInfoReturnObjResultResponse struct {
	EdgeName        *string `json:"edgeName"`        /*  设备名称  */
	DeviceModel     *string `json:"deviceModel"`     /*  本参数表示设备型号<br/><br/>取值范围：<br/>low:基础版<br/>medium:标准版<br/>high:豪华版<br/>vcpe:虚拟智能网关  */
	Status          *string `json:"status"`          /*  本参数表示设备状态<br/><br/>取值范围：<br/>online:在线<br/>offline:下线  */
	CurrentVersion  *string `json:"currentVersion"`  /*  软件版本号  */
	HardwareVersion *string `json:"hardwareVersion"` /*  硬件版本号  */
}
