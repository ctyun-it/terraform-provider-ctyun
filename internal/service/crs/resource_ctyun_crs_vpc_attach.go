package crs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	_ resource.Resource                = &ctyunCrsVpcAttach{}
	_ resource.ResourceWithConfigure   = &ctyunCrsVpcAttach{}
	_ resource.ResourceWithImportState = &ctyunCrsVpcAttach{}
)

type ctyunCrsVpcAttach struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
}

func NewCtyunCrsVpcAttach() resource.Resource {
	return &ctyunCrsVpcAttach{}
}

func (c *ctyunCrsVpcAttach) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_crs_vpc_attach"
}

type CtyunCrsVpcAttachConfig struct {
	ID       types.String `tfsdk:"id"`
	VpcID    types.String `tfsdk:"vpc_id"`
	SubnetID types.String `tfsdk:"subnet_id"`
	RegionID types.String `tfsdk:"region_id"`
}

func (c *ctyunCrsVpcAttach) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10007018/10007025`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				Optional:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
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
		},
	}
}

func (c *ctyunCrsVpcAttach) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCrsVpcAttachConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.attach(ctx, plan)
	if err != nil {
		if strings.Contains(err.Error(), "请稍后重新发起请求") {
			time.Sleep(10 * time.Second)
			err = c.attach(ctx, plan)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunCrsVpcAttach) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCrsVpcAttachConfig
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

func (c *ctyunCrsVpcAttach) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {

}

func (c *ctyunCrsVpcAttach) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCrsVpcAttachConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.detach(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDetach(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunCrsVpcAttach) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [vpcID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunCrsVpcAttachConfig
	var vpcID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		vpcID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &vpcID, &regionID)
		if err != nil {
			return
		}
	}

	if vpcID == "" {
		err = fmt.Errorf("vpcID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	cfg.RegionID = types.StringValue(regionID)
	cfg.VpcID = types.StringValue(vpcID)
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &cfg)...)
}

func (c *ctyunCrsVpcAttach) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
}

// attach vpc接入
func (c *ctyunCrsVpcAttach) attach(ctx context.Context, plan CtyunCrsVpcAttachConfig) (err error) {
	params := &crs.CrsCreateInstanceVpceLinkedVpcsV2Request{
		RegionId: plan.RegionID.ValueString(),
		VpcList: []*crs.CrsCreateInstanceVpceLinkedVpcsV2VpcListRequest{{
			VpcId:    plan.VpcID.ValueString(),
			SubnetId: plan.SubnetID.ValueStringPointer(),
		}},
	}
	resp, err := c.meta.Apis.SdkCrsApis.CrsCreateInstanceVpceLinkedVpcsV2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if len(resp.ReturnObj) == 0 || len(resp.ReturnObj[0].VpcLinkResult) == 0 {
		err = common.InvalidReturnObjError
	} else if *resp.ReturnObj[0].VpcLinkResult[0].State != "ACTIVE" {
		err = fmt.Errorf("API return error. Message: %s", *resp.ReturnObj[0].VpcLinkResult[0].Msg)
	}
	return
}

// checkBeforeAttach 接入前检查
func (c *ctyunCrsVpcAttach) checkBeforeAttach(ctx context.Context, plan CtyunCrsVpcAttachConfig) error {
	vpcID, subnetID := plan.VpcID.ValueString(), plan.SubnetID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpcID, plan.RegionID.ValueString(), "")
	if err != nil {
		return err
	}
	if len(subnets) == 0 {
		return fmt.Errorf("%s 至少要有一个子网", vpcID)
	}

	if _, ok := subnets[subnetID]; !ok {
		return fmt.Errorf("子网 %s 不属于 %s", subnetID, vpcID)
	}
	return nil
}

// disAttach vpc取消接入
func (c *ctyunCrsVpcAttach) detach(ctx context.Context, plan CtyunCrsVpcAttachConfig) (err error) {
	params := &crs.CrsDeleteInstanceVpceLinkedVpcsV2Request{
		RegionId:  plan.RegionID.ValueString(),
		VpcIdList: []string{plan.VpcID.ValueString()},
	}
	resp, err := c.meta.Apis.SdkCrsApis.CrsDeleteInstanceVpceLinkedVpcsV2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
	} else if len(resp.ReturnObj) == 0 || len(resp.ReturnObj[0].VpcLinkResult) == 0 {
		err = common.InvalidReturnObjError
	} else if *resp.ReturnObj[0].VpcLinkResult[0].State != "INACTIVE" {
		err = fmt.Errorf("API return error. Message: %s", *resp.ReturnObj[0].VpcLinkResult[0].Msg)
	}
	return
}

// checkAfterDetach 取消接入后检查
func (c *ctyunCrsVpcAttach) checkAfterDetach(ctx context.Context, plan CtyunCrsVpcAttachConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 10)
	retryer.Start(
		func(currentTime int) bool {
			var state string
			state, err = c.getAttach(ctx, plan)
			if err != nil {
				return false
			}
			if state != "INACTIVE" {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("%s取消接入失败", plan.VpcID.ValueString())
	}
	return
}

// getAttach 查询vpc接入状态
func (c *ctyunCrsVpcAttach) getAttach(ctx context.Context, plan CtyunCrsVpcAttachConfig) (state string, err error) {
	params := &crs.CrsGetInstanceVpceLinkedVpcsV2Request{
		RegionId:  plan.RegionID.ValueString(),
		VpcIdList: plan.VpcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCrsApis.CrsGetInstanceVpceLinkedVpcsV2Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if len(resp.ReturnObj) == 0 || len(resp.ReturnObj[0].VpcLinkResult) == 0 {
		err = common.InvalidReturnObjError
		return
	}
	state = *resp.ReturnObj[0].VpcLinkResult[0].State
	return
}

// getAndMerge 查询并合并
func (c *ctyunCrsVpcAttach) getAndMerge(ctx context.Context, plan *CtyunCrsVpcAttachConfig) (err error) {
	state, err := c.getAttach(ctx, *plan)
	if err != nil {
		return
	}
	if state != "ACTIVE" {
		err = common.ResourceNotExistError
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.VpcID.ValueString(), plan.RegionID.ValueString()))
	return
}
