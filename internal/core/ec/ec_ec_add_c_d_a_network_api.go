package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcAddCDANetworkApi
/* 添加CDA网络实例 */
type EcEcAddCDANetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcAddCDANetworkApi(client *core.CtyunClient) *EcEcAddCDANetworkApi {
	return &EcEcAddCDANetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cda-instance/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcAddCDANetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcAddCDANetworkRequest) (*EcEcAddCDANetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcAddCDANetworkRequest
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
	var resp EcEcAddCDANetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcAddCDANetworkRequest struct {
	EcID          string    `json:"ecID"`                    /*  云间高速实例ID  */
	CgwID         string    `json:"cgwID"`                   /*  55450584-fe8f-49b0-9e09-ef915caeee10  */
	CdaID         string    `json:"cdaID"`                   /*  55450584-fe8f-49b0-9e09-ef915caeee10  */
	CdaName       string    `json:"cdaName"`                 /*  云间高速侧显示的云专线名称（建议保持和cda创建时的名称一致）  */
	CdaCidrV4List []string  `json:"cdaCidrV4List"`           /*  已选择的V4子网列表，值类型为String  */
	CdaCidrV6List []*string `json:"cdaCidrV6List,omitempty"` /*  已选择的V6子网列表，值类型为String  */
	RtbID         string    `json:"rtbID"`                   /*  路由表ID  */
	CdaInfo       string    `json:"cdaInfo"`                 /*  云专线信息，json格式的字符串  */
	Account       *string   `json:"account,omitempty"`       /*  天翼云客户邮箱  */
	DcID          *string   `json:"dcID,omitempty"`          /*  资源池ID  */
	Weights       *int32    `json:"weights,omitempty"`       /*  权重，专线默认50，无冗余实例则不传  */
	RouteLearn    *int32    `json:"routeLearn,omitempty"`    /*  路由学习开关，开启后云网关自动学习网络实例路由<br/>取值范围:<br/>1:学习<br/>0:不学习<br/>默认学习  */
	RouteSync     *int32    `json:"routeSync,omitempty"`     /*  路由同步开关，开启后云网关路由自动同步到网络实例<br/>取值范围:<br/>1:同步<br/>0:不同步<br/>默认同步  */
}

type EcEcAddCDANetworkResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcAddCDANetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcAddCDANetworkReturnObjResponse struct {
	OplogID *string `json:"oplogID"` /*  操作日志id  */
}
