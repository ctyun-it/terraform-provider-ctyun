package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetObjectAclApi
/* 获取对象的访问权限控制列表（ACL）。
 */type ZosGetObjectAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetObjectAclApi(client *core.CtyunClient) *ZosGetObjectAclApi {
	return &ZosGetObjectAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-object-acl",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetObjectAclApi) Do(ctx context.Context, credential core.Credential, req *ZosGetObjectAclRequest) (*ZosGetObjectAclResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("key", req.Key)
	if req.VersionID != "" {
		ctReq.AddParam("versionID", req.VersionID)
	}
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetObjectAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetObjectAclRequest struct {
	Bucket    string /*  桶名  */
	Key       string /*  对象名  */
	VersionID string /*  版本ID，在开启多版本时可使用  */
	RegionID  string /*  区域 ID  */
}

type ZosGetObjectAclResponse struct {
	StatusCode  int64                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                            `json:"message,omitempty"`     /*  状态描述  */
	Description string                            `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetObjectAclReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                            `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetObjectAclReturnObjResponse struct {
	Owner  *ZosGetObjectAclReturnObjOwnerResponse    `json:"owner"`  /*  所有者  */
	Grants []*ZosGetObjectAclReturnObjGrantsResponse `json:"grants"` /*  授权信息  */
}

type ZosGetObjectAclReturnObjOwnerResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名称  */
	ID          string `json:"ID,omitempty"`          /*  用户名  */
}

type ZosGetObjectAclReturnObjGrantsResponse struct {
	Grantee    *ZosGetObjectAclReturnObjGrantsGranteeResponse `json:"grantee"`              /*  被授予权限的人的容器  */
	Permission string                                         `json:"permission,omitempty"` /*  权限，为 WRITE, WRITE_ACP, FULL_CONTROL, READ, READ_ACP 之中的值<br>WRITE：向桶中写对象的权限<br>WRITE_ACP：修改桶的访问控制权限的能力<br>READ：读取桶中文件列表的能力<br>READ_ACP：获取桶的访问控制权限的能力<br>FULL_CONTROL：同桶的所属者相同的权限，以上能力都具有  */
}

type ZosGetObjectAclReturnObjGrantsGranteeResponse struct {
	EmailAddress string `json:"emailAddress,omitempty"` /*  邮箱地址  */
	RawType      string `json:"type,omitempty"`         /*  用户类型， CanonicalUser，Group 二者之一  */
	DisplayName  string `json:"displayName,omitempty"`  /*  展示名称  */
	ID           string `json:"ID,omitempty"`           /*  用户名  */
	URI          string `json:"URI,omitempty"`          /*  URI，不存在时为 null  */
}
