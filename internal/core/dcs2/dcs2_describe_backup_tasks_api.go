package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeBackupTasksApi
/* 查询备份任务执行情况
 */type Dcs2DescribeBackupTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeBackupTasksApi(client *core.CtyunClient) *Dcs2DescribeBackupTasksApi {
	return &Dcs2DescribeBackupTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/redisMgr/describeBackupTasks",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeBackupTasksApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeBackupTasksRequest) (*Dcs2DescribeBackupTasksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("restoreName", req.RestoreName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeBackupTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeBackupTasksRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId  string /*  实例ID  */
	RestoreName string /*  备份名<br><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7771&isNormal=1&vid=270">手动备份数据</a> restoreName字段  */
}

type Dcs2DescribeBackupTasksResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeBackupTasksReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeBackupTasksReturnObjResponse struct {
	CreateTime  string `json:"createTime,omitempty"`  /*  创建时间（格式：yyyy-MM-dd HH:mm:ss）  */
	RestoreName string `json:"restoreName,omitempty"` /*  备份名  */
	Status      string `json:"status,omitempty"`      /*  节点状态<li>success：成功<li>processing：进行中<li>fail：失败。  */
}
