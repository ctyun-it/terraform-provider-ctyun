package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcAclCloneApi
/* 克隆 Acl,仅实现acl的规则复制，不包括关联资源和相关属性
 */type CtvpcAclCloneApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcAclCloneApi(client *core.CtyunClient) *CtvpcAclCloneApi {
	return &CtvpcAclCloneApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/acl/clone",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcAclCloneApi) Do(ctx context.Context, credential core.Credential, req *CtvpcAclCloneRequest) (*CtvpcAclCloneResponse, error) {
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
	var resp CtvpcAclCloneResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcAclCloneRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	DestRegionID string `json:"destRegionID,omitempty"` /*  目标资源池，仅支持从4.0资源池复制到4.0资源池  */
	SrcAclID     string `json:"srcAclID,omitempty"`     /*  源aclID  */
	VpcID        string `json:"vpcID,omitempty"`        /*  目标资源池得到的acl归属的vpc  */
	Name         string `json:"name,omitempty"`         /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
}

type CtvpcAclCloneResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcAclCloneReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcAclCloneReturnObjResponse struct {
	AclID *string `json:"aclID,omitempty"` /*  acl id  */
	Name  *string `json:"name,omitempty"`  /*  acl 名称  */
}
