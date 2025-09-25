package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPauseZMSMigrationTaskApi
/* 对象存储迁移任务暂停
 */type ZosPauseZMSMigrationTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPauseZMSMigrationTaskApi(client *core.CtyunClient) *ZosPauseZMSMigrationTaskApi {
	return &ZosPauseZMSMigrationTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/zms/stop-migration",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPauseZMSMigrationTaskApi) Do(ctx context.Context, credential core.Credential, req *ZosPauseZMSMigrationTaskRequest) (*ZosPauseZMSMigrationTaskResponse, error) {
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
	var resp ZosPauseZMSMigrationTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPauseZMSMigrationTaskRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池 ID  */
	MigrationID string `json:"migrationID,omitempty"` /*  迁移任务ID  */
}

type ZosPauseZMSMigrationTaskResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
