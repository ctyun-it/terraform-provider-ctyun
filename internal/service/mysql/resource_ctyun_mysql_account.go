package mysql

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunMysqlAccount{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlAccount{}
	_ resource.ResourceWithImportState = &CtyunMysqlAccount{}
)

type CtyunMysqlAccount struct {
	meta         *common.CtyunMetadata
	name         string
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlAccount) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_account"
}
func NewCtyunMysqlAccount() resource.Resource {
	return &CtyunMysqlAccount{}
}

func (c *CtyunMysqlAccount) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.mysqlService = business.NewMysqlService(meta)
}

func (c *CtyunMysqlAccount) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] [instance_id],[name],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunMysqlAccountConfig
	var regionID, instID, name string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instID, &name)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instID, &name, &regionID)
		if err != nil {
			return
		}
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	if instID == "" {
		err = fmt.Errorf("instance_id不能为空")
		return
	}
	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}
	config.InstID = types.StringValue(instID)
	config.RegionID = types.StringValue(regionID)
	config.Name = types.StringValue(name)
	err = c.getAndMergeMysqlAccount(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunMysqlAccount) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理MySQL实例的账户", "关系数据库MySQL版", "https://www.ctyun.cn/document/10033813/10133363"),
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
				Optional:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "企业项目ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "数据库账号名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "数据库账号的密码，创建时填写，导入时不填，支持更新。不建议使用弱密码，长度[8,20]位",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 20),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "备注，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schema_privilege_list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "数据库权限配置列表，支持更新。",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"grant_schema": schema.StringAttribute{
							Required:    true,
							Description: "授权数据库名称，支持更新。",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"privilege": schema.StringAttribute{
							Required:    true,
							Description: "数据库权限，支持更新。取值范围：read_only, ddl, dml, rw",
							Validators: []validator.String{
								stringvalidator.OneOf(business.MysqlSchemaPrivileges...),
							},
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunMysqlAccount) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	if plan.Password.ValueString() == "" {
		err = fmt.Errorf("创建时密码必须填写")
		return
	}
	err = c.mysqlService.WaitInstanceStatus(
		ctx,
		plan.InstID.ValueString(),
		plan.ProjectID.ValueString(),
		plan.RegionID.ValueString(),
		business.MysqlRunningStatusStarted,
		business.MysqlOrderStatusStarted,
	)
	if err != nil {
		return
	}
	// 开始创建新用户
	err = c.createMysqlAccount(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlAccount(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAccount) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlAccountConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlAccount(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAccount) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlAccountConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updateMysqlAccount(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlAccount(ctx, &state)
	if err != nil {
		return
	}
	state.ProjectID = plan.ProjectID
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAccount) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunMysqlAccountConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &mysql.TeledbDeleteAccountRequest{
		OuterProdInstId: config.InstID.ValueString(),
		AccountName:     config.Name.ValueString(),
	}
	header := &mysql.TeledbDeleteAccountRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDeleteAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("删除mysql实例id=%s的%s用户失败，接口返回nil，具体原因请联系研发确认！", config.InstID.ValueString(), config.Name.ValueString())
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("delete mysql user failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	}
	if err != nil {
		return
	}
}

// createMysqlAccount 创建mysql账号
func (c *CtyunMysqlAccount) createMysqlAccount(ctx context.Context, config *CtyunMysqlAccountConfig) error {
	params := &mysql.TeledbCreateAccountRequest{
		OuterProdInstId: config.InstID.ValueString(),
		AccountName:     config.Name.ValueString(),
		AccountPassword: c.encodeBase64(config.Password.ValueString()),
	}
	header := &mysql.TeledbCreateAccountRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.SchemaPrivilegeList.IsNull() {
		var privilegeList []MysqlSchemaPrivilegeModel
		diags := config.SchemaPrivilegeList.ElementsAs(ctx, &privilegeList, false)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		var schemaPrivilegeList []mysql.SchemaPrivilegeVO
		for _, privilegeItem := range privilegeList {
			var privilege mysql.SchemaPrivilegeVO
			privilege.GrantSchema = privilegeItem.GrantSchema.ValueString()
			c.judgeSchemaPrivilege(&privilege, privilegeItem.Privilege.ValueString())
			schemaPrivilegeList = append(schemaPrivilegeList, privilege)
		}
		params.SchemaPrivilegeVOList = schemaPrivilegeList
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateAccountApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("为mysql实例(id=%s)创建用户%s时失败，接口返回为nil。请与研发联系确认问题原因。", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("create mysql account failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	// 若用户增加注释，需要处理更新remark
	if config.Description.ValueString() != "" {
		err = c.updateRemark(ctx, config, config.Description.ValueString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunMysqlAccount) encodeBase64(password string) string {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	return encodedPassword
}

// getAndMergeMysqlAccount 查询mysql账号
func (c *CtyunMysqlAccount) getAndMergeMysqlAccount(ctx context.Context, config *CtyunMysqlAccountConfig) error {
	config.ID = types.StringValue(config.InstID.ValueString() + "," + config.Name.ValueString())
	respPrivilege, err := c.getMysqlAccountInfo(ctx, config)
	if err != nil {
		return err
	}
	config.Description = types.StringValue(respPrivilege.Remark)
	privileges := []MysqlSchemaPrivilegeModel{}
	for _, privilegeItem := range respPrivilege.SchemaPrivilegeVOList {
		var privilege MysqlSchemaPrivilegeModel
		privilege.GrantSchema = types.StringValue(privilegeItem.GrantSchema)
		schemaPrivilege := c.getPrivilege(privilegeItem)
		if schemaPrivilege == "" {
			continue
		}
		privilege.Privilege = types.StringValue(schemaPrivilege)
		privileges = append(privileges, privilege)
	}
	var diags diag.Diagnostics
	config.SchemaPrivilegeList, diags = types.SetValueFrom(ctx, utils.StructToTFObjectTypes(MysqlSchemaPrivilegeModel{}), &privileges)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	return nil
}

func (c *CtyunMysqlAccount) getMysqlAccountInfo(ctx context.Context, config *CtyunMysqlAccountConfig) (*mysql.TeledbGetAccountInfoResponseReturnObj, error) {
	params := &mysql.TeledbGetAccountInfoRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetAccountInfoRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetAccountInfoApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例(id=%s)的用户信息失败", config.InstID.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		// MYSQL_10002=实例id不存在 || message: accountName not exist
		if strings.Contains(*resp.Error, "MYSQL_10002") || strings.Contains(resp.Message, "accountName not exist") {
			err = common.ResourceNotExistError
			return nil, err
		}
		err = fmt.Errorf("get mysql account failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	for _, accountPrivilege := range resp.ReturnObj {
		Name := accountPrivilege.AccountName
		if Name == config.Name.ValueString() {
			return &accountPrivilege, nil
		}
	}
	return nil, common.ResourceNotExistError
}

func (c *CtyunMysqlAccount) updateMysqlAccount(ctx context.Context, state *CtyunMysqlAccountConfig, plan *CtyunMysqlAccountConfig) error {
	// 更新密码
	err := c.updatePassword(ctx, state, plan)
	if err != nil {
		return err
	}
	// 更新权限
	// 增加和删除权限
	// 1. 处理新增与更新权限操作
	err = c.updatePrivilege(ctx, state, plan)
	if err != nil {
		return err
	}
	// 2. 更新备注
	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		err = c.updateRemark(ctx, state, plan.Description.ValueString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunMysqlAccount) updatePassword(ctx context.Context, state *CtyunMysqlAccountConfig, plan *CtyunMysqlAccountConfig) error {
	if plan.Password.Equal(state.Password) {
		return nil
	}
	params := &mysql.TeledbResetPasswordRequest{
		OuterProdInstId: state.InstID.ValueString(),
		AccountName:     state.Name.ValueString(),
		AccountPassword: c.encodeBase64(plan.Password.ValueString()),
	}
	header := &mysql.TeledbResetPasswordRequestHeader{
		InstID:   state.InstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbResetPasswordApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("修改mysql实例id=%s密码失败，接口返回nil。具体原因请联系研发确认。", state.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("reset password failed, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	state.Password = plan.Password
	return nil
}

func (c *CtyunMysqlAccount) judgeSchemaPrivilege(privilegeVo *mysql.SchemaPrivilegeVO, privilege string) {
	trueValue := true
	switch privilege {
	case business.MysqlSchemaPrivilegeRw:
		privilegeVo.ReadAndWrite = &trueValue
	case business.MysqlSchemaPrivilegeDDL:
		privilegeVo.DDLPrivilege = &trueValue
	case business.MysqlSchemaPrivilegeDML:
		privilegeVo.DMLPrivilege = &trueValue
	case business.MysqlSchemaPrivilegeReadOnly:
		privilegeVo.ReadOnly = &trueValue
	}
}

func (c *CtyunMysqlAccount) getPrivilege(item mysql.SchemaPrivilegeVO) string {
	if *item.DDLPrivilege {
		return business.MysqlSchemaPrivilegeDDL
	} else if *item.ReadAndWrite {
		return business.MysqlSchemaPrivilegeRw
	} else if *item.DMLPrivilege {
		return business.MysqlSchemaPrivilegeDML
	} else if *item.ReadOnly {
		return business.MysqlSchemaPrivilegeReadOnly
	}
	return ""
}

func (c *CtyunMysqlAccount) updatePrivilege(ctx context.Context, state *CtyunMysqlAccountConfig, plan *CtyunMysqlAccountConfig) error {
	grantPrivilegeMap := make(map[string]string)
	revokePrivilegeMap := make(map[string]string)
	// 将state和plan中表-权限 映射成map
	statePrivilegeMap, err := c.getPrivilegeMap(ctx, state)
	if err != nil {
		return err
	}
	planPrivilegeMap, err := c.getPrivilegeMap(ctx, plan)
	if err != nil {
		return err
	}
	// 先处理更新、新增权限，利用plan的列表轮询，与state列表做对比
	for schemaName, privilege := range planPrivilegeMap {
		// 若plan有schema权限记录，但是state没有，触发新增
		if _, ok := statePrivilegeMap[schemaName]; !ok {
			grantPrivilegeMap[schemaName] = privilege
		} else {
			if privilege != statePrivilegeMap[schemaName] {
				grantPrivilegeMap[schemaName] = privilege
			}
		}
	}

	err = c.requestGrantAndUpdateSchemaPrivilege(ctx, state, grantPrivilegeMap)
	if err != nil {
		return err
	}
	// 处理删除权限，利用state的列表轮询，与plan列表做对比
	for schemaName, privilege := range statePrivilegeMap {
		// state中记录的schema，但是plan阶段没有该schema触发revoke操作
		if _, ok := planPrivilegeMap[schemaName]; !ok {
			revokePrivilegeMap[schemaName] = privilege
		}
	}
	err = c.revokeSchemaPrivilege(ctx, state, revokePrivilegeMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunMysqlAccount) getPrivilegeMap(ctx context.Context, config *CtyunMysqlAccountConfig) (map[string]string, error) {
	privilegeMap := make(map[string]string)
	if !config.SchemaPrivilegeList.IsNull() {
		var schemaPrivilegeList []MysqlSchemaPrivilegeModel
		diags := config.SchemaPrivilegeList.ElementsAs(ctx, &schemaPrivilegeList, false)
		if diags.HasError() {
			err := fmt.Errorf(diags[0].Detail())

			return nil, err
		}
		for _, schemaPrivilege := range schemaPrivilegeList {
			schemaName := schemaPrivilege.GrantSchema.ValueString()
			privilege := schemaPrivilege.Privilege.ValueString()
			privilegeMap[schemaName] = privilege
		}
	}
	return privilegeMap, nil
}

func (c *CtyunMysqlAccount) requestGrantAndUpdateSchemaPrivilege(ctx context.Context, state *CtyunMysqlAccountConfig, privilegeMap map[string]string) error {
	if len(privilegeMap) <= 0 {
		return nil
	}
	params := &mysql.TeledbGrantPrivilegeRequest{
		OuterProdInstId: state.InstID.ValueString(),
		AccountName:     state.Name.ValueString(),
	}
	header := &mysql.TeledbGrantPrivilegeRequestHeader{
		InstID:   state.InstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}

	var schemaPrivilegeVOList []mysql.SchemaPrivilegeVO
	for schemaName, privilege := range privilegeMap {
		var schemaPrivilegeVo mysql.SchemaPrivilegeVO
		schemaPrivilegeVo.GrantSchema = schemaName
		c.judgeSchemaPrivilege(&schemaPrivilegeVo, privilege)
		schemaPrivilegeVOList = append(schemaPrivilegeVOList, schemaPrivilegeVo)
	}
	params.SchemaPrivilegeVOList = schemaPrivilegeVOList

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGrantPrivilegeApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("为mysql实例(id=%s)添加用户权限失败。添加列表为：%#v", state.InstID.ValueString(), privilegeMap)
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunMysqlAccount) revokeSchemaPrivilege(ctx context.Context, state *CtyunMysqlAccountConfig, privilegeMap map[string]string) error {
	if len(privilegeMap) <= 0 {
		return nil
	}
	params := &mysql.TeledbRevokeSchemaRequest{
		OuterProdInstId: state.InstID.ValueString(),
		AccountName:     state.Name.ValueString(),
		DatabaseVOList:  nil,
	}
	var schemaList []mysql.DatabaseVO
	for schemaName, _ := range privilegeMap {
		var schemaItem mysql.DatabaseVO
		schemaItem.RevokeSchema = schemaName
		schemaList = append(schemaList, schemaItem)
	}
	params.DatabaseVOList = schemaList
	header := &mysql.TeledbRevokeSchemaRequestHeader{
		InstID:   state.InstID.ValueString(),
		RegionID: state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbRevokeSchemaApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("撤销mysql实例id=%s的%s账户%#v库权限失败，请联系研发确认问题原因。", state.InstID.ValueString(), state.Name.ValueString(), privilegeMap)
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunMysqlAccount) updateRemark(ctx context.Context, config *CtyunMysqlAccountConfig, desc string) error {
	params := &mysql.TeledbUpdateAccountRemarkRequest{
		OuterProdInstId: config.InstID.ValueString(),
		AccountName:     config.Name.ValueString(),
		Remark:          &desc,
	}
	header := &mysql.TeledbUpdateAccountRemarkRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateAccountRemarkApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新mysql数据库(id=%s)的%s用户备注失败，接口返回nil。请与研发联系确认问题原因", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

type CtyunMysqlAccountConfig struct {
	InstID              types.String `tfsdk:"instance_id"`
	ProjectID           types.String `tfsdk:"project_id"`
	RegionID            types.String `tfsdk:"region_id"`
	Name                types.String `tfsdk:"name"`
	Password            types.String `tfsdk:"password"`
	SchemaPrivilegeList types.Set    `tfsdk:"schema_privilege_list"`
	Description         types.String `tfsdk:"description"`
	ID                  types.String `tfsdk:"id"`
}

type MysqlSchemaPrivilegeModel struct {
	GrantSchema types.String `tfsdk:"grant_schema"` // 数据库schema
	Privilege   types.String `tfsdk:"privilege"`
}
