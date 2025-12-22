package nat

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctnat"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &ctyunPrivateSnatResource{}
	_ resource.ResourceWithConfigure   = &ctyunPrivateSnatResource{}
	_ resource.ResourceWithImportState = &ctyunPrivateSnatResource{}
)

type ctyunPrivateSnatResource struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateSnatResource() resource.Resource {
	return &ctyunPrivateSnatResource{}
}

func (c *ctyunPrivateSnatResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[natGateWayID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrivateSnatConfig
	var ID, projectID, regionID, natGateWayID string
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		projectID = c.meta.GetExtraIfEmpty(projectID, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &natGateWayID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &natGateWayID, &projectID, &regionID)
		if err != nil {
			return
		}
	}
	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if natGateWayID == "" {
		err = fmt.Errorf("natGateWayID不能为空")
	}
	config.ID = types.StringValue(ID)
	config.SNatID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionID)
	config.NatGatewayID = types.StringValue(natGateWayID)
	err = c.getAndMergeSnat(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *ctyunPrivateSnatResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_nat_snat"
}

func (c *ctyunPrivateSnatResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026759/10378381`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID，同snat_id",
			},
			"snat_id": schema.StringAttribute{
				Computed:      true,
				Description:   "Snat规则的id",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
				Required:    true,
				Description: "子网ID，需要和NAT网关同属一个VPC，支持更新",
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"addresses": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "中转IP地址，必须在中转网段指定的网络范围内 支持更新",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SNAT描述 支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
					validator2.DescNotStartWithHttp()},
			},
			"source_vpc_name": schema.StringAttribute{
				Computed:      true,
				Description:   "源vpc名称",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			//"source_subnet_name": schema.StringAttribute{
			//	Computed:      true,
			//	Description:   "源Subnet名称",
			//	PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			//},
			"state": schema.StringAttribute{
				Computed:      true,
				Description:   "SNAT状态: running代表运行中, freeze代表已冻结, expired代表已到期",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (c *ctyunPrivateSnatResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPrivateSnatConfig
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
	createResp, err := c.createSnat(ctx, plan)
	if err != nil {
		return
	}

	plan.SNatID = types.StringValue(createResp.ReturnObj.SnatID)
	plan.ID = plan.SNatID

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

func (c *ctyunPrivateSnatResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunPrivateSnatConfig
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

func (c *ctyunPrivateSnatResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunPrivateSnatConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunPrivateSnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新snat基础信息
	err = c.updateSnatInfo(ctx, state, plan)
	if err != nil {
		return
	}

	state.Addresses = plan.Addresses
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

func (c *ctyunPrivateSnatResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var state CtyunPrivateSnatConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err = c.meta.Apis.SdkCtNatApis.CtnatDeletePrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatDeletePrivatenatSnatRequest{
		RegionID: state.RegionID.ValueString(),
		SnatID:   state.SNatID.ValueString(),
	})
	if err != nil {
		return
	}

	// 轮询查询直到删除成功
	err = c.DeleteLoop(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunPrivateSnatResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunPrivateSnatResource) getAndMergeSnat(ctx context.Context, config *CtyunPrivateSnatConfig) (err error) {
	config.ID = config.SNatID

	// 查询SNAT列表，因为没有单独的查询接口
	resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatQueryPrivatenatSnatRequest{
		RegionID:     config.RegionID.ValueString(),
		NatGatewayID: config.NatGatewayID.ValueString(),
		SnatID:       config.SNatID.ValueString(),
		PageNo:       1,
		PageSize:     10,
	})
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
		err = fmt.Errorf("snat not found")
		return
	}

	snat := resp.ReturnObj[0]
	config.SNatID = types.StringValue(snat.SnatID)
	config.Description = types.StringValue(snat.Description)
	config.SourceVpcName = types.StringValue(snat.SrcVpcName)
	//config.SourceSubnetName = types.StringValue(snat.SrcSubnetName)
	config.State = types.StringValue(snat.State)

	snatIps, diags := types.SetValueFrom(ctx, types.StringType, snat.Addresses)
	if diags.HasError() {
		err = fmt.Errorf("failed to convert addresses to list")
		return
	}
	config.Addresses = snatIps
	// 确保source_subnet_id保持不变
	config.SourceSubnetID = types.StringValue(snat.SrcSubnetID)

	return nil
}

func (c *ctyunPrivateSnatResource) checkBeforeCreateSnat(ctx context.Context, plan CtyunPrivateSnatConfig) (err error) {
	// 检查source_subnet_id是否为空
	if plan.SourceSubnetID.ValueString() == "" {
		err = fmt.Errorf("source_subnet_id不能为空")
		return
	}

	// 检查eip
	var snatIps []string
	diag := plan.Addresses.ElementsAs(ctx, &snatIps, false)
	if diag.HasError() {
		err = fmt.Errorf(diag.Errors()[0].Detail())
		return
	}

	// 检查私网nat网关是否存在
	_, err = business.NewPrivateNatService(c.meta).GetPrivateNatByID(ctx, plan.NatGatewayID.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		err = fmt.Errorf("校验私网nat网关失败: %s", err.Error())
		return
	}

	return
}

// createSnat创建Snat
func (c *ctyunPrivateSnatResource) createSnat(ctx context.Context, plan CtyunPrivateSnatConfig) (resp *ctnat.CtnatCreatePrivatenatSnatResponse, err error) {
	regionId := plan.RegionID.ValueString()
	natGatewayId := plan.NatGatewayID.ValueString()
	sourceSubnetId := plan.SourceSubnetID.ValueString()
	var snatIps []string
	plan.Addresses.ElementsAs(ctx, &snatIps, true)

	params := &ctnat.CtnatCreatePrivatenatSnatRequest{
		RegionID:       regionId,
		NatGatewayID:   natGatewayId,
		SourceSubnetID: sourceSubnetId,
		SnatIPs:        snatIps,
		Description:    plan.Description.ValueString(),
	}

	resp, err = c.meta.Apis.SdkCtNatApis.CtnatCreatePrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

func (c *ctyunPrivateSnatResource) updateSnatInfo(ctx context.Context, state CtyunPrivateSnatConfig, plan CtyunPrivateSnatConfig) (err error) {
	if plan.SourceSubnetID.Equal(state.SourceSubnetID) &&
		plan.Description.Equal(state.Description) &&
		plan.Addresses.Equal(state.Addresses) {
		return
	}

	var snatIps []string
	plan.Addresses.ElementsAs(ctx, &snatIps, true)

	params := &ctnat.CtnatModifyPrivatenatSnatRequest{
		RegionID:       state.RegionID.ValueString(),
		SnatID:         state.SNatID.ValueString(),
		SourceSubnetID: plan.SourceSubnetID.ValueString(),
		SnatIps:        snatIps,
		Description:    plan.Description.ValueString(),
	}

	_, err = c.meta.Apis.SdkCtNatApis.CtnatModifyPrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	}

	// 轮询请求查看是否更新成功
	err = c.updateLoop(ctx, &state, params)
	return
}

func (c *ctyunPrivateSnatResource) updateLoop(ctx context.Context, state *CtyunPrivateSnatConfig, updatedParams *ctnat.CtnatModifyPrivatenatSnatRequest, loopCount ...int) (err error) {
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
			resp, err2 := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatQueryPrivatenatSnatRequest{
				RegionID:     state.RegionID.ValueString(),
				NatGatewayID: state.NatGatewayID.ValueString(),
				SnatID:       state.SNatID.ValueString(),
				PageNo:       1,
				PageSize:     10,
			})
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return false
			} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
				err = common.InvalidReturnObjError
				return false
			}

			snatInfo := resp.ReturnObj[0]
			if updatedParams.SourceSubnetID != "" && updatedParams.SourceSubnetID != snatInfo.SrcSubnetID {
				return true
			}
			if updatedParams.Description != "" && updatedParams.Description != snatInfo.Description {
				return true
			}
			return false
		})

	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更新完成!Snat: " + state.SNatID.ValueString())
	}
	return
}

func (c *ctyunPrivateSnatResource) DeleteLoop(ctx context.Context, state CtyunPrivateSnatConfig) (err error) {
	var respErr error
	retryer, err := business.NewRetryer(time.Second*5, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err := c.meta.Apis.SdkCtNatApis.CtnatQueryPrivatenatSnatApi.Do(ctx, c.meta.SdkCredential, &ctnat.CtnatQueryPrivatenatSnatRequest{
				RegionID:     state.RegionID.ValueString(),
				NatGatewayID: state.NatGatewayID.ValueString(),
				SnatID:       state.SNatID.ValueString(),
				PageNo:       1,
				PageSize:     10,
			})
			if err != nil {
				respErr = err
				return false
			} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
				// SNAT已删除
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

type CtyunPrivateSnatConfig struct {
	ID             types.String `tfsdk:"id"`
	RegionID       types.String `tfsdk:"region_id"`        //区域id
	NatGatewayID   types.String `tfsdk:"nat_gateway_id"`   //NAT网关ID
	SourceSubnetID types.String `tfsdk:"source_subnet_id"` //子网id
	Addresses      types.Set    `tfsdk:"addresses"`        //IP地址，必须在中转网段指定的网络范围内
	Description    types.String `tfsdk:"description"`      //支持拉丁字母、中文、数字, 特殊字符
	SNatID         types.String `tfsdk:"snat_id"`          //snat id
	SourceVpcName  types.String `tfsdk:"source_vpc_name"`  //源vpc名称
	//SourceSubnetName types.String `tfsdk:"source_subnet_name"` //源Subnet名称

	State types.String `tfsdk:"state"` //SNAT状态
}
