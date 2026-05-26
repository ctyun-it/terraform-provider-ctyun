package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetProvisionConfigApi
/* 查询函数版本预留实例配置 */
type CfGetProvisionConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetProvisionConfigApi(client *core.CtyunClient) *CfGetProvisionConfigApi {
	return &CfGetProvisionConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/resources/functions/{functionName}/provision-config",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetProvisionConfigApi) Do(ctx context.Context, credential core.Credential, req *CfGetProvisionConfigRequest) (*CfGetProvisionConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.Qualifier != nil && *req.Qualifier != "" {
		ctReq.AddParam("qualifier", *req.Qualifier)
	}
	if req.Version != nil && *req.Version != "" {
		ctReq.AddParam("version", *req.Version)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetProvisionConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetProvisionConfigRequest struct {
	FunctionName string  `json:"functionName"`        /*  函数名称  */
	RegionId     string  `json:"regionId"`            /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	Qualifier    *string `json:"qualifier,omitempty"` /*  函数别名。函数版本、函数别名必须填其中一个  */
	Version      *string `json:"version,omitempty"`   /*  函数版本。函数版本、函数别名必须填其中一个  */
}

type CfGetProvisionConfigResponse struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfGetProvisionConfigReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetProvisionConfigReturnObjResponse struct {
	Current          *int32                                                   `json:"current"`          /*  实际资源个数  */
	CurrentError     *string                                                  `json:"currentError"`     /*  预留实例创建失败时的错误信息  */
	FunctionId       *string                                                  `json:"functionId"`       /*  函数id  */
	FunctionName     *string                                                  `json:"functionName"`     /*  函数名称  */
	KsvcName         *string                                                  `json:"ksvcName"`         /*  ksvc名称  */
	Version          *string                                                  `json:"version"`          /*  函数版本  */
	Target           *int32                                                   `json:"target"`           /*  目标资源个数  */
	IdleMode         *bool                                                    `json:"idleMode"`         /*  是否为闲置模式  */
	Qualifier        *string                                                  `json:"qualifier"`        /*  函数别名  */
	ScheduledActions []*CfGetProvisionConfigReturnObjScheduledActionsResponse `json:"scheduledActions"` /*  定时策略配置  */
}

type CfGetProvisionConfigReturnObjScheduledActionsResponse struct {
	ScheduleExpression *string `json:"scheduleExpression"` /*  定时配置表达式  */
	Name               *string `json:"name"`               /*  定时策略名称  */
	StartTime          *string `json:"startTime"`          /*  策略生效时间  */
	EndTime            *string `json:"endTime"`            /*  策略失效时间  */
	Target             *int32  `json:"target"`             /*  预留目标实例数  */
	Timezone           *string `json:"timezone"`           /*  时区  */
}
