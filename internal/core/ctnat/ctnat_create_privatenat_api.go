package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtnatCreatePrivatenatApi
/* 创建私网NAT
 */type CtnatCreatePrivatenatApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtnatCreatePrivatenatApi(client *core.CtyunClient) *CtnatCreatePrivatenatApi {
	return &CtnatCreatePrivatenatApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/privatenat/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtnatCreatePrivatenatApi) Do(ctx context.Context, credential core.Credential, req *CtnatCreatePrivatenatRequest) (*CtnatCreatePrivatenatResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtnatCreatePrivatenatRequest
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
	var resp CtnatCreatePrivatenatResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtnatCreatePrivatenatRequest struct {
	RegionID        string `json:"regionID,omitempty"`        /*  区域 id  */
	ProjectID       string `json:"projectID,omitempty"`       /*  项目ID  */
	VpcID           string `json:"vpcID,omitempty"`           /*  需要创建 私网NAT 网关的 VPC 的ID  */
	SubnetID        string `json:"subnetID,omitempty"`        /*  需要创建私网NBAT网关的Subnet的ID  */
	Spec            string `json:"spec,omitempty"`            /*  规格, small 表示小型, medium 表示中型, large 表示大型, xlarge 表示超大型  */
	Name            string `json:"name,omitempty"`            /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description     string `json:"description,omitempty"`     /*  <支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&\*()\_\-+= <>?:\"{}\&#124;,.\/;'[\]·~！@#￥%……&\*（） ——\-+={} &#124; 《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	ClientToken     string `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType       string `json:"cycleType,omitempty"`       /*  订购类型：month / year / on_demand  */
	CycleCount      int32  `json:"cycleCount,omitempty"`      /*  订购时长  */
	AzName          string `json:"azName,omitempty"`          /*  可用区名称  */
	AutoRenew       *bool  `json:"autoRenew"`                 /*  是否自动续约  */
	PayVoucherPrice string `json:"payVoucherPrice,omitempty"` /*  代金券金额，支持到小数点后两位，仅包周期支持代金券  */
}

type CtnatCreatePrivatenatResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码（800 为成功，900 为失败  */
	Message     string                                  `json:"message"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description string                                  `json:"description"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   string                                  `json:"errorCode"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   *CtnatCreatePrivatenatReturnObjResponse `json:"returnObj"`   /*  object  */
}

type CtnatCreatePrivatenatReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID"`        /*  订单 id。  */
	MasterOrderNO        string `json:"masterOrderNO"`        /*  订单编号, 可以为 null。  */
	MasterResourceStatus string `json:"masterResourceStatus"` /*  资源状态: started（启用） / renewed（续订） / refunded（退订） / destroyed（销毁） / failed（失败） / starting（正在启用） / changed（变配）/ expired（过期）/ unknown（未知）  */
	MasterResourceID     string `json:"masterResourceID"`     /*  可以为 null。  */
	RegionID             string `json:"regionID"`             /*  可用区 id。  */
	NatGatewayID         string `json:"natGatewayID"`         /*  nat 网关 ID  */
}
