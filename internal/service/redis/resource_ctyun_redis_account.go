package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunRedisAccount{}
	_ resource.ResourceWithConfigure   = &ctyunRedisAccount{}
	_ resource.ResourceWithImportState = &ctyunRedisAccount{}
)

type ctyunRedisAccount struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisAccount() resource.Resource {
	return &ctyunRedisAccount{}
}

func (c *ctyunRedisAccount) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_account"
}

type CtyunRedisAccountConfig struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	InstanceId  types.String `tfsdk:"instance_id"`
	RegionId    types.String `tfsdk:"region_id"`
	Password    types.String `tfsdk:"password"`
	Description types.String `tfsdk:"description"`
	Privilege   types.String `tfsdk:"privilege"`
}

func (c *ctyunRedisAccount) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10403139`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "账户名称，规则如下：\n以英文字母、数字、下划线开头，且只能由英文字母、数字、句点、中划线、下划线组成。\n长度3-64。\n名称不可重复。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$`),
						"必须以英文字母、数字、下划线开头，只能包含英文字母、数字、句点、中划线、下划线",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID。",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "密码，规则如下：\n长度8-26字符。\n必须同时包含大写字母、小写字母、数字和英文格式特殊符号(@%^*_+!$-=.)中的至少三种类型。\n不能有空格。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 26),
				},
				Sensitive: true,
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "账户描述信息。支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(255),
				},
			},
			"privilege": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "账户权限，可选值：ro(只读)、rw(读写) 支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("ro", "rw"),
				},
			},
		},
	}
}

func (c *ctyunRedisAccount) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建账户
	err = c.create(ctx, plan)
	if err != nil {
		return
	}

	// 查询创建结果
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisAccount) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisAccount) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// tf文件中的
	var plan CtyunRedisAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunRedisAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新账户信息
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}

	state.Password = plan.Password
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisAccount) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除账户
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunRedisAccount) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisAccount) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [账户名称],[实例ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRedisAccountConfig

	var accountName, instanceId, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &accountName, &instanceId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &accountName, &instanceId, &regionId)
		if err != nil {
			return
		}
	}

	if accountName == "" {
		err = fmt.Errorf("账户名称不能为空")
		return
	}
	if instanceId == "" {
		err = fmt.Errorf("实例ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	cfg.RegionId = types.StringValue(regionId)
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.Name = types.StringValue(accountName)

	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建账户
func (c *ctyunRedisAccount) create(ctx context.Context, plan CtyunRedisAccountConfig) (err error) {
	params := &ctgdcs2.Dcs2CreateAccountRequest{
		RegionId:           plan.RegionId.ValueString(),
		ProdInstId:         plan.InstanceId.ValueString(),
		AccountName:        plan.Name.ValueString(),
		AccountPassword:    plan.Password.ValueString(),
		AccountDescription: plan.Description.ValueString(),
		AccountPrivilege:   plan.Privilege.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2CreateAccountApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// update 更新账户信息
func (c *ctyunRedisAccount) update(ctx context.Context, plan, state CtyunRedisAccountConfig) (err error) {
	// 如果密码变更，使用密码修改API
	if plan.Password.ValueString() != state.Password.ValueString() {
		passwordParams := &ctgdcs2.Dcs2ModifyAccountPasswordRequest{
			RegionId:           state.RegionId.ValueString(),
			ProdInstId:         state.InstanceId.ValueString(),
			AccountName:        plan.Name.ValueString(),
			AccountPassword:    plan.Password.ValueString(),
			OldAccountPassword: state.Password.ValueString(),
		}

		passwordResp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifyAccountPasswordApi.Do(ctx, c.meta.SdkCredential, passwordParams)
		if err != nil {
			return err
		} else if passwordResp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", passwordResp.Message)
		} else if passwordResp.ReturnObj == nil {
			return common.InvalidReturnObjError
		}
	}

	// 如果描述变更，使用描述修改API
	if plan.Description.ValueString() != state.Description.ValueString() {
		descriptionParams := &ctgdcs2.Dcs2ModifyAccountDescriptionRequest{
			RegionId:           state.RegionId.ValueString(),
			ProdInstId:         state.InstanceId.ValueString(),
			AccountName:        plan.Name.ValueString(),
			AccountDescription: plan.Description.ValueString(),
		}

		descriptionResp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifyAccountDescriptionApi.Do(ctx, c.meta.SdkCredential, descriptionParams)
		if err != nil {
			return err
		} else if descriptionResp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", descriptionResp.Message)
		} else if descriptionResp.ReturnObj == nil {
			return common.InvalidReturnObjError
		}
	}

	// 如果权限变更，使用权限修改API
	if plan.Privilege.ValueString() != state.Privilege.ValueString() {
		privilegeParams := &ctgdcs2.Dcs2GrantAccountPrivilegeRequest{
			RegionId:         state.RegionId.ValueString(),
			ProdInstId:       state.InstanceId.ValueString(),
			AccountName:      plan.Name.ValueString(),
			AccountPrivilege: plan.Privilege.ValueString(),
		}

		privilegeResp, err := c.meta.Apis.SdkDcs2Apis.Dcs2GrantAccountPrivilegeApi.Do(ctx, c.meta.SdkCredential, privilegeParams)
		if err != nil {
			return err
		} else if privilegeResp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s", privilegeResp.Message)
		} else if privilegeResp.ReturnObj == nil {
			return common.InvalidReturnObjError
		}
	}

	return
}

// destroy 删除账户
func (c *ctyunRedisAccount) destroy(ctx context.Context, plan CtyunRedisAccountConfig) (err error) {
	params := &ctgdcs2.Dcs2DeleteAccountRequest{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		AccountName: plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DeleteAccountApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getAndMerge 从远端查询账户信息
func (c *ctyunRedisAccount) getAndMerge(ctx context.Context, plan *CtyunRedisAccountConfig) (err error) {
	params := &ctgdcs2.Dcs2DescribeAccountsRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeAccountsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Rows == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 查找匹配的账户信息
	var accountData *ctgdcs2.Dcs2DescribeAccountsReturnObjRowsResponse
	for _, account := range resp.ReturnObj.Rows {
		if account.Name == plan.Name.ValueString() {
			accountData = account
			break
		}
	}

	if accountData == nil {
		err = fmt.Errorf("account %s not found", plan.Name.ValueString())
		return
	}

	// 设置权限信息
	if accountData.RawType != "" {
		plan.Privilege = types.StringValue(accountData.RawType)
	} else {
		plan.Privilege = types.StringValue("rw") // 默认权限
	}

	// 设置描述信息
	if accountData.AccountDescription == "" {
		plan.Description = types.StringNull()
	} else {
		plan.Description = types.StringValue(accountData.AccountDescription)
	}

	id := fmt.Sprintf("%s,%s,%s", plan.Name.ValueString(), plan.InstanceId.ValueString(), plan.RegionId.ValueString())
	plan.Id = types.StringValue(id)
	return
}
