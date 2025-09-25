package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketEncryptionApi
/* 为指定桶设置服务端加密。
 */type ZosPutBucketEncryptionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketEncryptionApi(client *core.CtyunClient) *ZosPutBucketEncryptionApi {
	return &ZosPutBucketEncryptionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-encryption",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketEncryptionApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketEncryptionRequest) (*ZosPutBucketEncryptionResponse, error) {
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
	var resp ZosPutBucketEncryptionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketEncryptionRequest struct {
	Bucket                            string                                                          `json:"bucket,omitempty"`                  /*  桶名  */
	RegionID                          string                                                          `json:"regionID,omitempty"`                /*  区域 ID  */
	ServerSideEncryptionConfiguration *ZosPutBucketEncryptionServerSideEncryptionConfigurationRequest `json:"serverSideEncryptionConfiguration"` /*  加密配置  */
}

type ZosPutBucketEncryptionServerSideEncryptionConfigurationRequest struct {
	Rules []*ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesRequest `json:"rules"` /*  规则  */
}

type ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesRequest struct {
	ApplyServerSideEncryptionByDefault *ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesApplyServerSideEncryptionByDefaultRequest `json:"applyServerSideEncryptionByDefault"` /*  加密配置  */
}

type ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesApplyServerSideEncryptionByDefaultRequest struct {
	SSEAlgorithm   string `json:"SSEAlgorithm,omitempty"`   /*  加密算法，仅支持 AES256 或 aws:kms，若传入 AES256，将自动生成 KMSMasterKeyID，若传入 aws:kms，需用户预先通过密钥管理服务创建密钥  */
	KMSMasterKeyID string `json:"KMSMasterKeyID,omitempty"` /*  当且仅当 SSEAlgorithm 为 aws:kms 时需要填写此参数，参数格式为"{密钥管理服务处的密钥ID}:::{regionID}:{userID}"  */
}

type ZosPutBucketEncryptionResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
