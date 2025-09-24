package ecs

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
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

/*
云主机备份策略绑定存储库
*/

func NewCtyunEcsBackupPolicyBindRepo() resource.Resource {
	return &ctyunEcsBackupPolicyBindRepo{}
}

type ctyunEcsBackupPolicyBindRepo struct {
	meta *common.CtyunMetadata
}

func (c *ctyunEcsBackupPolicyBindRepo) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ecs_backup_policy_bind_repo"
}

type CtyunEcsBackupPolicyBindRepoConfig struct {
	ID           types.String `tfsdk:"id"`
	PolicyID     types.String `tfsdk:"policy_id"`
	RegionID     types.String `tfsdk:"region_id"`
	RepositoryID types.String `tfsdk:"repository_id"`
}

func (c *ctyunEcsBackupPolicyBindRepo) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026751/10235038**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机备份策略id",
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

			"repository_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机备份存储库ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.UUID(),
				},
			},
		},
	}
}

func (c *ctyunEcsBackupPolicyBindRepo) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunEcsBackupPolicyBindRepoConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeBindRepo(ctx, plan)
	if err != nil {
		return
	}

	// 实际创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	err = c.checkAfterBindRepo(ctx, plan)
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

func (c *ctyunEcsBackupPolicyBindRepo) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunEcsBackupPolicyBindRepo) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindRepoConfig
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

func (c *ctyunEcsBackupPolicyBindRepo) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunEcsBackupPolicyBindRepoConfig
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

func (c *ctyunEcsBackupPolicyBindRepo) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// create 创建
func (c *ctyunEcsBackupPolicyBindRepo) create(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (err error) {

	params := &ctecs2.CtecsInstanceBackupPolicyBindRepoRequest{
		RegionID:     plan.RegionID.ValueString(),
		PolicyID:     plan.PolicyID.ValueString(),
		RepositoryID: plan.RepositoryID.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsInstanceBackupPolicyBindRepoApi.Do(ctx, c.meta.SdkCredential, params)
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

func (c *ctyunEcsBackupPolicyBindRepo) checkBeforeBindRepo(ctx context.Context, cfg CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	params := &ctecs2.CtecsListInstanceBackupPolicyRequest{
		RegionID: cfg.RegionID.ValueString(),
		PolicyID: cfg.PolicyID.ValueString(),
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsListInstanceBackupPolicyApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	} else if resp.ReturnObj.CurrentCount != 1 {
		return fmt.Errorf("备份策略必须存在")
	}

	if cfg.RepositoryID.ValueString() != "" {
		//1.使用限制，本接口只支持在拉萨3、上海7、广州6、郴州2、长沙3、北京5、内蒙6、南京3、重庆2、合肥2、成都4、晋中、昆明2、乌鲁木齐27、福州25、衡阳3、长沙37、张家界2、华北2、央企北京1、华东1、上海32、上海33、上海36资源池进行公测
		//2.备份策略与存储库必须存在
		params := &ctecs2.CtecsListInstanceBackupRepoRequest{
			RegionID:     cfg.RegionID.ValueString(),
			RepositoryID: cfg.RepositoryID.ValueString(),
		}
		// 调用API
		respRepo, err := c.meta.Apis.SdkCtEcsApis.CtecsListInstanceBackupRepoApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if respRepo.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return err
		} else if respRepo.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return err
		} else if len(respRepo.ReturnObj.Results) != 1 || respRepo.ReturnObj.Results[0].RepositoryID == "" {
			return fmt.Errorf("存储库必须存在")
		}
		//3.备份策略与存储库处于相同的企业项目下
		if respRepo.ReturnObj.Results[0].ProjectID != resp.ReturnObj.PolicyList[0].ProjectID {
			return fmt.Errorf("备份策略与存储库需处于相同的企业项目下")
		}
		//4.备份策略在此之前未绑定存储库，且该存储库未到期
		if len(resp.ReturnObj.PolicyList[0].RepositoryList) != 0 {
			return fmt.Errorf("备份策略已绑定存储库，不可再绑定")
		}
		if respRepo.ReturnObj.Results[0].Expired {
			return fmt.Errorf("该存储库已到期")
		}
	}

	return
}

// checkAfterBindRepo 绑定后检查
func (c *ctyunEcsBackupPolicyBindRepo) checkAfterBindRepo(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			if plan.RepositoryID.ValueString() != "" {
				hasBind, err := c.getBindingRepos(ctx, plan)
				if err != nil {
					return false
				}
				if hasBind {
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
		err = fmt.Errorf("云主机备份策略 %s 和存储库 %s 未关联  regionID： %s", plan.PolicyID.String(), plan.RepositoryID.ValueString(), plan.RegionID.ValueString())
	}
	return nil
}

// checkBeforeDissociate 解绑前检查
func (c *ctyunEcsBackupPolicyBindRepo) checkBeforeDissociate(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	hasBind, err := c.getBindingRepos(ctx, plan)
	if err != nil {
		return
	}
	if !hasBind {
		err = fmt.Errorf("云主机备份策略 %s 和存储库 %s 未关联", plan.PolicyID.String(), plan.RepositoryID.ValueString())
		return
	}
	return
}

// checkAfterDissociation 解绑后检查
func (c *ctyunEcsBackupPolicyBindRepo) checkAfterDissociation(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			hasBind, err := c.getBindingRepos(ctx, plan)
			if err != nil {
				return false
			}
			if !hasBind {
				executeSuccessFlag = true
				return false
			}
			return true
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		return fmt.Errorf("云主机备份策略 %s 和存储库%s  解绑失败", plan.PolicyID.ValueString(), plan.RepositoryID.ValueString())
	}
	return nil
}

// dissociate 解绑
func (c *ctyunEcsBackupPolicyBindRepo) delete(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	params := &ctecs2.CtecsInstanceBackupPolicyUnbindRepoRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
	}

	// 创建实例
	resp, err := c.meta.Apis.SdkCtEcsApis.CtecsInstanceBackupPolicyUnbindRepoApi.Do(ctx, c.meta.SdkCredential, params)
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

func (c *ctyunEcsBackupPolicyBindRepo) getBindingRepos(ctx context.Context, plan CtyunEcsBackupPolicyBindRepoConfig) (hasBind bool, err error) {

	params := &ctecs2.CtecsListInstanceBackupPolicyRequest{
		RegionID: plan.RegionID.ValueString(),
		PolicyID: plan.PolicyID.ValueString(),
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
	}

	// 解析返回值
	for _, policy := range resp.ReturnObj.PolicyList {
		// 处理备份库列表
		repositoryList := make([]repositoryListModel, len(policy.RepositoryList))
		for i, repo := range policy.RepositoryList {
			repoItem := repositoryListModel{
				RepositoryID:   types.StringValue(repo.RepositoryID),
				RepositoryName: types.StringValue(repo.RepositoryName),
			}
			repositoryList[i] = repoItem
			if plan.RepositoryID.ValueString() == repo.RepositoryID {
				return true, nil
			}
		}
	}

	return false, nil

}

// getAndMerge 查询绑定关系
func (c *ctyunEcsBackupPolicyBindRepo) getAndMerge(ctx context.Context, plan *CtyunEcsBackupPolicyBindRepoConfig) (err error) {
	policyId, repositoryID, regionID := plan.PolicyID.ValueString(), plan.RepositoryID.ValueString(), plan.RegionID.ValueString()
	hasBind, err := c.getBindingRepos(ctx, *plan)

	if err != nil {
		return
	}
	if !hasBind {
		err = fmt.Errorf("云主机备份策略 %s 和存储库 %s 未关联  regionID： %s", policyId, repositoryID, regionID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", policyId, repositoryID, regionID))
	return
}

// 导入命令：terraform import [配置标识].[导入配置名称] [policyID],[repositoryID],[regionID]
func (c *ctyunEcsBackupPolicyBindRepo) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunEcsBackupPolicyBindRepoConfig
	var repositoryID, policyID, regionID string
	err = terraform_extend.Split(request.ID, &policyID, &repositoryID, &regionID)
	if err != nil {
		return
	}

	cfg.RepositoryID = types.StringValue(repositoryID)
	cfg.PolicyID = types.StringValue(policyID)
	cfg.RegionID = types.StringValue(regionID)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}
