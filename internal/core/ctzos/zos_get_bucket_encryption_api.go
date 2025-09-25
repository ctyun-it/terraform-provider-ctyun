package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketEncryptionApi
/* 查询指定桶的服务端加密配置。
 */type ZosGetBucketEncryptionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketEncryptionApi(client *core.CtyunClient) *ZosGetBucketEncryptionApi {
	return &ZosGetBucketEncryptionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-encryption",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketEncryptionApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketEncryptionRequest) (*ZosGetBucketEncryptionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketEncryptionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketEncryptionRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketEncryptionResponse struct {
	StatusCode  int64                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                   `json:"message,omitempty"`     /*  状态描述  */
	Description string                                   `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketEncryptionReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketEncryptionReturnObjResponse struct {
	ServerSideEncryptionConfiguration *ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationResponse `json:"serverSideEncryptionConfiguration"` /*  加密配置  */
}

type ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationResponse struct {
	Rules []*ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationRulesResponse `json:"rules"` /*  规则  */
}

type ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationRulesResponse struct {
	ApplyServerSideEncryptionByDefault *ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationRulesApplyServerSideEncryptionByDefaultResponse `json:"applyServerSideEncryptionByDefault"` /*  加密配置规则  */
}

type ZosGetBucketEncryptionReturnObjServerSideEncryptionConfigurationRulesApplyServerSideEncryptionByDefaultResponse struct {
	SSEAlgorithm   string `json:"SSEAlgorithm,omitempty"`   /*  加密算法，仅支持 AES256 或 aws:kms  */
	KMSMasterKeyID string `json:"KMSMasterKeyID,omitempty"` /*  KMS密钥ID  */
}
