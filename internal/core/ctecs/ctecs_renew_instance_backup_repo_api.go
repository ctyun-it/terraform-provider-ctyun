package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsRenewInstanceBackupRepoApi
/* 续订云主机备份存储库。<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />&emsp;&emsp;计费模式：确认云主机备份存储库的计费模式，详细查看<a href="https://www.ctyun.cn/document/10051003/10100892">计费模式</a><br />&emsp;&emsp;代金券：只支持预付费用户抵扣包周期资源的金额，且不可超过包周期资源的金额
 */type CtecsRenewInstanceBackupRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsRenewInstanceBackupRepoApi(client *core.CtyunClient) *CtecsRenewInstanceBackupRepoApi {
	return &CtecsRenewInstanceBackupRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-repo/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsRenewInstanceBackupRepoApi) Do(ctx context.Context, credential core.Credential, req *CtecsRenewInstanceBackupRepoRequest) (*CtecsRenewInstanceBackupRepoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsRenewInstanceBackupRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsRenewInstanceBackupRepoRequest struct {
	RegionID        string  `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>    */
	RepositoryID    string  `json:"repositoryID,omitempty"` /*  云主机备份存储库ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033742">产品定义-存储库</a>来了解存储库<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6909&data=87&isNormal=1&vid=81">查询存储库列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6910&data=87&isNormal=1&vid=81">创建存储库</a>  */
	CycleCount      int32   `json:"cycleCount,omitempty"`   /*  订购时长，该参数需要与cycleType一同使用<br />注：最长订购周期为60个月（5年）  */
	CycleType       string  `json:"cycleType,omitempty"`    /*  订购周期类型 ，取值范围: <br />MONTH表示按月订购,<br />YEAR表示按年订购  */
	ClientToken     string  `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，则代表为同一个请求。保留时间为24小时  */
	PayVoucherPrice float32 `json:"payVoucherPrice"`        /*  代金券，满足以下规则：两位小数，不足两位自动补0，超过两位小数无效；不可为负数；字段为0时表示不使用代金券  */
}

type CtecsRenewInstanceBackupRepoResponse struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）。  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsRenewInstanceBackupRepoReturnObjResponse `json:"returnObj"`             /*  返回参数。  */
}

type CtecsRenewInstanceBackupRepoReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  主订单ID。调用方在拿到masterOrderID之后，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单号  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID  */
}
