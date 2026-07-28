package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ListTaskInfoApi
/* 查询数据迁移任务列表 */
type Dcs2ListTaskInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ListTaskInfoApi(client *core.CtyunClient) *Dcs2ListTaskInfoApi {
	return &Dcs2ListTaskInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/transfer/listTaskInfo",
			ContentType:  "",
		},
	}
}

func (a *Dcs2ListTaskInfoApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ListTaskInfoRequest) (*Dcs2ListTaskInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("pageNum", req.PageNum)
	ctReq.AddParam("pageSize", req.PageSize)
	if *req.Status != "" {
		ctReq.AddParam("status", *req.Status)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2ListTaskInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Dcs2ListTaskInfoRequest struct {
	RegionId string  `json:"regionId"`         /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	PageNum  string  `json:"pageNum"`          /*  页码（范围：> 0）。  */
	PageSize string  `json:"pageSize"`         /*  数量（范围：[1,100]）。  */
	Status   *string `json:"status,omitempty"` /*  查询指定的任务状态，可选值：<li>0：所有状态（默认）。<li>1：运行中。<li>2：成功。<li>3：失败。  */
}

type Dcs2ListTaskInfoResponse struct {
	StatusCode int32                              `json:"statusCode"`          /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    *string                            `json:"message,omitempty"`   /*  响应信息。  */
	ReturnObj  *Dcs2ListTaskInfoReturnObjResponse `json:"returnObj,omitempty"` /*  返回数据对象，数据见returnObj。  */
	RequestId  *string                            `json:"requestId,omitempty"` /*  请求 ID。  */
	Code       *string                            `json:"code,omitempty"`      /*  响应码，仅表示请求是否执行。  */
	Error      *string                            `json:"error,omitempty"`     /*  错误码，参见错误码说明。  */
}

type Dcs2ListTaskInfoReturnObjResponse struct {
	Total int32                                    `json:"total"`          /*  总数。  */
	Size  int32                                    `json:"size"`           /*  当前页数量。  */
	List  []*Dcs2ListTaskInfoReturnObjListResponse `json:"list,omitempty"` /*  迁移任务列表。  */
}

type Dcs2ListTaskInfoReturnObjListResponse struct {
	UserId          *string                                      `json:"userId,omitempty"`          /*  用户ID。  */
	TenantId        *string                                      `json:"tenantId,omitempty"`        /*  租户ID。  */
	TaskId          *string                                      `json:"taskId,omitempty"`          /*  任务ID。  */
	SourceSpuInstId *string                                      `json:"sourceSpuInstId,omitempty"` /*  源库实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	TargetSpuInstId *string                                      `json:"targetSpuInstId,omitempty"` /*  目标库实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	RawType         int32                                        `json:"type"`                      /*  类型。  */
	Status          int32                                        `json:"status"`                    /*  任务状态。<li>0：初始态。<li> 1：运行中。<li>2：成功。<li>3：失败。  */
	RunStep         int32                                        `json:"runStep"`                   /*  迁移进度。<li>1：全量同步中。<li>2：增量同步中。  */
	CreateTime      int64                                        `json:"createTime"`                /*  创建时间戳（秒）。  */
	CompleteTime    int64                                        `json:"completeTime"`              /*  结束时间戳（秒，=-1时表示时间未知）。  */
	Detail          *Dcs2ListTaskInfoReturnObjListDetailResponse `json:"detail,omitempty"`          /*  详情。  */
}

type Dcs2ListTaskInfoReturnObjListDetailResponse struct {
	SourceDbInfo *Dcs2ListTaskInfoReturnObjListDetailSourceDbInfoResponse `json:"sourceDbInfo,omitempty"` /*  源库信息。  */
	TargetDbInfo *Dcs2ListTaskInfoReturnObjListDetailTargetDbInfoResponse `json:"targetDbInfo,omitempty"` /*  目标库信息。  */
}

type Dcs2ListTaskInfoReturnObjListDetailSourceDbInfoResponse struct {
	SpuInstId       *string `json:"spuInstId,omitempty"`       /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	IpAddr          *string `json:"ipAddr,omitempty"`          /*  连接地址。  */
	OriginalCluster *bool   `json:"originalCluster,omitempty"` /*  是否是原生cluster集群。  */
	AccountName     *string `json:"accountName,omitempty"`     /*  数据库账号。  */
	Password        *string `json:"password,omitempty"`        /*  数据库密码。  */
	InstanceName    *string `json:"instanceName,omitempty"`    /*  实例名称。  */
}

type Dcs2ListTaskInfoReturnObjListDetailTargetDbInfoResponse struct {
	SpuInstId       *string `json:"spuInstId,omitempty"`       /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	IpAddr          *string `json:"ipAddr,omitempty"`          /*  连接地址。  */
	OriginalCluster *bool   `json:"originalCluster,omitempty"` /*  是否是原生cluster集群。  */
	AccountName     *string `json:"accountName,omitempty"`     /*  数据库账号。  */
	Password        *string `json:"password,omitempty"`        /*  数据库密码。  */
	InstanceName    *string `json:"instanceName,omitempty"`    /*  实例名称。  */
}
