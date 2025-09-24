package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2AsyncExportInstancesApi
/* 异步导出实例列表。
 */type Dcs2AsyncExportInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2AsyncExportInstancesApi(client *core.CtyunClient) *Dcs2AsyncExportInstancesApi {
	return &Dcs2AsyncExportInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instanceManageMgrServant/asyncExportInstances",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2AsyncExportInstancesApi) Do(ctx context.Context, credential core.Credential, req *Dcs2AsyncExportInstancesRequest) (*Dcs2AsyncExportInstancesResponse, error) {
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
	var resp Dcs2AsyncExportInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2AsyncExportInstancesRequest struct {
	RegionId      string   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProjectId     string   `json:"projectId,omitempty"`     /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br>说明：默认值为"0"  */
	InstanceName  string   `json:"instanceName,omitempty"`  /*  实例名称  */
	Capacity      string   `json:"capacity,omitempty"`      /*  实例规格容量，单位G  */
	ProdInstId    string   `json:"prodInstId,omitempty"`    /*  实例ID  */
	Vip           string   `json:"vip,omitempty"`           /*  实例虚拟IP  */
	Status        int32    `json:"status,omitempty"`        /*  实例状态<li> 0：有效<li>1：开通中<li>2：暂停<li>3：变更中<li>4：开通失败<li>5：停止中<li>6：已停止<li>8：已退订;  */
	EngineVersion string   `json:"engineVersion,omitempty"` /*  引擎版本  */
	PayType       int32    `json:"payType,omitempty"`       /*  付费类型<li>0：包年/包月<li>1: 按需  */
	CpuArchType   string   `json:"cpuArchType,omitempty"`   /*  cpu架构<li>x86<li>arm  */
	LabelIds      []string `json:"labelIds"`                /*  标签ID  */
}

type Dcs2AsyncExportInstancesResponse struct {
	StatusCode int32                                      `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                     `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2AsyncExportInstancesReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                     `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                     `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2AsyncExportInstancesReturnObjResponse struct {
	TaskId string `json:"taskId,omitempty"` /*  任务ID  */
}
