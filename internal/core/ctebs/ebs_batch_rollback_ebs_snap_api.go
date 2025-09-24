package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsBatchRollbackEbsSnapApi
/* 当发生误操作或系统故障等问题时，您可以使用已创建的快照来回滚数据，使云硬盘的数据恢复至创建快照的时刻，实现云硬盘数据的恢复。
 */ /* 本接口支持快照批量回滚。
 */ /* 注意：
 */ /* 1、快照回滚为不可逆操作，从快照创建时刻到回滚操作开始时刻这段时间内的数据会被删除。为避免误操作，建议您在回滚之前为云硬盘创建一次快照进行数据备份。
 */ /* 2、如果您在创建云硬盘快照之后对该盘进行了扩容，当您回滚到该快照后，您之前扩容的分区和文件系统会丢失。因此，您需要在回滚后登录云主机重新进行扩展分区和文件系统的操作，以恢复到扩容后的状态。
 */type EbsBatchRollbackEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsBatchRollbackEbsSnapApi(client *core.CtyunClient) *EbsBatchRollbackEbsSnapApi {
	return &EbsBatchRollbackEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/batch-rollback-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsBatchRollbackEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsBatchRollbackEbsSnapRequest) (*EbsBatchRollbackEbsSnapResponse, error) {
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
	var resp EbsBatchRollbackEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsBatchRollbackEbsSnapRequest struct {
	SnapshotList string `json:"snapshotList,omitempty"` /*  快照ID列表,以逗号分隔，最多8个。全部快照都存在才可进行批量重置。  */
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID。  */
	AzName       string `json:"azName,omitempty"`       /*  多可用区资源池下，必须指定可用区名称。  */
}

type EbsBatchRollbackEbsSnapResponse struct {
	StatusCode  int32                                     `json:"statusCode"`            /*  返回状态码(800为成功，900为失败/处理中)。  */
	Message     *string                                   `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                                   `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsBatchRollbackEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                                   `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsBatchRollbackEbsSnapReturnObjResponse struct{}
