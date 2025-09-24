package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunVpcRouteTable{}
	_ resource.ResourceWithConfigure   = &ctyunVpcRouteTable{}
	_ resource.ResourceWithImportState = &ctyunVpcRouteTable{}
)

type ctyunVpcRouteTable struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcRouteTable() resource.Resource {
	return &ctyunVpcRouteTable{}
}

func (c *ctyunVpcRouteTable) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_route_table"
}

type CtyunVpcRouteTableConfig struct {
	ID           types.String `tfsdk:"id"`
	RouteTableID types.String `tfsdk:"route_table_id"`
	RegionID     types.String `tfsdk:"region_id"`
	VpcID        types.String `tfsdk:"vpc_id"`
	Name         types.String `tfsdk:"name"`
	ProjectID    types.String `tfsdk:"project_id"`
}

func (c *ctyunVpcRouteTable) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027724**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，值与路由表ID相同",
			},
			"route_table_id": schema.StringAttribute{
				Computed:    true,
				Description: "路由表ID",
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "名称，支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}-]+$"), "名称不符合规则"),
				},
			},
		},
	}
}

func (c *ctyunVpcRouteTable) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpcRouteTableConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	routeTableID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.RouteTableID = types.StringValue(routeTableID)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpcRouteTable) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcRouteTableConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpcRouteTable) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpcRouteTableConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpcRouteTableConfig
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

func (c *ctyunVpcRouteTable) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcRouteTableConfig
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

func (c *ctyunVpcRouteTable) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [routeTableID],[regionID]
func (c *ctyunVpcRouteTable) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpcRouteTableConfig
	var routeTableID, regionID string
	err = terraform_extend.Split(request.ID, &routeTableID, &regionID)
	if err != nil {
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

// checkBeforeCreate 创建前检查
func (c *ctyunVpcRouteTable) checkBeforeCreate(ctx context.Context, plan CtyunVpcRouteTableConfig) (err error) {
	vpcID, regionID, projectID := plan.VpcID.ValueString(), plan.RegionID.ValueString(), plan.ProjectID.ValueString()
	err = business.NewVpcService(c.meta).MustExist(ctx, vpcID, regionID, projectID)
	return
}

// create 创建路由表
func (c *ctyunVpcRouteTable) create(ctx context.Context, plan CtyunVpcRouteTableConfig) (routeTableID string, err error) {
	vpcID, regionID, projectID := plan.VpcID.ValueString(), plan.RegionID.ValueString(), plan.ProjectID.ValueString()
	params := &ctvpc.CtvpcCreateRouteTableRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionID,
		VpcID:       vpcID,
		Name:        plan.Name.ValueString(),
	}
	if projectID != "" {
		params.ProjectID = &projectID
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	routeTableID = *resp.ReturnObj.Id
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpcRouteTable) getAndMerge(ctx context.Context, plan *CtyunVpcRouteTableConfig) (err error) {
	routeTableID, regionID := plan.RouteTableID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcShowRouteTableRequest{
		RegionID:     regionID,
		RouteTableID: routeTableID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	plan.VpcID = utils.SecStringValue(resp.ReturnObj.VpcID)
	plan.Name = utils.SecStringValue(resp.ReturnObj.Name)
	plan.ID = plan.RouteTableID
	return
}

// update 更新路由表
func (c *ctyunVpcRouteTable) update(ctx context.Context, plan, state CtyunVpcRouteTableConfig) (err error) {
	if plan.Name.Equal(state.Name) {
		return
	}
	routeTableID, regionID, name := state.RouteTableID.ValueString(), state.RegionID.ValueString(), plan.Name.ValueString()
	params := &ctvpc.CtvpcUpdateRouteTableAttributeRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     regionID,
		RouteTableID: routeTableID,
		Name:         &name,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateRouteTableAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// delete 删除路由表
func (c *ctyunVpcRouteTable) delete(ctx context.Context, plan CtyunVpcRouteTableConfig) (err error) {
	routeTableID, regionID := plan.RouteTableID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteRouteTableRequest{
		RegionID:     regionID,
		RouteTableID: routeTableID,
		ClientToken:  uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteRouteTableApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
