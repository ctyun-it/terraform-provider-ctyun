package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupCreateRepoApi
/* 创建云硬盘备份存储库，该接口会涉及计费<br />
 */ /* <b>准备工作：</b><br />
 */ /* 构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />
 */ /* 认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */ /* 计费模式：确认创建云硬盘备份存储库的计费模式，详细查看<a href="https://www.ctyun.cn/document/10026730/10030877">计费模式</a><br />
 */ /* 地域选择：选择创建云硬盘备份存储库的资源池，详细查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a><br />
 */ /* 产品选型：创建云硬盘备份存储库前，请先阅读<a href="https://www.ctyun.cn/document/10026752/10037454">入门流程</a>了解云硬盘备份的基本信息，以及操作步骤
 */type EbsbackupCreateRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupCreateRepoApi(client *core.CtyunClient) *EbsbackupCreateRepoApi {
	return &EbsbackupCreateRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/repo/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupCreateRepoApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupCreateRepoRequest) (*EbsbackupCreateRepoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupCreateRepoRequest
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
	var resp EbsbackupCreateRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupCreateRepoRequest struct {
	ClientToken    string `json:"clientToken,omitempty"`    /*  用于保证订单幂等性。要求单个云平台账户内唯一。使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求。  */
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  云硬盘备份存储库名称，长度为 2~32 个字符，只能由数字、字母、- 组成，不能以数字、- 开头，且不能以 - 结尾。  */
	Size           int32  `json:"size,omitempty"`           /*  云硬盘备份存储库容量，单位为GB，取值范围：100-1024000，默认为100。  */
	OnDemand       string `json:"onDemand,omitempty"`       /*  是否按需下单，取值范围：
	true：是
	false：否
	默认为false。  */
	CycleType string `json:"cycleType,omitempty"` /*  表示订购周期类型 ，取值范围：
	MONTH：按月
	YEAR：按年
	onDemand为false时，必须指定该参数。  */
	CycleCount int32 `json:"cycleCount,omitempty"` /*  订购时长，最大为3年（36个月），与cycleType配合使用。
	cycleType为MONTH时，单位为月，cycleType为YEAR时，单位为年。  */
	AutoRenewStatus int32 `json:"autoRenewStatus,omitempty"` /*  本参数表示是否自动续订 ，只有onDemand为false时生效。取值范围：
	0：不续费
	1：自动续费
	默认不自动续费，如果选择自动续费：
	按月购买：自动续订周期为3个月。
	按年购买：自动续订周期为1年。  */
	ProjectID string `json:"projectID,omitempty"` /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupCreateRepoResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                `json:"message"`     /*  错误信息的英文描述。  */
	Description string                                `json:"description"` /*  错误信息的本地化描述（中文）。  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务错误细分码，发生错误时返回。  */
	ReturnObj   *EbsbackupCreateRepoReturnObjResponse `json:"returnObj"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupCreateRepoReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单号，可根据订单号查询具体资源ID。  */
	RegionID      string `json:"regionID"`      /*  资源所属资源池ID。  */
}
