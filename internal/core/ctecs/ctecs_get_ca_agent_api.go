package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsGetCaAgentApi
/* 调用此接口可以查询一台或多台弹性云主机、物理机内是否安装了云助手agent
 */ /* 说明：仅支持批量查询弹性云主机或物理机是否安装了云助手agent，不支持混合查询
 */type CtecsGetCaAgentApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsGetCaAgentApi(client *core.CtyunClient) *CtecsGetCaAgentApi {
	return &CtecsGetCaAgentApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/cloud-assistant/get-ca-agent",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsGetCaAgentApi) Do(ctx context.Context, credential core.Credential, req *CtecsGetCaAgentRequest) (*CtecsGetCaAgentResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsGetCaAgentResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsGetCaAgentRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID  */
	InstanceIDs string `json:"instanceIDs,omitempty"` /*  待执行命令的云主机、物理机ID列表, 使用英文 , 分割  */
	PageNo      int32  `json:"pageNo,omitempty"`      /*  当前页码，默认值为1  */
	PageSize    int32  `json:"pageSize,omitempty"`    /*  分页查询时设置的每页行数，最大值为100，默认为10  */
}

type CtecsGetCaAgentResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800 为成功，900 为失败）  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码，详见错误码说明  */
	Message     string                            `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                            `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsGetCaAgentReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsGetCaAgentReturnObjResponse struct {
	CaAgentStatusSet []*CtecsGetCaAgentReturnObjCaAgentStatusSetResponse `json:"caAgentStatusSet"`     /*  状态列表  */
	TotalCount       int32                                               `json:"totalCount,omitempty"` /*  命令总个数  */
	PageNo           int32                                               `json:"pageNo,omitempty"`     /*  当前页码  */
	PageSize         int32                                               `json:"pageSize,omitempty"`   /*  每页行数  */
}

type CtecsGetCaAgentReturnObjCaAgentStatusSetResponse struct {
	Status     string `json:"status,omitempty"`     /*  agent当前状态，取值范围：<br />Error：异常；<br />Running：运行中；  */
	AgentName  string `json:"agentName,omitempty"`  /*  agent名称  */
	Version    string `json:"version,omitempty"`    /*  版本  */
	InstanceID string `json:"instanceID,omitempty"` /*  实例ID  */
}
