package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbShowAccessControlApi
/* 查询策略地址组详情，访问控制采用黑、白名单方式实现，此接口为查询黑、白名单的地址组。
 */type CtelbShowAccessControlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbShowAccessControlApi(client *core.CtyunClient) *CtelbShowAccessControlApi {
	return &CtelbShowAccessControlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/show-access-control",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbShowAccessControlApi) Do(ctx context.Context, credential core.Credential, req *CtelbShowAccessControlRequest) (*CtelbShowAccessControlResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.Id != "" {
		ctReq.AddParam("id", req.Id)
	}
	if req.AccessControlID != "" {
		ctReq.AddParam("accessControlID", req.AccessControlID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbShowAccessControlResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbShowAccessControlRequest struct {
	RegionID        string /*  区域ID  */
	Id              string /*  访问控制ID, 该字段后续废弃  */
	AccessControlID string /*  访问控制ID, 推荐使用该字段, 当同时使用 id 和 accessControlID 时，优先使用 accessControlID  */
}

type CtelbShowAccessControlResponse struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                   `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbShowAccessControlReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                   `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbShowAccessControlReturnObjResponse struct {
	AzName      string   `json:"azName,omitempty"`      /*  可用区名称  */
	ProjectID   string   `json:"projectID,omitempty"`   /*  项目ID  */
	ID          string   `json:"ID,omitempty"`          /*  访问控制ID  */
	Name        string   `json:"name,omitempty"`        /*  访问控制名称  */
	Description string   `json:"description,omitempty"` /*  描述  */
	SourceIps   []string `json:"sourceIps"`             /*  IP地址的集合或者CIDR  */
	CreateTime  string   `json:"createTime,omitempty"`  /*  创建时间，为UTC格式  */
}
