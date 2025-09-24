package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbQuerySlaApi
/* 查看规格列表
 */type CtelbQuerySlaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbQuerySlaApi(client *core.CtyunClient) *CtelbQuerySlaApi {
	return &CtelbQuerySlaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/query-sla",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbQuerySlaApi) Do(ctx context.Context, credential core.Credential, req *CtelbQuerySlaRequest) (*CtelbQuerySlaResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbQuerySlaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbQuerySlaRequest struct {
	RegionID string /*  区域ID  */
}

type CtelbQuerySlaResponse struct {
	StatusCode  int32                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbQuerySlaReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       string                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbQuerySlaReturnObjResponse struct {
	RegionID    string `json:"regionID,omitempty"`    /*  区域ID  */
	AzName      string `json:"azName,omitempty"`      /*  az名称  */
	ProjectID   string `json:"projectID,omitempty"`   /*  项目ID  */
	ID          string `json:"ID,omitempty"`          /*  规格ID  */
	Name        string `json:"name,omitempty"`        /*  规格名称  */
	Description string `json:"description,omitempty"` /*  规格描述  */
	Spec        string `json:"spec,omitempty"`        /*  规格类型: 标准型I / 标准型II / 增强型I / 增强型II / 高阶型I / 高阶型II / 存量 /免费型  */
}
