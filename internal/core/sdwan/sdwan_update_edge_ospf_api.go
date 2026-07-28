package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanUpdateEdgeOspfApi
/* 修改智能网关ospf */
type SdwanUpdateEdgeOspfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanUpdateEdgeOspfApi(client *core.CtyunClient) *SdwanUpdateEdgeOspfApi {
	return &SdwanUpdateEdgeOspfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/edge-ospf/update",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanUpdateEdgeOspfApi) Do(ctx context.Context, credential core.Credential, req *SdwanUpdateEdgeOspfRequest) (*SdwanUpdateEdgeOspfResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanUpdateEdgeOspfRequest
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
	var resp SdwanUpdateEdgeOspfResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanUpdateEdgeOspfRequest struct {
	EdgeID     string `json:"edgeID"`     /*  智能网关ID  */
	OspfIntfIP string `json:"ospfIntfIP"` /*  ospf接口IP  */
	AuthEnable string `json:"authEnable"` /*  是否开启md5认证  */
	KeyID      string `json:"keyID"`      /*  md5认证的key-id  */
	MdsKey     string `json:"mdsKey"`     /*  md5认证的mds-key  */
	HelloTime  string `json:"helloTime"`  /*  上传Hello数据包的时间  */
	DeadTime   string `json:"deadTime"`   /*  等待接收Hello数据包的时间  */
	AreaID     string `json:"areaID"`     /*  区域ID  */
	RouterID   string `json:"routerID"`   /*  ospf业务标识  */
}

type SdwanUpdateEdgeOspfResponse struct {
	StatusCode  int32     `json:"statusCode"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   *string   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   []*string `json:"returnObj"`   /*  返回参数  */
	Error       *string   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}
