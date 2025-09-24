package vpce

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
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &ctyunVpceService{}
	_ resource.ResourceWithConfigure   = &ctyunVpceService{}
	_ resource.ResourceWithImportState = &ctyunVpceService{}
)

type ctyunVpceService struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceService() resource.Resource {
	return &ctyunVpceService{}
}

func (c *ctyunVpceService) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service"
}

type CtyunVpceServiceConfig struct {
	ID             types.String `tfsdk:"id"`
	RegionID       types.String `tfsdk:"region_id"`
	VpcID          types.String `tfsdk:"vpc_id"`
	Type           types.String `tfsdk:"type"`
	Name           types.String `tfsdk:"name"`
	InstanceType   types.String `tfsdk:"instance_type"`
	InstanceID     types.String `tfsdk:"instance_id"`
	SubnetID       types.String `tfsdk:"subnet_id"`
	AutoConnection types.Bool   `tfsdk:"auto_connection"`
	Rules          types.Set    `tfsdk:"rules"`
	WhitelistEmail types.Set    `tfsdk:"whitelist_email"`
	whitelist      []string
	rules          []CtyunVpceServiceRule
}

type CtyunVpceServiceRule struct {
	Protocol     types.String `tfsdk:"protocol"`
	ServerPort   types.Int32  `tfsdk:"server_port"`
	EndpointPort types.Int32  `tfsdk:"endpoint_port"`
}

func (c *ctyunVpceService) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10042658/10217013**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
			"type": schema.StringAttribute{
				Required:    true,
				Description: "接口还是反向，interface:接口，reverse:反向",
				Validators: []validator.String{
					stringvalidator.OneOf(business.VpceServiceTypeInterface, business.VpceServiceTypeReverse),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "支持拉丁字母、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}-]+$"), "终端节点服务名称不符合规则"),
				},
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "服务后端实例类型，vm:虚机类型,bm:物理机,vip:vip类型,lb:负载均衡类型,当type为interface时必填。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf("vm", "bm", "vip", "lb"),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeInterface),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeReverse),
					),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "服务后端实例ID，当type为interface时必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeInterface),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeReverse),
					),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "服务后端子网id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"auto_connection": schema.BoolAttribute{
				Required:    true,
				Description: "是否自动连接，true表示自动链接，false表示非自动链接，支持更新",
			},
			"whitelist_email": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "白名单邮箱，最多支持10个，支持更新",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validator2.Email()),
					setvalidator.SizeAtMost(10),
				},
			},
			"rules": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Set{
					validator2.AlsoRequiresEqualSet(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeInterface),
					),
					validator2.ConflictsWithEqualSet(
						path.MatchRoot("type"),
						types.StringValue(business.VpceServiceTypeReverse),
					),
				},
				Description: "节点服务规则，当type为interface时必填，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "协议，TCP:TCP协议,UDP:UDP协议，支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("TCP", "UDP"),
							},
						},
						"server_port": schema.Int32Attribute{
							Required:    true,
							Description: "服务端口(用于创建backend传入)(1-65535)，支持更新",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
						"endpoint_port": schema.Int32Attribute{
							Required:    true,
							Description: "节点端口(用于创建rule传入)(1-65535)，支持更新",
							Validators: []validator.Int32{
								int32validator.Between(1, 65535),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpceService) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServiceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	endpointServiceID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(endpointServiceID)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)

	err = c.checkAfterCreate(ctx, plan)
	if err != nil {
		return
	}
	err = c.calcWhitelist(ctx, &plan)
	if err != nil {
		return
	}
	err = c.addWhitelist(ctx, plan)
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

func (c *ctyunVpceService) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpceService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpceServiceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpceServiceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	plan.ID, plan.RegionID = state.ID, state.RegionID
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 更新白名单
	err = c.updateWhitelist(ctx, plan, state)
	if err != nil {
		return
	}

	if plan.Type.ValueString() == business.VpceServiceTypeInterface {
		// 更新规则
		err = c.updateRule(ctx, plan, state)
		if err != nil {
			return
		}

		// 更新后端服务
		err = c.updateBackend(ctx, plan, state)
		if err != nil {
			return
		}
	}

	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpceService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunVpceService) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [endpointServiceID],[regionID]
func (c *ctyunVpceService) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceServiceConfig
	var endpointServiceID, regionID string
	err = terraform_extend.Split(request.ID, &endpointServiceID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(endpointServiceID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpceService) create(ctx context.Context, plan CtyunVpceServiceConfig) (endpointSeverID string, err error) {
	params := &ctvpc.CtvpcCreateEndpointServiceRequest{
		ClientToken:    uuid.NewString(),
		RegionID:       plan.RegionID.ValueString(),
		VpcID:          plan.VpcID.ValueString(),
		Name:           plan.Name.ValueString(),
		RawType:        plan.Type.ValueStringPointer(),
		InstanceType:   plan.InstanceType.ValueStringPointer(),
		InstanceID:     plan.InstanceID.ValueStringPointer(),
		SubnetID:       plan.SubnetID.ValueStringPointer(),
		AutoConnection: plan.AutoConnection.ValueBool(),
	}

	err = c.calcRule(ctx, &plan)
	if err != nil {
		return
	}
	for _, r := range plan.rules {
		params.Rules = append(params.Rules, &ctvpc.CtvpcCreateEndpointServiceRulesRequest{
			Protocol:     r.Protocol.ValueString(),
			ServerPort:   r.ServerPort.ValueInt32(),
			EndpointPort: r.EndpointPort.ValueInt32(),
		})
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.EndpointService == nil {
		err = common.InvalidReturnObjError
		return
	}
	endpointSeverID = utils.SecString(resp.ReturnObj.EndpointService.EndpointServiceID)
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunVpceService) checkAfterCreate(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	endpointService, err := c.show(ctx, plan)
	if err != nil {
		return
	}
	// 后端信息不正确
	if len(endpointService.Backends) == 0 && plan.Type.ValueString() == business.VpceServiceTypeInterface {
		err = fmt.Errorf("终端节点服务创建成功，但后端服务不正确，请使用正确instance_id重新创建")
		return
	}
	return
}

// buildRules 构造终端节点服务规则
func (c *ctyunVpceService) buildRules(ctx context.Context, plan CtyunVpceServiceConfig) (rulesReq []*ctvpc.CtvpcCreateEndpointServiceRulesRequest, err error) {
	if plan.Rules.IsUnknown() || plan.Rules.IsNull() {
		return
	}
	var rules []CtyunVpceServiceRule
	diags := plan.Rules.ElementsAs(ctx, &rules, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	for _, r := range rules {
		item := &ctvpc.CtvpcCreateEndpointServiceRulesRequest{
			Protocol:     r.Protocol.ValueString(),
			ServerPort:   r.ServerPort.ValueInt32(),
			EndpointPort: r.EndpointPort.ValueInt32(),
		}
		rulesReq = append(rulesReq, item)
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpceService) getAndMerge(ctx context.Context, plan *CtyunVpceServiceConfig) (err error) {
	endpointService, err := c.show(ctx, *plan)
	if err != nil {
		return
	}
	plan.VpcID = utils.SecStringValue(endpointService.VpcID)
	plan.Name = utils.SecStringValue(endpointService.Name)
	plan.Type = utils.SecStringValue(endpointService.RawType)
	plan.AutoConnection = utils.SecBoolValue(endpointService.AutoConnection)

	if len(endpointService.Backends) != 0 {
		backend := endpointService.Backends[0]
		plan.InstanceType = utils.SecStringValue(backend.InstanceType)
		plan.InstanceID = utils.SecStringValue(backend.InstanceID)
	} else {
		plan.InstanceType = types.StringNull()
		plan.InstanceID = types.StringNull()
	}

	err = c.mergeRules(ctx, plan, endpointService)
	if err != nil {
		return
	}

	err = c.mergeWhitelist(ctx, plan)
	if err != nil {
		return
	}

	return
}

// update 更新
func (c *ctyunVpceService) update(ctx context.Context, plan, state CtyunVpceServiceConfig) (err error) {
	endpointServiceID, regionID := state.ID.ValueString(), state.RegionID.ValueString()
	params := &ctvpc.CtvpcModifyEndpointServiceRequest{
		RegionID:          regionID,
		EndpointServiceID: endpointServiceID,
		Name:              plan.Name.ValueStringPointer(),
		AutoConnection:    plan.AutoConnection.ValueBoolPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcModifyEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// delete 删除
func (c *ctyunVpceService) delete(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	endpointServiceID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteEndpointServiceRequest{
		RegionID:    regionID,
		ID:          endpointServiceID,
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// show 查询VPCEs详情
func (c *ctyunVpceService) show(ctx context.Context, plan CtyunVpceServiceConfig) (endpointService ctvpc.CtvpcShowEndpointServiceReturnObjResponse, err error) {
	endpointServiceID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcShowEndpointServiceRequest{
		RegionID:          regionID,
		EndpointServiceID: endpointServiceID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowEndpointServiceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	endpointService = *resp.ReturnObj
	return
}

// calcWhitelist 将types.Set类型的白名单转换为[]string
func (c *ctyunVpceService) calcWhitelist(ctx context.Context, plan *CtyunVpceServiceConfig) (err error) {
	if plan.WhitelistEmail.IsNull() || plan.WhitelistEmail.IsUnknown() {
		return
	}
	plan.whitelist = []string{}
	diags := plan.WhitelistEmail.ElementsAs(ctx, &plan.whitelist, true)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
	}
	return
}

// addWhitelist 添加白名单
func (c *ctyunVpceService) addWhitelist(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	for i, _ := range plan.whitelist {
		email := plan.whitelist[i]
		params := &ctvpc.CtvpcCreateEndpointServiceWhitelistRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Email:             &email,
		}
		var resp *ctvpc.CtvpcCreateEndpointServiceWhitelistResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceWhitelistApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delWhitelist 删除白名单
func (c *ctyunVpceService) delWhitelist(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	for i, _ := range plan.whitelist {
		email := plan.whitelist[i]
		params := &ctvpc.CtvpcDeleteEndpointServiceWhitelistRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Email:             &email,
		}
		var resp *ctvpc.CtvpcDeleteEndpointServiceWhitelistResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceWhitelistApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// updateWhitelist 更新白名单
func (c *ctyunVpceService) updateWhitelist(ctx context.Context, plan, state CtyunVpceServiceConfig) (err error) {
	err = c.calcWhitelist(ctx, &plan)
	if err != nil {
		return
	}
	err = c.calcWhitelist(ctx, &state)
	if err != nil {
		return
	}

	add, del := utils.DifferenceStrArray(plan.whitelist, state.whitelist)
	plan.whitelist = del
	err = c.delWhitelist(ctx, plan)
	if err != nil {
		return
	}
	plan.whitelist = add
	err = c.addWhitelist(ctx, plan)
	if err != nil {
		return
	}
	return
}

// mergeWhitelist 查询当前白名单
func (c *ctyunVpceService) mergeWhitelist(ctx context.Context, plan *CtyunVpceServiceConfig) (err error) {
	params := ctvpc.CtvpcNewEndpointServiceWhiteListRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewEndpointServiceWhiteListApi.Do(ctx, c.meta.SdkCredential, &params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	whitelist := []string{}
	for _, item := range resp.ReturnObj.Whitelist {
		if item != nil && item.Email != nil {
			whitelist = append(whitelist, *item.Email)
		}
	}
	t, diags := types.SetValueFrom(ctx, types.StringType, whitelist)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	plan.WhitelistEmail = t
	plan.whitelist = whitelist
	return
}

// addRule 新增端口映射
func (c *ctyunVpceService) addRule(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	for _, rule := range plan.rules {
		params := &ctvpc.CtvpcCreateEndpointServiceRuleRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Protocol:          rule.Protocol.ValueString(),
			EndpointPort:      rule.EndpointPort.ValueInt32(),
			ServerPort:        rule.ServerPort.ValueInt32(),
		}
		var resp *ctvpc.CtvpcCreateEndpointServiceRuleResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceRuleApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delRule 删除端口映射
func (c *ctyunVpceService) delRule(ctx context.Context, plan CtyunVpceServiceConfig) (err error) {
	for _, rule := range plan.rules {
		params := &ctvpc.CtvpcDeleteEndpointServiceRuleRequest{
			ClientToken:       uuid.NewString(),
			RegionID:          plan.RegionID.ValueString(),
			EndpointServiceID: plan.ID.ValueString(),
			Protocol:          rule.Protocol.ValueString(),
			EndpointPort:      rule.EndpointPort.ValueInt32(),
			ServerPort:        rule.ServerPort.ValueInt32(),
		}
		var resp *ctvpc.CtvpcDeleteEndpointServiceRuleResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceRuleApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		}
	}
	return
}

// delRule 删除端口映射
func (c *ctyunVpceService) calcRule(ctx context.Context, plan *CtyunVpceServiceConfig) (err error) {
	if plan.Rules.IsUnknown() || plan.Rules.IsNull() {
		return
	}
	plan.rules = []CtyunVpceServiceRule{}
	diags := plan.Rules.ElementsAs(ctx, &plan.rules, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	return
}

// updateRule 更新端口映射
func (c *ctyunVpceService) updateRule(ctx context.Context, plan, state CtyunVpceServiceConfig) (err error) {
	err = c.calcRule(ctx, &plan)
	if err != nil {
		return
	}
	err = c.calcRule(ctx, &state)
	if err != nil {
		return
	}

	add, del := utils.DifferenceStructArray[CtyunVpceServiceRule](plan.rules, state.rules)
	plan.rules = del
	err = c.delRule(ctx, plan)
	if err != nil {
		return
	}
	plan.rules = add
	err = c.addRule(ctx, plan)
	if err != nil {
		return
	}
	return
}

// mergeRules 计算当前rules
func (c *ctyunVpceService) mergeRules(ctx context.Context, plan *CtyunVpceServiceConfig, endpointService ctvpc.CtvpcShowEndpointServiceReturnObjResponse) (err error) {
	rules := []CtyunVpceServiceRule{}
	for _, item := range endpointService.Rules {
		if item != nil {
			rules = append(rules, CtyunVpceServiceRule{
				Protocol:     utils.SecStringValue(item.Protocol),
				EndpointPort: types.Int32Value(item.EndpointPort),
				ServerPort:   types.Int32Value(item.ServerPort),
			})
		}
	}
	ruleObj := utils.StructToTFObjectTypes(CtyunVpceServiceRule{})
	t, diags := types.SetValueFrom(ctx, ruleObj, rules)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	plan.Rules = t
	plan.rules = rules
	return
}

// updateBackend 更新后端服务
func (c *ctyunVpceService) updateBackend(ctx context.Context, plan, state CtyunVpceServiceConfig) (err error) {
	if plan.InstanceType.Equal(state.InstanceType) && plan.InstanceID.Equal(state.InstanceID) {
		return
	}

	endpointServiceID, regionID := state.ID.ValueString(), state.RegionID.ValueString()
	p := plan.InstanceType.ValueString()
	if p == "lb" {
		p = "elb"
	}
	params := &ctvpc.CtvpcVpceUpdateBackendRequest{
		RegionID:          regionID,
		EndpointServiceID: endpointServiceID,
		InstanceID:        plan.InstanceID.ValueString(),
		InstanceType:      p,
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcVpceUpdateBackendApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
