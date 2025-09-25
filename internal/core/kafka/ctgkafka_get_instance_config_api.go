package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaGetInstanceConfigApi
/* 获取实例配置。
 */type CtgkafkaGetInstanceConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaGetInstanceConfigApi(client *core.CtyunClient) *CtgkafkaGetInstanceConfigApi {
	return &CtgkafkaGetInstanceConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instances/getInstanceConfig",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaGetInstanceConfigApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaGetInstanceConfigRequest) (*CtgkafkaGetInstanceConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaGetInstanceConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaGetInstanceConfigRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaGetInstanceConfigResponse struct {
	StatusCode string                                      `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                      `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaGetInstanceConfigReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaGetInstanceConfigReturnObjResponse struct {
	Data []*CtgkafkaGetInstanceConfigReturnObjDataResponse `json:"data"` /*  返回数据  */
}

type CtgkafkaGetInstanceConfigReturnObjDataResponse struct {
	Name          string `json:"name,omitempty"`          /*  配置名称  */
	Value         string `json:"value,omitempty"`         /*  当前配置值  */
	Valid_values  string `json:"valid_values,omitempty"`  /*  配置有效值。  */
	Default_value string `json:"default_value,omitempty"` /*  配置默认值。  */
	VarType       string `json:"varType,omitempty"`       /*  值类型。  */
	Config_type   string `json:"config_type,omitempty"`   /*  配置类型：<br><li>static：静态配置<br><li>dynamic：动态配置 <br>说明：静态配置修改后需重启实例方可生效，动态配置无需重启实例。  */
	Desc          string `json:"desc,omitempty"`          /*  配置说明。  */
}
