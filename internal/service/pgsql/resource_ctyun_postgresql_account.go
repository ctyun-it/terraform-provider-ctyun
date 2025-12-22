package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunPostgresqlAccount{}
	_ resource.ResourceWithConfigure   = &CtyunPostgresqlAccount{}
	_ resource.ResourceWithImportState = &CtyunPostgresqlAccount{}
)

type CtyunPostgresqlAccount struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPostgresqlAccount) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_account"
}
func NewCtyunPostgresqlAccount() resource.Resource {
	return &CtyunPostgresqlAccount{}
}

func (c *CtyunPostgresqlAccount) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPostgresqlAccount) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[projectID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPostgresqlAccountConfig
	var regionID, projectID, name, instID string
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
	config.ID = types.StringValue(instID + "-" + name)
	config.InstID = types.StringValue(instID)
	config.RegionID = types.StringValue(regionID)
	config.ProjectID = types.StringValue(projectID)
	config.Name = types.StringValue(name)
	err = c.getAndMergePostgresqlAccount(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPostgresqlAccount) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10161317",
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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
				Description: "数据库账号名称，,格式限制: 1.名称唯一；2. 以字母开头，以字母或数字结尾；3.由小写字母、数字或下划线组成；4. 长度：2~63个字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 63),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "数据库账号密码，支持更新。由大写字母、小写字母、特殊字符、数字中三种或者三种以上组成(特殊字符：@!#$%^&*()_-=)",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 32),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备注，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"user_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(business.PgsqlAccountTypeNormal),
				Description: "账号类型，取值范围：normal-普通账号，advanced-高权限账号。默认为普通账号",
				Validators: []validator.String{
					stringvalidator.OneOf(business.PgsqlAccountTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_lock": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "判断账号是否需要锁定，取值范围：true-锁定账号，false-解锁账号。默认为false。支持更新",
			},
			"schema_privilege_list": schema.SetNestedAttribute{
				Optional:    true,
				Description: "账号需要授权的数据库列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"grant_schema": schema.StringAttribute{
							Required:    true,
							Description: "需要授权的数据库名称，支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"privilege": schema.StringAttribute{
							Required:    true,
							Description: "授权数据库的权限，取值范围：readwrite-读写，readonly-只读，支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("readwrite", "readonly"),
							},
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "实例账户id",
			},
			"rol_super": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否具有超级用户权限",
			},
			"rol_inherit": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否自动继承其所属角色的权限",
			},
			"rol_create_role": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否支持创建其他子用户",
			},
			"rol_create_db": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否可以创建数据库",
			},
			"rol_can_login": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否可以登录数据库",
			},
			"rol_conn_limit": schema.Int32Attribute{
				Computed:    true,
				Description: "用户连接实例的最大并发连接数。-1表示没有限制",
			},
			"rol_by_pass_rls": schema.BoolAttribute{
				Computed:    true,
				Description: "用户是否绕过每个行级安全策略",
			},
		},
	}
}

func (c *CtyunPostgresqlAccount) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建新用户
	err = c.CreatePostgresqlAccount(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePostgresqlAccount(ctx, &plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(plan.InstID.ValueString() + "-" + plan.Name.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlAccount) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlAccountConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePostgresqlAccount(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlAccount) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPostgresqlAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPostgresqlAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updatePgsqlAccount(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergePostgresqlAccount(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlAccount) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPostgresqlAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &pgsql.PgsqlDeleteAccountRequest{
		ProdInstId: config.InstID.ValueString(),
		Username:   config.Name.ValueString(),
	}
	header := &pgsql.PgsqlDeleteAccountRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除mysql实例id=%s的%s用户失败，接口返回nil，具体原因请联系研发确认！", config.InstID.ValueString(), config.Name.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("delete mysql user failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	}
	if err != nil {
		return
	}
}

func (c *CtyunPostgresqlAccount) CreatePostgresqlAccount(ctx context.Context, config *CtyunPostgresqlAccountConfig) error {
	params := &pgsql.PgsqlCreateAccountRequest{
		ProdInstId: config.InstID.ValueString(),
		Username:   config.Name.ValueString(),
		Password:   config.Password.ValueString(),
		UserType:   config.UserType.ValueString(),
	}

	header := &pgsql.PgsqlCreateAccountRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)创建用户(account_name=%s)失败，接口返回nil。请与研发联系确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	// 授权数据库
	if !config.SchemaPrivilegeList.IsNull() {
		grantPrivilegeMap := make(map[string]string)
		var privilegeList []PostgresqlSchemaPrivilegeModel
		diags := config.SchemaPrivilegeList.ElementsAs(ctx, &privilegeList, true)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return err
		}
		for _, privilege := range privilegeList {
			schemaName := privilege.GrantSchema.ValueString()
			if _, exist := grantPrivilegeMap[schemaName]; !exist {
				grantPrivilegeMap[schemaName] = privilege.Privilege.ValueString()
			} else {
				err = fmt.Errorf("输入授权表(schema_privilege_list)信息有误！存在重复表权限！")
				return err
			}
		}
		err = c.grantSchemaPrivilege(ctx, grantPrivilegeMap, config)
		if err != nil {
			return err
		}
	}

	if !config.Description.IsNull() {
		err = c.updateInstanceDescription(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CtyunPostgresqlAccount) updateInstanceDescription(ctx context.Context, config *CtyunPostgresqlAccountConfig) error {
	params := &pgsql.PgsqlUpdateAccountRemarkRequest{
		ProdInstId:  config.InstID.ValueString(),
		Username:    config.Name.ValueString(),
		Description: config.Description.ValueStringPointer(),
	}
	header := &pgsql.PgsqlUpdateAccountRemarkRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateAccountRemarkApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新postgresql实例(id=%s)的备注失败，接口返回nil。请与研发联系确认问题原因！", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunPostgresqlAccount) getAndMergePostgresqlAccount(ctx context.Context, config *CtyunPostgresqlAccountConfig) error {
	resp, err := c.getPgsqlAccountDetail(ctx, config)
	if err != nil {
		return err
	}
	config.RolSuper = types.BoolValue(resp.RolSuper)
	config.RolInherit = types.BoolValue(resp.RolInherit)
	config.RolCreateRole = types.BoolValue(resp.RolCreateRole)
	config.RolCanLogin = types.BoolValue(resp.RolCanLogin)
	config.RolCreateDB = types.BoolValue(resp.RolCreateDB)
	config.RolConnLimit = types.Int32Value(resp.RolConnLimit)
	config.RolByPassRls = types.BoolValue(resp.RolByPassRls)
	if config.SchemaPrivilegeList.IsNull() {
		var privilege []PostgresqlSchemaPrivilegeModel
		privilegeList, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(PostgresqlSchemaPrivilegeModel{}), privilege)
		if diags.HasError() {
			err = errors.New(diags[0].Detail())
			return err
		}
		config.SchemaPrivilegeList = privilegeList
	}
	return nil
}

func (c *CtyunPostgresqlAccount) getPgsqlAccountDetail(ctx context.Context, config *CtyunPostgresqlAccountConfig) (*pgsql.PgsqlGetAccountListResponseReturnObj, error) {
	params := &pgsql.PgsqlGetAccountListRequest{
		ProdInstId: config.InstID.ValueString(),
		PageNum:    1,
		PageSize:   100,
		Username:   config.Name.ValueStringPointer(),
	}
	header := &pgsql.PgsqlGetAccountListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetAccountListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询postgresql实例(id=%s)的账户信息(account_name=%s)失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(resp.ReturnObj.List) <= 0 {
		err = fmt.Errorf("查询postgresql实例(id=%s)的账户信息失败，接口未查询到任何(account_name=%s)信息。", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	}
	if len(resp.ReturnObj.List) > 1 {
		err = fmt.Errorf("查询postgresql实例(id=%s)的账户信息失败，接口未查询到多条(account_name=%s)信息", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	}
	return &resp.ReturnObj.List[0], nil
}

func (c *CtyunPostgresqlAccount) updatePgsqlAccount(ctx context.Context, state *CtyunPostgresqlAccountConfig, plan *CtyunPostgresqlAccountConfig) error {
	// 更新remark
	if !plan.Description.IsNull() && !state.Description.Equal(plan.Description) {
		err := c.updateInstanceDescription(ctx, plan)
		if err != nil {
			return err
		}
		state.Description = plan.Description
	}
	// 新增/删除表权限
	err := c.updateSchemaPrivilege(ctx, state, plan)
	if err != nil {
		return err
	}
	// 修改密码
	if !plan.Password.Equal(state.Password) {
		err = c.updatePassword(ctx, state, plan)
		if err != nil {
			return err
		}
		state.Password = plan.Password
	}
	// 锁定/解锁用户
	if !plan.IsLock.IsNull() && !state.IsLock.Equal(plan.IsLock) {
		if plan.IsLock.ValueBool() {
			err = c.lockAccount(ctx, state)
			if err != nil {
				return err
			}
		} else {
			err = c.unlockAccount(ctx, state)
			if err != nil {
				return err
			}
		}
		state.IsLock = plan.IsLock
	}

	return nil
}

func (c *CtyunPostgresqlAccount) updateSchemaPrivilege(ctx context.Context, state *CtyunPostgresqlAccountConfig, plan *CtyunPostgresqlAccountConfig) error {
	if plan.SchemaPrivilegeList.IsNull() || plan.SchemaPrivilegeList.IsUnknown() {
		return nil
	}
	// 创建两个map，确定schema是否为新增/删除？权限是否需要发生变化
	grantPrivilegeMap := make(map[string]string)
	revokePrivilegeMap := make(map[string]string)

	statePrivilegeMap, err := c.getPrivilegeMap(ctx, state)
	if err != nil || statePrivilegeMap == nil {
		return err
	}
	planPrivilegeMap, err := c.getPrivilegeMap(ctx, plan)
	if err != nil || planPrivilegeMap == nil {
		return err
	}
	// 遍历plan节点，如果plan中有，state中没有，说明需要新增。
	for schemaName, privilege := range planPrivilegeMap {
		// 如果plan和state中都有该权限，进一步需要判断，如果权限相同则不变，权限不相同，需要重新授权
		if value, exist := statePrivilegeMap[schemaName]; exist {
			if value != privilege {
				grantPrivilegeMap[schemaName] = privilege
			}
		} else {
			// 如果plan中有，state中没有该schema权限，需要授权
			grantPrivilegeMap[schemaName] = privilege
		}
	}
	// 遍历state节点，plan中没有，state中说明需要删除
	for schemaName, privilege := range statePrivilegeMap {
		if _, exist := planPrivilegeMap[schemaName]; !exist {
			revokePrivilegeMap[schemaName] = privilege
		}
	}
	err = c.grantSchemaPrivilege(ctx, grantPrivilegeMap, state)
	if err != nil {
		return err
	}
	err = c.revokeSchemaPrivilege(ctx, revokePrivilegeMap, state)
	if err != nil {
		return err
	}
	state.SchemaPrivilegeList = plan.SchemaPrivilegeList
	return nil
}

func (c *CtyunPostgresqlAccount) getPrivilegeMap(ctx context.Context, state *CtyunPostgresqlAccountConfig) (map[string]string, error) {
	privilegeMap := make(map[string]string)
	var privilegeList []PostgresqlSchemaPrivilegeModel
	diags := state.SchemaPrivilegeList.ElementsAs(ctx, &privilegeList, true)
	if diags.HasError() {
		err := errors.New(diags[0].Detail())
		return nil, err
	}
	for _, privilege := range privilegeList {
		if _, exist := privilegeMap[privilege.GrantSchema.ValueString()]; !exist {
			privilegeMap[privilege.GrantSchema.ValueString()] = privilege.Privilege.ValueString()
		}
	}
	return privilegeMap, nil
}

func (c *CtyunPostgresqlAccount) grantSchemaPrivilege(ctx context.Context, privilegeMap map[string]string, config *CtyunPostgresqlAccountConfig) error {
	for schemaName, privilege := range privilegeMap {
		params := &pgsql.PgsqlGrantPrivilegeRequest{
			ProdInstId:    config.InstID.ValueString(),
			DbName:        schemaName,
			Username:      config.Name.ValueString(),
			UserPrivilege: privilege,
		}
		header := &pgsql.PgsqlGrantPrivilegeRequestHeader{
			RegionID: config.RegionID.ValueString(),
		}
		if !config.ProjectID.IsNull() {
			header.ProjectID = config.ProjectID.ValueStringPointer()
		}
		resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGrantPrivilegeApi.Do(ctx, c.meta.Credential, params, header)
		if err != nil {
			return err
		} else if resp == nil {
			err = fmt.Errorf("postgresql实例(id=%s)授权account=%s的%s权限失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString(), privilege)
			return err
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
			return err
		}
	}
	return nil
}

func (c *CtyunPostgresqlAccount) revokeSchemaPrivilege(ctx context.Context, privilegeMap map[string]string, config *CtyunPostgresqlAccountConfig) error {
	for schemaName, _ := range privilegeMap {
		params := &pgsql.PgsqlRevokePrivilegeRequest{
			ProdInstId: config.InstID.ValueString(),
			DbName:     schemaName,
			Username:   config.Name.ValueString(),
		}
		header := &pgsql.PgsqlRevokePrivilegeRequestHeader{
			RegionID: config.RegionID.ValueString(),
		}
		if !config.ProjectID.IsNull() {
			header.ProjectID = config.ProjectID.ValueStringPointer()
		}
		resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlRevokePrivilegeApi.Do(ctx, c.meta.Credential, params, header)
		if err != nil {
			return err
		} else if resp == nil {
			err = fmt.Errorf("postgresql实例(id=%s)撤销account=%s权限失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
			return err
		} else if resp.StatusCode != common.NormalStatusCode {
			err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
			return err
		}
	}
	return nil
}

func (c *CtyunPostgresqlAccount) lockAccount(ctx context.Context, config *CtyunPostgresqlAccountConfig) error {
	params := &pgsql.PgsqlLockAccountRequest{
		ProdInstId: config.InstID.ValueString(),
		Username:   config.Name.ValueString(),
	}
	header := &pgsql.PgsqlLockAccountRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlLockAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)冻结账户(account_name=%s)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunPostgresqlAccount) unlockAccount(ctx context.Context, config *CtyunPostgresqlAccountConfig) error {
	params := &pgsql.PgsqlUnLockAccountRequest{
		ProdInstId: config.InstID.ValueString(),
		Username:   config.Name.ValueString(),
	}
	header := &pgsql.PgsqlUnLockAccountRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUnLockAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)解冻账户(account_name=%s)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunPostgresqlAccount) updatePassword(ctx context.Context, state *CtyunPostgresqlAccountConfig, plan *CtyunPostgresqlAccountConfig) error {
	params := &pgsql.PgsqlResetPasswordRequest{
		ProdInstId: state.InstID.ValueString(),
		Username:   state.Name.ValueString(),
		Password:   plan.Password.ValueString(),
	}
	header := &pgsql.PgsqlResetPasswordRequestHeader{
		RegionID: state.RegionID.ValueString(),
	}
	if !state.ProjectID.IsNull() {
		header.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlResetPasswordApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)重置用户密码失败，接口返回nil，请联系研发确认问题原因！", state.InstID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

type CtyunPostgresqlAccountConfig struct {
	RegionID            types.String `tfsdk:"region_id"`
	ProjectID           types.String `tfsdk:"project_id"`
	InstID              types.String `tfsdk:"instance_id"`
	Name                types.String `tfsdk:"name"`
	Password            types.String `tfsdk:"password"`
	UserType            types.String `tfsdk:"user_type"`
	IsLock              types.Bool   `tfsdk:"is_lock"` // 是否锁定
	SchemaPrivilegeList types.Set    `tfsdk:"schema_privilege_list"`
	Description         types.String `tfsdk:"description"`
	ID                  types.String `tfsdk:"id"`
	RolSuper            types.Bool   `tfsdk:"rol_super"`
	RolInherit          types.Bool   `tfsdk:"rol_inherit"`
	RolCreateRole       types.Bool   `tfsdk:"rol_create_role"`
	RolCreateDB         types.Bool   `tfsdk:"rol_create_db"`
	RolCanLogin         types.Bool   `tfsdk:"rol_can_login"`
	RolConnLimit        types.Int32  `tfsdk:"rol_conn_limit"`
	RolByPassRls        types.Bool   `tfsdk:"rol_by_pass_rls"`
}

type PostgresqlSchemaPrivilegeModel struct {
	GrantSchema types.String `tfsdk:"grant_schema"` // 数据库schema
	Privilege   types.String `tfsdk:"privilege"`
}
