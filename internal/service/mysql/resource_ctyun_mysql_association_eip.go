package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
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
	_ resource.Resource                = &CtyunMysqlAssociationEip{}
	_ resource.ResourceWithConfigure   = &CtyunMysqlAssociationEip{}
	_ resource.ResourceWithImportState = &CtyunMysqlAssociationEip{}
)

type CtyunMysqlAssociationEip struct {
	meta       *common.CtyunMetadata
	eipService *business.EipService
}

func (c *CtyunMysqlAssociationEip) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_association_eip"
}
func NewCtyunMysqlAssociationEip() resource.Resource {
	return &CtyunMysqlAssociationEip{}
}

func (c *CtyunMysqlAssociationEip) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10033813/10033927**`,
		Attributes: map[string]schema.Attribute{
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性IP的id",
				Validators: []validator.String{
					validator2.EipValidate(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"eip_status": schema.Int32Attribute{
				Computed:    true,
				Description: " 弹性ip状态 0->unbind，1->bind,2->binding",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "eip绑定状态，与eip_status一致",
			},
		},
	}
}

func (c *CtyunMysqlAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 绑定eip之前，轮询确认eip未被绑定
	err = c.preCheckEipUnbound(ctx, &plan)
	if err != nil {
		return
	}
	// 实例绑定弹性IP
	err = c.MysqlBindEip(ctx, &plan)
	if err != nil {
		return
	}
	// 轮询查看绑定状态
	err = c.BindLoop(ctx, &plan, business.EipStatusBind, business.MysqlBindEipStatusACTIVE)
	if err != nil {
		return
	}
	// 查询实例详情，确认是否绑定成功
	err = c.getAndMergeBindEip(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunMysqlAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunAssociationEipConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeBindEip(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "is not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunMysqlAssociationEip) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	//暂无可更新内容
}

func (c *CtyunMysqlAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	eip, err := c.eipService.GetEipAddressByEipID(ctx, state.EipID.ValueString(), state.RegionID.ValueString())
	if err != nil {
		return
	}

	unbindParams := &mysql.TeledbUnbindEipRequest{
		EipID:  state.EipID.ValueString(),
		Eip:    *eip.EipAddress,
		InstID: state.InstID.ValueString(),
	}
	unbindHeader := &mysql.TeledbUnbindEipRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		unbindHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbUnbindEipApi.Do(ctx, c.meta.Credential, unbindParams, unbindHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确定解绑成功
	err = c.BindLoop(ctx, &state, business.EipStatusUnbind, business.MysqlBindEipStatusDOWN)
	if err != nil {
		return
	}
	return
}

func (c *CtyunMysqlAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// todo
}

func (c *CtyunMysqlAssociationEip) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.eipService = business.NewEipService(c.meta)
}

func (c *CtyunMysqlAssociationEip) MysqlBindEip(ctx context.Context, config *CtyunAssociationEipConfig) (err error) {
	// 绑定前，确定实例状态为running
	err = c.StartedLoop(ctx, config, 60)
	if err != nil {
		return err
	}

	params := &mysql.TeledbBindEipRequest{
		EipID:  config.EipID.ValueString(),
		Eip:    config.eipAddress,
		InstID: config.InstID.ValueString(),
	}
	header := &mysql.TeledbBindEipRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbBindEipApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	}
	return
}

func (c *CtyunMysqlAssociationEip) getAndMergeBindEip(ctx context.Context, config *CtyunAssociationEipConfig) (err error) {
	//detailParams := &mysql.TeledbQueryDetailRequest{
	//	OuterProdInstId: config.InstID.ValueString(),
	//}
	//header := &mysql.TeledbQueryDetailRequestHeaders{
	//	InstID:   config.InstID.ValueString(),
	//	RegionID: config.RegionID.ValueString(),
	//}
	//if config.ProjectID.ValueString() != "" {
	//	header.ProjectID = config.ProjectID.ValueStringPointer()
	//}
	//
	//resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, header)
	//if err != nil {
	//	return err
	//} else if resp.StatusCode != 0 {
	//	err = fmt.Errorf("API return error. Message: %s", resp.Message)
	//	return
	//} else if resp.ReturnObj == nil {
	//	err = common.InvalidReturnObjError
	//	return
	//}
	//returnObj := resp.ReturnObj

	params := &mysql.TeledbBoundEipListRequest{
		RegionID: config.RegionID.ValueString(),
		EipID:    config.EipID.ValueStringPointer(),
	}

	headers := &mysql.TeledbBoundEipListRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		headers.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbBoundEipListApi.Do(ctx, c.meta.Credential, params, headers)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	eipInfo := resp.ReturnObj.Data[0]
	config.EipStatus = types.Int32Value(eipInfo.BindStatus)
	config.Status = types.StringValue(eipInfo.Status)
	return
}

func (c *CtyunMysqlAssociationEip) BindLoop(ctx context.Context, config *CtyunAssociationEipConfig, finalBindStatus int32, finalStatus string, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			params := &mysql.TeledbBoundEipListRequest{
				RegionID: config.RegionID.ValueString(),
				EipID:    config.EipID.ValueStringPointer(),
			}

			headers := &mysql.TeledbBoundEipListRequestHeader{}
			if config.ProjectID.ValueString() != "" {
				headers.ProjectID = config.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbBoundEipListApi.Do(ctx, c.meta.Credential, params, headers)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 200 {
				err = fmt.Errorf("API return error. Message: %s ", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			// 解析返回的绑定eip列表
			returnObj := resp.ReturnObj.Data
			if len(returnObj) > 1 {
				err = errors.New("根据eip id 查询到多个eip详情，返回有误！")
				return false
			} else if len(returnObj) < 1 {
				err = errors.New("eip 有误，未查询到该eip绑定信息")
				return false
			}
			if returnObj[0].Status == finalStatus && returnObj[0].BindStatus == finalBindStatus {
				return false
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，eip仍未绑定/解绑成功！")
	}
	return
}

func (c *CtyunMysqlAssociationEip) StartedLoop(ctx context.Context, state *CtyunAssociationEipConfig, loopCount ...int) (err error) {
	count := 30
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			// 获取实例详情
			detailParams := &mysql.TeledbQueryDetailRequest{
				OuterProdInstId: state.InstID.ValueString(),
			}
			detailHeaders := &mysql.TeledbQueryDetailRequestHeaders{
				InstID:   state.InstID.ValueString(),
				RegionID: state.RegionID.ValueString(),
			}
			if state.ProjectID.ValueString() != "" {
				detailHeaders.ProjectID = state.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 0 {
				err = fmt.Errorf("API return error. Message: %s", resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			runningStatus := resp.ReturnObj.ProdRunningStatus
			orderStatus := resp.ReturnObj.ProdOrderStatus
			// 若变配前，发现数据库已冻结，将其恢复
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("当前数据库状态为暂停状态，请启用后再进行绑定")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
				return false
			}
			if orderStatus == business.MysqlOrderStatusPause {
				err = errors.New("订单处于暂停状态，不可进行变更操作")
				return false
			}
			if runningStatus == business.MysqlRunningStatusStopping || runningStatus == business.MysqlRunningStatusStopped {
				err = errors.New("主机处于关机状态，不可进行变更操作")
				return false
			}

			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未到达启动状态！")
	}
	return
}

func (c *CtyunMysqlAssociationEip) preCheckEipUnbound(ctx context.Context, state *CtyunAssociationEipConfig) (err error) {
	eip, err := c.eipService.GetEipAddressByEipID(ctx, state.EipID.ValueString(), state.RegionID.ValueString())
	if err != nil {
		return
	}
	if *eip.Status != "DOWN" {
		return errors.New("eip 已经被绑定，无法再次绑定")
	}
	state.eipAddress = *eip.EipAddress
	return
}

type CtyunAssociationEipConfig struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性id
	InstID    types.String `tfsdk:"inst_id"`    //实例id
	ProjectID types.String `tfsdk:"project_id"` //项目id
	RegionID  types.String `tfsdk:"region_id"`  //区域Id
	EipStatus types.Int32  `tfsdk:"eip_status"` //弹性ip状态 0->unbind，1->bind,2->binding
	Status    types.String `tfsdk:"status"`     //弹性ip可读状态

	eipAddress string
}
