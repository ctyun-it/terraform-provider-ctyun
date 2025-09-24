package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbListAccessControlApi
/* 查询策略地址组，访问控制采用黑、白名单方式实现，此接口为查询黑、白名单的地址组。
 */type CtelbListAccessControlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbListAccessControlApi(client *core.CtyunClient) *CtelbListAccessControlApi {
	return &CtelbListAccessControlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/elb/list-access-control",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbListAccessControlApi) Do(ctx context.Context, credential core.Credential, req *CtelbListAccessControlRequest) (*CtelbListAccessControlResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	for _, ele0 := range req.IDs {
		if ele0 != "" {
			ctReq.AddParam("IDs", ele0)
		}
	}
	if req.Name != "" {
		ctReq.AddParam("name", req.Name)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtelbListAccessControlResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbListAccessControlRequest struct {
	RegionID string   /*  区域ID  */
	IDs      []string /*  访问控制ID列表  */
	Name     string   /*  访问控制名称,只能由数字，字母，-组成不能以数字和-开头，最大长度32  */
}

type CtelbListAccessControlResponse struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtelbListAccessControlReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                     `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbListAccessControlReturnObjResponse struct {
	AzName      string   `json:"azName,omitempty"`      /*  可用区名称  */
	ProjectID   string   `json:"projectID,omitempty"`   /*  项目ID  */
	ID          string   `json:"ID,omitempty"`          /*  访问控制ID  */
	Name        string   `json:"name,omitempty"`        /*  访问控制名称  */
	Description string   `json:"description,omitempty"` /*  描述  */
	SourceIps   []string `json:"sourceIps"`             /*  IP地址的集合或者CIDR  */
	CreateTime  string   `json:"createTime,omitempty"`  /*  创建时间，为UTC格式  */
}
