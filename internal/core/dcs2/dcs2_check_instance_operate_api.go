package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CheckInstanceOperateApi
/* 查询分布式缓存Redis实例是否可以扩容。
 */type Dcs2CheckInstanceOperateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CheckInstanceOperateApi(client *core.CtyunClient) *Dcs2CheckInstanceOperateApi {
	return &Dcs2CheckInstanceOperateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/check/checkInstanceOperate",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CheckInstanceOperateApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CheckInstanceOperateRequest) (*Dcs2CheckInstanceOperateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2CheckInstanceOperateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CheckInstanceOperateRequest struct {
	RegionId     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	Operate      string `json:"operate,omitempty"`      /*  操作类型，可选值：upgrade（扩容）  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例ID  */
	ShardMemSize int32  `json:"shardMemSize,omitempty"` /*  目标实例的分片规格（单位：G）<br>目标实例为基础版与增强版时（OriginalMultipleReadLvs除外）必填<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 使用表SeriesInfo中shardMemSizeItems(单分片内存可选值)，若shardMemSizeItems为空则无需填写  */
	ShardCount   int32  `json:"shardCount,omitempty"`   /*  目标实例的分片数量。<li>DirectClusterSingle: 3-256<li>DirectCluster: 3-256<li>ClusterOriginalProxy: 3-64<br><br>其他实例类型无需填写此参数  */
	MemSize      int32  `json:"memSize,omitempty"`      /*  目标实例的总容量（单位：G）<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7726&isNormal=1&vid=270">资源池可创建规格</a> 计算方式：单分片内存shardMemSize × 分片数量shardCount 或 使用表SeriesInfo中memSizeItems(内存可选值)  */
	CopiesCount  int32  `json:"copiesCount,omitempty"`  /*  目标实例的副本数；目标实例为OriginalMultipleReadLvs时必填，取值范围为[2,6]  */
}

type Dcs2CheckInstanceOperateResponse struct {
	StatusCode int32                                      `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                     `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2CheckInstanceOperateReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                     `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                     `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2CheckInstanceOperateReturnObjResponse struct {
	IsSupport        *bool  `json:"isSupport"`                  /*  是否支持当前操作  */
	NotSupportReason string `json:"notSupportReason,omitempty"` /*  不支持的原因  */
}
