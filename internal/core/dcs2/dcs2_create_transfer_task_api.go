package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CreateTransferTaskApi
/* 创建数据迁移任务。
 */type Dcs2CreateTransferTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CreateTransferTaskApi(client *core.CtyunClient) *Dcs2CreateTransferTaskApi {
	return &Dcs2CreateTransferTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/transfer/createTransferTask",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CreateTransferTaskApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CreateTransferTaskRequest) (*Dcs2CreateTransferTaskResponse, error) {
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
	var resp Dcs2CreateTransferTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateTransferTaskRequest struct {
	RegionId     string                                     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	SourceDbInfo *Dcs2CreateTransferTaskSourceDbInfoRequest `json:"sourceDbInfo"`           /*  源数据库  */
	TargetDbInfo *Dcs2CreateTransferTaskTargetDbInfoRequest `json:"targetDbInfo"`           /*  目标数据库  */
	SyncMode     int32                                      `json:"syncMode,omitempty"`     /*  同步模式<li>1： 全量同步+增量同步<li>2：全量同步  */
	ConflictMode int32                                      `json:"conflictMode,omitempty"` /*  数据冲突时的处理办法<li>1：中断迁移<li>2：跳过目标key，继续执行<li>3：覆盖目标key，继续执行  */
}

type Dcs2CreateTransferTaskSourceDbInfoRequest struct {
	SpuInstId       string `json:"spuInstId,omitempty"`   /*  实例ID  */
	IpAddr          string `json:"ipAddr,omitempty"`      /*  连接地址  */
	OriginalCluster *bool  `json:"originalCluster"`       /*  是否是原生cluster集群，输入实例id可不填，否则必填  */
	AccountName     string `json:"accountName,omitempty"` /*  数据库账号  */
	Password        string `json:"password,omitempty"`    /*  数据库密码  */
}

type Dcs2CreateTransferTaskTargetDbInfoRequest struct {
	SpuInstId       string `json:"spuInstId,omitempty"`   /*  实例ID  */
	IpAddr          string `json:"ipAddr,omitempty"`      /*  连接地址  */
	OriginalCluster *bool  `json:"originalCluster"`       /*  是否是原生cluster集群，输入实例id可不填，否则必填  */
	AccountName     string `json:"accountName,omitempty"` /*  数据库账号  */
	Password        string `json:"password,omitempty"`    /*  数据库密码  */
}

type Dcs2CreateTransferTaskResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                   `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2CreateTransferTaskReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                   `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                   `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2CreateTransferTaskReturnObjResponse struct{}
