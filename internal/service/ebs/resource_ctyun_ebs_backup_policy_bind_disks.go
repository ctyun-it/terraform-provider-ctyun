package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
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
云硬盘备份策略绑定云硬盘
*/

func NewCtyunEcsBackupPolicyBindDisks() resource.Resource {
	return &ctyunEcsBackupPolicyBindDisks{}
}

type ctyunEcsBackupPolicyBindDisks struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEcsBackupPolicyBindDisks) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_backup_policy_bind_disks"
}

type CtyunEcsBackupPolicyBindDisksConfig struct {
	ID         types.String `tfsdk:"id"`
	PolicyID   types.String `tfsdk:"policy_id"`
	RegionID   types.String `tfsdk:"region_id"`
	DiskIDList types.String `tfsdk:"disk_id_list"`
}

func (c *ctyunEcsBackupPolicyBindDisks) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026752/10037452**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘备份策略id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
			"disk_id_list": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘ID列表，多台使用英文逗号分割",
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

func (c *ctyunEcsBackupPolicyBindDisks) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsBackupPolicyBindDisksConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeBindDisks(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	err = c.checkAfterBindDisks(ctx, plan)
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

func (c *ctyunEcsBackupPolicyBindDisks) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEcsBackupPolicyBindDisks) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindDisksConfig
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

func (c *ctyunEcsBackupPolicyBindDisks) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindDisksConfig
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

func (c *ctyunEcsBackupPolicyBindDisks) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEcsBackupPolicyBindDisks) create(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (err error) {

	params := &ctebsbackup.EbsbackupEbsBackupPolicyBindDisksRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
		DiskIDs:  plan.DiskIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupEbsBackupPolicyBindDisksApi.Do(ctx, c.meta.SdkCredential, params)
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

func (c *ctyunEcsBackupPolicyBindDisks) checkBeforeBindDisks(ctx context.Context, cfg CtyunEcsBackupPolicyBindDisksConfig) (err error) {
	params := &ctebsbackup.EbsbackupListBackupPolicyRequest{
		RegionID: cfg.RegionID.ValueString(),
		PolicyID: cfg.PolicyID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
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

	if cfg.DiskIDList.ValueString() != "" {
		// 拆分实例ID列表
		diskIds := strings.Split(cfg.DiskIDList.ValueString(), ",")
		// 对每个实例ID进行校验
		for _, diskId := range diskIds {
			//云硬盘只能绑定一个备份策略，如果一个云硬盘已经绑定了备份策略，再次调用该接口绑定新的备份策略，旧的备份策略会自动解绑。

			//查询云硬盘
			_, err2 := c.meta.Apis.CtEbsApis.EbsShowApi.Do(ctx, c.meta.Credential, &ctebs.EbsShowRequest{
				RegionId: cfg.RegionID.ValueString(),
				DiskId:   diskId,
			})
			if err2 != nil {
				// 实例已经被退订的情况
				if err2.ErrorCode() == common.EcsInstanceNotFound {
					return nil
				}
				return err2
			}

		}

	}

	return
}

// checkAfterBindDisks 绑定后检查
func (c *ctyunEcsBackupPolicyBindDisks) checkAfterBindDisks(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			if plan.DiskIDList.ValueString() != "" {
				bindID, err = c.getBindingDisks(ctx, plan)
				if err != nil {
					return false
				}
				if bindID == plan.DiskIDList.ValueString() {
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
		err = fmt.Errorf("云硬盘策略 %s 和云硬盘 %s 未关联  regionID： %s", plan.PolicyID.String(), plan.DiskIDList.ValueString(), plan.RegionID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEcsBackupPolicyBindDisks) checkBeforeDissociate(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (err error) {

	bindID, err := c.getBindingDisks(ctx, plan)
	if err != nil {
		return
	}
	if bindID != plan.DiskIDList.ValueString() {
		err = fmt.Errorf("云硬盘策略 %s 和云硬盘 %s 未关联", plan.PolicyID.String(), plan.DiskIDList.ValueString())
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEcsBackupPolicyBindDisks) checkAfterDissociation(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (err error) {
	var executeSuccessFlag bool
	var bindID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			bindID, err = c.getBindingDisks(ctx, plan)
			if err != nil {
				return false
			}
			if bindID != plan.DiskIDList.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云硬盘策略 %s 和云硬盘%s  解绑失败", plan.PolicyID.ValueString(), plan.DiskIDList.ValueString())
	}
	return nil
}

// dissociate 解绑
func (c *ctyunEcsBackupPolicyBindDisks) delete(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (err error) {
	params := &ctebsbackup.EbsbackupEbsBackupPolicyUnbindDisksRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
		DiskIDs:  plan.DiskIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupEbsBackupPolicyUnbindDisksApi.Do(ctx, c.meta.SdkCredential, params)
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

func (c *ctyunEcsBackupPolicyBindDisks) getBindingDisks(ctx context.Context, plan CtyunEcsBackupPolicyBindDisksConfig) (instanceIdList string, err error) {
	// 组装请求体
	params := &ctebsbackup.EbsbackupListEbsBackupPolicyDisksRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.CtEbsBackupApis.EbsbackupListEbsBackupPolicyDisksApi.Do(ctx, c.meta.SdkCredential, params)
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
	var diskIds []string
	for _, policy := range resp.ReturnObj.DiskList {
		diskIds = append(diskIds, policy.DiskID)
	}
	instanceIdList = strings.Join(diskIds, ",")
	return

}

// getAndMerge 查询绑定关系
func (c *ctyunEcsBackupPolicyBindDisks) getAndMerge(ctx context.Context, plan *CtyunEcsBackupPolicyBindDisksConfig) (err error) {
	policyId, diskIDList, regionID := plan.PolicyID.ValueString(), plan.DiskIDList.ValueString(), plan.RegionID.ValueString()
	bindID, err := c.getBindingDisks(ctx, *plan)
	if err != nil {
		return
	}
	if bindID != diskIDList {
		err = fmt.Errorf("云硬盘策略 %s 和云硬盘 %s 未关联  regionID： %s", policyId, diskIDList, regionID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", policyId, diskIDList, regionID))
	return
}

// 导入命令：terraform import [配置标识].[导入配置名称] [policyID],[diskIDList],[regionID]
func (c *ctyunEcsBackupPolicyBindDisks) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsBackupPolicyBindDisksConfig
	var diskIDList, policyID, regionID string
	err = terraform_extend.Split(request.ID, &policyID, &diskIDList, &regionID)
	if err != nil {
		return
	}

	cfg.DiskIDList = types.StringValue(diskIDList)
	cfg.PolicyID = types.StringValue(policyID)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
