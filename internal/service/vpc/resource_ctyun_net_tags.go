package vpc

import (
	"context"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunNetTags{}
	_ resource.ResourceWithConfigure   = &ctyunNetTags{}
	_ resource.ResourceWithImportState = &ctyunNetTags{}
)

type ctyunNetTags struct {
	meta        *common.CtyunMetadata
	tagsService *business.TagsService
}

func NewCtyunNetTagsResource() resource.Resource {
	return &ctyunNetTags{}
}

func (c *ctyunNetTags) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_net_tags"
}

func (c *ctyunNetTags) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026753/10219975`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID，值",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "资源类型，resourceType only support vpc / subnet / acl / security_group / route_table / havip / port / multicast_domain / vpc_peer / vpce_endpoint / vpce_endpoint_service / ipv6_gateway / elb / private_nat / nat / eip / bandwidth /ipv6_bandwidth",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.NetResourceTypes...),
				},
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "资源 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"tags": schema.SetNestedAttribute{
				Optional:    true,
				Description: "标签列表。最多10个标签，标签键不可重复，键值长度1~32字符，不能换行或以空格开头/结尾。",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "标签id。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
						"key": schema.StringAttribute{
							Required:    true,
							Description: "标签键。支持更新",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "标签值。支持更新",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
					},
				},
			},
		},
	}
}

// Create 创建nat
func (c *ctyunNetTags) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunNetTagsConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.tagsService.BindTags(ctx, plan.RegionID.ValueString(), plan.ResourceType.ValueString(), plan.ResourceID.ValueString(), &plan.Tags)
	if err != nil {
		return
	}
	//ID 由RegionID  ResourceType ResourceID 组合
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.RegionID.ValueString(), plan.ResourceType.ValueString(), plan.ResourceID.ValueString()))
	// 创建后反查创建后的nat信息
	err = c.getAndMergeNat(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunNetTags) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunNetTagsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergeNat(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)

}

func (c *ctyunNetTags) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunNetTagsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunNetTagsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.tagsService.UpdateBind(ctx, state.RegionID.ValueString(), state.ResourceType.ValueString(), state.ResourceID.ValueString(), &plan.Tags)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeNat(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunNetTags) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 获取state
	var state CtyunNetTagsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.tagsService.Unbind(ctx, state.RegionID.ValueString(), state.ResourceType.ValueString(), state.ResourceID.ValueString(), &state.Tags)
	if err != nil {
		return
	}

	// 私网NAT删除API没有返回MasterResourceStatus，所以跳过状态检查
	return
}

func (c *ctyunNetTags) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.tagsService = business.NewTagsService(c.meta)
}

func (c *ctyunNetTags) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [resourceType],[resourceID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var config CtyunNetTagsConfig
	var regionID, resourceType, resourceID string
	// 根据分隔符数量判断是否输入了regionID,
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &resourceType, &resourceID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &resourceType, &resourceID, &regionID)
		if err != nil {
			return
		}
	}

	err = terraform_extend.Split(request.ID, &regionID, &resourceType, &resourceID)
	if err != nil {
		return
	}
	config.RegionID = types.StringValue(regionID)
	config.ResourceType = types.StringValue(resourceType)
	config.ResourceID = types.StringValue(resourceID)
	err = c.getAndMergeNat(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *ctyunNetTags) getAndMergeNat(ctx context.Context, plan *CtyunNetTagsConfig) (err error) {
	tags, err := c.tagsService.QueryAll(ctx, plan.RegionID.ValueString(), plan.ResourceType.ValueString(), plan.ResourceID.ValueString())
	if err != nil {
		return
	}
	plan.Tags = tags

	return
}

type CtyunNetTagsConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionID     types.String `tfsdk:"region_id"`     //区域id
	ResourceID   types.String `tfsdk:"resource_id"`   //需要创建 NAT 网关的 VPC 的 ID
	ResourceType types.String `tfsdk:"resource_type"` //需要创建 NAT 网关的 VPC 的 ID
	Tags         types.Set    `tfsdk:"tags"`
}
