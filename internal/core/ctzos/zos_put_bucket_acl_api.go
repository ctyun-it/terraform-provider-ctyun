package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketAclApi
/* 设置桶的ACL，管理账户对存储桶的访问权限和ACL权限。
 */type ZosPutBucketAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketAclApi(client *core.CtyunClient) *ZosPutBucketAclApi {
	return &ZosPutBucketAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-acl",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketAclApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketAclRequest) (*ZosPutBucketAclResponse, error) {
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
	var resp ZosPutBucketAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketAclRequest struct {
	Bucket              string                                     `json:"bucket,omitempty"`    /*  桶名  */
	RegionID            string                                     `json:"regionID,omitempty"`  /*  区域 ID  */
	ACL                 string                                     `json:"ACL,omitempty"`       /*  ACL 配置，允许的值为 private，public-read。此参数不能与 accessControlPolicy 共同使用  */
	AccessControlPolicy *ZosPutBucketAclAccessControlPolicyRequest `json:"accessControlPolicy"` /*  访问控制策略，不能与 ACL 共同使用  */
}

type ZosPutBucketAclAccessControlPolicyRequest struct {
	Owner  *ZosPutBucketAclAccessControlPolicyOwnerRequest    `json:"owner"`  /*  所有者  */
	Grants []*ZosPutBucketAclAccessControlPolicyGrantsRequest `json:"grants"` /*  授权信息数组  */
}

type ZosPutBucketAclAccessControlPolicyOwnerRequest struct {
	DisplayName string `json:"displayName,omitempty"` /*  用户名  */
	ID          string `json:"ID,omitempty"`          /*  用户 ID  */
}

type ZosPutBucketAclAccessControlPolicyGrantsRequest struct {
	Grantee    *ZosPutBucketAclAccessControlPolicyGrantsGranteeRequest `json:"grantee"`              /*  授权信息  */
	Permission string                                                  `json:"permission,omitempty"` /*  权限，为 WRITE, WRITE_ACP, FULL_CONTROL, READ, READ_ACP 之中的值<br>WRITE：向桶中写对象的权限<br>WRITE_ACP：修改桶的访问控制权限的能力<br>READ：读取桶中文件列表的能力<br>READ_ACP：获取桶的访问控制权限的能力<br>FULL_CONTROL：同桶的所属者相同的权限，以上能力都具有  */
}

type ZosPutBucketAclAccessControlPolicyGrantsGranteeRequest struct {
	EmailAddress string `json:"emailAddress,omitempty"` /*  邮箱地址  */
	RawType      string `json:"type,omitempty"`         /*  用户类型， CanonicalUser，AmazonCustomerByEmail，Group 三者之一。type 为 CanonicalUser 时，必填 ID；为 AmazonCustomerByEmail，必填 emailAddress；为 Group 必填URI。另外，使用 AmazonCustomerByEmail 时，将会保存其指向到的 CanonicalUser 类型的用户  */
	DisplayName  string `json:"displayName,omitempty"`  /*  用户名  */
	ID           string `json:"ID,omitempty"`           /*  用户 ID  */
	URI          string `json:"URI,omitempty"`          /*  URI，不存在时为 null  */
}

type ZosPutBucketAclResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
