package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanCheckSDWANCidrApi
/* 校验同网段 */
type SdwanCheckSDWANCidrApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanCheckSDWANCidrApi(client *core.CtyunClient) *SdwanCheckSDWANCidrApi {
	return &SdwanCheckSDWANCidrApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/check-cidr",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanCheckSDWANCidrApi) Do(ctx context.Context, credential core.Credential, req *SdwanCheckSDWANCidrRequest) (*SdwanCheckSDWANCidrResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanCheckSDWANCidrRequest
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
	var resp SdwanCheckSDWANCidrResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanCheckSDWANCidrRequest struct {
	EdgeID    string    `json:"edgeID"`             /*  智能网关站点ID  */
	CheckType string    `json:"checkType"`          /*  本参数表示校验类型<br/><br/>取值范围:<br/>sdwan:校验edge绑定SDWAN<br/>edge:绑定sdwan的情况下修改盒子私网网段<br/>snat:校验snat对外服务IP  */
	SdwanID   *string   `json:"sdwanID,omitempty"`  /*  SD-WAN ID ,checkType为sdwan/edge必填  */
	Ipv4Cidr  []*string `json:"ipv4Cidr,omitempty"` /*  设备ipv4私网网段，值类型为string, checkType为edge必填  */
	Ipv6Cidr  []*string `json:"ipv6Cidr,omitempty"` /*  设备ipv6私网网段，值类型为string  */
	SnatIP    *string   `json:"snatIP,omitempty"`   /*  snat对外服务IP，checkType为snat必填  */
}

type SdwanCheckSDWANCidrResponse struct {
	StatusCode  int32   `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)，默认值:800  */
	ErrorCode   *string `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	EdgeID      *string `json:"edgeID"`      /*  智能网关站点ID  */
	EdgeName    *string `json:"edgeName"`    /*  智能网关站点名称  */
	DstEdgeName *string `json:"dstEdgeName"` /*  冲突智能网关站点名称  */
	DstEdgeID   *string `json:"dstEdgeID"`   /*  冲突智能网关站点ID  */
	Status      int32   `json:"status"`      /*  本参数表示网段状态<br/><br/>取值范围:<br/>1:重复<br/>2:包含  */
	Error       *string `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
