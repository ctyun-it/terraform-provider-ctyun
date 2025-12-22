package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ resource.Resource              = &CtyunMysqlRdsParameterTemplate{}
	_ resource.ResourceWithConfigure = &CtyunMysqlRdsParameterTemplate{}
)

type CtyunMysqlRdsParameterTemplate struct {
	meta         *common.CtyunMetadata
	mysqlService *business.MysqlService
}

func (c *CtyunMysqlRdsParameterTemplate) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_rds_parameter_template"
}
func NewCtyunMysqlRdsParameterTemplate() resource.Resource {
	return &CtyunMysqlRdsParameterTemplate{}
}

func (c *CtyunMysqlRdsParameterTemplate) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.mysqlService = business.NewMysqlService(meta)
}

func (c *CtyunMysqlRdsParameterTemplate) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10033813/10035295",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql数据库实例ID，为该实例管理只读实例",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"template_id": schema.Int64Attribute{
				Optional:    true,
				Description: "参数模板id，当mysql实例应用参数模板时必填。参数模板id可以根据data.ctyun_mysql_param_templates获取",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "参数模板管理id",
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "要修改的参数对。传入该参数，则无需传入template_id，当前mysql实例的参数可根据data.ctyun_mysql_parameters获取。",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (c *CtyunMysqlRdsParameterTemplate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunMysqlRdsParameterTemplateConfig
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
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}

	//err = c.getAndMergeBindEip(ctx, &plan)
	//if err != nil {
	//	return
	//}
	plan.ID = types.StringValue(uuid.NewString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlRdsParameterTemplate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	return
}

func (c *CtyunMysqlRdsParameterTemplate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunMysqlRdsParameterTemplateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunMysqlRdsParameterTemplateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlRdsParameterTemplate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunMysqlRdsParameterTemplate) mysqlApplyTemplate(ctx context.Context, config *CtyunMysqlRdsParameterTemplateConfig) error {
	params := &mysql.TeledbUpdateRdsTemplateParameterRequest{
		OuterProdInstId: config.InstID.ValueString(),
		ID:              config.TemplateID.ValueInt64Pointer(),
	}
	header := &mysql.TeledbUpdateRdsTemplateParameterRequestHeader{
		RegionID: config.RegionID.ValueString(),
		InstID:   config.InstID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() && config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateRdsTemplateParameterApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("mysql实例(id=%s)应用参数模板(id=%d)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.TemplateID.ValueInt64())
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}

	// 应用请求发出后，轮询确认
	err = c.applyLoop(ctx, config, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunMysqlRdsParameterTemplate) create(ctx context.Context, config *CtyunMysqlRdsParameterTemplateConfig) error {
	if !config.TemplateID.IsNull() {
		// mysql实例应用参数模板
		err := c.mysqlApplyTemplate(ctx, config)
		if err != nil {
			return err
		}
	}
	if !config.Parameters.IsNull() {
		// mysql实例修改参数值
		err := c.mysqlUpdateParameters(ctx, config)
		return err
	}
	return nil
}

func (c *CtyunMysqlRdsParameterTemplate) mysqlUpdateParameters(ctx context.Context, config *CtyunMysqlRdsParameterTemplateConfig) error {
	parameters := make(map[string]string)
	diags := config.Parameters.ElementsAs(ctx, &parameters, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params := &mysql.TeledbUpdateRdsTemplateParameterRequest{
		OuterProdInstId: config.InstID.ValueString(),
		Parameters:      &parameters,
	}
	header := &mysql.TeledbUpdateRdsTemplateParameterRequestHeader{
		RegionID: config.RegionID.ValueString(),
		InstID:   config.InstID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() && config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUpdateRdsTemplateParameterApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("mysql实例(id=%s)应用参数模板(id=%d)失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString(), config.TemplateID.ValueInt64())
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunMysqlRdsParameterTemplate) update(ctx context.Context, state *CtyunMysqlRdsParameterTemplateConfig, plan *CtyunMysqlRdsParameterTemplateConfig) error {
	if !plan.TemplateID.IsNull() {
		// mysql实例应用参数模板
		err := c.mysqlApplyTemplate(ctx, plan)
		if err != nil {
			return err
		}
		state.TemplateID = plan.TemplateID
	}
	if !plan.Parameters.IsNull() {
		// mysql实例修改参数值
		err := c.mysqlUpdateParameters(ctx, plan)
		if err != nil {
			return err
		}
		state.Parameters = plan.Parameters

	}
	return nil
}

func (c *CtyunMysqlRdsParameterTemplate) applyLoop(ctx context.Context, config *CtyunMysqlRdsParameterTemplateConfig, loopCount ...int) error {
	count := 60
	cnt := 2
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: config.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   config.InstID.ValueString(),
				RegionID: config.RegionID.ValueString(),
			}
			if config.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}

			status := resp.ReturnObj.ProdRunningStatus
			if status == business.MysqlRunningStatusApplying {
				return true
			}
			if status == business.MysqlRunningStatusStarted {
				if cnt > 0 {
					cnt--
					return true
				}
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return err
}

type CtyunMysqlRdsParameterTemplateConfig struct {
	InstID     types.String `tfsdk:"instance_id"`
	ProjectID  types.String `tfsdk:"project_id"`
	RegionID   types.String `tfsdk:"region_id"`
	TemplateID types.Int64  `tfsdk:"template_id"`
	ID         types.String `tfsdk:"id"`
	Parameters types.Map    `tfsdk:"parameters"`
}
