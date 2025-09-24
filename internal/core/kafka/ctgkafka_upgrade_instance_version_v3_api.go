package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUpgradeInstanceVersionV3Api
/* 实例升级，支持小版本升级和大版本升级。
 */type CtgkafkaUpgradeInstanceVersionV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUpgradeInstanceVersionV3Api(client *core.CtyunClient) *CtgkafkaUpgradeInstanceVersionV3Api {
	return &CtgkafkaUpgradeInstanceVersionV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/upgradeInstanceVersion",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUpgradeInstanceVersionV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUpgradeInstanceVersionV3Request) (*CtgkafkaUpgradeInstanceVersionV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUpgradeInstanceVersionV3Request
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
	var resp CtgkafkaUpgradeInstanceVersionV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUpgradeInstanceVersionV3Request struct {
	RegionId      string `json:"regionId,omitempty"`      /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId    string `json:"prodInstId,omitempty"`    /*  实例ID。  */
	UpdateVersion string `json:"updateVersion,omitempty"` /*  升级目标版本号，可<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382">查询可升级实例版本</a>获取。  */
}

type CtgkafkaUpgradeInstanceVersionV3Response struct {
	StatusCode string                                             `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                             `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaUpgradeInstanceVersionV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                             `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUpgradeInstanceVersionV3ReturnObjResponse struct {
	Data *CtgkafkaUpgradeInstanceVersionV3ReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type CtgkafkaUpgradeInstanceVersionV3ReturnObjDataResponse struct{}
