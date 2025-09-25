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
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunSnatResource{}
	_ resource.ResourceWithConfigure   = &ctyunSnatResource{}
	_ resource.ResourceWithImportState = &ctyunSnatResource{}
)

type ctyunSnatResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunSnatResource() resource.Resource {
	return &ctyunSnatResource{}
}

func (c *ctyunSnatResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunSnatConfig
	var id, regionID string
	err = terraform_extend.Split(request.ID, &id, &regionID)
	if err != nil {
		return
	}
	config.RegionID = types.StringValue(regionID)
	config.SNatID = types.StringValue(id)
	err = c.getAndMergeSnat(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *ctyunSnatResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nat_snat"
}

func (c *ctyunSnatResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10166496`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，同snat_id",
			},
			"snat_id": schema.StringAttribute{
				Computed:    true,
				Description: "Snat规则的id",
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
			"nat_gateway_id": schema.StringAttribute{
				Required:    true,
				Description: "NAT网关Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"source_subnet_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "子网ID，需要和NAT网关同属一个VPC，与source_cidr有且只能填写一个，支持更新",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("source_cidr")),
				},
			},
			"source_cidr": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "自定义网段，与source_subnet_id有且只能填写一个，支持更新",
				Validators: []validator.String{
					validator2.Cidr(),
					stringvalidator.ConflictsWith(path.MatchRoot("source_subnet_id")),
				},
			},
			"snat_ips": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "弹性公网IP集合，每个元素为eipID，至少输入1个，最多5个，支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(5),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SNAT描述，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"subnet_type": schema.Int32Attribute{
				Computed:    true,
				Description: "子网类型：1-使用子网ID，2-使用自定义网段",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间",
			},
			"eips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "绑定的 eip 信息",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"eip_id": schema.StringAttribute{
							Computed:    true,
							Description: "弹性 IP id",
						},
						"ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "弹性 IP 地址",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunSnatResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunSnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前检查
	err = c.checkBeforeCreateSnat(ctx, plan)
	if err != nil {
		return
	}

	// 创建
	id, err := c.createSnat(ctx, plan)
	if err != nil {
		return
	}
	plan.SNatID = types.StringValue(id)
	time.Sleep(5 * time.Second)
	// 添加description
	err = c.updateDescription(ctx, plan)
	if err != nil {
		return
	}

	// 反查信息
	err = c.getAndMergeSnat(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *ctyunSnatResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunSnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 远端查询
	err = c.getAndMergeSnat(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(request.State.Set(ctx, &state)...)
}

func (c *ctyunSnatResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunSnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunSnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新snat基础信息
	err = c.updateSnatInfo(ctx, state, plan)
	if err != nil {
		return
	}

	err = c.addAndDeleteSnatEips(ctx, state, plan)
	if err != nil {
		return
	}
	state.SnatIps = plan.SnatIps
	// 查询远端信息并更新本地
	err = c.getAndMergeSnat(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSnatResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunSnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteSnatEntryApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDeleteSnatEntryRequest{
		RegionID:    state.RegionID.ValueString(),
		SNatID:      state.SNatID.ValueString(),
		ClientToken: uuid.NewString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	err = c.DeleteLoop(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunSnatResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunSnatResource) getAndMergeSnat(ctx context.Context, config *CtyunSnatConfig) (err error) {
	config.ID = config.SNatID
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowSnatApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowSnatRequest{
		RegionID: config.RegionID.ValueString(),
		SNatID:   config.SNatID.ValueString(),
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

	snat := resp.ReturnObj
	config.SNatID = utils.SecStringValue(snat.SNatID)
	config.Description = utils.SecStringValue(snat.Description)
	config.SourceCIDR = utils.SecStringValue(snat.SubnetCidr)

	if c.inRange(business.SNatSubnetTypes, snat.SubnetType) {
		config.SubnetType = types.Int32Value(snat.SubnetType)
	} else {
		err = fmt.Errorf("the value of SubnetType %d is invaild ", snat.SubnetType)
		return
	}
	config.CreateTime = utils.SecStringValue(snat.CreationTime)
	config.SourceSubnetID = utils.SecStringValue(snat.SubnetID)

	var eipList []CtyunSnatEipsList
	var eipIds []types.String
	for _, eip := range snat.Eips {
		eipIds = append(eipIds, utils.SecStringValue(eip.EipID))
		eipItem := CtyunSnatEipsList{
			EipID:     utils.SecStringValue(eip.EipID),
			IpAddress: utils.SecStringValue(eip.IpAddress),
		}
		eipList = append(eipList, eipItem)
	}
	config.SnatIps, _ = types.SetValueFrom(ctx, types.StringType, eipIds)

	eipObj := utils.StructToTFObjectTypes(CtyunSnatEipsList{})

	config.Eips, _ = types.ListValueFrom(ctx, eipObj, eipList)
	return nil
}

func (c *ctyunSnatResource) checkBeforeCreateSnat(ctx context.Context, plan CtyunSnatConfig) (err error) {
	if plan.SourceCIDR.ValueString() == "" && plan.SourceSubnetID.ValueString() == "" {
		err = fmt.Errorf("子网ID和自定义网段不能都为空")
		return
	}
	nat, err := business.NewNatService(c.meta).GetNatByID(ctx, plan.NatGatewayID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		err = fmt.Errorf("校验nat网关失败" + err.Error())
		return
	}
	subnets, err := business.NewVpcService(c.meta).GetVpcSubnet(ctx, *nat.VpcID, plan.RegionID.ValueString(), "")
	if err != nil {
		return err
	}

	// 传了子网，但不在该vpc内
	if plan.SourceSubnetID.ValueString() != "" && subnets[plan.SourceSubnetID.ValueString()].SubnetId == "" {
		err = fmt.Errorf("子网 %s 不属于 %s 所在的 %s", plan.SourceSubnetID.ValueString(), plan.NatGatewayID.ValueString(), *nat.VpcID)
	}
	// 检查eip
	var snatIps []string
	diag := plan.SnatIps.ElementsAs(ctx, &snatIps, false)
	if diag.HasError() {
		err = fmt.Errorf(diag.Errors()[0].Detail())
		return
	}
	for _, eipID := range snatIps {
		err = business.NewEipService(c.meta).MustExist(ctx, eipID, plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}
	return
}

// createSnat创建Snat
func (c *ctyunSnatResource) createSnat(ctx context.Context, plan CtyunSnatConfig) (id string, err error) {
	regionId := plan.RegionID.ValueString()
	natGatewayId := plan.NatGatewayID.ValueString()
	sourceSubnetId := plan.SourceSubnetID.ValueString()
	sourceCIDR := plan.SourceCIDR.ValueString()
	var snatIps []string
	plan.SnatIps.ElementsAs(ctx, &snatIps, true)
	params := &ctvpc.CtvpcCreateSnatEntryRequest{
		RegionID:     regionId,
		NatGatewayID: natGatewayId,
		SnatIps:      snatIps,
		ClientToken:  uuid.NewString(),
	}
	if sourceSubnetId != "" {
		params.SourceSubnetID = &sourceSubnetId
	}
	if sourceCIDR != "" {
		params.SourceCIDR = &sourceCIDR
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateSnatEntryApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Snat == nil || resp.ReturnObj.Snat.SnatID == nil {
		err = common.InvalidReturnObjError
		return
	}
	id = *resp.ReturnObj.Snat.SnatID
	return
}

func (c *ctyunSnatResource) updateSnatInfo(ctx context.Context, state CtyunSnatConfig, plan CtyunSnatConfig) (err error) {
	if plan.SourceSubnetID.Equal(state.SourceSubnetID) &&
		plan.SourceCIDR.Equal(state.SourceCIDR) &&
		plan.Description.Equal(state.Description) {
		return
	}

	params := &ctvpc.CtvpcUpdateSnatEntryAttributeRequest{
		RegionID:    state.RegionID.ValueString(),
		SNatID:      state.SNatID.ValueString(),
		ClientToken: uuid.NewString(),
		Description: plan.Description.ValueStringPointer(),
	}
	if plan.SourceSubnetID.ValueString() != "" {
		params.SourceSubnetID = plan.SourceSubnetID.ValueStringPointer()
	}
	if plan.SourceCIDR.ValueString() != "" {
		params.SourceCIDR = plan.SourceCIDR.ValueStringPointer()
	}

	resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateSnatEntryAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err2 != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	// 轮询请求查看是否更新成功
	err = c.updateLoop(ctx, &state, params)

	return
}

func (c *ctyunSnatResource) addAndDeleteSnatEips(ctx context.Context, state CtyunSnatConfig, plan CtyunSnatConfig) (err error) {

	var planSnatIps []string
	var stateSnatIps []string
	diags := plan.SnatIps.ElementsAs(ctx, &planSnatIps, true)
	if diags.HasError() {
		return
	}
	diags = state.SnatIps.ElementsAs(ctx, &stateSnatIps, true)
	if diags.HasError() {
		return
	}

	// 先筛选出哪些是需要删除的，哪些是需要添加的Eip
	addIpAddressIds, deleteIpAddressIds := utils.DifferenceStrArray(planSnatIps, stateSnatIps)

	//snat添加eip
	if len(addIpAddressIds) > 0 {
		addResp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcAssociateEipsToSnatApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcAssociateEipsToSnatRequest{
			RegionID:     state.RegionID.ValueString(),
			NatGatewayID: state.NatGatewayID.ValueString(),
			SNatID:       state.SNatID.ValueString(),
			IpAddressIds: addIpAddressIds,
			ClientToken:  uuid.NewString(),
		})
		if err2 != nil {
			err = err2
			return
		} else if addResp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *addResp.Message, *addResp.Description)
			return
		}
	}
	if len(deleteIpAddressIds) > 0 {
		deleteResp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcDisassociateEipsFromSnatApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDisassociateEipsFromSnatRequest{
			RegionID:     state.RegionID.ValueString(),
			NatGatewayID: state.NatGatewayID.ValueString(),
			SNatID:       state.SNatID.ValueString(),
			IpAddressIds: deleteIpAddressIds,
			ClientToken:  uuid.NewString(),
		})
		if err2 != nil {
			err = err2
			return
		} else if deleteResp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *deleteResp.Message, *deleteResp.Description)
			return
		}
	}
	time.Sleep(5 * time.Second)
	return nil
}

func (c *ctyunSnatResource) updateLoop(ctx context.Context, state *CtyunSnatConfig, updatedParams *ctvpc.CtvpcUpdateSnatEntryAttributeRequest, loopCount ...int) (err error) {
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
			resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcShowSnatApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowSnatRequest{
				RegionID: state.RegionID.ValueString(),
				SNatID:   state.SNatID.ValueString(),
			})
			if err != nil {
				err = err2
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			snatInfo := resp.ReturnObj
			if updatedParams.SourceSubnetID != nil && *updatedParams.SourceSubnetID != *snatInfo.SubnetID {
				return true
			}
			if updatedParams.SourceCIDR != nil && *updatedParams.SourceCIDR != *snatInfo.SubnetCidr {
				return true
			}
			if updatedParams.Description != nil && *updatedParams.Description != *snatInfo.Description {
				return true
			}
			return false
		})

	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新完成!Snat: " + state.SNatID.ValueString())
	}
	return
}

func (c *ctyunSnatResource) DeleteLoop(ctx context.Context, state CtyunSnatConfig) (err error) {
	var respErr error
	retryer, err := business.NewRetryer(time.Second*5, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowSnatApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcShowSnatRequest{
				RegionID: state.RegionID.ValueString(),
				SNatID:   state.SNatID.ValueString(),
			})
			if err != nil {
				respErr = err
				return false
			} else if resp.ReturnObj == nil {
				return false
			} else {
				//如果仍能查询到snat信息，说明未删除完成
				return true
			}
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未删除!Snat: " + state.SNatID.ValueString())
	}
	return respErr
}

func (c *ctyunSnatResource) inRange(ranges []int32, num int32) bool {
	for _, v := range ranges {
		if v == num {
			return true
		}
	}
	return false
}

// 因为create 接口没有description参数，通过创建后，update接口为description赋值
func (c *ctyunSnatResource) updateDescription(ctx context.Context, plan CtyunSnatConfig) (err error) {
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		return
	}
	if plan.SNatID.IsNull() || plan.SNatID.IsUnknown() {
		err = errors.New("snat未创建成功，snat id缺失，导致无法添加description")
		return
	}
	params := &ctvpc.CtvpcUpdateSnatEntryAttributeRequest{
		RegionID:    plan.RegionID.ValueString(),
		SNatID:      plan.SNatID.ValueString(),
		ClientToken: uuid.NewString(),
		Description: plan.Description.ValueStringPointer(),
	}
	if !plan.SourceSubnetID.IsNull() && !plan.SourceSubnetID.IsUnknown() {
		params.SourceSubnetID = plan.SourceSubnetID.ValueStringPointer()
	} else if !plan.SourceCIDR.IsNull() && !plan.SourceCIDR.IsUnknown() {
		params.SourceCIDR = plan.SourceCIDR.ValueStringPointer()
	} else {
		err = errors.New("source_subnet_id 和source_cidr同时只能填写一个")
		return err
	}

	resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateSnatEntryAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err2 != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	// 轮询请求查看是否更新成功
	err = c.updateLoop(ctx, &plan, params)

	return
}

type CtyunSnatConfig struct {
	ID             types.String `tfsdk:"id"`
	RegionID       types.String `tfsdk:"region_id"`        //区域id
	NatGatewayID   types.String `tfsdk:"nat_gateway_id"`   //NAT网关ID
	SourceSubnetID types.String `tfsdk:"source_subnet_id"` //子网id，【非自定义情况必传 sourceCIDR和sourceSubnetID二选一必传】｜ 5fe30709-93ef-522f-a1a0-d8c8f6803e0d
	SourceCIDR     types.String `tfsdk:"source_cidr"`      //自定义输入VPC、交换机或ECS实例的网段，还可以输入任意网段。【自定义子网信息必传】】
	SnatIps        types.Set    `tfsdk:"snat_ips"`         //弹性公网IP集合
	Description    types.String `tfsdk:"description"`      //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&()_-+= <>?:"{},./;'[]·~！@#￥%……&（） ——-+={}
	SNatID         types.String `tfsdk:"snat_id"`          //snat id
	SubnetType     types.Int32  `tfsdk:"subnet_type"`      //子网类型：1-有vpcID的子网，0-自定义
	CreateTime     types.String `tfsdk:"create_time"`      //创建时间
	Eips           types.List   `tfsdk:"eips"`             //绑定的 eip 信息
}

type CtyunSnatEipsList struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性 IP id
	IpAddress types.String `tfsdk:"ip_address"` //弹性 IP 地址
}

type SNatLoopCreateResponse struct {
	SNatID types.String
	Status types.String
}
