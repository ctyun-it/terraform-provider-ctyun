package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &ctyunRabbitmqVhost{}
	_ resource.ResourceWithConfigure   = &ctyunRabbitmqVhost{}
	_ resource.ResourceWithImportState = &ctyunRabbitmqVhost{}
)

type ctyunRabbitmqVhost struct {
	meta            *common.CtyunMetadata
	rabbitmqService *business.RabbitmqService
}

func NewCtyunRabbitmqVhost() resource.Resource {
	return &ctyunRabbitmqVhost{}
}

func (c *ctyunRabbitmqVhost) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rabbitmq_vhost"
}

type CtyunRabbitmqVhostConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	Name       types.String `tfsdk:"name"`
	RegionID   types.String `tfsdk:"region_id"`
}

func (c *ctyunRabbitmqVhost) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
				Description: "vhost名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "vhost名称不符合规则"),
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
		},
	}
}

func (c *ctyunRabbitmqVhost) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRabbitmqVhostConfig
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

func (c *ctyunRabbitmqVhost) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqVhostConfig
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

func (c *ctyunRabbitmqVhost) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunRabbitmqVhost) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRabbitmqVhostConfig
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

func (c *ctyunRabbitmqVhost) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.rabbitmqService = business.NewRabbitmqService(meta)
}

func (c *ctyunRabbitmqVhost) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [name],[instanceID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRabbitmqVhostConfig
	var name, instanceID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &name, &instanceID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instanceID, &regionID)
		if err != nil {
			return
		}
	}
	if name == "" {
		err = fmt.Errorf("name不能为空")
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
	cfg.InstanceID = types.StringValue(instanceID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunRabbitmqVhost) checkBeforeCreate(ctx context.Context, plan CtyunRabbitmqVhostConfig) (err error) {
	name, instanceID, regionID := plan.Name.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	exist, err := c.rabbitmqService.CheckVhostExist(ctx, name, instanceID, regionID)
	if err != nil {
		return
	}
	if exist {
		err = fmt.Errorf("rabbitmq vhost %s already exist", name)
	}
	return
}

// create 创建
func (c *ctyunRabbitmqVhost) create(ctx context.Context, plan CtyunRabbitmqVhostConfig) (err error) {
	params := &amqp.AmqpVhostCreateV3Request{
		RegionId:   plan.RegionID.ValueString(),
		Name:       plan.Name.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpVhostCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// delete 创建
func (c *ctyunRabbitmqVhost) delete(ctx context.Context, plan CtyunRabbitmqVhostConfig) (err error) {
	params := &amqp.AmqpVhostDeleteV3Request{
		RegionId:   plan.RegionID.ValueString(),
		Name:       plan.Name.ValueString(),
		ProdInstId: plan.InstanceID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkAmqpApis.AmqpVhostDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunRabbitmqVhost) getAndMerge(ctx context.Context, plan *CtyunRabbitmqVhostConfig) (err error) {
	name, instanceID, regionID := plan.Name.ValueString(), plan.InstanceID.ValueString(), plan.RegionID.ValueString()
	exist, err := c.rabbitmqService.CheckVhostExist(ctx, name, instanceID, regionID)
	if err != nil {
		return
	}
	if !exist {
		err = common.ResourceNotExistError
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", name, instanceID, regionID))
	return
}
