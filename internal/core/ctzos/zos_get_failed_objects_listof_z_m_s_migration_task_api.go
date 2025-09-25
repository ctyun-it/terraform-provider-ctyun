package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosGetFailedObjectsListofZMSMigrationTaskApi
/* 查询迁移任务的失败对象列表
 */type ZosGetFailedObjectsListofZMSMigrationTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetFailedObjectsListofZMSMigrationTaskApi(client *core.CtyunClient) *ZosGetFailedObjectsListofZMSMigrationTaskApi {
	return &ZosGetFailedObjectsListofZMSMigrationTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/zms/list-migration-failed-detail",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetFailedObjectsListofZMSMigrationTaskApi) Do(ctx context.Context, credential core.Credential, req *ZosGetFailedObjectsListofZMSMigrationTaskRequest) (*ZosGetFailedObjectsListofZMSMigrationTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("migrationID", req.MigrationID)
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetFailedObjectsListofZMSMigrationTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetFailedObjectsListofZMSMigrationTaskRequest struct {
	RegionID    string /*  资源池 ID  */
	MigrationID string /*  迁移任务ID  */
	PageSize    int32  /*  页大小，默认值 10，取值范围 1~50  */
	PageNo      int32  /*  页码，默认值为1  */
}

type ZosGetFailedObjectsListofZMSMigrationTaskResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
