package business

import (
	"context"
	"errors"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-core"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctecs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

const (
	JobStatusExecuting = 0
	JobStatusSuccess   = 1
	JobStatusFail      = 2
)

type GeneralJobHelper struct {
	api *ctecs.JobShowApi
}

func NewGeneralJobHelper(api *ctecs.JobShowApi) *GeneralJobHelper {
	return &GeneralJobHelper{api: api}
}

func (helper GeneralJobHelper) JobLoop(ctx context.Context, credential ctyunsdk.Credential, regionId string, jobId string) (*ctecs.JobShowResponse, error) {
	var resp ctecs.JobShowResponse
	var respError error
	// 暂时就不实现批量
	retryer, _ := NewRetryer(time.Second*5, 60)
	result := retryer.Start(
		func(currentTime int) bool {
			do, requestError := helper.api.Do(ctx, credential, &ctecs.JobShowRequest{
				RegionId: regionId,
				JobId:    jobId,
			})
			if requestError != nil {
				tflog.Error(ctx, "轮询通用任务状态发生异常", map[string]interface{}{"loopTime": currentTime + 1, "err": requestError})
			}
			switch do.Status {
			case JobStatusExecuting:
				return true
			case JobStatusSuccess:
				resp = *do
				return false
			case JobStatusFail:
				resp = *do
				tflog.Error(ctx, "轮询通用任务状态失败", map[string]interface{}{"status": do.Status, "jobStatus": do.JobStatus, "jobId": jobId, "regionId": regionId})
				respError = errors.New("轮询通用任务状态失败，轮询到的任务状态为：" + do.JobStatus)
				return false
			default:
				tflog.Error(ctx, "轮询通用任务状态失败,未知状态", map[string]interface{}{"status": do.Status, "jobStatus": do.JobStatus, "jobId": jobId, "regionId": regionId})
				respError = errors.New("轮询通用任务状态失败，轮询到的任务状态为：" + do.JobStatus)
				return false
			}
		},
	)
	if result.ReturnReason == ReachMaxLoopTime {
		// 这里出来的全都是异常的
		tflog.Info(ctx, "轮询通用任务状态失败，已超过最大轮询次数", map[string]interface{}{"jobId": jobId})
		return nil, errors.New("轮询通用任务状态失败，已超过最大轮询次数，jobId：" + jobId)
	}
	return &resp, respError
}
