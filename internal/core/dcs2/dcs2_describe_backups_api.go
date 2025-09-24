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
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeBackupsResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeBackupsReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeBackupsReturnObjResponse struct {
	Total int32                                       `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeBackupsReturnObjRowsResponse `json:"rows"`            /*  备份文件信息集合，见Restore  */
}

type Dcs2DescribeBackupsReturnObjRowsResponse struct {
	RestoreName    string `json:"restoreName,omitempty"`    /*  备份名  */
	CreateTime     string `json:"createTime,omitempty"`     /*  创建时间（格式：yyyy-MM-dd HH:mm:ss）  */
	Status         string `json:"status,omitempty"`         /*  节点状态<li>success：成功<li>processing：进行中<li>fail：失败。  */
	RecoveryStatus string `json:"recoveryStatus,omitempty"` /*  备份恢复状态<li>success：成功<li>processing：进行中<li>fail：失败<li>create：备份点创建。  */
	RawType        int32  `json:"type,omitempty"`           /*  备份类型<li>0：手动备份<li>1：自动备份  */
}
