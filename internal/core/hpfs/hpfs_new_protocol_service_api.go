package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsNewProtocolServiceApi
/* 创建协议服务
 */type HpfsNewProtocolServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsNewProtocolServiceApi(client *core.CtyunClient) *HpfsNewProtocolServiceApi {
	return &HpfsNewProtocolServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/hpfs/new-protocol-service",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsNewProtocolServiceApi) Do(ctx context.Context, credential core.Credential, req *HpfsNewProtocolServiceRequest) (*HpfsNewProtocolServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*HpfsNewProtocolServiceRequest
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
	var resp HpfsNewProtocolServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsNewProtocolServiceRequest struct {
	RegionID           string `json:"regionID,omitempty"`           /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	SfsUID             string `json:"sfsUID,omitempty"`             /*  并行文件唯一ID  */
	ProtocolSpec       string `json:"protocolSpec,omitempty"`       /*  协议服务规格，目前仅支持general（通用型）  */
	ProtocolType       string `json:"protocolType,omitempty"`       /*  协议服务的协议类型，目前仅支持nfs  */
	VpcID              string `json:"vpcID,omitempty"`              /*  虚拟网 ID  */
	SubnetID           string `json:"subnetID,omitempty"`           /*  子网 ID，3.0资源池必填，4.0资源池若isVpce为true则必填  */
	IsVpce             *bool  `json:"isVpce"`                       /*  是否创建终端节点，默认false，仅4.0资源池生效  */
	IpVersion          int32  `json:"ipVersion,omitempty"`          /*  终端节点的类型，0:ipv4,1:ipv6,2:双栈，默认为0，仅isVpce为true时生效，仅4.0资源池生效  */
	ProtocolDescrption string `json:"protocolDescrption,omitempty"` /*  协议服务的描述，最高支持128字符  */
}

type HpfsNewProtocolServiceResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码(800为成功，900为处理中/失败，详见errorCode)  */
	Message     string                                   `json:"message"`     /*  响应描述  */
	Description string                                   `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsNewProtocolServiceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsNewProtocolServiceReturnObjResponse struct {
	RegionID  string                                              `json:"regionID"`  /*  资源所属资源池 ID  */
	Resources []*HpfsNewProtocolServiceReturnObjResourcesResponse `json:"resources"` /*  资源明细  */
}

type HpfsNewProtocolServiceReturnObjResourcesResponse struct {
	SfsUID            string `json:"sfsUID"`            /*  并行文件内部唯一 ID  */
	ProtocolServiceID string `json:"protocolServiceID"` /*  协议服务唯一ID  */
}
