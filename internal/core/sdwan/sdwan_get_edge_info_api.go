package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SdwanGetEdgeInfoApi
/* 查找智能网关信息 */
type SdwanGetEdgeInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanGetEdgeInfoApi(client *core.CtyunClient) *SdwanGetEdgeInfoApi {
	return &SdwanGetEdgeInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sdwan/edge-info/list",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanGetEdgeInfoApi) Do(ctx context.Context, credential core.Credential, req *SdwanGetEdgeInfoRequest) (*SdwanGetEdgeInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.EdgeID != nil && *req.EdgeID != "" {
		ctReq.AddParam("edgeID", *req.EdgeID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SdwanGetEdgeInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanGetEdgeInfoRequest struct {
	EdgeID   *string `json:"edgeID,omitempty"` /*  智能网关ID  */
	PageNo   int32   `json:"pageNo"`           /*  页数  */
	PageSize int32   `json:"pageSize"`         /*  页大小  */
	Search   *string `json:"search,omitempty"` /*  模糊查询，可以查resourceId  */
}

type SdwanGetEdgeInfoResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanGetEdgeInfoReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type SdwanGetEdgeInfoReturnObjResponse struct {
	Result       []*SdwanGetEdgeInfoReturnObjResultResponse `json:"result"`       /*  查询edge 信息  */
	TotalCount   int32                                      `json:"totalCount"`   /*  总数量  */
	CurrentCount int32                                      `json:"currentCount"` /*  当前页数量  */
}

type SdwanGetEdgeInfoReturnObjResultResponse struct {
	EdgeID           *string                                                  `json:"edgeID"`           /*  智能网关ID  */
	EdgeName         *string                                                  `json:"edgeName"`         /*  设备名称  */
	Status           *string                                                  `json:"status"`           /*  本参数表示设备状态<br/><br/>取值范围:<br/>online:在线<br/>offline:下线  */
	EdgeCIDR         *string                                                  `json:"edgeCIDR"`         /*  edge子网  */
	UseType          *string                                                  `json:"useType"`          /*  本参数表示使用方式<br/><br/>取值范围:<br/>singleNode:单机<br/>activeStandby:双机  */
	DeployMode       *string                                                  `json:"deployMode"`       /*  本参数表示接入方式<br/><br/>取值范围:<br/>inline-mode:串接<br/>dual-arm-mode:双臂旁挂<br/>single-arm-mode:单臂旁挂  */
	SerialNumberDict *SdwanGetEdgeInfoReturnObjResultSerialNumberDictResponse `json:"serialNumberDict"` /*  智能网关SN信息  */
	BizStatus        *string                                                  `json:"bizStatus"`        /*  本参数表示业务状态<br/><br/>取值范围：<br/>lock:已锁定<br/>processing:操作中<br/>ordered:已下单<br/>active:已激活  */
	ServiceList      []*SdwanGetEdgeInfoReturnObjResultServiceListResponse    `json:"serviceList"`      /*  增值服务信息  */
	CloudService     *SdwanGetEdgeInfoReturnObjResultCloudServiceResponse     `json:"cloudService"`     /*  所属云服务  */
	TotalSize        *string                                                  `json:"totalSize"`        /*  总带宽  */
	ExpireTime       *string                                                  `json:"expireTime"`       /*  到期时间  */
	Description      *string                                                  `json:"description"`      /*  edge描述  */
}

type SdwanGetEdgeInfoReturnObjResultSerialNumberDictResponse struct {
	MasterSN *string `json:"masterSN"` /*  主edge的sn  */
	SlaveSN  *string `json:"slaveSN"`  /*  备edge的sn  */
}

type SdwanGetEdgeInfoReturnObjResultServiceListResponse struct {
	ServiceName *string `json:"serviceName"` /*  本参数表示增值服务名称<br/>取值范围：<br/>5G<br/>WIFI<br/>LTE<br/>extOnceOutlay<br/>extOncService<br/>extOutlayService  */
	CreateTime  *string `json:"createTime"`  /*  创建时间  */
}

type SdwanGetEdgeInfoReturnObjResultCloudServiceResponse struct {
	SdwanID   *string `json:"sdwanID"`   /*  sdwan id  */
	SdwanName *string `json:"sdwanName"` /*  sdwan名称  */
}
