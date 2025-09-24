package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

/*
云主机备份策略绑定云主机
*/

func NewCtyunEcsBackupPolicyBindInstances() resource.Resource {
	return &ctyunEcsBackupPolicyBindInstances{}
}

type ctyunEcsBackupPolicyBindInstances struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEcsBackupPolicyBindInstances) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_backup_policy_bind_instances"
}

type CtyunEcsBackupPolicyBindInstancesConfig struct {
	ID             types.String `tfsdk:"id"`
	PolicyID       types.String `tfsdk:"policy_id"`
	RegionID       types.String `tfsdk:"region_id"`
	InstanceIDList types.String `tfsdk:"instance_id_list"`
}

func (c *ctyunEcsBackupPolicyBindInstances) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026751/10033775**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机备份策略id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"instance_id_list": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID列表，多台使用英文逗号分割",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1), // 至少包含一个字符
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9\-_,]+$`),
						"必须是由逗号分隔的UUID列表",
					),
				},
			},
		},
	}
}

func (c *ctyunEcsBackupPolicyBindInstances) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsBackupPolicyBindInstancesConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeBindInstances(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	err = c.checkAfterBindInstances(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunEcsBackupPolicyBindInstances) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEcsBackupPolicyBindInstances) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindInstancesConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未关联") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunEcsBackupPolicyBindInstances) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindInstancesConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeDissociate(ctx, state)
	if err != nil {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociation(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunEcsBackupPolicyBindInstances) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEcsBackupPolicyBindInstances) create(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (err error) {

	params := &ctecs2.CtecsInstanceBackupPolicyBindInstancesRequest{
		RegionID:       plan.RegionID.ValueString(),
		PolicyID:       plan.PolicyID.ValueString(),
		InstanceIDList: plan.InstanceIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsInstanceBackupPolicyBindInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEcsBackupPolicyBindInstances) checkBeforeBindInstances(ctx context.Context, cfg CtyunEcsBackupPolicyBindInstancesConfig) (err error) {
	params := &ctecs2.CtecsListInstanceBackupPolicyRequest{
		RegionID: cfg.RegionID.ValueString(),
		PolicyID: cfg.PolicyID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListInstanceBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.CurrentCount != 1 {
		return fmt.Errorf("备份策略必须存在")
	}

	if cfg.InstanceIDList.ValueString() != "" {
		// 拆分实例ID列表
		instanceIDs := strings.Split(cfg.InstanceIDList.ValueString(), ",")
		// 对每个实例ID进行校验
		for _, instanceID := range instanceIDs {
			//1.使用限制，本接口只支持在拉萨3、上海7、广州6、郴州2、长沙3、北京5、内蒙6、南京3、重庆2、合肥2、成都4、晋中、昆明2、乌鲁木齐27、福州25、衡阳3、长沙37、张家界2、华北2、央企北京1、华东1、上海32、上海33、上海36资源池进行公测
			//3.云主机备份策略可绑定的云主机有配额限制，请在配额限制范围内绑定
			//5.云主机盘限制：不可含有本地盘、共享盘、ISCSI磁盘模式盘

			//查询云主机
			instance_details_resp, err2 := c.meta.Apis.CtEcsApis.EcsInstanceDetailsApi.Do(ctx, c.meta.Credential, &ctecs.EcsInstanceDetailsRequest{
				RegionId:   cfg.RegionID.ValueString(),
				InstanceId: instanceID,
			})
			if err2 != nil {
				// 实例已经被退订的情况
				if err2.ErrorCode() == common.EcsInstanceNotFound {
					return nil
				}
				return err2
			}

			status := instance_details_resp.InstanceStatus
			allowedStatuses := map[string]bool{
				business.EcsStatusRunning: true,
				business.EcsStatusStopped: true,
			}
			//2.备份策略与云主机处于相同的企业项目下，才可进行绑定
			if instance_details_resp.ProjectId != resp.ReturnObj.PolicyList[0].ProjectID {
				return fmt.Errorf("备份策略与云主机处于相同的企业项目下，才可进行绑定")
			}

			//4.云主机存在且状态为运行中或关机，不可重复绑定
			if !allowedStatuses[status] {
				return fmt.Errorf("云主机状态无效(当前:%s)，仅允许在%s或%s状态下绑定备份策略",
					status, business.EcsStatusRunning, business.EcsStatusStopped)
			}
		}

	}

	return
}

// checkAfterBindInstances 绑定后检查
func (c *ctyunEcsBackupPolicyBindInstances) checkAfterBindInstances(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			if plan.InstanceIDList.ValueString() != "" {
				bindID, err = c.getBindingInstances(ctx, plan)
				if err != nil {
					return false
				}
				if bindID == plan.InstanceIDList.ValueString() {
					executeSuccessFlag = true
					return false
				}

			}

			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("云主机策略 %s 和云主机 %s 未关联  regionID： %s", plan.PolicyID.String(), plan.InstanceIDList.ValueString(), plan.RegionID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEcsBackupPolicyBindInstances) checkBeforeDissociate(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (err error) {

	bindID, err := c.getBindingInstances(ctx, plan)
	if err != nil {
		return
	}
	if bindID != plan.InstanceIDList.ValueString() {
		err = fmt.Errorf("云主机策略 %s 和云主机 %s 未关联", plan.PolicyID.String(), plan.InstanceIDList.ValueString())
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEcsBackupPolicyBindInstances) checkAfterDissociation(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			bindID, err = c.getBindingInstances(ctx, plan)
			if err != nil {
				return false
			}
			if bindID != plan.InstanceIDList.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云主机策略 %s 和云主机%s  解绑失败", plan.PolicyID.ValueString(), plan.InstanceIDList.ValueString())
	}
	return nil
}

// dissociate 解绑
func (c *ctyunEcsBackupPolicyBindInstances) delete(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (err error) {
	params := &ctecs2.CtecsInstanceBackupPolicyUnbindInstancesRequest{
		RegionID:       plan.RegionID.ValueString(),
		PolicyID:       plan.PolicyID.ValueString(),
		InstanceIDList: plan.InstanceIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsInstanceBackupPolicyUnbindInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEcsBackupPolicyBindInstances) getBindingInstances(ctx context.Context, plan CtyunEcsBackupPolicyBindInstancesConfig) (instanceIdList string, err error) {
	// 组装请求体
	params := &ctecs2.CtecsListInstanceBackupPolicyBindInstancesRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListInstanceBackupPolicyBindInstancesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	var instanceIds []string
	for _, policy := range resp.ReturnObj.InstancePolicies {
		instanceIds = append(instanceIds, policy.InstanceID)
	}
	instanceIdList = strings.Join(instanceIds, ",")
	return

}

// getAndMerge 查询绑定关系
func (c *ctyunEcsBackupPolicyBindInstances) getAndMerge(ctx context.Context, plan *CtyunEcsBackupPolicyBindInstancesConfig) (err error) {
	policyId, instanceIDList, regionID := plan.PolicyID.ValueString(), plan.InstanceIDList.ValueString(), plan.RegionID.ValueString()
	bindID, err := c.getBindingInstances(ctx, *plan)
	if err != nil {
		return
	}
	if bindID != instanceIDList {
		err = fmt.Errorf("云主机策略 %s 和云主机 %s 未关联  regionID： %s", policyId, instanceIDList, regionID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", policyId, instanceIDList, regionID))
	return
}

// 导入命令：terraform import [配置标识].[导入配置名称] [policyID],[instanceIDList],[regionID]
func (c *ctyunEcsBackupPolicyBindInstances) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsBackupPolicyBindInstancesConfig
	var instanceIDList, policyID, regionID string
	err = terraform_extend.Split(request.ID, &policyID, &instanceIDList, &regionID)
	if err != nil {
		return
	}

	cfg.InstanceIDList = types.StringValue(instanceIDList)
	cfg.PolicyID = types.StringValue(policyID)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
