package vpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	vpcSdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = &ctyunSecurityGroup{}
	_ resource.ResourceWithConfigure   = &ctyunSecurityGroup{}
	_ resource.ResourceWithImportState = &ctyunSecurityGroup{}
)

func NewCtyunSecurityGroup() resource.Resource {
	return &ctyunSecurityGroup{}
}

type ctyunSecurityGroup struct {
	meta       *common.CtyunMetadata
	name       string
	vpcService *business.VpcService
}

func (c *ctyunSecurityGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_security_group"
	c.name = response.TypeName
}

func (c *ctyunSecurityGroup) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理安全组", "虚拟私有云（Virtual Private Cloud，VPC）", "https://www.ctyun.cn/document/10026730/10225459"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "id",
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "vpcId",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "安全组名称，字母、中文、数字，下划线，连字符，中文/英文字母开头，长度2-32，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}-]+$"), "安全组名称不符合规则"),
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "描述，长度最大为128，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(128),
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
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
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
		},
	}
}

func (c *ctyunSecurityGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunSecurityGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.checkCreate(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	createResponse, err := c.meta.Apis.CtVpcApis.SecurityGroupCreateApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupCreateRequest{
		RegionId:    regionId,
		ProjectId:   projectId,
		ClientToken: uuid.NewString(),
		VpcId:       plan.VpcId.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(createResponse.SecurityGroupId)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, ctyunRequestError := c.getAndMergeSecurityGroup(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunSecurityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeSecurityGroup(ctx, state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunSecurityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	var plan CtyunSecurityGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	_, err := c.meta.Apis.CtVpcApis.SecurityGroupModifyAttributionApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupModifyAttributionRequest{
		SecurityGroupId: state.Id.ValueString(),
		RegionId:        state.RegionId.ValueString(),
		ProjectId:       state.ProjectId.ValueString(),
		ClientToken:     uuid.NewString(),
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		Enabled:         true,
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, ctyunRequestError := c.getAndMergeSecurityGroup(ctx, state)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunSecurityGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtVpcApis.SecurityGroupDeleteApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupDeleteRequest{
		SecurityGroupId: state.Id.ValueString(),
		RegionId:        state.RegionId.ValueString(),
		ProjectId:       state.ProjectId.ValueString(),
		ClientToken:     uuid.NewString(),
	})
	if err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic(err.Error(), err.Error()))
		return
	}
}

func (c *ctyunSecurityGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],<region_id>", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunSecurityGroupConfig
	var securityGroupId, regionId string
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		securityGroupId = request.ID
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &securityGroupId, &regionId)
		if err != nil {
			return
		}
	}

	if securityGroupId == "" {
		err = fmt.Errorf("security_group_id不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}

	cfg.Id = types.StringValue(securityGroupId)
	cfg.RegionId = types.StringValue(regionId)
	instance, err := c.getAndMergeSecurityGroup(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunSecurityGroup) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
}

// getAndMergeSecurityGroup 查询安全组
func (c *ctyunSecurityGroup) getAndMergeSecurityGroup(ctx context.Context, cfg CtyunSecurityGroupConfig) (*CtyunSecurityGroupConfig, error) {
	resp, err := c.meta.Apis.CtVpcApis.SecurityGroupDescribeAttributeApi.Do(ctx, c.meta.Credential, &ctvpc.SecurityGroupDescribeAttributeRequest{
		RegionId:        cfg.RegionId.ValueString(),
		SecurityGroupId: cfg.Id.ValueString(),
		Direction:       "all",
	})
	if err != nil {
		if err.ErrorCode() == common.OpenapiSecurityGroupNotFound {
			return nil, common.ResourceNotExistError
		}
		return nil, err
	}
	cfg.Id = types.StringValue(resp.Id)
	cfg.VpcId = types.StringValue(resp.VpcId)
	cfg.Name = types.StringValue(resp.SecurityGroupName)
	cfg.CreateTime = types.StringValue(resp.CreationTime)
	cfg.Description = types.StringValue(resp.Description)

	// 获取project id
	err2 := c.getProjectIdByName(ctx, &cfg)
	if err2 != nil {
		return nil, err2
	}
	return &cfg, nil
}

// checkCreate 校验创建动作是否能执行
func (c *ctyunSecurityGroup) checkCreate(ctx context.Context, plan CtyunSecurityGroupConfig) error {
	return c.vpcService.MustExist(ctx, plan.VpcId.ValueString(), plan.RegionId.ValueString())
}

func (c *ctyunSecurityGroup) getProjectIdByName(ctx context.Context, cfg *CtyunSecurityGroupConfig) error {
	params := &vpcSdk.CtvpcListSecurityGroupsRequest{
		RegionID: cfg.RegionId.ValueString(),
		VpcID:    cfg.VpcId.ValueStringPointer(),
		//InstanceID:   cfg.Id.ValueStringPointer(),
		QueryContent: cfg.Name.ValueStringPointer(),
		PageNo:       1,
		PageSize:     10,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListSecurityGroupsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("查询安全组失败，resp为空！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
		err = common.InvalidReturnObjError
		return err
	}
	cfg.ProjectId = types.StringValue(*resp.ReturnObj[0].ProjectID)
	return nil
}

type CtyunSecurityGroupConfig struct {
	Id          types.String `tfsdk:"id"`
	VpcId       types.String `tfsdk:"vpc_id"`
	Name        types.String `tfsdk:"name"`
	CreateTime  types.String `tfsdk:"create_time"`
	Description types.String `tfsdk:"description"`
	ProjectId   types.String `tfsdk:"project_id"`
	RegionId    types.String `tfsdk:"region_id"`
}
