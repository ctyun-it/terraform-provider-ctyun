package mongodb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &CtyunMongodbAccount{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbAccount{}
	_ resource.ResourceWithImportState = &CtyunMongodbAccount{}
)

func NewCtyunMongodbAccount() resource.Resource {
	return &CtyunMongodbAccount{}
}

type CtyunMongodbAccount struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMongodbAccount) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_account"
}

func (c *CtyunMongodbAccount) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 {inst_id},{account_name}",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MongoDB实例ID",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "账号名称，以字母开头，由字母、数字和下划线组成，长度2-16个字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 16),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"), "实例名称不符合规范"),
				},
			},
			// 实现一个validator方法
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "实例密码，长度为8~26个字符，支持更新，必须包含大写字母、小写字母、数字和特殊字符~!@#%^*_=+ 支持更新",
				Validators: []validator.String{
					validator2.DBPassword(
						8,
						26,
						4,
						"MongoDB",
						"~!@#%^*_=+",
					),
				},
			},
			"database": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "数据库名称，默认为admin 支持更新",
				Default:     stringdefault.StaticString("admin"),
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"roles": schema.ListNestedAttribute{
				Required:    true,
				Description: "角色列表，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"db": schema.StringAttribute{
							Required:    true,
							Description: "数据库名称 支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"role": schema.StringAttribute{
							Required:    true,
							Description: "角色，可选值：read、readWrite、readWriteAnyDatabase等，默认为readWrite 支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("read", "readWrite", "readAnyDatabase", "readWriteAnyDatabase", "dbAdmin", "dbAdminAnyDatabase", "userAdmin", "userAdminAnyDatabase", "clusterAdmin"),
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMongodbAccount) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbAccount) Create(ctx context.Context, req resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan MongodbAccountConfig
	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, &plan)
	if err != nil {
		return
	}
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	// 查询账号信息
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunMongodbAccount) Read(ctx context.Context, req resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state MongodbAccountConfig
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunMongodbAccount) Update(ctx context.Context, req resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan, state MongodbAccountConfig
	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 如果密码变更，更新密码
	if !plan.Password.IsNull() && !plan.Password.Equal(state.Password) {
		err = c.updateAccountPassword(ctx, &plan)
		if err != nil {
			return
		}
	}
	// 如果权限变更，更新权限
	// 检查roles块或者database字段是否变更
	rolesChanged := len(plan.Roles) != len(state.Roles)
	if !rolesChanged {
		for i := range plan.Roles {
			if !plan.Roles[i].Database.Equal(state.Roles[i].Database) ||
				!plan.Roles[i].Role.Equal(state.Roles[i].Role) {
				rolesChanged = true
				break
			}
		}
	}

	if rolesChanged || !plan.Database.Equal(state.Database) {
		err = c.updateAccountPermission(ctx, &plan)
		if err != nil {
			response.Diagnostics.AddError("更新MongoDB账号权限失败", err.Error())
			return
		}
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunMongodbAccount) Delete(ctx context.Context, req resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state MongodbAccountConfig
	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunMongodbAccount) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config MongodbAccountConfig
	var regionID, projectID, instID, name string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		projectID = c.meta.GetExtraIfEmpty(projectID, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &name, &instID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instID, &projectID, &regionID)
		if err != nil {
			return
		}
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if instID == "" {
		err = fmt.Errorf("instdID不能为空")
		return
	}
	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}
	config.ID = types.StringValue(fmt.Sprintf("%s,%s", instID, name))
	config.InstanceID = types.StringValue(instID)
	config.RegionID = types.StringValue(regionID)
	config.ProjectID = types.StringValue(projectID)
	config.Name = types.StringValue(name)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)

}
func (c *CtyunMongodbAccount) create(ctx context.Context, plan MongodbAccountConfig) (err error) {
	// 创建账号权限列表
	var roles []mongodb.MongodbAccountRole

	// 如果定义了roles块，则使用roles块中的权限

	for _, role := range plan.Roles {
		roles = append(roles, mongodb.MongodbAccountRole{
			DB:   role.Database.ValueString(),
			Role: role.Role.ValueString(),
		})
	}

	// 对密码进行Base64编码
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(plan.Password.ValueString()))

	createReq := &mongodb.MongodbCreateAccountRequest{
		AccountName:     plan.Name.ValueString(),
		AccountPassword: encodedPassword,
		Roles:           &roles,
	}

	// 只有当database字段被设置时才添加到请求中
	if !plan.Database.IsNull() && !plan.Database.IsUnknown() {
		createReq.DatabaseName = plan.Database.ValueStringPointer()
	}

	headers := &mongodb.MongodbCreateAccountRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "开始创建MongoDB账号", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbCreateAccountApi.Do(ctx, c.meta.Credential, createReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// getAndMerge 查询账号信息并更新到配置中
func (c *CtyunMongodbAccount) getAndMerge(ctx context.Context, plan *MongodbAccountConfig) (err error) {
	// 解析ID获取inst_id和account_name
	instanceID := plan.InstanceID.ValueString()

	describeReq := &mongodb.MongodbDescribeAccountsRequest{
		ProdInstId: instanceID,
		PageNow:    1,
		PageSize:   1000,
	}

	headers := &mongodb.MongodbDescribeAccountsRequestHeaders{
		RegionID: plan.RegionID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbDescribeAccountsApi.Do(ctx, c.meta.Credential, describeReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	}
	// 这里获取账号信息, 并更新到配置中  循环list 获取账号信息user字段 ==plan.user
	for _, account := range resp.ReturnObj.List {
		if account.User == plan.Name.ValueString() {
			// 如果数据库没有设置，则使用API返回的值
			if plan.Database.IsNull() || plan.Database.IsUnknown() {
				plan.Database = types.StringValue(account.DB)
			}

			// 如果roles块为空，则根据API返回的角色信息填充
			if len(plan.Roles) == 0 && len(account.Roles) > 0 {
				plan.Roles = make([]MongodbAccountRole, len(account.Roles))
				for i, role := range account.Roles {
					plan.Roles[i] = MongodbAccountRole{
						Database: types.StringValue(role.DB),
						Role:     types.StringValue(role.Role),
					}
				}
			}
			plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.InstanceID.ValueString(), plan.Name.ValueString()))
			return nil
		}
	}

	return common.ResourceNotExistError
}

// updateAccountPassword 更新账户密码
func (c *CtyunMongodbAccount) updateAccountPassword(ctx context.Context, plan *MongodbAccountConfig) error {
	// 对密码进行Base64编码
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(plan.Password.ValueString()))

	updatePasswordReq := &mongodb.MongodbUpdateAccountPasswordRequest{
		AccountName:     plan.Name.ValueString(),
		AccountPassword: encodedPassword,
		Database:        plan.Database.ValueString(),
	}

	headers := &mongodb.MongodbUpdateAccountPasswordRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "更新MongoDB账号密码", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	updatePasswordResp, err := c.meta.Apis.SdkMongodbApis.MongodbUpdateAccountPasswordApi.Do(ctx, c.meta.Credential, updatePasswordReq, headers)
	if err != nil {
		return err
	}

	if updatePasswordResp.StatusCode != 800 {
		return fmt.Errorf("API返回错误: %s", *updatePasswordResp.Message)
	}

	return nil
}

// updateAccountPermission 更新账户权限
func (c *CtyunMongodbAccount) updateAccountPermission(ctx context.Context, plan *MongodbAccountConfig) error {
	// 创建权限对象数组
	var roles []mongodb.MongodbAccountRole

	// 如果定义了roles块，则使用roles块中的权限
	for _, role := range plan.Roles {
		roles = append(roles, mongodb.MongodbAccountRole{
			DB:   role.Database.ValueString(),
			Role: role.Role.ValueString(),
		})
	}

	modifyPermissionReq := &mongodb.MongodbModifyAccountPermissionRequest{
		AccountName: plan.Name.ValueString(),
		Roles:       &roles,
	}

	// 只有当database字段被设置时才添加到请求中
	if !plan.Database.IsNull() && !plan.Database.IsUnknown() {
		modifyPermissionReq.DatabaseName = plan.Database.ValueStringPointer()
	}

	headers := &mongodb.MongodbModifyAccountPermissionRequestHeaders{
		RegionID:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}
	if !plan.ProjectID.IsNull() {
		headers.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "更新MongoDB账号权限", map[string]interface{}{
		"instance_id":  plan.InstanceID.ValueString(),
		"account_name": plan.Name.ValueString(),
	})

	modifyPermissionResp, err := c.meta.Apis.SdkMongodbApis.MongodbModifyAccountPermissionApi.Do(ctx, c.meta.Credential, modifyPermissionReq, headers)
	if err != nil {
		return err
	}

	if modifyPermissionResp.StatusCode != 800 {
		return fmt.Errorf("API返回错误: %s", *modifyPermissionResp.Message)
	}

	return nil
}
func (c *CtyunMongodbAccount) delete(ctx context.Context, state MongodbAccountConfig) (err error) {
	deleteReq := &mongodb.MongodbDeleteAccountRequest{
		AccountName:  state.Name.ValueString(),
		DatabaseName: state.Database.ValueString(),
	}

	headers := &mongodb.MongodbDeleteAccountRequestHeaders{
		RegionID:   state.RegionID.ValueString(),
		ProdInstId: state.InstanceID.ValueString(),
	}
	if !state.ProjectID.IsNull() {
		headers.ProjectID = state.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "删除MongoDB账号", map[string]interface{}{
		"instance_id":  state.InstanceID.ValueString(),
		"account_name": state.Name.ValueString(),
	})

	resp, err := c.meta.Apis.SdkMongodbApis.MongodbDeleteAccountApi.Do(ctx, c.meta.Credential, deleteReq, headers)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

func (c *CtyunMongodbAccount) checkBeforeCreate(ctx context.Context, m *MongodbAccountConfig) (err error) {
	return
}

type MongodbAccountRole struct {
	Database types.String `tfsdk:"db"`
	Role     types.String `tfsdk:"role"`
}

type MongodbAccountConfig struct {
	ID         types.String         `tfsdk:"id"`
	InstanceID types.String         `tfsdk:"instance_id"`
	RegionID   types.String         `tfsdk:"region_id"`
	ProjectID  types.String         `tfsdk:"project_id"`
	Name       types.String         `tfsdk:"name"`
	Password   types.String         `tfsdk:"password"`
	Database   types.String         `tfsdk:"database"`
	Roles      []MongodbAccountRole `tfsdk:"roles"`
}
