package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunPgsqlDatabase{}
	_ resource.ResourceWithConfigure   = &CtyunPgsqlDatabase{}
	_ resource.ResourceWithImportState = &CtyunPgsqlDatabase{}
)

type CtyunPgsqlDatabase struct {
	meta *common.CtyunMetadata
}

func (c *CtyunPgsqlDatabase) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_database"
}
func NewCtyunPgsqlDatabase() resource.Resource {
	return &CtyunPgsqlDatabase{}
}

func (c *CtyunPgsqlDatabase) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlDatabase) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[projectID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPostgresqlDatabaseConfig
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
	err = c.getAndMergePgsqlDatabase(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPgsqlDatabase) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10159978",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "pgsql实例id",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "数据库名称,mysql库名限制建议:长度为2~63个字符，以字母开头，以字母或数字结尾，由小写字母、数字、下划线或中划线组成，数据库名称在实例内必须是唯一的，禁用关键字",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 63),
					validator2.PgsqlDatabaseName(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "数据库描述，支持更新。长度为2~256个字符。以中文、英文字母开头，可以包含数字、中文、英文、下划线（_）、短横线（-）",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 256),
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
			"charset_collate": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "字符串排序规则，charset_name为utf8不传，其他的字符集必须传入",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.ConflictsWithEqualString(
						path.MatchRoot("charset_name"),
						types.StringValue("utf8"),
					),
				},
			},
			"charset_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "字符分类，charset_name为utf8不传，其他的字符集必须传入",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.ConflictsWithEqualString(
						path.MatchRoot("charset_name"),
						types.StringValue("utf8"),
					),
				},
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "数据库所有者",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "postgresql实例数据库id",
			},
		},
	}
}

func (c *CtyunPgsqlDatabase) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlDatabaseConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建数据库
	err = c.createPgsqlDatabase(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePgsqlDatabase(ctx, &plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(plan.InstID.ValueString() + "-" + plan.Name.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlDatabase) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlDatabaseConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePgsqlDatabase(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlDatabase) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPostgresqlDatabaseConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPostgresqlDatabaseConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.updatePgsqlDatabase(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergePgsqlDatabase(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlDatabase) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPostgresqlDatabaseConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &pgsql.PgsqlDeleteDatabaseRequest{
		ProdInstId: config.InstID.ValueString(),
		DBName:     config.Name.ValueString(),
	}
	header := &pgsql.PgsqlDeleteDatabaseRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDeleteDatabaseApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)删除数据库失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return
	}
	return
}

func (c *CtyunPgsqlDatabase) createPgsqlDatabase(ctx context.Context, config *CtyunPostgresqlDatabaseConfig) error {
	params := &pgsql.PgsqlCreateDatabaseRequest{
		ProdInstId: config.InstID.ValueString(),
		DBName:     config.Name.ValueString(),
		DBEncoding: config.CharSetName.ValueString(),
	}
	if !config.CharSetCollate.IsNull() && !config.CharSetCollate.IsUnknown() {
		params.DBCollate = config.CharSetCollate.ValueStringPointer()
	}
	if !config.CharSetType.IsNull() && !config.CharSetType.IsUnknown() {
		params.DBType = config.CharSetType.ValueStringPointer()
	}
	if !config.Owner.IsNull() && !config.Owner.IsUnknown() {
		params.DBOwner = config.Owner.ValueStringPointer()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.DBDescription = config.Description.ValueStringPointer()
	}
	header := &pgsql.PgsqlCreateDatabaseRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateDatabaseApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)创建数据库失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil

}

func (c *CtyunPgsqlDatabase) getAndMergePgsqlDatabase(ctx context.Context, config *CtyunPostgresqlDatabaseConfig) error {
	databaseDetail, err := c.getDatabaseDetail(ctx, config)
	if err != nil {
		return err
	}
	config.CharSetName = types.StringValue(databaseDetail.DBEncoding)
	config.CharSetCollate = types.StringValue(databaseDetail.DBCollate)
	config.CharSetType = types.StringValue(databaseDetail.DbType)
	config.Owner = types.StringValue(databaseDetail.DBOwner)
	config.Description = types.StringValue(databaseDetail.DBDescription)
	return nil
}

func (c *CtyunPgsqlDatabase) getDatabaseDetail(ctx context.Context, config *CtyunPostgresqlDatabaseConfig) (detail pgsql.PgsqlGetDatabaseSchemaResponseReturnObj, err error) {
	params := &pgsql.PgsqlGetDatabaseSchemaRequest{
		ProdInstId: config.InstID.ValueString(),
		DBName:     config.Name.ValueStringPointer(),
	}

	header := &pgsql.PgsqlGetDatabaseSchemaRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetDatabaseSchemaApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)查询数据库详情，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
	} else if len(resp.ReturnObj) == 0 {
		err = common.ResourceNotExistError
	} else if len(resp.ReturnObj) > 1 {
		for _, v := range resp.ReturnObj {
			if v.DBName == config.Name.ValueString() {
				detail = v
			}
			return
		}
		err = common.ResourceNotExistError
	} else if len(resp.ReturnObj) == 1 {
		detail = resp.ReturnObj[0]
	}
	return

}

func (c *CtyunPgsqlDatabase) updatePgsqlDatabase(ctx context.Context, state *CtyunPostgresqlDatabaseConfig, plan *CtyunPostgresqlDatabaseConfig) error {
	// 修改备注
	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		err := c.updateRemark(ctx, plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunPgsqlDatabase) updateRemark(ctx context.Context, config *CtyunPostgresqlDatabaseConfig) error {
	params := &pgsql.PgsqlUpdateDatabaseRemarkRequest{
		ProdInstId:  config.InstID.ValueString(),
		DBName:      config.Name.ValueString(),
		Description: config.Description.ValueStringPointer(),
	}
	header := &pgsql.PgsqlUpdateDatabaseRemarkRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateDatabaseRemarkApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例(id=%s)修改数据库备注失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

type CtyunPostgresqlDatabaseConfig struct {
	InstID         types.String `tfsdk:"instance_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	RegionID       types.String `tfsdk:"region_id"`
	Name           types.String `tfsdk:"name"`
	CharSetName    types.String `tfsdk:"charset_name"`
	CharSetCollate types.String `tfsdk:"charset_collate"`
	CharSetType    types.String `tfsdk:"charset_type"`
	Owner          types.String `tfsdk:"owner"`
	Description    types.String `tfsdk:"description"`
	ID             types.String `tfsdk:"id"`
}
