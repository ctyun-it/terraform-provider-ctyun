package mysql

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunMysqlDatabase{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlDatabase{}
	_ resource.ResourceWithImportState = &CtyunMysqlDatabase{}
)

type CtyunMysqlDatabase struct {
	meta         *common.CtyunMetadata
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlDatabase) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_database"
}
func NewCtyunMysqlDatabase() resource.Resource {
	return &CtyunMysqlDatabase{}
}

func (c *CtyunMysqlDatabase) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.mysqlService = business.NewMysqlService(meta)
}

func (c *CtyunMysqlDatabase) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[projectID],s[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunMysqlDatabaseConfig
	var name, regionID, projectID, instID string
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
	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if instID == "" {
		err = fmt.Errorf("instID不能为空")
		return
	}
	config.ID = types.StringValue(fmt.Sprintf("%s", instID+"-"+name))
	config.Name = types.StringValue(name)
	config.InstID = types.StringValue(instID)
	config.RegionID = types.StringValue(regionID)
	config.ProjectID = types.StringValue(projectID)
	err = c.getAndMergeMysqlDatabase(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunMysqlDatabase) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10140487",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql实例id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
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
			//todo validator
			"name": schema.StringAttribute{
				Required:    true,
				Description: "数据库名称,mysql库名限制建议:以小写字母开头，且以小写字母或数字结尾，可包含数字或下划线，不含其他特殊字符",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					validator2.MysqlDatabaseName(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "备注，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"charset_name": schema.StringAttribute{
				Required:    true,
				Description: "字符集",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"user_grant_privilege": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_name": schema.StringAttribute{
							Computed:    true,
							Description: "数据库用户名",
						},
						"read_only": schema.BoolAttribute{
							Computed:    true,
							Description: "用户读写权限",
						},
						"select_priv": schema.BoolAttribute{
							Computed:    true,
							Description: "查询权限",
						},
						"insert_priv": schema.BoolAttribute{
							Computed:    true,
							Description: "写入权限",
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据库id",
			},
		},
	}
}

func (c *CtyunMysqlDatabase) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunMysqlDatabaseConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
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
	// 开始创建数据库
	err = c.createMysqlDatabase(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergeMysqlDatabase(ctx, &plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(plan.InstID.ValueString() + "-" + plan.Name.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlDatabase) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunMysqlDatabaseConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeMysqlDatabase(ctx, &state)
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

func (c *CtyunMysqlDatabase) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlDatabaseConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlDatabaseConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updateMysqlDatabase(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeMysqlDatabase(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlDatabase) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunMysqlDatabaseConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.deleteMysqlDatabase(ctx, config)
	if err != nil {
		return
	}
}

// createMysqlDatabase 创建数据库
func (c *CtyunMysqlDatabase) createMysqlDatabase(ctx context.Context, config *CtyunMysqlDatabaseConfig) error {
	// 创建前，确认db name是否可用
	//err := c.checkDBName(ctx, config)
	//if err != nil {
	//	return err
	//}
	params := &mysql.TeledbCreateDatabaseRequest{
		OuterProdInstId: config.InstID.ValueString(),
		DBName:          config.Name.ValueString(),
		CharSetName:     config.CharSetName.ValueString(),
	}
	header := &mysql.TeledbCreateDatabaseRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCreateDatabaseApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("mysql实例(id=%s)创建数据库失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("create mysql database error, API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}

	// 如果description不为空，调用接口更新description
	if !config.Description.IsNull() {
		err = c.updateDescription(ctx, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunMysqlDatabase) updateDescription(ctx context.Context, config *CtyunMysqlDatabaseConfig) error {
	params := &mysql.TeledbUpdateDatabaseRemarkRequest{
		OuterProdInstId: config.InstID.ValueString(),
		DatabaseName:    config.Name.ValueString(),
		Remark:          config.Description.ValueStringPointer(),
	}
	header := &mysql.TeledbUpdateDatabaseRemarkRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateDatabaseRemarkApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("mysql实例(id=%s)更新database(name=%s)的备注失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunMysqlDatabase) getAndMergeMysqlDatabase(ctx context.Context, config *CtyunMysqlDatabaseConfig) error {
	resp, err := c.getMysqlDatabaseInfo(ctx, config)
	if err != nil {
		return err
	}

	var userGrantPrivilegeList []CtyunMysqlDatabaseGrantPrivilegeModel
	for _, privilegeItem := range resp.UserVOList {
		var privilege CtyunMysqlDatabaseGrantPrivilegeModel
		privilege.AccountName = types.StringValue(privilegeItem.AccountName)
		privilege.InsertPriv = types.BoolValue(business.PrivilegeMap[privilegeItem.InsertPriv])
		privilege.SelectPriv = types.BoolValue(business.PrivilegeMap[privilegeItem.SelectPriv])
		privilege.ReadOnly = types.BoolValue(privilegeItem.ReadOnly)
		userGrantPrivilegeList = append(userGrantPrivilegeList, privilege)
	}
	userGrantPrivilegeListTmp, diag := types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunMysqlDatabaseGrantPrivilegeModel{}), &userGrantPrivilegeList)
	if diag.HasError() {
		err = errors.New(diag[0].Detail())
		return err
	}
	config.UserGrantPrivilege = userGrantPrivilegeListTmp
	return nil
}

func (c *CtyunMysqlDatabase) getMysqlDatabaseInfo(ctx context.Context, config *CtyunMysqlDatabaseConfig) (*mysql.TeledbGetDatabaseSchemaResponseReturnObj, error) {
	params := &mysql.TeledbGetDatabaseSchemaRequest{
		OuterProdInstId: config.InstID.ValueString(),
	}
	header := &mysql.TeledbGetDatabaseSchemaRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetDatabaseSchemaApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询mysql实例(id=%s)的database schema(name=%s)失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	for _, schemaItem := range resp.ReturnObj {
		schemaName := schemaItem.GrantSchema
		if schemaName == config.Name.ValueString() {
			return &schemaItem, nil
		}
	}
	return nil, fmt.Errorf("未查询到db_name=%s，mysql实例(id=%s)的schema 信息", config.Name.ValueString(), config.InstID.ValueString())
}

func (c *CtyunMysqlDatabase) checkDBName(ctx context.Context, config *CtyunMysqlDatabaseConfig) error {
	params := &mysql.TeledbCheckDatabaseNameAvailableRequest{
		OuterProdInstId: config.InstID.ValueString(),
		DBName:          config.Name.ValueString(),
	}
	header := &mysql.TeledbCheckDatabaseNameAvailableRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbCheckDatabaseNameAvailableApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("mysql实例(id=%s)验证database(name=%s)有效性失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode == 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	if !resp.ReturnObj.Available {
		err = fmt.Errorf("database(name=%s)不可用", config.Name.ValueString())
		return err
	}
	return nil
}

func (c *CtyunMysqlDatabase) updateMysqlDatabase(ctx context.Context, state *CtyunMysqlDatabaseConfig, plan *CtyunMysqlDatabaseConfig) error {
	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		state.Description = plan.Description
		err := c.updateDescription(ctx, state)
		if err != nil {
			return err
		}
		state.Description = plan.Description
	}
	return nil
}

func (c *CtyunMysqlDatabase) deleteMysqlDatabase(ctx context.Context, config CtyunMysqlDatabaseConfig) error {
	params := &mysql.TeledbDeleteDatabaseRequest{
		OuterProdInstId: config.InstID.ValueString(),
		DBName:          config.Name.ValueString(),
	}
	header := &mysql.TeledbDeleteDatabaseRequestHeader{
		InstID:   config.InstID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbDeleteDatabaseApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除mysql实例(id=%s)的database(name=%s)失败，接口返回nil。请联系研发确认问题原因！", config.InstID.ValueString(), config.Name.ValueString())
		return err
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

type CtyunMysqlDatabaseConfig struct {
	InstID             types.String `tfsdk:"instance_id"`
	ProjectID          types.String `tfsdk:"project_id"`
	RegionID           types.String `tfsdk:"region_id"`
	Name               types.String `tfsdk:"name"`
	CharSetName        types.String `tfsdk:"charset_name"`
	Description        types.String `tfsdk:"description"`
	UserGrantPrivilege types.List   `tfsdk:"user_grant_privilege"`
	ID                 types.String `tfsdk:"id"`
}

type CtyunMysqlDatabaseGrantPrivilegeModel struct {
	AccountName types.String `tfsdk:"account_name"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	SelectPriv  types.Bool   `tfsdk:"select_priv"`
	InsertPriv  types.Bool   `tfsdk:"insert_priv"`
}
