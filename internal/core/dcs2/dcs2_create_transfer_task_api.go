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
	_, err := ctReq.WriteJson(struct {
		*Dcs2CreateTransferTaskRequest
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
	var resp Dcs2CreateTransferTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateTransferTaskRequest struct {
	RegionId     string                                     `json:"regionId,omitempty"`     /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	SourceDbInfo *Dcs2CreateTransferTaskSourceDbInfoRequest `json:"sourceDbInfo"`           /*  源数据库。  */
	TargetDbInfo *Dcs2CreateTransferTaskTargetDbInfoRequest `json:"targetDbInfo"`           /*  目标数据库。  */
	SyncMode     int32                                      `json:"syncMode,omitempty"`     /*  同步模式，可选值：<li>1： 全量同步+增量同步。<li>2：全量同步。  */
	ConflictMode int32                                      `json:"conflictMode,omitempty"` /*  数据冲突时的处理办法，可选值：<li>1：中断迁移。<li>2：跳过目标key，继续执行。<li>3：覆盖目标key，继续执行。  */
}

type Dcs2CreateTransferTaskSourceDbInfoRequest struct {
	SpuInstId       string `json:"spuInstId,omitempty"`   /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	IpAddr          string `json:"ipAddr,omitempty"`      /*  连接地址。  */
	OriginalCluster *bool  `json:"originalCluster"`       /*  是否是原生cluster集群，输入实例ID可不填，否则必填。  */
	AccountName     string `json:"accountName,omitempty"` /*  数据库账号。  */
	Password        string `json:"password,omitempty"`    /*  数据库密码。  */
}

type Dcs2CreateTransferTaskTargetDbInfoRequest struct {
	SpuInstId       string `json:"spuInstId,omitempty"`   /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	IpAddr          string `json:"ipAddr,omitempty"`      /*  连接地址。  */
	OriginalCluster *bool  `json:"originalCluster"`       /*  是否是原生cluster集群，输入实例ID可不填，否则必填。  */
	AccountName     string `json:"accountName,omitempty"` /*  数据库账号。  */
	Password        string `json:"password,omitempty"`    /*  数据库密码。  */
}

type Dcs2CreateTransferTaskResponse struct {
	StatusCode int32                                    `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                   `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2CreateTransferTaskReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                   `json:"requestId"`  /*  请求 ID。  */
	Code       string                                   `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                   `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2CreateTransferTaskReturnObjResponse struct {
	TaskId string `json:"taskId"`
}
