package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowAclApi
/* 查看 Acl 的详细信息
 */type CtvpcShowAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowAclApi(client *core.CtyunClient) *CtvpcShowAclApi {
	return &CtvpcShowAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/acl/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowAclApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowAclRequest) (*CtvpcShowAclResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	ctReq.AddParam("aclID", req.AclID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowAclRequest struct {
	RegionID  string  /*  资源池ID  */
	ProjectID *string /*  企业项目 ID，默认为'0'  */
	AclID     string  /*  aclID  */
}

type CtvpcShowAclResponse struct {
	StatusCode  int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcShowAclReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                          `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcShowAclReturnObjResponse struct {
	AclID       *string   `json:"aclID,omitempty"`       /*  id  */
	Name        *string   `json:"name,omitempty"`        /*  名称  */
	Description *string   `json:"description,omitempty"` /*  描述  */
	VpcID       *string   `json:"vpcID,omitempty"`       /*  VPC  */
	Enabled     *string   `json:"enabled,omitempty"`     /*  disable,enable  */
	InPolicyID  []*string `json:"inPolicyID"`            /*  入规则id数组  */
	OutPolicyID []*string `json:"outPolicyID"`           /*  出规则id数组  */
	CreatedAt   *string   `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt   *string   `json:"updatedAt,omitempty"`   /*  更新时间  */
	SubnetIDs   []*string `json:"subnetIDs"`             /*  acl 绑定的子网 id  */
}
