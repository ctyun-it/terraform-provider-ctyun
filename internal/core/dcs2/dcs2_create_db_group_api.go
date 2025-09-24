package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2CreateDbGroupApi
/* 新增分组，分组命名与db一一对应。
 */type Dcs2CreateDbGroupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2CreateDbGroupApi(client *core.CtyunClient) *Dcs2CreateDbGroupApi {
	return &Dcs2CreateDbGroupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/groupManageMgrServant/createDbGroup",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2CreateDbGroupApi) Do(ctx context.Context, credential core.Credential, req *Dcs2CreateDbGroupRequest) (*Dcs2CreateDbGroupResponse, error) {
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
	var resp Dcs2CreateDbGroupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2CreateDbGroupRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	GroupId    string `json:"groupId,omitempty"`    /*  分组名称，返回格式为：group.实例名.  test  */
	Db         string `json:"db,omitempty"`         /*  数据库编号，取值范围：0-255  */
}

type Dcs2CreateDbGroupResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2CreateDbGroupReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2CreateDbGroupReturnObjResponse struct {
	GroupName      string `json:"groupName,omitempty"`      /*  分组名称  */
	GroupInfo      string `json:"groupInfo,omitempty"`      /*  分组信息  */
	LastestTime    string `json:"lastestTime,omitempty"`    /*  更新时间  */
	Dborder        string `json:"dborder,omitempty"`        /*  分组对应的db  */
	RedisSetName   string `json:"redisSetName,omitempty"`   /*  集群名  */
	Title          string `json:"title,omitempty"`          /*  扩展，不使用  */
	Nodes          string `json:"nodes,omitempty"`          /*  扩展，不使用  */
	Link           string `json:"link,omitempty"`           /*  扩展，不使用  */
	UserName       string `json:"userName,omitempty"`       /*  实例名称  */
	UserPwd        string `json:"userPwd,omitempty"`        /*  用户密码  */
	UseClientCache string `json:"useClientCache,omitempty"` /*  是否使用客户端缓存  */
	MaxObjects     string `json:"maxObjects,omitempty"`     /*  扩展，不使用  */
	MaxLife        string `json:"maxLife,omitempty"`        /*  扩展，不使用  */
}
