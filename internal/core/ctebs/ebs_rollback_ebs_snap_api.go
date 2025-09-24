package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsRollbackEbsSnapApi
/* 当发生误操作或系统故障等问题时，您可以使用已创建的快照来回滚数据，使云硬盘的数据恢复至创建快照的时刻，实现云硬盘数据的恢复。
 */ /* 注意：
 */ /* 1、快照回滚为不可逆操作，从快照创建时刻到回滚操作开始时刻这段时间内的数据会被删除。为避免误操作，建议您在回滚之前为云硬盘创建一次快照进行数据备份。
 */ /* 2、如果您在创建云硬盘快照之后对该盘进行了扩容，当您回滚到该快照后，您之前扩容的分区和文件系统会丢失。因此，您需要在回滚后登录云主机重新进行扩展分区和文件系统的操作，以恢复到扩容后的状态。
 */type EbsRollbackEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsRollbackEbsSnapApi(client *core.CtyunClient) *EbsRollbackEbsSnapApi {
	return &EbsRollbackEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs_snapshot/rollback-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsRollbackEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsRollbackEbsSnapRequest) (*EbsRollbackEbsSnapResponse, error) {
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
	var resp EbsRollbackEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsRollbackEbsSnapRequest struct {
	SnapshotID string `json:"snapshotID,omitempty"` /*  快照ID。  */
	DiskID     string `json:"diskID,omitempty"`     /*  快照所属的云硬盘ID。  */
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID。  */
	AzName     string `json:"azName,omitempty"`     /*  快照所属的可用区名称。  */
}

type EbsRollbackEbsSnapResponse struct {
	StatusCode  int32                                `json:"statusCode"`            /*  返回状态码（800为成功，900为失败/处理中）。  */
	Message     *string                              `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                              `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsRollbackEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                              `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                              `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                              `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsRollbackEbsSnapReturnObjResponse struct {
	SnapshotJobID *string `json:"snapshotJobID,omitempty"` /*  异步任务ID，可通过公共接口/v4/job/info查询该jobID来查看异步任务最终执行结果。  */
}
