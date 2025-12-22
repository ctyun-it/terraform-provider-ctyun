package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/amqp"

	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
)

var (
	_ resource.Resource                = &ctyunRabbitmqQueue{}
	_ resource.ResourceWithConfigure   = &ctyunRabbitmqQueue{}
	_ resource.ResourceWithImportState = &ctyunRabbitmqQueue{}
)

type ctyunRabbitmqQueue struct {
	meta            *common.CtyunMetadata
	rabbitmqService *business.RabbitmqService
}

func NewCtyunRabbitmqQueue() resource.Resource {
	return &ctyunRabbitmqQueue{}
}

func (c *ctyunRabbitmqQueue) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_queue"
}

type CtyunRabbitmqQueueConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	RegionID   types.String `tfsdk:"region_id"`
	Durable    types.Bool   `tfsdk:"durable"`     // 是否持久化，可选
	AutoDelete types.Bool   `tfsdk:"auto_delete"` // 是否自动删除，可选
	//Node                  types.String `tfsdk:"node"`                      // 队列所在节点，可选
	XExpires              types.Int64  `tfsdk:"x_expires"`                 // 队列过期时间(ms)，可选
	XDeadLetterExchange   types.String `tfsdk:"x_dead_letter_exchange"`    // 死信交换器名称，可选
	XDeadLetterRoutingKey types.String `tfsdk:"x_dead_letter_routing_key"` // 死信路由键，可选
	XMessageTTL           types.Int64  `tfsdk:"x_message_ttl"`             // 消息过期时间(ms)，可选
	XMaxLength            types.Int64  `tfsdk:"x_max_length"`              // 队列最大消息长度，可选
	XMaxLengthBytes       types.Int64  `tfsdk:"x_max_length_bytes"`        // 队列消息总字节数上限，可选
	XOverflow             types.String `tfsdk:"x_overflow"`                // 队列消息处理策略，可选
	XMaxPriority          types.Int64  `tfsdk:"x_max_priority"`            // 队列最大优先级，可选
	XQueueMode            types.String `tfsdk:"x_queue_mode"`              // 队列模式，可选
}

func (c *ctyunRabbitmqQueue) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10000118/10220893`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "队列名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "队列名称不符合规则"),
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
			//"node": schema.StringAttribute{
			//	Description: "队列所在节点，默认为实例随机节点",
			//	Optional:    true,
			//	Computed:    true,
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.RequiresReplace(),
			//		stringplanmodifier.UseStateForUnknown(),
			//	},
			//},
			"x_expires": schema.Int64Attribute{
				Description: "队列过期时间，过期后队列自动删除，单位为ms",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"x_dead_letter_exchange": schema.StringAttribute{
				Description: "死信交换器名称。消息被拒绝或过期时将重新发布到该交换器",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"x_dead_letter_routing_key": schema.StringAttribute{
				Description: "死信路由键。死信交换器会发送死信消息到绑定对应路由键的队列上",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"x_message_ttl": schema.Int64Attribute{
				Description: "消息过期时间，发布到队列的消息在被丢弃之前可以存活的时间，单位为ms",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"x_max_length": schema.Int64Attribute{
				Description: "队列最大消息长度",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"x_max_length_bytes": schema.Int64Attribute{
				Description: "队列消息的总字节数上限",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"x_overflow": schema.StringAttribute{
				Description: "队列消息处理策略。可选值：drop-head（默认，丢弃最早消息）、reject-publish（拒绝接收新消息）",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("drop-head", "reject-publish"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString("drop-head"),
			},
			"x_max_priority": schema.Int64Attribute{
				Description: "队列最大优先级，范围为0~255，默认0",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 255),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
				Default: int64default.StaticInt64(0),
			},

			"x_queue_mode": schema.StringAttribute{
				Description: "队列模式。可选值：default（默认模式）、lazy（惰性模式）",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "lazy"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString("default"),
			},
		},
	}
}

func (c *ctyunRabbitmqQueue) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRabbitmqQueueConfig
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

func (c *ctyunRabbitmqQueue) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqQueueConfig
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

func (c *ctyunRabbitmqQueue) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunRabbitmqQueue) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqQueueConfig
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

func (c *ctyunRabbitmqQueue) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.rabbitmqService = business.NewRabbitmqService(meta)
}

func (c *ctyunRabbitmqQueue) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[vhost],[instanceID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRabbitmqQueueConfig
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

func (c *ctyunRabbitmqQueue) checkBeforeCreate(ctx context.Context, plan CtyunRabbitmqQueueConfig) (err error) {
	name, vhost, instanceID, regionID := plan.Name.ValueString(), plan.Vhost.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	// 确保vhost 存在
	exist, err := c.rabbitmqService.CheckVhostExist(ctx, vhost, instanceID, regionID)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("rabbitmq vhost %s not exist", vhost)
	}
	// 确保死信交换机存在
	if plan.XDeadLetterExchange.ValueString() != "" {
		exist, err = c.rabbitmqService.CheckExchangeExist(ctx, plan.XDeadLetterExchange.ValueString(), vhost, instanceID, regionID)
		if err != nil {
			return
		}
		if !exist {
			return fmt.Errorf("rabbitmq x-dead-letter exchange %s not exist", plan.XDeadLetterExchange.ValueString())
		}
	}
	// 确保队列名称没有被占用
	exist, err = c.rabbitmqService.CheckQueueExist(ctx, name, vhost, instanceID, regionID)
	if err != nil {
		return
	}
	if exist {
		return fmt.Errorf("rabbitmq exchange %s already exist", name)
	}
	return
}

// create 创建
func (c *ctyunRabbitmqQueue) create(ctx context.Context, plan CtyunRabbitmqQueueConfig) (err error) {
	params := &amqp.AmqpQueueCreateV3Request{
		RegionId:              plan.RegionID.ValueString(),
		ProdInstId:            plan.InstanceID.ValueString(),
		Vhost:                 plan.Vhost.ValueString(),
		Name:                  plan.Name.ValueString(),
		Durable:               plan.Durable.ValueBoolPointer(),
		Auto_delete:           plan.AutoDelete.ValueBoolPointer(),
		XExpires:              plan.XExpires.ValueInt64Pointer(),
		XDeadLetterExchange:   plan.XDeadLetterExchange.ValueStringPointer(),
		XDeadLetterRoutingKey: plan.XDeadLetterRoutingKey.ValueStringPointer(),
		XMessageTTL:           plan.XMessageTTL.ValueInt64Pointer(),
		XMaxLength:            plan.XMaxLength.ValueInt64Pointer(),
		XMaxLengthBytes:       plan.XMaxLengthBytes.ValueInt64Pointer(),
		XOverflow:             plan.XOverflow.ValueStringPointer(),
		XMaxPriority:          plan.XMaxPriority.ValueInt64Pointer(),
		XQueueMode:            plan.XQueueMode.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpQueueCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// delete 创建
func (c *ctyunRabbitmqQueue) delete(ctx context.Context, plan CtyunRabbitmqQueueConfig) (err error) {
	params := &amqp.AmqpQueueDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		Name:       plan.Name.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		Vhost:      plan.Vhost.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpQueueDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getQueueByName 根据名称查询队列
func (c *ctyunRabbitmqQueue) getQueueByName(ctx context.Context, plan CtyunRabbitmqQueueConfig) (queue *amqp.AmqpQueueQueryV3ReturnObjDataItem, err error) {
	params := &amqp.AmqpQueueQueryV3Request{
		RegionId:   plan.RegionID.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
		Vhost:      plan.Vhost.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpQueueQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
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
	queue = resp.ReturnObj.Data.Items[0]
	return
}

// getAndMerge 从远端查询
func (c *ctyunRabbitmqQueue) getAndMerge(ctx context.Context, plan *CtyunRabbitmqQueueConfig) (err error) {
	queue, err := c.getQueueByName(ctx, *plan)
	if err != nil {
		return
	}
	plan.Durable = types.BoolValue(queue.Durable)
	plan.Vhost = types.StringValue(queue.Vhost)
	plan.AutoDelete = types.BoolValue(queue.AutoDelete)
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s,%s",
		plan.Name.ValueString(),
		plan.Vhost.ValueString(),
		plan.InstanceID.ValueString(),
		plan.RegionID.ValueString()))
	return
}
