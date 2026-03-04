package zos

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctzos"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type ctyunZosSubscribe struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
	orderLooper   *business.OrderLooper
}

func (c *ctyunZosSubscribe) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_subscribe"
}

func (c *ctyunZosSubscribe) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)
	c.orderLooper = business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
}

func NewctyunZosSubscribe() resource.Resource {
	return &ctyunZosSubscribe{}
}

func (c *ctyunZosSubscribe) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("订阅对象存储服务", "对象存储（CT-ZOS，Zettabyte Object Storage）", "https://www.ctyun.cn/document/10027350"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			//"id": schema.StringAttribute{
			//	Computed:    true,
			//	Description: "zos ID",
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.UseStateForUnknown(),
			//	},
			//},
		}}
}

func (c *ctyunZosSubscribe) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunZosSubscribeConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	flag := false
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		if strings.Contains(err.Error(), "该资源池已开通zos服务，请勿重复开通！") {
			err = nil
			flag = true
			response.Diagnostics.AddWarning("该资源池已开通zos服务，请勿重复开通！", "因该资源池已开通zos服务，直接跳过开通步骤。")
		} else {
			return
		}
	}
	if !flag {
		err = c.subscribe(ctx, &plan)
		if err != nil {
			return
		}
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunZosSubscribe) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *ctyunZosSubscribe) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *ctyunZosSubscribe) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *ctyunZosSubscribe) checkBeforeCreate(ctx context.Context, plan CtyunZosSubscribeConfig) error {
	// 根据region id 查看是否支持zos
	regionListResp, err := c.meta.Apis.SdkCtZosApis.ZosListRegionsApi.Do(ctx, c.meta.SdkCredential, &ctzos.ZosListRegionsRequest{})
	if err != nil {
		return err
	} else if regionListResp == nil {
		err = fmt.Errorf("查询所有支持zos资源池失败，接口返回nil。可联系研发确认原因")
		return err
	} else if regionListResp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", regionListResp.Message, regionListResp.Description)
		return err
	} else if regionListResp.ReturnObj == nil || len(regionListResp.ReturnObj) == 0 {
		err = fmt.Errorf("查询所有支持zos资源池失败，接口返回为空。可联系研发确认原因")
		return err
	}
	flag := false
	for _, region := range regionListResp.ReturnObj {
		if region.RegionID == plan.RegionID.ValueString() {
			flag = true
			break
		}
	}
	if !flag {
		err = fmt.Errorf("资源池ID %s 不支持zos服务", plan.RegionID.ValueString())
		return err
	}
	// 根据region id 确认是否已经开通zos服务
	status, err := c.getZosOrderStatus(ctx, plan)
	if err != nil {
		return err
	}
	// 如果开通了，则提示用户无需重复开通
	if status {
		err = fmt.Errorf("该资源池已开通zos服务，请勿重复开通！")
		return err
	}
	return nil
}

func (c *ctyunZosSubscribe) subscribe(ctx context.Context, plan *CtyunZosSubscribeConfig) error {
	resp, err := c.requestSubscribe(ctx, plan)
	if err != nil {
		return err
	}
	// 判断是否开通成功
	err = c.checkSubscribeResult(ctx, plan, resp.ReturnObj.MasterOrderID)
	if err != nil {
		return err
	}
	// 开通成功后，再通过clientToken 获取zos id
	//resp, err = c.requestSubscribe(ctx, plan)
	//if err != nil {
	//	return err
	//}
	//plan.ID = types.StringValue(resp.ReturnObj.Resources.ZosID)
	return nil
}

func (c *ctyunZosSubscribe) checkSubscribeResult(ctx context.Context, plan *CtyunZosSubscribeConfig, masterOrderID string) error {
	// 	// 根据order id 轮询订单，确认订单是否结束
	err := c.orderLooper.WaitOrderFinish(ctx, c.meta.Credential, masterOrderID)
	if err != nil {
		return err
	}

	// 查询开通状态
	err = c.statusLoop(ctx, *plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunZosSubscribe) getZosOrderStatus(ctx context.Context, plan CtyunZosSubscribeConfig) (bool, error) {
	params := &ctzos.ZosGetOssServiceStatusRequest{
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetOssServiceStatusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return false, err
	} else if resp == nil {
		err = fmt.Errorf("查询资源池 %s 是否开通zos服务失败，接口返回nil。可联系研发确认原因", plan.RegionID.ValueString())
		return false, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return false, err
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("查询资源池 %s 是否开通zos服务失败，接口返回为空。可联系研发确认原因", plan.RegionID.ValueString())
		return false, err
	}
	if resp.ReturnObj.State == "true" {
		return true, nil
	}
	return false, nil

}

func (c *ctyunZosSubscribe) statusLoop(ctx context.Context, plan CtyunZosSubscribeConfig) error {
	retryer, err := business.NewRetryer(time.Second*10, 20)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			status, err2 := c.getZosOrderStatus(ctx, plan)
			if err2 != nil {
				err = err2
				return false
			}
			if status {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，zos仍未开通成功！")
	}
	return err
}

func (c *ctyunZosSubscribe) requestSubscribe(ctx context.Context, plan *CtyunZosSubscribeConfig) (*ctzos.ZosNewOssResponse, error) {
	params := &ctzos.ZosNewOssRequest{
		RegionID:    plan.RegionID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosNewOssApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("开通资源池 %s zos服务失败，接口返回nil。可联系研发确认原因", plan.RegionID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode && resp.ErrorCode != "oss.order.inProgress" {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	}
	return resp, err
}

type CtyunZosSubscribeConfig struct {
	RegionID types.String `tfsdk:"region_id"`
}
