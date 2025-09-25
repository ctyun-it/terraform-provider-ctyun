package oceanfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// OceanfsNewSfsApi
/* 创建文件系统
 */type OceanfsNewSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOceanfsNewSfsApi(client *core.CtyunClient) *OceanfsNewSfsApi {
	return &OceanfsNewSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oceanfs/new-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *OceanfsNewSfsApi) Do(ctx context.Context, credential core.Credential, req *OceanfsNewSfsRequest) (*OceanfsNewSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OceanfsNewSfsRequest
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
	var resp OceanfsNewSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OceanfsNewSfsRequest struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	ProjectID   string `json:"projectID,omitempty"`   /*  资源所属企业项目 ID，默认为"0"  */
	SfsType     string `json:"sfsType,omitempty"`     /*  海量文件类型，massive  */
	SfsProtocol string `json:"sfsProtocol,omitempty"` /*  协议类型，nfs/cifs, nfs 适用于 linux，cifs 适用于 windows  */
	SfsName     string `json:"sfsName,omitempty"`     /*  海量文件名。单账户单资源池下，命名需唯一  */
	SfsSize     int32  `json:"sfsSize,omitempty"`     /*  大小，单位 GB，最小 100GB  */
	OnDemand    *bool  `json:"onDemand"`              /*  是否按需下单。true/false，默认为 true。如果是按实际使用量付费功能的白名单用户，无须传此参数  */
	CycleType   string `json:"cycleType,omitempty"`   /*  包周期（subscription）类型，year/month。onDemand 为 false 时，必须指定  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  包周期数。onDemand 为 false 时必须指定。周期最大长度不能超过 3 年  */
	AzName      string `json:"azName,omitempty"`      /*  多可用区资源池下，必须指定可用区。4.0资源池必填  */
	Vpc         string `json:"vpc,omitempty"`         /*  虚拟网 ID  */
	Subnet      string `json:"subnet,omitempty"`      /*  子网 ID  */
	IsVpce      *bool  `json:"isVpce"`                /*  是否创建走VPCE网络的文件系统，默认为false  */
}

type OceanfsNewSfsResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码(800为成功，900为失败/订单处理中)  */
	Message     string                            `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                            `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                            `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
	ErrorDetail *OceanfsNewSfsErrorDetailResponse `json:"errorDetail"` /*  错误明细。一般情况下，会对订单侧(bss)的海量文件订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。其他订单侧(bss)的错误，以sfs.order.procFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
}

type OceanfsNewSfsErrorDetailResponse struct{}
