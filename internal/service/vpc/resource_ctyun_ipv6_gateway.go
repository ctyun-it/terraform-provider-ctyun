package vpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
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
)

var (
	_ resource.Resource                = &ctyunIPv6GatewayService{}
	_ resource.ResourceWithConfigure   = &ctyunIPv6GatewayService{}
	_ resource.ResourceWithImportState = &ctyunIPv6GatewayService{}
)

type ctyunIPv6GatewayService struct {
	meta *common.CtyunMetadata
	name string
}

func NewCtyunIPv6GatewayService() resource.Resource {
	return &ctyunIPv6GatewayService{}
}

func (c *ctyunIPv6GatewayService) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ipv6_gateway"
	c.name = response.TypeName
}

type CtyunIPv6GatewayConfig struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	RegionID  types.String `tfsdk:"region_id"`
	VpcID     types.String `tfsdk:"vpc_id"`
	ProjectID types.String `tfsdk:"project_id"`
}

func (c *ctyunIPv6GatewayService) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理IPv6带宽", "IPv6带宽（IPv6 Bandwidth）", "https://www.ctyun.cn/document/10026753/10037269"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "IPv6网关ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "IPv6网关名称",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					explanmodifier.Project(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
		},
	}
}

func (c *ctyunIPv6GatewayService) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunIPv6GatewayConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	clientToken := uuid.NewString()
	err = c.create(ctx, clientToken, &plan)
	if err != nil {
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunIPv6GatewayService) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIPv6GatewayConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunIPv6GatewayService) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (c *ctyunIPv6GatewayService) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIPv6GatewayConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除
	err = c.delete(ctx, &state)
	if err != nil {
		return
	}

	return
}

func (c *ctyunIPv6GatewayService) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunIPv6GatewayService) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import %s.[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunIPv6GatewayConfig

	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionId
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

// create 创建
func (c *ctyunIPv6GatewayService) create(ctx context.Context, clientToken string, plan *CtyunIPv6GatewayConfig) (err error) {
	params := &ctvpc.CtvpcCreateIPv6GatewayRequest{
		ClientToken: clientToken,
		RegionID:    plan.RegionID.ValueString(),
		VpcID:       plan.VpcID.ValueString(),
		ProjectID:   plan.ProjectID.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateIPv6GatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}

func (c *ctyunIPv6GatewayService) getAndMerge(ctx context.Context, plan *CtyunIPv6GatewayConfig) (err error) {
	regionID := plan.RegionID.ValueString()
	params := &ctvpc.CtvpcListIPv6GatewayRequest{
		RegionID: regionID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListIPv6GatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	} else if len(resp.ReturnObj) == 0 {
		err = common.ResourceNotExistError
		return
	}

	for _, gateway := range resp.ReturnObj {
		if *gateway.VpcID == plan.VpcID.ValueString() || *gateway.Ipv6GatewayID == plan.ID.ValueString() {
			plan.ID = types.StringValue(utils.SecString(gateway.Ipv6GatewayID))
			plan.Name = types.StringValue(utils.SecString(gateway.Name))
			plan.VpcID = types.StringValue(utils.SecString(gateway.VpcID))
			plan.ProjectID = types.StringValue(utils.SecString(gateway.ProjectIdEcs))
			break
		}
	}

	return
}

func (c *ctyunIPv6GatewayService) delete(ctx context.Context, plan *CtyunIPv6GatewayConfig) (err error) {
	clientToken := uuid.NewString()
	ID, regionID := plan.ID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteIPv6GatewayRequest{
		RegionID:      regionID,
		Ipv6GatewayID: ID,
		ClientToken:   clientToken,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteIPv6GatewayApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}

	return
}
