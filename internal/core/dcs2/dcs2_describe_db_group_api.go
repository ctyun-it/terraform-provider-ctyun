package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2DescribeDbGroupApi
/* 查询分组列表。
 */type Dcs2DescribeDbGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2DescribeDbGroupApi(client *core.CtyunClient) *Dcs2DescribeDbGroupApi {
	return &Dcs2DescribeDbGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/groupManageMgrServant/describeDbGroup",
			ContentType:  "",
		},
	}
}

func (a *Dcs2DescribeDbGroupApi) Do(ctx context.Context, credential core.Credential, req *Dcs2DescribeDbGroupRequest) (*Dcs2DescribeDbGroupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2DescribeDbGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2DescribeDbGroupRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string /*  实例ID  */
}

type Dcs2DescribeDbGroupResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2DescribeDbGroupReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2DescribeDbGroupReturnObjResponse struct {
	Total int32                                       `json:"total,omitempty"` /*  数量  */
	Rows  []*Dcs2DescribeDbGroupReturnObjRowsResponse `json:"rows"`            /*  GroupNode数据集List  */
}

type Dcs2DescribeDbGroupReturnObjRowsResponse struct {
	GroupName      string `json:"groupName,omitempty"`      /*  分组名称  */
	GroupInfo      string `json:"groupInfo,omitempty"`      /*  分组信息  */
	LatestTime     string `json:"latestTime,omitempty"`     /*  更新时间， UTC格式  */
	Dborder        string `json:"dborder,omitempty"`        /*  分组对应的db  */
	RedisSetName   string `json:"redisSetName,omitempty"`   /*  redis集群名  */
	Title          string `json:"title,omitempty"`          /*  拓展，未使用  */
	Nodes          string `json:"nodes,omitempty"`          /*  拓展，未使用  */
	Link           string `json:"link,omitempty"`           /*  拓展，未使用  */
	UserName       string `json:"userName,omitempty"`       /*  实例名称  */
	UserPwd        string `json:"userPwd,omitempty"`        /*  用户密码  */
	UseClientCache string `json:"useClientCache,omitempty"` /*  拓展，未使用  */
	MaxObjects     string `json:"maxObjects,omitempty"`     /*  拓展，未使用  */
	MaxLife        string `json:"maxLife,omitempty"`        /*  拓展，未使用  */
}
