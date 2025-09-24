package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbAcl{}
	_ resource.ResourceWithConfigure   = &CtyunElbAcl{}
	_ resource.ResourceWithImportState = &CtyunElbAcl{}
)

type CtyunElbAcl struct {
	meta *common.CtyunMetadata
}

func (c *CtyunElbAcl) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbAcl) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbAcl) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_acl"
}

func (c *CtyunElbAcl) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10032777**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
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
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_+= <>?:,.,/;'[]·！@#￥%……&*（） ——+={}，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					validator2.Desc(),
				},
			},
			"source_ips": schema.SetAttribute{
				Required:    true,
				Description: "IP地址的集合或者CIDR, 单次最多添加 10 条数据，支持更新",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(10),
					setvalidator.ValueStringsAre(stringvalidator.UTF8LengthAtLeast(1)),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "访问控制ID",
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				// az时候有必要设定默认值
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, false),
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
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunElbAcl) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunElbAclConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.createAcl(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后反查创建后的nat信息
	err = c.getAndMergeAcl(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbAcl) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbAclConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeAcl(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbAcl) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunElbAclConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}
	// 更新ACL基本信息
	err = c.updateAclInfo(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeAcl(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbAcl) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbAclConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &ctelb.CtelbDeleteAccessControlRequest{
		ClientToken:     uuid.NewString(),
		RegionID:        state.RegionID.ValueString(),
		AccessControlID: state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteAccessControlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
}

func (c *CtyunElbAcl) createAcl(ctx context.Context, config *CtyunElbAclConfig) (err error) {
	if config.RegionID.ValueString() == "" {
		err = fmt.Errorf("region id不能为空")
		return
	}

	params := &ctelb.CtelbCreateAccessControlRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		Name:        config.Name.ValueString(),
	}
	if !config.Description.IsNull() {
		params.Description = config.Description.ValueString()
	}
	var sourceIps []string
	if !config.SourceIps.IsNull() {
		diag := config.SourceIps.ElementsAs(ctx, &sourceIps, false)
		if diag.HasError() {
			return
		}
		params.SourceIps = sourceIps
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateAccessControlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.ID = types.StringValue(resp.ReturnObj.ID)
	return
}

func (c *CtyunElbAcl) getAndMergeAcl(ctx context.Context, config *CtyunElbAclConfig) (err error) {
	params := &ctelb.CtelbShowAccessControlRequest{
		RegionID:        config.RegionID.ValueString(),
		AccessControlID: config.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowAccessControlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	returnObj := resp.ReturnObj
	//解析acl详情
	config.Name = types.StringValue(returnObj.Name)
	config.Description = types.StringValue(returnObj.Description)
	config.CreateTime = types.StringValue(returnObj.CreateTime)
	// 解析sourceIps
	var sourceIps []types.String
	for _, sourceIp := range returnObj.SourceIps {
		sourceIps = append(sourceIps, types.StringValue(sourceIp))
	}
	config.SourceIps, _ = types.SetValueFrom(ctx, types.StringType, sourceIps)
	return
}

func (c *CtyunElbAcl) updateAclInfo(ctx context.Context, state *CtyunElbAclConfig, plan *CtyunElbAclConfig) (err error) {
	params := &ctelb.CtelbUpdateAccessControlRequest{
		ClientToken:     uuid.NewString(),
		RegionID:        state.RegionID.ValueString(),
		SourceIps:       nil,
		AccessControlID: state.ID.ValueString(),
	}
	if !plan.Name.IsNull() && plan.Name.ValueString() != state.Name.ValueString() {
		params.Name = plan.Name.ValueString()
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != state.Description.ValueString() {
		params.Description = plan.Description.ValueString()
	}
	var planSourceIps []string
	var stateSourceIps []string
	plan.SourceIps.ElementsAs(ctx, &planSourceIps, false)
	state.SourceIps.ElementsAs(ctx, &stateSourceIps, false)
	// 判断state和plan的source ip是否相同，如果不同则更新
	if !c.compareStringList(planSourceIps, stateSourceIps) {
		params.SourceIps = planSourceIps
	}

	if params.Name == "" && params.Description == "" && params.SourceIps == nil {
		return
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateAccessControlApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// 将切片转换为 set
func (c *CtyunElbAcl) sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, v := range slice {
		set[v] = struct{}{} // 空结构体占位
	}
	return set
}

func (c *CtyunElbAcl) compareStringList(arr1 []string, arr2 []string) bool {
	set1 := c.sliceToSet(arr1)
	set2 := c.sliceToSet(arr2)
	// 首先，检查两个 set 的长度是否相同
	if len(set1) != len(set2) {
		return false
	}
	// 然后检查 set 中的元素是否完全相同
	for key := range set1 {
		if _, exists := set2[key]; !exists {
			return false
		}
	}
	return true
}

func NewCtyunElbAcl() resource.Resource {
	return &CtyunElbAcl{}
}

type CtyunElbAclConfig struct {
	RegionID    types.String `tfsdk:"region_id"`   //区域ID
	Name        types.String `tfsdk:"name"`        //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description types.String `tfsdk:"description"` //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_+= <>?:,.,/;'[]·！@#￥%……&*（） ——+={}
	SourceIps   types.Set    `tfsdk:"source_ips"`  //IP地址的集合或者CIDR, 单次最多添加 10 条数据
	ID          types.String `tfsdk:"id"`          //访问控制ID
	AzName      types.String `tfsdk:"az_name"`     //可用区名称
	ProjectID   types.String `tfsdk:"project_id"`  //项目ID
	CreateTime  types.String `tfsdk:"create_time"` //创建时间，为UTC格式
}
