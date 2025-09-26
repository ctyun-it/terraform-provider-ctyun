package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	_ resource.Resource                = &CtyunPgsqlAssociationEip{}
	_ resource.ResourceWithConfigure   = &CtyunPgsqlAssociationEip{}
	_ resource.ResourceWithImportState = &CtyunPgsqlAssociationEip{}
)

type CtyunPgsqlAssociationEip struct {
	meta       *common.CtyunMetadata
	eipService *business.EipService
}

func (c *CtyunPgsqlAssociationEip) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_association_eip"
}
func NewCtyunMysqlAssociationEip() resource.Resource {
	return &CtyunPgsqlAssociationEip{}
}

func (c *CtyunPgsqlAssociationEip) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10174601`,
		Attributes: map[string]schema.Attribute{
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性id",
				Validators: []validator.String{
					validator2.EipValidate(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"eip": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip地址",
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
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
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量",
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
				Validators: []validator.Int32{
					int32validator.Between(0, 2),
				},
			},
		},
	}
}

func (c *CtyunPgsqlAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPgsqlAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 实例绑定弹性IP
	err = c.PgsqlBindEip(ctx, &plan)
	// 查询eip状态，确认是否被绑定
	// 轮询查看绑定状态
	err = c.BindLoop(ctx, &plan, business.EipStatusBind)
	if err != nil {
		return
	}
	err = c.getAndMergeBindEip(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPgsqlAssociationEipConfig
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

func (c *CtyunPgsqlAssociationEip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunPgsqlAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPgsqlAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	unbindParams := &pgsql.PgsqlUnBindEipRequest{
		EipID:  state.EipID.ValueString(),
		Eip:    state.Eip.ValueString(),
		InstID: state.InstID.ValueString(),
	}
	unbindHeader := &pgsql.PgsqlUnBindEipRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		unbindHeader.ProjectId = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUnBindEipApi.Do(ctx, c.meta.Credential, unbindParams, unbindHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	// 轮询确认eip是否解绑成功
	err = c.BindLoop(ctx, &state, business.EipStatusUnbind)
	if err != nil {
		return
	}
	state.EipStatus = types.Int32Value(business.EipStatusUnbind)

	return
}

func (c *CtyunPgsqlAssociationEip) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.eipService = business.NewEipService(c.meta)
}
func (c *CtyunPgsqlAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
}

func (c *CtyunPgsqlAssociationEip) PgsqlBindEip(ctx context.Context, config *CtyunPgsqlAssociationEipConfig) (err error) {
	eip, err := c.eipService.GetEipAddressByEipID(ctx, config.EipID.ValueString(), config.RegionID.ValueString())
	if err != nil {
		return
	}
	config.Eip = types.StringValue(*eip.EipAddress)
	params := &pgsql.PgsqlBindEipRequest{
		EipID:  config.EipID.ValueString(),
		Eip:    *eip.EipAddress,
		InstID: config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlBindEipRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectId = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlBindEipApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	}
	return
}

func (c *CtyunPgsqlAssociationEip) BindLoop(ctx context.Context, config *CtyunPgsqlAssociationEipConfig, bindStatus int32, loopCount ...int) (err error) {
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
			header := &mysql.TeledbBoundEipListRequestHeader{}
			if config.ProjectID.ValueString() != "" {
				header.ProjectID = config.ProjectID.ValueStringPointer()
			}
			resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbBoundEipListApi.Do(ctx, c.meta.Credential, params, header)
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
			if len(returnObj) != 1 {
				err = fmt.Errorf("eip获取数量有误！")
				return false
			}
			if returnObj[0].BindStatus == bindStatus {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，eip仍未绑定/解绑成功！")
	}
	return
}

func (c *CtyunPgsqlAssociationEip) getAndMergeBindEip(ctx context.Context, config *CtyunPgsqlAssociationEipConfig) (err error) {
	params := &mysql.TeledbBoundEipListRequest{
		RegionID: config.RegionID.ValueString(),
		EipID:    config.EipID.ValueStringPointer(),
	}
	header := &mysql.TeledbBoundEipListRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbBoundEipListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s ", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回的绑定eip列表
	returnObj := resp.ReturnObj.Data
	if len(returnObj) != 1 {
		err = fmt.Errorf("eip获取数量有误！")
		return
	}

	config.EipStatus = types.Int32Value(returnObj[0].BindStatus)
	return
}

type CtyunPgsqlAssociationEipConfig struct {
	EipID     types.String `tfsdk:"eip_id"`     //弹性id
	Eip       types.String `tfsdk:"eip"`        //弹性ip
	InstID    types.String `tfsdk:"inst_id"`    //实例id
	ProjectID types.String `tfsdk:"project_id"` //项目id
	RegionID  types.String `tfsdk:"region_id"`  //区域Id
	EipStatus types.Int32  `tfsdk:"eip_status"` //弹性ip状态 0->unbind，1->bind
}
