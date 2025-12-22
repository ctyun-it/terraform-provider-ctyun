package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeBackupsApi
/* 查询备份文件信息
 */type Dcs2DescribeBackupsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeBackupsApi(client *core.CtyunClient) *Dcs2DescribeBackupsApi {
	return &Dcs2DescribeBackupsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisMgr/describeBackups",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeBackupsApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeBackupsRequest) (*Dcs2DescribeBackupsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.RestoreName != "" {
		ctReq.AddParam("restoreName", req.RestoreName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeBackupsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeBackupsRequest struct {
	RegionId    string `json:"regionId,omitempty"`    /*  资源池ID。获取方法如下：<br>方法一：通过查看附录文档<a  target="_blank" rel="noopener noreferrer" href="https://www.ctyun.cn/document/10029420/11067697">分布式缓存服务Redis资源池</a>获取资源池ID。<br>方法二：可调用  <a  target="_blank" rel="noopener noreferrer" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a> 接口获取resPoolCode字段。  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID。获取方法如下：<br>方法一：可登录分布式缓存控制台在实例列表复制实例ID。<br>方法二：可调用<a target="_blank" href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7728&isNormal=1&vid=270"> 查询实例列表 </a>接口获取prodInstId字段。  */
	RestoreName string `json:"restoreName,omitempty"` /*  备份恢复记录。  */
}

type Dcs2DescribeBackupsResponse struct {
	StatusCode int32                                 `json:"statusCode"` /*  响应状态码。<li>800：成功。<li>900：失败。  */
	Message    string                                `json:"message"`    /*  响应信息。  */
	ReturnObj  *Dcs2DescribeBackupsReturnObjResponse `json:"returnObj"`  /*  返回数据对象，数据见returnObj。  */
	RequestId  string                                `json:"requestId"`  /*  请求 ID。  */
	Code       string                                `json:"code"`       /*  响应码，仅表示请求是否执行。  */
	Error      string                                `json:"error"`      /*  错误码，参见错误码说明。  */
}

type Dcs2DescribeBackupsReturnObjResponse struct {
	Total int32                                       `json:"total"` /*  数量。  */
	Rows  []*Dcs2DescribeBackupsReturnObjRowsResponse `json:"rows"`  /*  备份文件信息集合，见Restore。  */
}

type Dcs2DescribeBackupsReturnObjRowsResponse struct {
	RestoreName    string `json:"restoreName"`    /*  备份名。  */
	CreateTime     string `json:"createTime"`     /*  创建时间（格式：yyyy-MM-dd HH:mm:ss）。  */
	Status         string `json:"status"`         /*  备份状态。<li>success：成功。<li>processing：进行中。<li>fail：失败。  */
	RecoveryStatus string `json:"recoveryStatus"` /*  备份恢复状态。<li>success：成功。<li>processing：进行中。<li>fail：失败。<li>create：备份点创建。  */
	RawType        int32  `json:"type"`           /*  备份类型。<li>0：手动备份。<li>1：自动备份。  */
	Remark         string `json:"remark"`         /*  备注信息。  */
}
