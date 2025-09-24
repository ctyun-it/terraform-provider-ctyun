package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsStatisticsInstanceBackupApi
/* 统计用户虚机盘总大小及备份总个数
 */type CtecsStatisticsInstanceBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsStatisticsInstanceBackupApi(client *core.CtyunClient) *CtecsStatisticsInstanceBackupApi {
	return &CtecsStatisticsInstanceBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-instance-resource",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsStatisticsInstanceBackupApi) Do(ctx context.Context, credential core.Credential, req *CtecsStatisticsInstanceBackupRequest) (*CtecsStatisticsInstanceBackupResponse, error) {
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
	var resp CtecsStatisticsInstanceBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsStatisticsInstanceBackupRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以调用<a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>查看最新的天翼云资源池列表  */
}

type CtecsStatisticsInstanceBackupResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsStatisticsInstanceBackupReturnObjResponse `json:"returnObj"`             /*  返回参数  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
}

type CtecsStatisticsInstanceBackupReturnObjResponse struct {
	TotalVolumeSize  int32 `json:"totalVolumeSize,omitempty"`  /*  虚机磁盘占用大小  */
	TotalBackupCount int32 `json:"totalBackupCount,omitempty"` /*  磁盘备份个数  */
}
