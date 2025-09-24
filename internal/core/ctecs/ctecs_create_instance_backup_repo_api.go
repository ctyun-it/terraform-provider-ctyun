package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateInstanceBackupRepoApi
/* 创建云主机备份存储库。<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />&emsp;&emsp;计费模式：确认云主机备份存储库的计费模式，详细查看<a href="https://www.ctyun.cn/document/10051003/10100892">计费模式</a><br /><b>注意事项：</b><br/>&emsp;&emsp;代金券：只支持预付费用户抵扣包周期资源的金额，且不可超过包周期资源的金额
 */type CtecsCreateInstanceBackupRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateInstanceBackupRepoApi(client *core.CtyunClient) *CtecsCreateInstanceBackupRepoApi {
	return &CtecsCreateInstanceBackupRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-repo/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateInstanceBackupRepoApi) Do(ctx context.Context, credential core.Credential, req *CtecsCreateInstanceBackupRepoRequest) (*CtecsCreateInstanceBackupRepoResponse, error) {
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
	var resp CtecsCreateInstanceBackupRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateInstanceBackupRepoRequest struct {
	RegionID        string  `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>    */
	ProjectID       string  `json:"projectID,omitempty"`       /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
	RepositoryName  string  `json:"repositoryName,omitempty"`  /*  云主机备份存储库名称，满足以下规则：长度为2-63字符，头尾不支持输入空格  */
	CycleCount      int32   `json:"cycleCount,omitempty"`      /*  订购时长，该参数需要与cycleType一同使用<br />注：最长订购周期为60个月（5年）  */
	CycleType       string  `json:"cycleType,omitempty"`       /*  订购周期类型，取值范围：<br />MONTH：按月<br />YEAR：按年<br />最长订购周期为5年  */
	ClientToken     string  `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，则代表为同一个请求。保留时间为24小时  */
	Size            int32   `json:"size,omitempty"`            /*  云主机备份存储库容量，单位GB，取值范围：[100-1024000]，默认值100  */
	AutoRenewStatus int32   `json:"autoRenewStatus,omitempty"` /*  是否自动续订，取值范围：<br />0（不续费），<br />1（自动续费），<br />注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年  */
	PayVoucherPrice float32 `json:"payVoucherPrice"`           /*  代金券，满足以下规则：两位小数，不足两位自动补0，超过两位小数无效；不可为负数；<br />注：字段为0时表示不使用代金券，默认不使用代金券。  */
}

type CtecsCreateInstanceBackupRepoResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）。  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsCreateInstanceBackupRepoReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsCreateInstanceBackupRepoReturnObjResponse struct {
	MasterOrderID    string `json:"masterOrderID,omitempty"`    /*  主订单ID。调用方在拿到masterOrderID之后，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO    string `json:"masterOrderNO,omitempty"`    /*  订单号  */
	RegionID         string `json:"regionID,omitempty"`         /*  资源池ID  */
	MasterResourceID string `json:"masterResourceID,omitempty"` /*  订单主资源ID  */
}
