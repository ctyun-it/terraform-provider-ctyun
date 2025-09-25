package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupRenewRepoApi
/* 续订云硬盘备份存储库，该接口会涉及计费<br />
 */ /* <b>准备工作：</b><br />
 */ /* 计费模式：确认续订存储库的计费模式，详细查看<a href="https://www.ctyun.cn/document/10026730/10030877">计费模式</a><br />
 */type EbsbackupRenewRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupRenewRepoApi(client *core.CtyunClient) *EbsbackupRenewRepoApi {
	return &EbsbackupRenewRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/repo/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupRenewRepoApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupRenewRepoRequest) (*EbsbackupRenewRepoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupRenewRepoRequest
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
	var resp EbsbackupRenewRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupRenewRepoRequest struct {
	ClientToken  string `json:"clientToken,omitempty"`  /*  用于保证订单幂等性。要求单个云平台账户内唯一。使用同一个ClientToken值，其他请求参数相同时，则代表为同一个请求  */
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	RepositoryID string `json:"repositoryID,omitempty"` /*  云硬盘备份存储库ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10039480">查询存储库列表</a>获取  */
	CycleType    string `json:"cycleType,omitempty"`    /*  本参数表示订购周期类型 ，取值范围：<br />MONTH：按月<br />YEAR：按年<br />最长订购周期为3年  */
	CycleCount   int32  `json:"cycleCount,omitempty"`   /*  订购时长，与cycleType配合，cycleType为Month时，单位为月，cycleType为YEAR时，单位为年  */
}

type EbsbackupRenewRepoResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                               `json:"message"`     /*  错误信息的英文描述  */
	Description string                               `json:"description"` /*  错误信息的本地化描述（中文）  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	ReturnObj   *EbsbackupRenewRepoReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见returnObj对象结构  */
	Error       string                               `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupRenewRepoReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO"` /*  订单号  */
	RegionID      string `json:"regionID"`      /*  资源所属资源池ID  */
}
