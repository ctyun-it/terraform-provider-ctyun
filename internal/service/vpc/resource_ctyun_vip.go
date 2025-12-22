package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunVip{}
	_ resource.ResourceWithConfigure   = &CtyunVip{}
	_ resource.ResourceWithImportState = &CtyunVip{}
)

func NewCtyunVip() resource.Resource {
	return &CtyunVip{}
}

type CtyunVip struct {
	meta *common.CtyunMetadata
}

func (c *CtyunVip) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vip"
}

func (c *CtyunVip) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026730/10224288`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "高可用虚IP的ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
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
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "VPC的ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"ip_address": schema.StringAttribute{
				Optional:    true,
				Description: "ip地址",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					validator2.Ip(),
				},
			},
			"vip_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "虚拟IP的类型，v4-IPv4类型虚IP，v6-IPv6类型虚IP。默认为v4",
				Default:     stringdefault.StaticString("v4"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("v4", "v6"),
				},
			},
			"ipv4_address": schema.StringAttribute{
				Computed:    true,
				Description: "高可用虚IP的IPv4地址",
			},
			"ipv6_address": schema.StringAttribute{
				Computed:    true,
				Description: "高可用虚IP的IPv6地址",
			},
		},
	}
}

func (c *CtyunVip) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunVip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunVipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *CtyunVip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunVipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunVip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// HaVip资源不支持更新操作，直接返回
	return
}

func (c *CtyunVip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunVipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.delete(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError("删除高可用虚拟IP失败", err.Error())
		return
	}
}

func (c *CtyunVip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id],[projectID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var state CtyunVipConfig
	var vipId, regionId string
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		vipId = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &vipId, &regionId)
		if err != nil {
			return
		}
	}

	if vipId == "" {
		err = fmt.Errorf("vipId不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	state.Id = types.StringValue(vipId)
	state.RegionId = types.StringValue(regionId)

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// create 创建havip
func (c *CtyunVip) create(ctx context.Context, plan *CtyunVipConfig) (err error) {
	// 准备请求参数
	req := &ctvpc.CtvpcCreateHavipRequest{
		ClientToken: uuid.NewString(),
		RegionID:    plan.RegionId.ValueString(),
		SubnetID:    plan.SubnetId.ValueString(),
	}

	// 设置可选参数
	if !plan.VpcId.IsNull() {
		vpcId := plan.VpcId.ValueString()
		req.NetworkID = &vpcId
	}

	if !plan.IpAddress.IsNull() {
		ipAddress := plan.IpAddress.ValueString()
		req.IpAddress = &ipAddress
	}

	if !plan.VipType.IsNull() {
		vipType := plan.VipType.ValueString()
		req.VipType = &vipType
	}

	// 调用API创建HaVip
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateHavipApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置资源ID和其他属性
	plan.Id = types.StringValue(*resp.ReturnObj.Uuid)

	if resp.ReturnObj.Ipv4 != nil {
		plan.Ipv4Address = types.StringValue(*resp.ReturnObj.Ipv4)
	}

	if resp.ReturnObj.Ipv6 != nil {
		plan.Ipv6Address = types.StringValue(*resp.ReturnObj.Ipv6)
	}

	return nil
}

// getAndMerge 查询havip并合并状态
func (c *CtyunVip) getAndMerge(ctx context.Context, state *CtyunVipConfig) (err error) {
	// 调用API获取HaVip详情
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowHavipApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowHavipRequest{
		RegionID: state.RegionId.ValueString(),
		HaVipID:  state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新状态
	returnObj := resp.ReturnObj

	if returnObj.Id != nil {
		state.Id = types.StringValue(*returnObj.Id)
	}

	if returnObj.Ipv4 != nil {
		state.Ipv4Address = types.StringValue(*returnObj.Ipv4)
	}

	//if returnObj.Ipv6 != nil {
	//	state.Ipv6Address = types.StringValue(*returnObj.Ipv6)
	//}

	if returnObj.VpcID != nil {
		state.VpcId = types.StringValue(*returnObj.VpcID)
	}

	if returnObj.SubnetID != nil {
		state.SubnetId = types.StringValue(*returnObj.SubnetID)
	}

	return
}

// delete 删除havip
func (c *CtyunVip) delete(ctx context.Context, state *CtyunVipConfig) (err error) {
	// 调用API删除HaVip
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteHavipApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDeleteHavipRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionId.ValueString(),
		HaVipID:     state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}

	return nil
}

type CtyunVipConfig struct {
	Id          types.String `tfsdk:"id"`
	RegionId    types.String `tfsdk:"region_id"`
	VpcId       types.String `tfsdk:"vpc_id"`
	SubnetId    types.String `tfsdk:"subnet_id"`
	IpAddress   types.String `tfsdk:"ip_address"`
	VipType     types.String `tfsdk:"vip_type"`
	Ipv4Address types.String `tfsdk:"ipv4_address"`
	Ipv6Address types.String `tfsdk:"ipv6_address"`
	ProjectId   types.String `tfsdk:"project_id"`
}
