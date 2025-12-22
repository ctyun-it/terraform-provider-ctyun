package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatQueryCreatePriceApi
/* 创建私网NAT询价。
 */type CtnatQueryCreatePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatQueryCreatePriceApi(client *core.CtyunClient) *CtnatQueryCreatePriceApi {
	return &CtnatQueryCreatePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/query-create-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatQueryCreatePriceApi) Do(ctx context.Context, credential core.Credential, req *CtnatQueryCreatePriceRequest) (*CtnatQueryCreatePriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatQueryCreatePriceRequest
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
	var resp CtnatQueryCreatePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatQueryCreatePriceRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  区域id  */
	VpcID       string `json:"vpcID,omitempty"`       /*  需要创建私网NAT的 VPC 的 ID  */
	SubnetID    string `json:"subnetID,omitempty"`    /*  需要创建私网NAT的 Subnet 的ID  */
	Name        string `json:"name,omitempty"`        /*  NAT 网关的名称。  */
	Spec        string `json:"spec,omitempty"`        /*  规格, small表示小型, medium表示中型, large表示大型, xlarge表示超大型  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订购类型：month（包月） / year（包年）/on_demand（按需）  */
	AzName      string `json:"azName,omitempty"`      /*  可用区名称  */
	ProjectID   string `json:"projectID,omitempty"`   /*  项目ID  */
}

type CtnatQueryCreatePriceResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                  `json:"message"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                  `json:"description"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtnatQueryCreatePriceReturnObjResponse `json:"returnObj"`   /*  业务数据  */
}

type CtnatQueryCreatePriceReturnObjResponse struct{}
