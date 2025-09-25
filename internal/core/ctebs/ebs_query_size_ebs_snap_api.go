package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsQuerySizeEbsSnapApi
/* 天翼云云硬盘快照服务以云硬盘的快照使用总容量为粒度进行统计并计费。您可以查询某个资源池下的快照使用量，快照使用总容量为当前云硬盘中所有快照的数据块所占用的存储空间之和。
 */type EbsQuerySizeEbsSnapApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsQuerySizeEbsSnapApi(client *core.CtyunClient) *EbsQuerySizeEbsSnapApi {
	return &EbsQuerySizeEbsSnapApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs_snapshot/query_size-ebs-snap",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsQuerySizeEbsSnapApi) Do(ctx context.Context, credential core.Credential, req *EbsQuerySizeEbsSnapRequest) (*EbsQuerySizeEbsSnapResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsQuerySizeEbsSnapResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsQuerySizeEbsSnapRequest struct {
	RegionID string /*  资源池ID。  */
}

type EbsQuerySizeEbsSnapResponse struct {
	StatusCode  int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败/处理中）。  */
	Message     *string                               `json:"message,omitempty"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description *string                               `json:"description,omitempty"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsQuerySizeEbsSnapReturnObjResponse `json:"returnObj"`             /*  返回数据结构体。  */
	Details     *string                               `json:"details,omitempty"`     /*  可忽略。  */
	ErrorCode   *string                               `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码。请参考错误码。  */
	Error       *string                               `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式大驼峰码。请参考错误码。  */
}

type EbsQuerySizeEbsSnapReturnObjResponse struct {
	SnapshotTotalSize  int32                                          `json:"snapshotTotalSize"`       /*  region下的所有云硬盘快照的使用量，单位为B。  */
	SnapshotTotalCount int32                                          `json:"snapshotTotalCount"`      /*  云硬盘快照数量。  */
	AzDisplayName      *string                                        `json:"azDisplayName,omitempty"` /*  可用区名称。  */
	Details            []*EbsQuerySizeEbsSnapReturnObjDetailsResponse `json:"details"`                 /*  不同类型的快照占用的使用量详情。  */
}

type EbsQuerySizeEbsSnapReturnObjDetailsResponse struct {
	VolumeType *string `json:"volumeType,omitempty"` /*  云硬盘类型，取值为：
	●SATA：普通IO。
	●SAS：高IO。
	●SSD：超高IO。
	●FAST-SSD：极速型SSD。
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
	SnapshotSize int32 `json:"snapshotSize"` /*  单个类型的云硬盘的快照使用量，单位为B。  */
}
