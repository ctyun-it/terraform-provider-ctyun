package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsGetEcsFlavorsApi
/* 查询资源池虚机规格信息
 */type CtecsGetEcsFlavorsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsGetEcsFlavorsApi(client *core.CtyunClient) *CtecsGetEcsFlavorsApi {
	return &CtecsGetEcsFlavorsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/common/get-ecs-flavors",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsGetEcsFlavorsApi) Do(ctx context.Context, credential core.Credential, req *CtecsGetEcsFlavorsRequest) (*CtecsGetEcsFlavorsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.Series != "" {
		ctReq.AddParam("series", req.Series)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsGetEcsFlavorsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsGetEcsFlavorsRequest struct {
	RegionID string /*  资源池ID  */
	AzName   string /*  多az可用区名称（4.0场景）  */
	Series   string /*  系列  */
}

type CtecsGetEcsFlavorsResponse struct {
	StatusCode  int32                                `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                               `json:"errorCode,omitempty"`   /*  具体错误码标志  */
	Message     string                               `json:"message,omitempty"`     /*  失败时的错误信息  */
	Description string                               `json:"description,omitempty"` /*  失败时的错误描述  */
	ReturnObj   *CtecsGetEcsFlavorsReturnObjResponse `json:"returnObj"`             /*  返回对象，成功时返回的数据  */
	Error       string                               `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsGetEcsFlavorsReturnObjResponse struct {
	TotalCount int32                                         `json:"totalCount,omitempty"` /*  总条数  */
	Results    []*CtecsGetEcsFlavorsReturnObjResultsResponse `json:"results"`              /*  规格列表  */
}

type CtecsGetEcsFlavorsReturnObjResultsResponse struct {
	FlavorID      string   `json:"flavorID,omitempty"`     /*  规格id  */
	SpecName      string   `json:"specName,omitempty"`     /*  规格名称  */
	FlavorType    string   `json:"flavorType,omitempty"`   /*  规格类型  */
	FlavorName    string   `json:"flavorName,omitempty"`   /*  规格类型名称  */
	CpuNum        int32    `json:"cpuNum,omitempty"`       /*  cpu核数  */
	MemSize       int32    `json:"memSize,omitempty"`      /*  内存大小  */
	MultiQueue    int32    `json:"multiQueue,omitempty"`   /*  网卡多队列数  */
	Pps           int32    `json:"pps,omitempty"`          /*  网络最大收发包能力  (万PPS)  */
	BandwidthBase float64  `json:"bandwidthBase"`          /*  基准带宽 (Gbps)  */
	BandwidthMax  float64  `json:"bandwidthMax"`           /*  最大带宽 (Gbps)  */
	CpuArch       string   `json:"cpuArch,omitempty"`      /*  cpu架构 （x86架构、arm架构）  */
	Series        string   `json:"series,omitempty"`       /*  系列  */
	AzList        []string `json:"azList"`                 /*  支持的az名称列表（4.0场景）(未传azName情况)  */
	NicCount      int32    `json:"nicCount,omitempty"`     /*  当前规格主机最大可挂载网卡数 (cpu类型的会返回，其他类型没有，且非必返回字段)  */
	CtLimitCount  int32    `json:"ctLimitCount,omitempty"` /*  最大连接数 (cpu类型的会返回，其他类型没有，且非必返回字段)  */
}
