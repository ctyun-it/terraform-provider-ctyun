package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/amqp"

	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
)

var (
	_ resource.Resource                = &ctyunRabbitmqExchange{}
	_ resource.ResourceWithConfigure   = &ctyunRabbitmqExchange{}
	_ resource.ResourceWithImportState = &ctyunRabbitmqExchange{}
)

type ctyunRabbitmqExchange struct {
	meta            *common.CtyunMetadata
	rabbitmqService *business.RabbitmqService
}

func NewCtyunRabbitmqExchange() resource.Resource {
	return &ctyunRabbitmqExchange{}
}

func (c *ctyunRabbitmqExchange) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_exchange"
}

type CtyunRabbitmqExchangeConfig struct {
	ID                types.String `tfsdk:"id"`
	Vhost             types.String `tfsdk:"vhost"`
	InstanceID        types.String `tfsdk:"instance_id"`
	RegionID          types.String `tfsdk:"region_id"`
	Name              types.String `tfsdk:"name"`
	Type              types.String `tfsdk:"type"`
	AutoDelete        types.Bool   `tfsdk:"auto_delete"`
	Durable           types.Bool   `tfsdk:"durable"`
	Internal          types.Bool   `tfsdk:"internal"`
	AlternateExchange types.String `tfsdk:"alternate_exchange"`
	XDelayedType      types.String `tfsdk:"x_delayed_type"`
}

func (c *ctyunRabbitmqExchange) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10000118/10001967`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "交换器名称，只能包含字母，数字，短横线-和下划线_",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "交换器名称不符合规则"),
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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "rabbitMq实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vhost": schema.StringAttribute{
				Required:    true,
				Description: "vhost名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "交换器类型，支持direct、topic、headers、fanout、x-delayed-message",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.RabbitMqExchangeTypes...),
				},
			},
			"auto_delete": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动删除，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"durable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否持久化，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"internal": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否内置，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"alternate_exchange": schema.StringAttribute{
				Optional:    true,
				Description: "备用交换机名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					//stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "备用交换器名称不符合规则"),
				},
			},
			"x_delayed_type": schema.StringAttribute{
				Optional:    true,
				Description: "当type为x-delayed-message时必填",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.RabbitMqExchangeTypeXDelayedMessage),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("type"),
						utils.StringArrayToValueArray(business.RabbitMqExchangeXDelayedTypes)...,
					),
					stringvalidator.OneOf(business.RabbitMqExchangeXDelayedTypes...),
				},
			},
		},
	}
}

func (c *ctyunRabbitmqExchange) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRabbitmqExchangeConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	err = c.create(ctx, plan)
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

func (c *ctyunRabbitmqExchange) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqExchangeConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRabbitmqExchange) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunRabbitmqExchange) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqExchangeConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 销毁
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunRabbitmqExchange) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.rabbitmqService = business.NewRabbitmqService(meta)
}

func (c *ctyunRabbitmqExchange) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[vhost],[instanceID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRabbitmqExchangeConfig
	var name, vhost, instanceID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 3 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &name, &vhost, &instanceID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &vhost, &instanceID, &regionID)
		if err != nil {
			return
		}
	}
	if name == "" {
		err = fmt.Errorf("name不能为空")
		return
	}
	if vhost == "" {
		err = fmt.Errorf("vhost不能为空")
		return
	}
	if instanceID == "" {
		err = fmt.Errorf("instanceID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Name = types.StringValue(name)
	cfg.Vhost = types.StringValue(vhost)
	cfg.InstanceID = types.StringValue(instanceID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunRabbitmqExchange) checkBeforeCreate(ctx context.Context, plan CtyunRabbitmqExchangeConfig) (err error) {
	name, vhost, instanceID, regionID := plan.Name.ValueString(), plan.Vhost.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	// 确保vhost 存在
	exist, err := c.rabbitmqService.CheckVhostExist(ctx, vhost, instanceID, regionID)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("rabbitmq vhost %s not exist", vhost)
	}
	// 确保备用交换机存在
	if plan.AlternateExchange.ValueString() != "" {
		exist, err = c.rabbitmqService.CheckExchangeExist(ctx, plan.AlternateExchange.ValueString(), vhost, instanceID, regionID)
		if err != nil {
			return
		}
		if !exist {
			return fmt.Errorf("rabbitmq alternate exchange %s not exist", plan.AlternateExchange.ValueString())
		}
	}
	// 确保交换器名称没有被占用
	exist, err = c.rabbitmqService.CheckExchangeExist(ctx, name, vhost, instanceID, regionID)
	if err != nil {
		return
	}
	if exist {
		return fmt.Errorf("rabbitmq exchange %s already exist", name)
	}
	return
}

// create 创建
func (c *ctyunRabbitmqExchange) create(ctx context.Context, plan CtyunRabbitmqExchangeConfig) (err error) {
	params := &amqp.AmqpExchangeCreateV3Request{
		RegionId:          plan.RegionID.ValueString(),
		ProdInstId:        plan.InstanceID.ValueString(),
		Vhost:             plan.Vhost.ValueString(),
		Name:              plan.Name.ValueString(),
		Auto_delete:       plan.AutoDelete.ValueBoolPointer(),
		RawType:           plan.Type.ValueString(),
		AlternateExchange: plan.AlternateExchange.ValueStringPointer(),
		XDelayedType:      plan.XDelayedType.ValueStringPointer(),
		Durable:           plan.Durable.ValueBoolPointer(),
		Internal:          plan.Internal.ValueBoolPointer(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpExchangeCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// delete 创建
func (c *ctyunRabbitmqExchange) delete(ctx context.Context, plan CtyunRabbitmqExchangeConfig) (err error) {
	params := &amqp.AmqpExchangeDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		Vhost:      plan.Vhost.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpExchangeDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getExchangeByName 根据名称查询交换器
func (c *ctyunRabbitmqExchange) getExchangeByName(ctx context.Context, plan CtyunRabbitmqExchangeConfig) (exchange *amqp.AmqpExchangeQueryV3ReturnObjDataItemsResponse, err error) {
	params := &amqp.AmqpExchangeQueryV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		Vhost:      plan.Vhost.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpExchangeQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		if resp.Message == "交换器不存在" {
			err = common.ResourceNotExistError
		} else {
			err = fmt.Errorf("API return error. Message: %s", resp.Message)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjResultsError
		return
	} else if len(resp.ReturnObj.Data.Items) == 0 {
		err = common.ResourceNotExistError
		return
	}
	for _, item := range resp.ReturnObj.Data.Items {
		if item.Name == plan.Name.ValueString() {
			exchange = item
			return
		}
	}

	return
}

// getAndMerge 从远端查询
func (c *ctyunRabbitmqExchange) getAndMerge(ctx context.Context, plan *CtyunRabbitmqExchangeConfig) (err error) {
	exchange, err := c.getExchangeByName(ctx, *plan)
	if err != nil {
		return
	}
	plan.AutoDelete = types.BoolValue(exchange.Auto_delete)
	plan.Durable = types.BoolValue(exchange.Durable)
	plan.Type = types.StringValue(exchange.RawType)
	if exchange.Argument.XDelayedType != "" {
		plan.XDelayedType = types.StringValue(exchange.Argument.XDelayedType)
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s,%s",
		plan.Name.ValueString(),
		plan.Vhost.ValueString(),
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString()))
	return
}
