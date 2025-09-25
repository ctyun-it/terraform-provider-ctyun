package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ModifyInstanceStrategyApi
/* 修改实例的淘汰策略配置，支持volatile-lru,allkeys-lru,volatile-lfu,allkeys-lfu,volatile-random,allkeys-random,volatile-ttl,noeviction。执行热key分析任务前应修改为volatile-lfu淘汰策略。
 */type Dcs2ModifyInstanceStrategyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ModifyInstanceStrategyApi(client *core.CtyunClient) *Dcs2ModifyInstanceStrategyApi {
	return &Dcs2ModifyInstanceStrategyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/keyAnalysisMgrServant/modifyInstanceStrategy",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ModifyInstanceStrategyApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ModifyInstanceStrategyRequest) (*Dcs2ModifyInstanceStrategyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ModifyInstanceStrategyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ModifyInstanceStrategyRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Policy     string `json:"policy,omitempty"`     /*  内存淘汰策略<li>volatile-lru：对设置了过期时间的 key 使用 LRU 淘汰<li>allkeys-lru：对所有 key 使用 LRU 淘汰<li>volatile-lfu：对设置了过期时间的 key 使用 LFU 淘汰<li>allkeys-lfu：对所有 key 使用 LFU 淘汰<li>volatile-random：对设置了过期时间的 key 随机淘汰<li>allkeys-random：对所有 key 随机淘汰<li>volatile-ttl：优先淘汰剩余时间短的 key<li>noeviction：不淘汰，内存满时返回错误（默认行为）  */
}

type Dcs2ModifyInstanceStrategyResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                       `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ModifyInstanceStrategyReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                       `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                       `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ModifyInstanceStrategyReturnObjResponse struct{}
