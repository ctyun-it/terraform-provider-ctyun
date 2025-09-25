package ebs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
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
将快照策略和云硬盘绑定
*/

func NewCtyunEbsSnapshotPolicyAssociation() resource.Resource {
	return &ctyunEbsSnapshotPolicyAssociation{}
}

type ctyunEbsSnapshotPolicyAssociation struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEbsSnapshotPolicyAssociation) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ebs_snapshot_policy_association"
}

type CtyunEbsSnapshotPolicyAssociationConfig struct {
	ID               types.String `tfsdk:"id"`
	SnapshotPolicyID types.String `tfsdk:"snapshot_policy_id"`
	RegionID         types.String `tfsdk:"region_id"`
	DiskIDList       types.String `tfsdk:"disk_id_list"`
}

func (c *ctyunEbsSnapshotPolicyAssociation) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027696/10118856**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"snapshot_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "云硬盘自动快照策略id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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

func (c *ctyunEbsSnapshotPolicyAssociation) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEbsSnapshotPolicyAssociationConfig
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

func (c *ctyunEbsSnapshotPolicyAssociation) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEbsSnapshotPolicyAssociation) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyAssociationConfig
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

func (c *ctyunEbsSnapshotPolicyAssociation) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEbsSnapshotPolicyAssociationConfig
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

func (c *ctyunEbsSnapshotPolicyAssociation) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEbsSnapshotPolicyAssociation) create(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {

	params := &ctebs2.EbsApplyPolicyEbsSnapRequest{
		RegionID:         plan.RegionID.ValueString(),
		SnapshotPolicyID: plan.SnapshotPolicyID.ValueString(),
		TargetDiskIDs:    plan.DiskIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsApplyPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) checkBeforeBindDisks(ctx context.Context, cfg CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	// 获取实例详情
	params := &ctebs2.EbsQueryPolicyEbsSnapRequest{
		RegionID:         cfg.RegionID.ValueString(),
		SnapshotPolicyID: cfg.SnapshotPolicyID.ValueStringPointer(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.SnapshotPolicyTotalCount != 1 {
		return fmt.Errorf("自动快照策略必须存在")
	}

	if cfg.DiskIDList.ValueString() != "" {
		// 拆分实例ID列表
		diskIds := strings.Split(cfg.DiskIDList.ValueString(), ",")
		// 对每个实例ID进行校验
		for _, diskId := range diskIds {
			//云硬盘只能绑定一个自动快照策略，如果一个云硬盘已经绑定了自动快照策略，再次调用该接口绑定新的自动快照策略，旧的自动快照策略会自动解绑。

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
func (c *ctyunEbsSnapshotPolicyAssociation) checkAfterBindDisks(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var snapshotPolicyID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			if plan.DiskIDList.ValueString() != "" {
				snapshotPolicyID, err = c.getBindingDisks(ctx, plan)
				if err != nil {
					return false
				}
				if snapshotPolicyID == plan.SnapshotPolicyID.ValueString() {
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
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联  regionID： %s", plan.SnapshotPolicyID.String(), plan.DiskIDList.ValueString(), plan.RegionID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEbsSnapshotPolicyAssociation) checkBeforeDissociate(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {

	snapshotPolicyID, err := c.getBindingDisks(ctx, plan)
	if err != nil {
		return
	}
	if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联", plan.SnapshotPolicyID.String(), plan.DiskIDList.ValueString())
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEbsSnapshotPolicyAssociation) checkAfterDissociation(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	var executeSuccessFlag bool
	var snapshotPolicyID string
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			snapshotPolicyID, err = c.getBindingDisks(ctx, plan)
			if err != nil {
				return false
			}
			if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘%s  解绑失败", plan.SnapshotPolicyID.ValueString(), plan.DiskIDList.ValueString())
	}
	return nil
}

// dissociate 解绑
func (c *ctyunEbsSnapshotPolicyAssociation) delete(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	params := &ctebs2.EbsCancelPolicyEbsSnapRequest{
		RegionID:      plan.RegionID.ValueString(),
		TargetDiskIDs: plan.DiskIDList.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEbsApis.EbsCancelPolicyEbsSnapApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.StatusCode == common.NormalStatusCode {
		return
	}

	return
}

func (c *ctyunEbsSnapshotPolicyAssociation) getBindingDisks(ctx context.Context, plan CtyunEbsSnapshotPolicyAssociationConfig) (snapshotPolicyID string, err error) {
	// 拆分磁盘ID列表
	diskIds := strings.Split(plan.DiskIDList.ValueString(), ",")
	var firstPolicyID string
	firstPolicyIDSet := false

	for _, diskId := range diskIds {
		diskId = strings.TrimSpace(diskId)
		if diskId == "" {
			continue
		}

		// 组装请求体
		params := &ctebs2.EbsQueryEbsByIDRequest{
			RegionID: plan.RegionID.ValueStringPointer(),
			DiskID:   diskId,
		}

		// 调用API
		resp, err := c.meta.Apis.SdkCtEbsApis.EbsQueryEbsByIDApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return "", err
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return "", err
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return "", err
		}

		// 安全处理SnapshotPolicyID
		var currentPolicyID string
		if resp.ReturnObj.SnapshotPolicyID != nil {
			currentPolicyID = *resp.ReturnObj.SnapshotPolicyID
		}
		// 如果resp.ReturnObj.SnapshotPolicyID是nil，currentPolicyID保持为空字符串

		// 设置参考策略ID
		if !firstPolicyIDSet {
			firstPolicyID = currentPolicyID
			firstPolicyIDSet = true
		} else if firstPolicyID != currentPolicyID {
			// 磁盘间绑定的策略不一致
			return currentPolicyID, nil
		}
	}

	return firstPolicyID, nil
}

// getAndMerge 查询绑定关系
func (c *ctyunEbsSnapshotPolicyAssociation) getAndMerge(ctx context.Context, plan *CtyunEbsSnapshotPolicyAssociationConfig) (err error) {
	policyId, diskIDList, regionID := plan.SnapshotPolicyID.ValueString(), plan.DiskIDList.ValueString(), plan.RegionID.ValueString()
	snapshotPolicyID, err := c.getBindingDisks(ctx, *plan)
	if err != nil {
		return
	}
	if snapshotPolicyID != plan.SnapshotPolicyID.ValueString() {
		err = fmt.Errorf("云硬盘自动快照策略 %s 和云硬盘 %s 未关联  regionID： %s", policyId, diskIDList, regionID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", policyId, diskIDList, regionID))
	return
}

// 导入命令：terraform import [配置标识].[导入配置名称] [policyID],[diskIDList],[regionID]
func (c *ctyunEbsSnapshotPolicyAssociation) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEbsSnapshotPolicyAssociationConfig
	var diskIDList, policyID, regionID string
	err = terraform_extend.Split(request.ID, &policyID, &diskIDList, &regionID)
	if err != nil {
		return
	}

	cfg.DiskIDList = types.StringValue(diskIDList)
	cfg.SnapshotPolicyID = types.StringValue(policyID)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
