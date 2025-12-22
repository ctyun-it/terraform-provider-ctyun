package vpc

import (
	"context"
	"fmt"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"regexp"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunDhcpOptionSet{}
	_ resource.ResourceWithConfigure   = &ctyunDhcpOptionSet{}
	_ resource.ResourceWithImportState = &ctyunDhcpOptionSet{}
)

func NewCtyunDhcpOptionSet() resource.Resource {
	return &ctyunDhcpOptionSet{}
}

type ctyunDhcpOptionSet struct {
	meta *common.CtyunMetadata
}

func (c *ctyunDhcpOptionSet) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunDhcpOptionSetConfig
	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionID,projectId
	if strings.Count(request.ID, common.ImportSeparator) == 0 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	config.Id = types.StringValue(ID)
	config.RegionId = types.StringValue(regionId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *ctyunDhcpOptionSet) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dhcpoptionset"
}

func (c *ctyunDhcpOptionSet) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "DHCP选项集ID",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "集合名，支持拉丁字母、中文、数字，下划线，连字符，必须以中文/英文字母开头，不能以数字、_和-、http:/https:开头，长度2-32 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z\\x{4e00}-\\x{9fa5}][0-9a-zA-Z_\\x{4e00}-\\x{9fa5}]+$"), "dhcp名称不符合规则"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "描述信息，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&**()_-+= <>?:\"{},./;'[**\r\n\r\n**]·~！@#￥%……&**（） —— -+={}《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128 支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "整个域名的总长度不能超过 255 个字符，每个子域名（包括顶级域名）的长度不能超过 63 个字符，域名中的字符集包括大写字母、小写字母、数字和连字符（减号），连字符不能位于域名的开头 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^-][a-zA-Z0-9\-]*(\.[a-zA-Z0-9\-]+)*$`), "连字符不能位于域名的开头，且格式应符合域名规范"),
				},
			},
			"dns_list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "服务ip地址列表，最多只能4个IP地址 支持更新",
				Validators: []validator.List{
					listvalidator.SizeAtMost(4),
				},
			},
		},
	}
}

func (c *ctyunDhcpOptionSet) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunDhcpOptionSet) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	// 关键：确保资源创建后正确设置所有必需属性
	// 特别是资源ID必须被设置
	if plan.Id.IsNull() || plan.Id.ValueString() == "" {
		response.Diagnostics.AddError(
			"Missing Resource ID",
			"Resource ID was not set after creation. This is a provider bug that should be reported.",
		)
		return
	}

	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSet) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *ctyunDhcpOptionSet) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunDhcpOptionSet) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunDhcpOptionSetConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := c.delete(ctx, &state)
	if err != nil {
		return
	}
}

// create 创建dhcpoptionset
func (c *ctyunDhcpOptionSet) create(ctx context.Context, plan *CtyunDhcpOptionSetConfig) (err error) {
	// 准备请求参数
	req := &ctvpc.CtvpcDhcpoptionsetscreateRequest{
		RegionID:   plan.RegionId.ValueString(),
		Name:       plan.Name.ValueString(),
		DomainName: plan.DomainName.ValueString(),
	}

	// 设置可选参数
	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	// 设置DNS列表
	var dnsList []string
	for _, dns := range plan.DnsList {
		dnsList = append(dnsList, dns.ValueString())
	}
	req.DnsList = dnsList

	// 调用API创建DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetscreateApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 设置资源ID
	plan.Id = utils.SecStringValue(resp.ReturnObj.DhcpOptionSetsID)

	return nil
}

// getAndMerge 查询dhcpoptionset并合并状态
func (c *ctyunDhcpOptionSet) getAndMerge(ctx context.Context, state *CtyunDhcpOptionSetConfig) (err error) {
	// 调用API获取DHCP选项集详情
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetsShowApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpoptionsetsShowRequest{
		RegionID:         state.RegionId.ValueString(),
		DhcpOptionSetsID: state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 更新状态
	state.Id = utils.SecStringValue(resp.ReturnObj.DhcpOptionSetsID)
	state.Name = utils.SecStringValue(resp.ReturnObj.Name)
	state.Description = utils.SecStringValue(resp.ReturnObj.Description)

	state.DomainName = utils.SecStringValue(resp.ReturnObj.DomainName)

	// 更新DNS列表
	state.DnsList = []types.String{}
	for _, dns := range resp.ReturnObj.DnsList {
		if dns != nil {
			state.DnsList = append(state.DnsList, types.StringValue(*dns))
		}
	}

	return
}

// update 更新dhcpoptionset
func (c *ctyunDhcpOptionSet) update(ctx context.Context, plan *CtyunDhcpOptionSetConfig) (err error) {
	// 准备请求参数
	req := &ctvpc.CtvpcUpdatedhcpoptionsetsRequest{
		RegionID:         plan.RegionId.ValueString(),
		DhcpOptionSetsID: plan.Id.ValueString(),
	}

	// 设置可选参数
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		req.Name = &name
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		req.Description = &description
	}

	if !plan.DomainName.IsNull() {
		domainName := plan.DomainName.ValueString()
		req.DomainName = &domainName
	}

	// 设置DNS列表
	var dnsList []*string
	for _, dns := range plan.DnsList {
		dnsValue := dns.ValueString()
		dnsList = append(dnsList, &dnsValue)
	}
	req.DnsList = dnsList

	// 调用API更新DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatedhcpoptionsetsApi.Do(ctx, c.meta.SdkCredential, req)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}

// delete 删除dhcpoptionset
func (c *ctyunDhcpOptionSet) delete(ctx context.Context, state *CtyunDhcpOptionSetConfig) (err error) {
	// 调用API删除DHCP选项集
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetsdeleteApi.Do(ctx, c.meta.SdkCredential, &ctvpc.CtvpcDhcpoptionsetsdeleteRequest{
		RegionID:         state.RegionId.ValueString(),
		DhcpOptionSetsID: state.Id.ValueString(),
	})
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return nil
}

type CtyunDhcpOptionSetConfig struct {
	Id          types.String   `tfsdk:"id"`
	RegionId    types.String   `tfsdk:"region_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	DomainName  types.String   `tfsdk:"domain_name"`
	DnsList     []types.String `tfsdk:"dns_list"`
}
