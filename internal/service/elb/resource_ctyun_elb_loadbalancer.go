package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunElbLoadBalancerResource{}
	_ resource.ResourceWithConfigure   = &CtyunElbLoadBalancerResource{}
	_ resource.ResourceWithImportState = &CtyunElbLoadBalancerResource{}
)

type CtyunElbLoadBalancerResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbLoadBalancer() resource.Resource {
	return &CtyunElbLoadBalancerResource{}
}
func (c *CtyunElbLoadBalancerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_loadbalancer"
}

func (c *CtyunElbLoadBalancerResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026756/10138703`,
		Attributes: map[string]schema.Attribute{
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
					stringvalidator.LengthBetween(0, 128),
				},
			},
			"eip_id": schema.StringAttribute{
				Optional:    true,
				Description: "弹性公网IP的ID。当resource_type=external为必填",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("resource_type"),
						types.StringValue(business.LbResourceTypeExternal),
					),
					validator2.EipValidate(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sla_name": schema.StringAttribute{
				Required:    true,
				Description: "lb的规格名称，支持:elb.s2.small（标准型Ⅰ），elb.s3.small（增强型Ⅰ），elb.s4.small（高阶型Ⅰ），elb.s5.small（超强型Ⅰ），elb.s2.large（标准型Ⅱ），elb.s3.large（增强型Ⅱ），elb.s4.large（高阶型Ⅱ），elb.s5.large（超强型Ⅱ），支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(append(business.ElbSlaNames, business.PgElbSlaNames...)...),
				},
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "资源类型。internal：内网负载均衡，external：公网负载均衡",
				Validators: []validator.String{
					stringvalidator.OneOf(business.LbResourceType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_ip_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "负载均衡的私有IP地址，不指定则自动分配",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "负载均衡Id",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				// az时候有必要设定默认值
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡实例默认创建port ID",
			},
			"ipv6_address": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡实例的IPv6地址",
			},
			"admin_status": schema.StringAttribute{
				Computed:    true,
				Description: "管理状态: DOWN / ACTIVE",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "负载均衡状态: DOWN / ACTIVE",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：year：按年，month：按月，on_demand：按需。当此值为month或year时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "订购时长, 当 cycleType = month, 支持订购 1 - 11 个月; 当 cycleType = year, 支持订购 1 - 3 年",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 3),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"pay_voucher_price": schema.StringAttribute{
				Optional:    true,
				Description: "代金券金额，支持到小数点后两位",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"eip_info": schema.ListNestedAttribute{
				Computed:    true,
				Description: "弹性公网IP信息",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"eip_id": schema.StringAttribute{
							Computed:    true,
							Description: "弹性公网IP的Id",
						},
						"bandwidth": schema.Float32Attribute{
							Computed:    true,
							Description: "弹性公网IP的带宽",
						},
						"is_talk_order": schema.BoolAttribute{
							Computed:    true,
							Description: "是否按需资源",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunElbLoadBalancerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查
	err = c.checkBeforeCreateElb(ctx, plan)
	if err != nil {
		return
	}

	returnObj, params, err := c.createPgElb(ctx, &plan)
	if err != nil {
		return
	}

	masterOrderId := returnObj.MasterOrderID
	// 创建保障型负载均衡为异步接口，需要轮询请求获取id
	loopResp, err := c.orderLoop(ctx, params, 600)
	if err != nil {
		return
	} else if loopResp == nil {
		err = common.InvalidReturnObjError
		return
	} else if loopResp.MasterOrderID != masterOrderId {
		err = fmt.Errorf("创建nat时订单ID和轮询订单ID不一致，创建时订单ID：%s, 轮询所得订单ID：%s", masterOrderId, loopResp.MasterOrderID)
	} else if loopResp.RegionID != plan.RegionID.ValueString() {
		err = fmt.Errorf("创建nat时regionId和轮询结果regionId不一致，计划的regionId：%s, 轮询所得regionId：%s", plan.RegionID.ValueString(), loopResp.RegionID)
	}
	// 将轮询所得elb id 存储plan中
	plan.ID = types.StringValue(loopResp.ElbID)

	// 创建后反查创建后的nat信息
	err = c.getAndMergeElb(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbLoadBalancerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbLoadBalancerConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeElb(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbLoadBalancerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新基本信息
	err = c.updateElbInfo(ctx, state, plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeElb(ctx, &state)
	if err != nil {
		return
	}

	// 变配
	err = c.modifyPgElbSpec(ctx, state, plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeElb(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbLoadBalancerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbLoadBalancerConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &ctelb.CtelbRefundPgelbRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ElbID:       state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbRefundPgelbApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 异步接口，需要轮询查看是否退订成功
	_, err = c.deleteLoop(ctx, params, 600)
	if err != nil {
		return
	}
}
func (c *CtyunElbLoadBalancerResource) ImportState(_ context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

}

func (c *CtyunElbLoadBalancerResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// createPgElb 创建保障型负载均衡
func (c *CtyunElbLoadBalancerResource) createPgElb(ctx context.Context, plan *CtyunElbLoadBalancerConfig) (returnObj ctelb.CtelbCreatePgelbReturnObjResponse, params *ctelb.CtelbCreatePgelbRequest, err error) {
	params = &ctelb.CtelbCreatePgelbRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     plan.RegionID.ValueString(),
		SubnetID:     plan.SubnetID.ValueString(),
		Name:         plan.Name.ValueString(),
		SlaName:      plan.SlaName.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		CycleType:    plan.CycleType.ValueString(),
		CycleCount:   int32(plan.CycleCount.ValueInt64()),
	}
	if !plan.ProjectID.IsNull() {
		params.ProjectID = plan.ProjectID.ValueString()
	}
	if !plan.VpcID.IsNull() {
		params.VpcID = plan.VpcID.ValueString()
	}

	if plan.ResourceType.ValueString() == business.LbResourceTypeExternal || !plan.EipID.IsNull() {
		params.EipID = plan.EipID.ValueString()
	}
	if !plan.PrivateIpAddress.IsNull() {
		params.PrivateIpAddress = plan.PrivateIpAddress.ValueString()
	}

	if !plan.PayVoucherPrice.IsNull() {
		params.PayVoucherPrice = plan.PayVoucherPrice.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreatePgelbApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	returnObj = *resp.ReturnObj
	return
}

func (c *CtyunElbLoadBalancerResource) checkBeforeCreateElb(_ context.Context, plan CtyunElbLoadBalancerConfig) error {
	// regionid不能为空，subnetID	(子网id)不能为空,name不能为空，slaName不能为空，resourceType不能为空
	regionId := plan.RegionID
	subnetId := plan.SubnetID
	slaName := plan.SlaName
	resourceType := plan.ResourceType
	name := plan.Name
	eipId := plan.EipID
	if regionId.IsNull() {
		return fmt.Errorf("regionID不能为空!")
	}
	if subnetId.IsNull() {
		return fmt.Errorf("subnetId-子网的ID不能为空!")
	}
	if slaName.IsNull() {
		return fmt.Errorf("slaName-lb的规格名称不能为空！")
	}
	if resourceType.IsNull() {
		return fmt.Errorf("resourceType-资源类型不能为空！")
	}
	if !c.isContains(resourceType.ValueString(), business.LbResourceType) {
		return fmt.Errorf("resourceType资源类型取值存在问题，resourceType取值范围为{internal：内网负载均衡，external：公网负载均衡}")
	}
	//当resourceType=external为必填, eipID不能为空
	if resourceType.ValueString() == business.LbResourceTypeExternal && eipId.IsNull() {
		return fmt.Errorf("当resourceType=external为必填, eipID不能为空")
	}

	if name.IsNull() {
		return fmt.Errorf("name不能为空")
	}
	return nil
}

func (c *CtyunElbLoadBalancerResource) getAndMergeElb(ctx context.Context, config *CtyunElbLoadBalancerConfig) (err error) {
	params := &ctelb.CtelbShowLoadBalancerRequest{
		RegionID: config.RegionID.ValueString(),
		ElbID:    config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	//解析resp.ReturnObj, 将最新的elb信息同步到state中
	if len(resp.ReturnObj) > 1 {
		err = fmt.Errorf("ReturnObj长度>1")
		return
	}
	elbObj := resp.ReturnObj[0]
	if config.RegionID.ValueString() != elbObj.RegionID {
		err = fmt.Errorf("elb详情regionid(%s)与plan的reigonid(%s)不一致！", elbObj.RegionID, config.RegionID.ValueString())
		return
	}
	if config.ID.ValueString() != elbObj.ID {
		err = fmt.Errorf("详情elb id(%s)与plan的elb id(%s)不一致！", elbObj.RegionID, config.RegionID.ValueString())
		return
	}
	config.Name = types.StringValue(elbObj.Name)
	config.Description = types.StringValue(elbObj.Description)
	config.VpcID = types.StringValue(elbObj.VpcID)
	config.SubnetID = types.StringValue(elbObj.SubnetID)
	config.PortID = types.StringValue(elbObj.PortID)
	config.PrivateIpAddress = types.StringValue(elbObj.PrivateIpAddress)
	config.Ipv6Address = types.StringValue(elbObj.Ipv6Address)
	config.SlaName = types.StringValue(elbObj.SlaName)
	config.AdminStatus = types.StringValue(elbObj.AdminStatus)
	config.Status = types.StringValue(elbObj.Status)
	config.ResourceType = types.StringValue(elbObj.ResourceType)
	config.CreatedTime = types.StringValue(elbObj.CreatedTime)
	config.UpdatedTime = types.StringValue(elbObj.UpdatedTime)
	EipInfoList := elbObj.EipInfo
	var eipInfos []EipInfoModel
	for _, eipItem := range EipInfoList {
		var eipInfo EipInfoModel
		eipInfo.EipID = types.StringValue(eipItem.EipID)
		eipInfo.Bandwidth = types.Float32Value(eipItem.Bandwidth)
		eipInfo.IsTalkOrder = utils.SecBoolValue(eipItem.IsTalkOrder)
		eipInfos = append(eipInfos, eipInfo)
	}

	eipInfoType := utils.StructToTFObjectTypes(EipInfoModel{})
	config.EipInfo, _ = types.ListValueFrom(ctx, eipInfoType, eipInfos)
	return
}

func (c *CtyunElbLoadBalancerResource) updateElbInfo(ctx context.Context, state CtyunElbLoadBalancerConfig, plan CtyunElbLoadBalancerConfig) (err error) {
	params := &ctelb.CtelbUpdateLoadBalancerRequest{
		ClientToken: uuid.NewString(),
		RegionID:    state.RegionID.ValueString(),
		ElbID:       state.ID.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueString()
	}
	if !plan.Name.IsNull() && !plan.Name.Equal(state.Name) {
		params.Name = plan.Name.ValueString()
	}
	if params.Name == "" && params.Description == "" {
		return
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)

	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// modifyPgElbSpec 性能保障型负载均衡变配
func (c *CtyunElbLoadBalancerResource) modifyPgElbSpec(ctx context.Context, state CtyunElbLoadBalancerConfig, plan CtyunElbLoadBalancerConfig) (err error) {
	if state.SlaName.Equal(plan.SlaName) {
		if !state.PayVoucherPrice.Equal(plan.PayVoucherPrice) {
			err = fmt.Errorf("当没有触发变配时，代金券金额修改无效")
		}
		return nil
	}
	params := &ctelb.CtelbModifyPgelbSpecRequest{
		ClientToken:     uuid.NewString(),
		RegionID:        state.RegionID.ValueString(),
		ElbID:           state.ID.ValueString(),
		SlaName:         plan.SlaName.ValueString(),
		PayVoucherPrice: plan.PayVoucherPrice.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbModifyPgelbSpecApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	// 变配后，需要轮询确认变配成功
	isModify, err := c.modifyLoop(ctx, state, plan, 600)
	if err != nil {
		return err
	}
	if !isModify {
		err = fmt.Errorf("变配失败，请重试！")
	}
	return
}

func (c *CtyunElbLoadBalancerResource) orderLoop(ctx context.Context, params *ctelb.CtelbCreatePgelbRequest, loopCount ...int) (loopResp *ctelb.CtelbCreatePgelbReturnObjResponse, err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			var resp *ctelb.CtelbCreatePgelbResponse
			resp, err = c.meta.Apis.SdkCtElbApis.CtelbCreatePgelbApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.Description == "订单已取消或撤单" {
				err = fmt.Errorf("订单已取消或撤单, 请检查参数或避免并发创建")
				return false
			}

			status := resp.ReturnObj.MasterResourceStatus
			switch status {
			case business.ElbStatusStarted:
				loopResp = resp.ReturnObj
				return false
			case business.ElbStatusInProgress:
				// 仍在开通，继续轮询
				return true
			case business.ElbStatusStarting:
				// 开通正在启动中，继续轮询
				return true
			default:
				// 开通过程中，非started和ingress状态，其他都为异常状态，将跳出轮询
				err = fmt.Errorf("创建保障型负载均衡期间，存在异常返回状态。当前状态为：" + status)
				return false
			}
		},
	)
	if err != nil {
		return
	}

	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建成功！")
	}

	return
}

func (c *CtyunElbLoadBalancerResource) deleteLoop(ctx context.Context, params *ctelb.CtelbRefundPgelbRequest, loopCount ...int) (loopResp *ctelb.CtelbRefundPgelbReturnObjResponse, err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtElbApis.CtelbRefundPgelbApi.Do(ctx, c.meta.SdkCredential, params)
			if err != nil {
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			}
			status := resp.ReturnObj.MasterResourceStatus
			switch status {
			case business.ElbStatusRefunded:
				loopResp = resp.ReturnObj
				return false
			case business.ElbStatusInProgress:
				return true
			case business.ElbStatusStarted:
				return true
			case business.ElbStatusUnknown:
				return true
			default:
				// 退订过程中，非refunded,ingress和started状态，其他都为异常状态，将跳出轮询
				err = fmt.Errorf("退订保障型负载均衡期间，存在异常返回状态。当前状态为：" + status)
				return false
			}
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未退订成功！")
	}

	return
}

// modifyLoop 实现轮询查elb修改信息，确认异步接口修改成功
func (c *CtyunElbLoadBalancerResource) modifyLoop(ctx context.Context, state CtyunElbLoadBalancerConfig, plan CtyunElbLoadBalancerConfig, loopCount ...int) (isModify bool, err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)

	params := &ctelb.CtelbShowLoadBalancerRequest{
		RegionID: state.RegionID.ValueString(),
		ElbID:    state.ID.ValueString(),
	}
	// 轮询调用查询elb详情
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtElbApis.CtelbShowLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 查询得出变配slaName已于预期相同时，退出轮询
			if resp.ReturnObj[0].SlaName == plan.SlaName.ValueString() {
				isModify = true
				return false
			}

			return true
		})

	if result.ReturnReason == business.ReachMaxLoopTime {
		return false, errors.New("轮询已达最大次数，资源仍未创建成功！")
	}

	return
}

func (c *CtyunElbLoadBalancerResource) isContains(value string, collect []string) bool {
	for _, v := range collect {
		if v == value {
			return true
		}
	}
	return false
}

type CtyunElbLoadBalancerConfig struct {
	RegionID         types.String `tfsdk:"region_id"`          //区域ID
	ProjectID        types.String `tfsdk:"project_id"`         //企业项目 ID，默认为'0'
	VpcID            types.String `tfsdk:"vpc_id"`             //vpc的ID
	SubnetID         types.String `tfsdk:"subnet_id"`          //子网的ID
	Name             types.String `tfsdk:"name"`               //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description      types.String `tfsdk:"description"`        //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
	EipID            types.String `tfsdk:"eip_id"`             //弹性公网IP的ID。当resourceType=external为必填
	SlaName          types.String `tfsdk:"sla_name"`           //lb的规格名称,支持elb.s1.small和elb.default，默认为elb.default
	ResourceType     types.String `tfsdk:"resource_type"`      //资源类型。internal：内网负载均衡，external：公网负载均衡
	PrivateIpAddress types.String `tfsdk:"private_ip_address"` //负载均衡的私有IP地址，不指定则自动分配
	ID               types.String `tfsdk:"id"`                 //负载均衡ID
	AzName           types.String `tfsdk:"az_name"`            //可用区名称
	PortID           types.String `tfsdk:"port_id"`            //负载均衡实例默认创建port ID
	Ipv6Address      types.String `tfsdk:"ipv6_address"`       //负载均衡实例的IPv6地址
	EipInfo          types.List   `tfsdk:"eip_info"`           //弹性公网IP信息
	AdminStatus      types.String `tfsdk:"admin_status"`       //管理状态: DOWN / ACTIVE
	Status           types.String `tfsdk:"status"`             //负载均衡状态: DOWN / ACTIVE
	CreatedTime      types.String `tfsdk:"created_time"`       //创建时间，为UTC格式
	UpdatedTime      types.String `tfsdk:"updated_time"`       //更新时间，为UTC格式
	// 升级保障型负载均衡字段
	CycleType       types.String `tfsdk:"cycle_type"`        //订购类型：month（包月） / year（包年）
	CycleCount      types.Int64  `tfsdk:"cycle_count"`       //订购时长, 当 cycleType = month, 支持订购 1 - 11 个月; 当 cycleType = year, 支持订购 1 - 3 年
	PayVoucherPrice types.String `tfsdk:"pay_voucher_price"` //代金券金额，支持到小数点后两位
}

type EipInfoModel struct {
	EipID       types.String  `tfsdk:"eip_id"`
	Bandwidth   types.Float32 `tfsdk:"bandwidth"`
	IsTalkOrder types.Bool    `tfsdk:"is_talk_order"`
}
