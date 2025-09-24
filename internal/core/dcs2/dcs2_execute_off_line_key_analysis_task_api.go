package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ExecuteOffLineKeyAnalysisTaskApi
/* 立即执行离线全量key分析任务
 */type Dcs2ExecuteOffLineKeyAnalysisTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ExecuteOffLineKeyAnalysisTaskApi(client *core.CtyunClient) *Dcs2ExecuteOffLineKeyAnalysisTaskApi {
	return &Dcs2ExecuteOffLineKeyAnalysisTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/keyAnalysisMgrServant/executeOffLineKeyAnalysisTask",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ExecuteOffLineKeyAnalysisTaskApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ExecuteOffLineKeyAnalysisTaskRequest) (*Dcs2ExecuteOffLineKeyAnalysisTaskResponse, error) {
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
	var resp Dcs2ExecuteOffLineKeyAnalysisTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ExecuteOffLineKeyAnalysisTaskRequest struct {
	RegionId    string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID  */
	NodeName    string `json:"nodeName,omitempty"`    /*  节点名称<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7752&isNormal=1&vid=270">获取redis节点名列表</a> node表nodeName字段  */
	RdbType     string `json:"rdbType,omitempty"`     /*  rdb文件类型<li>rdbNew：新建rdb文件<li>rdbBackup：使用历史rdb文件  */
	RestoreName string `json:"restoreName,omitempty"` /*  备份名<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7778&isNormal=1&vid=270">查询备份文件信息</a> Restore表restoreName字段  */
}

type Dcs2ExecuteOffLineKeyAnalysisTaskResponse struct {
	StatusCode int32                                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ExecuteOffLineKeyAnalysisTaskReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ExecuteOffLineKeyAnalysisTaskReturnObjResponse struct{}
