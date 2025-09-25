package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2QueryQuotaTotalAndUsedApi
/* 查询租户和用户默认可以创建的实例数和总内存的配额限制。不同的资源池配额可能不同。
 */type Dcs2QueryQuotaTotalAndUsedApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2QueryQuotaTotalAndUsedApi(client *core.CtyunClient) *Dcs2QueryQuotaTotalAndUsedApi {
	return &Dcs2QueryQuotaTotalAndUsedApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/quota/queryQuotaTotalAndUsed",
			ContentType:  "",
		},
	}
}

func (a *Dcs2QueryQuotaTotalAndUsedApi) Do(ctx context.Context, credential core.Credential, req *Dcs2QueryQuotaTotalAndUsedRequest) (*Dcs2QueryQuotaTotalAndUsedResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2QueryQuotaTotalAndUsedResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2QueryQuotaTotalAndUsedRequest struct {
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
}

type Dcs2QueryQuotaTotalAndUsedResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                       `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2QueryQuotaTotalAndUsedReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                       `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                       `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2QueryQuotaTotalAndUsedReturnObjResponse struct {
	QuotaCenterDataList []*Dcs2QueryQuotaTotalAndUsedReturnObjQuotaCenterDataListResponse `json:"quotaCenterDataList"` /*  配额列表  */
}

type Dcs2QueryQuotaTotalAndUsedReturnObjQuotaCenterDataListResponse struct {
	UserId                             string                                                                                              `json:"userId,omitempty"`                   /*  用户ID  */
	AccountId                          string                                                                                              `json:"accountId,omitempty"`                /*  账号ID  */
	TenantId                           string                                                                                              `json:"tenantId,omitempty"`                 /*  租户ID  */
	RegionId                           string                                                                                              `json:"regionId,omitempty"`                 /*  资源池ID  */
	QuotaCenterAccountResourceDataList []*Dcs2QueryQuotaTotalAndUsedReturnObjQuotaCenterDataListQuotaCenterAccountResourceDataListResponse `json:"quotaCenterAccountResourceDataList"` /*  配额信息  */
}

type Dcs2QueryQuotaTotalAndUsedReturnObjQuotaCenterDataListQuotaCenterAccountResourceDataListResponse struct {
	QuotaId    string `json:"quotaId,omitempty"`    /*  配置项编码  */
	Used       int32  `json:"used,omitempty"`       /*  已使用额度  */
	QuotaTotal int32  `json:"quotaTotal,omitempty"` /*  总配额  */
}
