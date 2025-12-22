package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ApplyTemplateToInstanceApi
/* 应用参数模板到对应实例。
 */type Dcs2ApplyTemplateToInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ApplyTemplateToInstanceApi(client *core.CtyunClient) *Dcs2ApplyTemplateToInstanceApi {
	return &Dcs2ApplyTemplateToInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/redisTemplate/applyTemplateToInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ApplyTemplateToInstanceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ApplyTemplateToInstanceRequest) (*Dcs2ApplyTemplateToInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Dcs2ApplyTemplateToInstanceRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ApplyTemplateToInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ApplyTemplateToInstanceRequest struct {
	RegionId    string   `json:"regionId,omitempty"`   /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	TemplateId  string   `json:"templateId,omitempty"` /*  模板ID。 可调用 <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=15104&data=81&isNormal=1&vid=270">查询参数模板列表</a> 接口，使用Template表id字段。  */
	ProdInstIds []string `json:"prodInstIds"`          /*  实例ID列表，多个实例ID用英文逗号分隔。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
}

type Dcs2ApplyTemplateToInstanceResponse struct {
	StatusCode int32                                         `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                        `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2ApplyTemplateToInstanceReturnObjResponse `json:"returnObj"`  /*  无返回数据，空对象。  */
	RequestId  string                                        `json:"requestId"`  /*  请求 ID。  */
	Code       string                                        `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                        `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2ApplyTemplateToInstanceReturnObjResponse struct{}
