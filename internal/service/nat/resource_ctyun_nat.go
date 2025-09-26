package nat

import (
	"context"
	"errors"
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
	_ resource.Resource                = &ctyunNat{}
	_ resource.ResourceWithConfigure   = &ctyunNat{}
	_ resource.ResourceWithImportState = &ctyunNat{}
)

type ctyunNat struct {
	meta *common.CtyunMetadata
}

func NewCtyunNatResource() resource.Resource {
	return &ctyunNat{}
}

func (c *ctyunNat) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat"
}

func (c *ctyunNat) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166493`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，值与nat_gateway_id相同",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
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
				Optional:    true,
				Computed:    true,
				Description: "需要创建 NAT 网关的 VPC 的 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"spec": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "规格 1~4, 1-表示小型, 2-表示中型, 3-表示大型, 4-表示超大型，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 4),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "nat名称，支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "nat描述，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:,'{},.,/;'[]·~！@#￥%……&*（） ——-+={}，支持更新",
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
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
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单id",
			},
			"nat_gateway_id": schema.StringAttribute{
				Computed:    true,
				Description: "网关id",
			},
			"vpc_name": schema.StringAttribute{
				Computed:    true,
				Description: "NAT所属的vpc专有网络名字",
			},
			"vpc_cidr": schema.StringAttribute{
				Computed:    true,
				Description: "当前网关所属的vpc cidr",
			},
			"creation_time": schema.StringAttribute{
				Computed:    true,
				Description: "NAT网关的创建时间",
			},
			"expired_time": schema.StringAttribute{
				Computed:    true,
				Description: "NAT网关实例的过期时间",
			},
		},
	}
}

// Create 创建nat
func (c *ctyunNat) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunNatConfig

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateNat(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	// NAT的创建，依赖于先有VPC
	returnObj, createParams, err := c.createNat(ctx, &plan)
	if err != nil {
		return
	}

	// 保存订单号
	masterOrderId := *returnObj.MasterOrderID
	plan.MasterOrderID = types.StringValue(masterOrderId)
	//response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	loopResponse, err := c.OrderLoop(ctx, createParams, 600)

	if err != nil {
		return
	} else if loopResponse == nil {
		err = common.InvalidReturnObjError
		return
	} else if loopResponse.MasterOrderId.ValueString() != masterOrderId {
		err = fmt.Errorf("创建nat时订单ID和轮询订单ID不一致，创建时订单ID：%s, 轮询所得订单ID：%s", masterOrderId, loopResponse.MasterOrderId)
	} else if loopResponse.RegionID.ValueString() != plan.RegionID.ValueString() {
		err = fmt.Errorf("创建nat时regionId和轮询结果regionId不一致，计划的regionId：%s, 轮询所得regionId：%s", plan.RegionID.ValueString(), loopResponse.RegionID.ValueString())
	}

	//plan.NatGatewayID = types.StringValue(loop.Uuid[0])
	if !loopResponse.NatGatewayId.IsNull() {
		plan.NatGatewayID = loopResponse.NatGatewayId
	}

	plan.ProjectID = utils.SecStringValue(createParams.ProjectID)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建后反查创建后的nat信息
	err = c.getAndMergeNat(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunNat) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunNatConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 通过订单号同步
	if !c.acquireAndSetIdIfOrderNotFinished(ctx, &state, response) {
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

func (c *ctyunNat) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunNatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunNatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新nat基础信息
	err = c.updateNatInfo(ctx, state, plan)
	if err != nil {
		return
	}
	//nat变配操作，规格(1-SMALL,2-MEDIUM,3-LARGE,4-XLARGE)的修改
	err = c.modifyNatSpec(ctx, state, plan)
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

func (c *ctyunNat) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunNatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteNatGatewayApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDeleteNatGatewayRequest{
		RegionID:     state.RegionID.ValueString(),
		NatGatewayID: state.NatGatewayID.ValueString(),
		ClientToken:  uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 根据返回值判断一下是否状态为退订状态(refunded)
	if len(*resp.ReturnObj.MasterResourceStatus) > 0 {
		// 若MasterResourceStatus不为空
		if !(*resp.ReturnObj.MasterResourceStatus == business.NatStatusRefunded) {
			err = fmt.Errorf("NatGatewayID :%s delete failed, MasterResourceStatus: %s", state.NatGatewayID, *resp.ReturnObj.MasterResourceStatus)
		}
	}
	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	err = helper.RefundLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID)
	if err != nil {
		return
	}

}

func (c *ctyunNat) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunNat) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunNatConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	config.RegionID = types.StringValue(regionID)
	config.NatGatewayID = types.StringValue(id)
	err = c.getAndMergeNat(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

// 在创建nat实例之前，进行检查
func (c *ctyunNat) checkBeforeCreateNat(ctx context.Context, plan CtyunNatConfig) error {
	return nil
}

func (c *ctyunNat) createNat(ctx context.Context, plan *CtyunNatConfig) (returnObj ctvpc.CtvpcCreateNatGatewayReturnObjResponse, createParams *ctvpc.CtvpcCreateNatGatewayRequest, err error) {
	regionID := plan.RegionID.ValueString()
	vpcID := plan.VpcID.ValueString()
	spec := plan.Spec.ValueInt32()
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	cycleType := plan.CycleType.ValueString()
	cycleCount := int32(plan.CycleCount.ValueInt64())
	azName := plan.AzName.ValueString()
	payVoucherPrice := plan.PayVoucherPrice.ValueString()
	projectID := plan.ProjectID.ValueString()

	params := &ctvpc.CtvpcCreateNatGatewayRequest{
		RegionID:    regionID,
		VpcID:       vpcID,
		Spec:        spec,
		Name:        name,
		Description: &description,
		ClientToken: uuid.NewString(),
		CycleType:   cycleType,
		AzName:      azName,
		ProjectID:   &projectID,
	}
	//if cycleType == business.OrderCycleTypeOnDemand {
	//	params.CycleCount = 1
	//	plan.CycleCount = types.Int64Value(1)
	//}
	if cycleType != business.OrderCycleTypeOnDemand && cycleCount > 0 {
		params.CycleCount = &cycleCount
	}
	if payVoucherPrice != "" {
		params.PayVoucherPrice = &payVoucherPrice
	}
	if projectID != "" {
		params.ProjectID = &projectID
	}

	// 调用创建接口
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateNatGatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	returnObj = *resp.ReturnObj
	createParams = params
	return
}

func (c *ctyunNat) getAndMergeNat(ctx context.Context, cfg *CtyunNatConfig) (err error) {
	cfg.ID = cfg.NatGatewayID
	//查看nat详情： ctvpc_show_nat_gateway_api.go
	params := &ctvpc.CtvpcShowNatGatewayRequest{
		RegionID:     cfg.RegionID.ValueString(),
		NatGatewayID: cfg.NatGatewayID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowNatGatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析resp.ReturnObj,将最新的nat信息同步到config中
	natObj := resp.ReturnObj
	cfg.VpcID = utils.SecStringValue(natObj.VpcID)
	//spec := utils.SecStringValue(natObj.Specs)
	spec := c.parseNatSpec(*natObj.Specs)
	if spec == 0 {
		err = errors.New("nat spec 返回值有误！当前值为：" + *natObj.Specs)
		return
	}
	cfg.Spec = types.Int32Value(spec)
	cfg.Name = utils.SecStringValue(natObj.Name)
	cfg.NatGatewayID = utils.SecStringValue(natObj.NatGatewayID)
	cfg.Description = utils.SecStringValue(natObj.Description)
	cfg.VpcName = utils.SecStringValue(natObj.VpcName)
	cfg.VpcCidr = utils.SecStringValue(natObj.VpcCidr)
	cfg.CreationTime = utils.SecStringValue(natObj.CreationTime)
	cfg.ExpiredTime = utils.SecStringValue(natObj.ExpiredTime)

	return nil
}

func (c *ctyunNat) acquireAndSetIdIfOrderNotFinished(ctx context.Context, state *CtyunNatConfig, response *resource.ReadResponse) bool {
	natGatewayId := state.NatGatewayID.ValueString()
	masterOrderId := state.MasterOrderID.ValueString()
	if natGatewayId != "" {
		return true
	}
	if masterOrderId == "" {
		// 没有受理的订单id，数据不可恢复，直接移除当前状态并返回
		response.State.RemoveResource(ctx)
		return false
	}

	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	resp, err := helper.OrderLoop(ctx, c.meta.Credential, masterOrderId)
	if err != nil || len(resp.Uuid) == 0 {
		// 报错，或受理没有返回数据的情况，表示单子未开通出来，且数据无法恢复
		response.State.RemoveResource(ctx)
		return false
	}
	//若成功的话，取出id
	state.NatGatewayID = types.StringValue(resp.Uuid[0])
	response.State.Set(ctx, state)
	return true
}
func (c *ctyunNat) modifyNatSpec(ctx context.Context, state CtyunNatConfig, plan CtyunNatConfig) (err error) {
	if plan.Spec.Equal(state.Spec) {
		if !state.PayVoucherPrice.Equal(plan.PayVoucherPrice) {
			err = fmt.Errorf("当没有触发变配时，代金券金额修改无效")
		}
		return
	}
	// 调用变配nat接口，规格(可传值：1-SMALL,2-MEDIUM,3-LARGE,4-XLARGE)
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcModifyNatSpecApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcModifyNatSpecRequest{
		RegionID:        state.RegionID.ValueString(),
		NatGatewayID:    state.NatGatewayID.ValueString(),
		Spec:            plan.Spec.ValueInt32(),
		ClientToken:     uuid.NewString(),
		PayVoucherPrice: plan.PayVoucherPrice.ValueStringPointer(),
	})
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	helper := business.NewOrderLooper(c.meta.Apis.CtEcsApis.EcsOrderQueryUuidApi)
	_, err = helper.OrderLoop(ctx, c.meta.Credential, *resp.ReturnObj.MasterOrderID, 600)
	if err != nil {
		return
	}

	return nil
}

func (c *ctyunNat) updateNatInfo(ctx context.Context, state CtyunNatConfig, plan CtyunNatConfig) (err error) {
	if plan.Name.Equal(state.Name) && plan.Description.Equal(state.Description) {
		return
	}
	params := &ctvpc.CtvpcUpdateNatGatewayAttributeRequest{
		RegionID:     plan.RegionID.ValueString(),
		NatGatewayID: state.NatGatewayID.ValueString(),
		Name:         plan.Name.ValueStringPointer(),
		Description:  plan.Description.ValueStringPointer(),
		ClientToken:  uuid.NewString(),
	}
	resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateNatGatewayAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	// 轮询详情接口，确认是否修改
	err = c.updateLoop(ctx, state, params, 30)
	if err != nil {
		return err
	}
	return
}

// 循环查询nat信息，确保更新成功
func (c *ctyunNat) updateLoop(ctx context.Context, state CtyunNatConfig, updatedParams *ctvpc.CtvpcUpdateNatGatewayAttributeRequest, loopCount ...int) (err error) {
	count := 10
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			params := &ctvpc.CtvpcShowNatGatewayRequest{
				RegionID:     state.RegionID.ValueString(),
				NatGatewayID: state.NatGatewayID.ValueString(),
			}
			resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcShowNatGatewayApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			if *resp.ReturnObj.Name == *updatedParams.Name && *resp.ReturnObj.Description == *updatedParams.Description {
				return false
			} else {
				return true
			}
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源信息仍未更新成功！")
	}
	return
}

func (c *ctyunNat) OrderLoop(ctx context.Context, params *ctvpc.CtvpcCreateNatGatewayRequest, loopCount ...int) (loopResponse *LoopOrderResponse, err error) {

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
			resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcCreateNatGatewayApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			} else if *resp.Description == "订单已取消或撤单" {
				err = fmt.Errorf("订单已取消或撤单, 请检查参数或避免并发创建")
				return false
			}

			status := *resp.ReturnObj.MasterResourceStatus
			switch status {
			case business.NatStatusStarted:
				// nat 已经started，跳出轮询，返回natGatewayID
				natGatewayId := utils.SecStringValue(resp.ReturnObj.NatGatewayID)
				masterOrderId := utils.SecStringValue(resp.ReturnObj.MasterOrderID)
				masterOrderNo := utils.SecStringValue(resp.ReturnObj.MasterOrderNO)
				masterResourceId := utils.SecStringValue(resp.ReturnObj.MasterResourceID)
				regionId := utils.SecStringValue(resp.ReturnObj.RegionID)
				loopResponse = &LoopOrderResponse{
					NatGatewayId:         natGatewayId,
					MasterOrderId:        masterOrderId,
					MasterOrderNO:        masterOrderNo,
					MasterResourceID:     masterResourceId,
					MasterResourceStatus: types.StringValue(status),
					RegionID:             regionId,
				}
				return false
			case business.NatStatusStarting:
				// 仍在开通，继续轮询
				return true
			case business.NatStatusInProgress:
				// 仍在开通，继续轮询
				return true
			default:
				// 在开通的时候，其他状态是异常的，因此抛出异常，并跳出轮询
				err = errors.New("Nat开通时，出现非starting 和 started的异常状态。当前状态为： " + status)
				return false
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建成功！")
	}

	return
}

func (c *ctyunNat) parseNatSpec(spec string) (specInt int32) {
	switch spec {
	case "small":
		specInt = 1
	case "medium":
		specInt = 2
	case "large":
		specInt = 3
	case "xlarge":
		specInt = 4
	default:
		specInt = 0
	}
	return specInt
}

type CtyunNatConfig struct {
	ID              types.String `tfsdk:"id"`
	RegionID        types.String `tfsdk:"region_id"`         //区域id
	VpcID           types.String `tfsdk:"vpc_id"`            //需要创建 NAT 网关的 VPC 的 ID
	Spec            types.Int32  `tfsdk:"spec"`              //规格 1~4, 1表示小型, 2表示中型, 3表示大型, 4表示超大型
	Name            types.String `tfsdk:"name"`              //支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description     types.String `tfsdk:"description"`       //支持拉丁字母、中文、数p字, 特殊字符：~!@#$%^&*()_-+= <>?:,'{},.,/;'[]·~！@#￥%……&*（） ——-+={}
	CycleType       types.String `tfsdk:"cycle_type"`        //订购类型：month（包月） / year（包年）/ on_demand（按需）
	CycleCount      types.Int64  `tfsdk:"cycle_count"`       //订购时长, 当 cycleType = month, 支持订购 1 - 11 个月; 当 cycleType = year, 支持订购 1 - 3 年
	AzName          types.String `tfsdk:"az_name"`           //可用区名称
	PayVoucherPrice types.String `tfsdk:"pay_voucher_price"` //代金券金额，支持到小数点后两位
	ProjectID       types.String `tfsdk:"project_id"`        //企业项目，不传默认为 0
	MasterOrderID   types.String `tfsdk:"master_order_id"`   //订单id
	NatGatewayID    types.String `tfsdk:"nat_gateway_id"`    //网关 ID
	VpcName         types.String `tfsdk:"vpc_name"`          //NAT所属的专有网络名字
	VpcCidr         types.String `tfsdk:"vpc_cidr"`          //当前网关所属的vpc cidr
	CreationTime    types.String `tfsdk:"creation_time"`     //NAT网关的创建时间
	ExpiredTime     types.String `tfsdk:"expired_time"`      //NAT网关实例的过期时间
}
type LoopOrderResponse struct {
	NatGatewayId         types.String
	MasterOrderId        types.String // 主订单id
	MasterOrderNO        types.String
	MasterResourceStatus types.String
	MasterResourceID     types.String
	RegionID             types.String
}
