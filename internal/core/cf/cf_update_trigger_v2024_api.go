package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfUpdateTriggerV2024Api
/* 修改触发器 */
type CfUpdateTriggerV2024Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfUpdateTriggerV2024Api(client *core.CtyunClient) *CfUpdateTriggerV2024Api {
	return &CfUpdateTriggerV2024Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/functions/{functionName}/triggers/{triggerName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfUpdateTriggerV2024Api) Do(ctx context.Context, credential core.Credential, req *CfUpdateTriggerV2024Request) (*CfUpdateTriggerV2024Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("triggerName", req.TriggerName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfUpdateTriggerV2024Request
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
		TriggerName  interface{} `json:"triggerName,omitempty"`
	}{
		req, nil, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfUpdateTriggerV2024Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfUpdateTriggerV2024Request struct {
	FunctionName  string `json:"functionName"`  /*  函数名称，函数必须已存在  */
	TriggerName   string `json:"triggerName"`   /*  触发器名称，触发器必须已存在  */
	RegionId      string `json:"regionId"`      /*  资源池id，标识不同的地区，如：华东1、西南1  */
	TriggerConfig string `json:"triggerConfig"` /*  触发器配置，JSON字符串 ，不同触发器类型配置不同，详情请查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=53&api=16023&data=42&isNormal=1&vid=40">创建触发器</a>  */
	Version       string `json:"version"`       /*  函数版本/别名，版本可以是1,2,...这种普通版本，也可以是特殊版本LATEST。  */
	Enable        bool   `json:"enable"`        /*  是否启用触发器。true：启用，false：禁用。  */
}

type CfUpdateTriggerV2024Response struct {
	StatusCode int32                                 `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      string                                `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    string                                `json:"message"`    /*  错误提示信息  */
	ReturnObj  *CfUpdateTriggerV2024ReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfUpdateTriggerV2024ReturnObjResponse struct {
	CreatedAt     *string `json:"createdAt"`     /*  创建时间  */
	UpdatedAt     *string `json:"updatedAt"`     /*  更新时间  */
	Creator       *int32  `json:"creator"`       /*  创建者ID  */
	Editor        *int32  `json:"editor"`        /*  编辑者ID  */
	TriggerName   *string `json:"triggerName"`   /*  触发器名称  */
	TriggerConfig *string `json:"triggerConfig"` /*  触发器配置，JSON字符串 ，不同触发器类型配置不同，详情请查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=53&api=16023&data=42&isNormal=1&vid=40">创建触发器</a>  */
	TriggerType   *string `json:"triggerType"`   /*  触发器类型  */
	Status        *int32  `json:"status"`        /*  触发器状态。1：启用；2：禁用；3：系统禁用  */
	Version       *string `json:"version"`       /*  别名或版本  */
	Region        *string `json:"region"`        /*  资源池id，标识不同的地区，如：华东1、西南1  */
	FunctionId    *int32  `json:"functionId"`    /*  函数ID  */
	UrlInternet   *string `json:"urlInternet"`   /*  外网URL ，非http触发器无值  */
	UrlIntranet   *string `json:"urlIntranet"`   /*  内网URL ，非http触发器无值  */
	IsVersion     *bool   `json:"isVersion"`     /*  是否是版本，区别于别名，1,2,...和LATEST属于版本  */
}
