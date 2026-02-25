package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
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
	_ resource.Resource                = &ctyunVpcRouteAssociationGatewayTable{}
	_ resource.ResourceWithConfigure   = &ctyunVpcRouteAssociationGatewayTable{}
	_ resource.ResourceWithImportState = &ctyunVpcRouteAssociationGatewayTable{}
)

type ctyunVpcRouteAssociationGatewayTable struct {
	meta *common.CtyunMetadata
	name string
}

func NewctyunVpcRouteAssociationGatewayTable() resource.Resource {
	return &ctyunVpcRouteAssociationGatewayTable{}
}

func (c *ctyunVpcRouteAssociationGatewayTable) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_route_association_gateway"
	c.name = response.TypeName
}

type ctyunVpcRouteAssociationGatewayTableConfig struct {
	ID           types.String `tfsdk:"id"`
	RouteTableID types.String `tfsdk:"route_table_id"`
	RegionID     types.String `tfsdk:"region_id"`
}

func (c *ctyunVpcRouteAssociationGatewayTable) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理虚拟私有云路由表与IPV4网关关联", "虚拟私有云（Virtual Private Cloud，VPC）", "https://www.ctyun.cn/document/10027724"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID，值与路由表ID相同",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"route_table_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
		},
	}
}

func (c *ctyunVpcRouteAssociationGatewayTable) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan ctyunVpcRouteAssociationGatewayTableConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
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

func (c *ctyunVpcRouteAssociationGatewayTable) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state ctyunVpcRouteAssociationGatewayTableConfig
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

func (c *ctyunVpcRouteAssociationGatewayTable) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan ctyunVpcRouteAssociationGatewayTableConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state ctyunVpcRouteAssociationGatewayTableConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpcRouteAssociationGatewayTable) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state ctyunVpcRouteAssociationGatewayTableConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunVpcRouteAssociationGatewayTable) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunVpcRouteAssociationGatewayTable) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg ctyunVpcRouteAssociationGatewayTableConfig
	var routeTableID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		routeTableID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &routeTableID, &regionID)
		if err != nil {
			return
		}
	}

	if routeTableID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	cfg.RegionID = types.StringValue(regionID)
	cfg.RouteTableID = types.StringValue(routeTableID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建路由表
func (c *ctyunVpcRouteAssociationGatewayTable) create(ctx context.Context, plan ctyunVpcRouteAssociationGatewayTableConfig) (err error) {

	ipv4GwID, err := c.getIPV4Id(ctx, &plan)

	params := &ctvpc.CtvpcIpv4GwBindRouteTableRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     plan.RegionID.ValueString(),
		Ipv4GwID:     ipv4GwID,
		RouteTableID: plan.RouteTableID.ValueString(),
	}
	var resp *ctvpc.CtvpcIpv4GwBindRouteTableResponse
	resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcIpv4GwBindRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}

// getAndMerge 从远端查询
func (c *ctyunVpcRouteAssociationGatewayTable) getAndMerge(ctx context.Context, plan *ctyunVpcRouteAssociationGatewayTableConfig) (err error) {
	//routeTableID, regionID := plan.RouteTableID.ValueString(), plan.RegionID.ValueString()

	ipv4GwID, err := c.getIPV4Id(ctx, plan)

	params2 := &ctvpc.CtvpcShowIPv4GwRequest{
		RegionID: plan.RegionID.ValueString(),
		Ipv4GwID: ipv4GwID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowIPv4GwApi.Do(ctx, c.meta.SdkCredential, params2)
	if err != nil {
		return
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiRouteTableAccessFailed {
		err = common.ResourceNotExistError
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	plan.ID = plan.RouteTableID
	return
}

func (c *ctyunVpcRouteAssociationGatewayTable) getIPV4Id(ctx context.Context, plan *ctyunVpcRouteAssociationGatewayTableConfig) (ipv4GwID string, err error) {
	routeTableID, regionID := plan.RouteTableID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcShowRouteTableRequest{
		RegionID:     regionID,
		RouteTableID: routeTableID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if utils.SecString(resp.ErrorCode) == common.OpenapiRouteTableAccessFailed {
		err = common.ResourceNotExistError
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	vpcID := utils.SecStringValue(resp.ReturnObj.VpcID)

	gatewaparams := &ctvpc.CtvpcListIPv4GwRequest{
		RegionID: regionID,
		VpcID:    vpcID.ValueStringPointer(),
	}
	var respGateway *ctvpc.CtvpcListIPv4GwResponse
	respGateway, err = c.meta.Apis.SdkCtVpcApis.CtvpcListIPv4GwApi.Do(ctx, c.meta.SdkCredential, gatewaparams)
	if err != nil {
		return
	} else if respGateway.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *respGateway.Message, *respGateway.Description)
		return
	} else if respGateway.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ipv4GwID = *respGateway.ReturnObj[0].Id
	return ipv4GwID, nil
}

// update 更新路由表
func (c *ctyunVpcRouteAssociationGatewayTable) update(ctx context.Context, plan, state ctyunVpcRouteAssociationGatewayTableConfig) (err error) {
	return
}

// delete 删除路由表
func (c *ctyunVpcRouteAssociationGatewayTable) delete(ctx context.Context, plan ctyunVpcRouteAssociationGatewayTableConfig) (err error) {
	ipv4GwID, err := c.getIPV4Id(ctx, &plan)

	params := &ctvpc.CtvpcIpv4GwUnbindRouteTableRequest{
		RegionID:    plan.RegionID.ValueString(),
		Ipv4GwID:    ipv4GwID,
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcIpv4GwUnbindRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
