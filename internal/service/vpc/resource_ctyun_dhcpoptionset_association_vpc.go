package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                = &ctyunDhcpOptionSetAssociationVpc{}
	_ resource.ResourceWithConfigure   = &ctyunDhcpOptionSetAssociationVpc{}
	_ resource.ResourceWithImportState = &ctyunDhcpOptionSetAssociationVpc{}
)

func NewCtyunDhcpOptionSetAssociationVpc() resource.Resource {
	return &ctyunDhcpOptionSetAssociationVpc{}
}

type ctyunDhcpOptionSetAssociationVpc struct {
	meta *common.CtyunMetadata
}

func (c *ctyunDhcpOptionSetAssociationVpc) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunDhcpOptionSetAssociationVpcConfig
	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionID,
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)

		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	config.DhcpOptionSetsId = types.StringValue(ID)
	config.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunDhcpOptionSetAssociationVpc) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dhcpoptionset_association_vpc"
}

func (c *ctyunDhcpOptionSetAssociationVpc) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10381274`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源唯一标识，格式为 region_id:dhcp_option_sets_id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"dhcp_option_sets_id": schema.StringAttribute{
				Required:    true,
				Description: "DHCP选项集ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "VPC ID列表 支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(0),
				},
			},
		},
	}
}

func (c *ctyunDhcpOptionSetAssociationVpc) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunDhcpOptionSetAssociationVpc) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {

	var plan CtyunDhcpOptionSetAssociationVpcConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.create(ctx, &plan)
	if err != nil {
		response.Diagnostics.AddError(
			"绑定DHCP选项集和VPC失败",
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSetAssociationVpc) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunDhcpOptionSetAssociationVpcConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMerge(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"读取DHCP选项集和VPC绑定关系失败",
			err.Error(),
		)
		return
	}

	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunDhcpOptionSetAssociationVpc) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunDhcpOptionSetAssociationVpcConfig
	var state CtyunDhcpOptionSetAssociationVpcConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.update(ctx, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"更新DHCP选项集和VPC绑定关系失败",
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSetAssociationVpc) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunDhcpOptionSetAssociationVpcConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.delete(ctx, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"解除DHCP选项集和VPC绑定关系失败",
			err.Error(),
		)
		return
	}
}

// create 绑定DHCP选项集和VPC
func (c *ctyunDhcpOptionSetAssociationVpc) create(ctx context.Context, plan *CtyunDhcpOptionSetAssociationVpcConfig) (err error) {
	// 遍历所有VPC ID进行绑定
	for _, vpcId := range plan.VpcIds {
		// 调用API绑定DHCP选项集和VPC
		resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpassociatevpcApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpassociatevpcRequest{
			RegionID:         plan.RegionId.ValueString(),
			DhcpOptionSetsID: plan.DhcpOptionSetsId.ValueString(),
			VpcID:            vpcId,
		})
		if err != nil {
			return err
		}

		if resp.StatusCode != 800 {
			return fmt.Errorf("API返回错误: %s (%s)", *resp.Message, *resp.Description)
		}
	}

	// 设置资源ID
	plan.Id = types.StringValue(fmt.Sprintf("%s:%s",
		plan.RegionId.ValueString(),
		plan.DhcpOptionSetsId.ValueString()))

	return nil
}

// getAndMerge 查询DHCP选项集和VPC绑定关系并合并状态
func (c *ctyunDhcpOptionSetAssociationVpc) getAndMerge(ctx context.Context, state *CtyunDhcpOptionSetAssociationVpcConfig) (*CtyunDhcpOptionSetAssociationVpcConfig, error) {
	// 调用API获取绑定的VPC列表
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcplistvpcApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcplistvpcRequest{
		RegionID:         state.RegionId.ValueString(),
		DhcpOptionSetsID: state.DhcpOptionSetsId.ValueString(),
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 800 {
		return nil, fmt.Errorf("API返回错误: %s (%s)", *resp.Message, *resp.Description)
	}

	if resp.ReturnObj == nil {
		return nil, nil
	}

	// 更新状态
	state.Id = types.StringValue(fmt.Sprintf("%s:%s",
		state.RegionId.ValueString(),
		state.DhcpOptionSetsId.ValueString()))

	// 更新VPC ID列表
	state.VpcIds = []string{}
	for _, vpc := range resp.ReturnObj.Results {
		if vpc.VpcID != nil {
			state.VpcIds = append(state.VpcIds, *vpc.VpcID)
		}
	}

	return state, nil
}

// update 更新DHCP选项集和VPC绑定关系
func (c *ctyunDhcpOptionSetAssociationVpc) update(ctx context.Context, plan *CtyunDhcpOptionSetAssociationVpcConfig, state *CtyunDhcpOptionSetAssociationVpcConfig) (err error) {
	// 计算需要新增和删除的VPC ID
	toAdd, toRemove := c.diffVpcIds(plan.VpcIds, state.VpcIds)

	// 解除不再需要的VPC绑定
	for _, vpcId := range toRemove {
		resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpdisassociatevpcApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpdisassociatevpcRequest{
			RegionID:         plan.RegionId.ValueString(),
			DhcpOptionSetsID: plan.DhcpOptionSetsId.ValueString(),
			VpcID:            vpcId,
		})
		if err != nil {
			return err
		}

		if resp.StatusCode != 800 {
			return fmt.Errorf("API返回错误: %s (%s)", *resp.Message, *resp.Description)
		}
	}

	// 添加新的VPC绑定
	for _, vpcId := range toAdd {
		resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpassociatevpcApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpassociatevpcRequest{
			RegionID:         plan.RegionId.ValueString(),
			DhcpOptionSetsID: plan.DhcpOptionSetsId.ValueString(),
			VpcID:            vpcId,
		})
		if err != nil {
			return err
		}

		if resp.StatusCode != 800 {
			return fmt.Errorf("API返回错误: %s (%s)", *resp.Message, *resp.Description)
		}
	}

	// 设置资源ID
	plan.Id = types.StringValue(fmt.Sprintf("%s:%s",
		plan.RegionId.ValueString(),
		plan.DhcpOptionSetsId.ValueString()))

	return nil
}

// delete 解除DHCP选项集和VPC绑定关系
func (c *ctyunDhcpOptionSetAssociationVpc) delete(ctx context.Context, state *CtyunDhcpOptionSetAssociationVpcConfig) (err error) {
	// 解除所有VPC的绑定
	for _, vpcId := range state.VpcIds {
		resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpdisassociatevpcApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpdisassociatevpcRequest{
			RegionID:         state.RegionId.ValueString(),
			DhcpOptionSetsID: state.DhcpOptionSetsId.ValueString(),
			VpcID:            vpcId,
		})
		if err != nil {
			// 忽略错误，因为VPC可能已经被删除
			continue
		}

		if resp.StatusCode != 800 {
			// 忽略错误，因为VPC可能已经被删除
			continue
		}
	}

	return nil
}

// diffVpcIds 计算VPC ID列表的差异
func (c *ctyunDhcpOptionSetAssociationVpc) diffVpcIds(planIds, stateIds []string) (toAdd, toRemove []string) {
	planSet := make(map[string]bool)
	stateSet := make(map[string]bool)

	// 构建计划和状态的集合
	for _, id := range planIds {
		planSet[id] = true
	}
	for _, id := range stateIds {
		stateSet[id] = true
	}

	// 找出需要添加的VPC ID
	for id := range planSet {
		if !stateSet[id] {
			toAdd = append(toAdd, id)
		}
	}

	// 找出需要移除的VPC ID
	for id := range stateSet {
		if !planSet[id] {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}

type CtyunDhcpOptionSetAssociationVpcConfig struct {
	Id               types.String `tfsdk:"id"`
	RegionId         types.String `tfsdk:"region_id"`
	DhcpOptionSetsId types.String `tfsdk:"dhcp_option_sets_id"`
	VpcIds           []string     `tfsdk:"vpc_ids"`
}
