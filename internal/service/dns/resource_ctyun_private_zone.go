package dns

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunPrivateZone{}
	_ resource.ResourceWithConfigure   = &CtyunPrivateZone{}
	_ resource.ResourceWithImportState = &CtyunPrivateZone{}
)

type CtyunPrivateZone struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunPrivateZone) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_zone"
}

func (c *CtyunPrivateZone) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunPrivateZone() resource.Resource {
	return &CtyunPrivateZone{}
}

func (c *CtyunPrivateZone) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrivateZoneConfig
	var ID, regionId string
	// 根据分隔符数量判断是否输入了regionID
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
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
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

func (c *CtyunPrivateZone) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026757/10033657",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "关联的vpc，最多同时支持5个VPC。支持更新",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 5),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "内网DNS名称。要求：由多个以点分隔的字符串组成，可包含字母、数字中划线、中划线不能在开头或末尾，单个字符串不超过63个字符，域名总长度不超过254个字符",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 254),
					validator2.DnsName(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "内网DNS描述，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"proxy_pattern": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("zone"),
				Description: "可选值：zone：当前可用区不进行递归解析，record：不完全劫持，进行递归解析代理，大小写不敏感。默认值zone，支持更新。",
				Validators: []validator.String{
					stringvalidator.OneOf("zone", "record"),
				},
			},
			"ttl": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(300),
				Description: "zone ttl, 单位秒。默认300，取值范围300-2147483647。支持更新",
				Validators: []validator.Int32{
					int32validator.Between(300, 2147483647),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "zone id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"tags": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "标签",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_id": schema.StringAttribute{
							Computed:    true,
							Description: "标签id",
						},
						"key": schema.StringAttribute{
							Required:    true,
							Description: "标签key，支持更新。",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "标签value，支持更新。",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
			},
		},
	}
}

func (c *CtyunPrivateZone) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPrivateZoneConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZone) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPrivateZoneConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZone) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPrivateZoneConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPrivateZoneConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZone) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPrivateZoneConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunPrivateZone) create(ctx context.Context, config *CtyunPrivateZoneConfig) error {
	var vpcIds []string
	diags := config.VpcIDList.ElementsAs(ctx, &vpcIds, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params := &ctvpc.CtvpcCreatePrivateZoneRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     config.RegionID.ValueString(),
		Name:         config.Name.ValueString(),
		TTL:          config.TTL.ValueInt32(),
		VpcIDList:    strings.Join(vpcIds, ","),
		ProxyPattern: config.ProxyPattern.ValueStringPointer(),
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.Description = config.Description.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreatePrivateZoneApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建内网DNS失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.ZoneID)
	time.Sleep(2 * time.Second)
	// 标签处理
	if !config.Tags.IsNull() && !config.Tags.IsUnknown() {
		err = c.addLabel(ctx, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunPrivateZone) getAndMerge(ctx context.Context, config *CtyunPrivateZoneConfig) error {
	detailResp, err := c.getPrivateZoneDetail(ctx, config)
	if err != nil {
		return err
	}
	config.Description = types.StringValue(*detailResp.Description)
	config.TTL = types.Int32Value(detailResp.TTL)
	config.Name = types.StringValue(*detailResp.Name)
	config.CreateTime = types.StringValue(*detailResp.CreatedAt)
	config.UpdateTime = types.StringValue(*detailResp.UpdatedAt)
	config.ProxyPattern = types.StringValue(*detailResp.ProxyPattern)
	var vpcIds []string
	for _, vpcItem := range detailResp.VpcAssociations {
		vpcIds = append(vpcIds, *vpcItem.VpcID)
	}
	vpcIdsTmp, diags := types.SetValueFrom(ctx, types.StringType, vpcIds)
	if diags.HasError() {
		return fmt.Errorf(diags[0].Detail())
	}
	config.VpcIDList = vpcIdsTmp

	var tags []CtyunPrivateZoneTagModel
	// 获取tags列表
	respTags, err := c.getTags(ctx, config)
	if err != nil {
		return err
	}
	for _, tagItem := range respTags {
		var tag CtyunPrivateZoneTagModel
		tag.TagID = types.StringValue(*tagItem.LabelID)
		tag.Key = types.StringValue(*tagItem.LabelKey)
		tag.Value = types.StringValue(*tagItem.LabelValue)
		tags = append(tags, tag)
	}
	tagsTmp, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(CtyunPrivateZoneTagModel{}), tags)
	if diags.HasError() {
		return fmt.Errorf(diags[0].Detail())
	}
	config.Tags = tagsTmp
	return nil
}

func (c *CtyunPrivateZone) getPrivateZoneDetail(ctx context.Context, config *CtyunPrivateZoneConfig) (*ctvpc.CtvpcShowPrivateZoneReturnObjResponse, error) {
	params := &ctvpc.CtvpcShowPrivateZoneRequest{
		RegionID: config.RegionID.ValueString(),
		ZoneID:   config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPrivateZoneApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取内网DNS详情失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj, nil
}

func (c *CtyunPrivateZone) update(ctx context.Context, state *CtyunPrivateZoneConfig, plan *CtyunPrivateZoneConfig) error {

	params := &ctvpc.CtvpcUpdatePrivateZoneAttributeRequest{
		RegionID:     state.RegionID.ValueString(),
		ZoneID:       state.ID.ValueString(),
		TTL:          plan.TTL.ValueInt32(),
		ProxyPattern: plan.ProxyPattern.ValueStringPointer(),
	}
	if !plan.VpcIDList.Equal(state.VpcIDList) {
		var vpcIds []string
		diags := plan.VpcIDList.ElementsAs(ctx, &vpcIds, false)
		if diags.HasError() {
			err := fmt.Errorf(diags[0].Detail())
			return err
		}
		vpcIdsStr := strings.Join(vpcIds, ",")
		params.VpcIDList = &vpcIdsStr
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		params.Description = plan.Description.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatePrivateZoneAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新内网DNS失败(id=%s)，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	// 处理标签更新
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && !plan.Tags.Equal(state.Tags) {
		err = c.updateTags(ctx, state, plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunPrivateZone) delete(ctx context.Context, config CtyunPrivateZoneConfig) error {
	params := &ctvpc.CtvpcDeletePrivateZoneRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		ZoneID:      config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeletePrivateZoneApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除内网DNS失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrivateZone) addLabel(ctx context.Context, config *CtyunPrivateZoneConfig) error {
	var tags []CtyunPrivateZoneTagModel
	diags := config.Tags.ElementsAs(ctx, &tags, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	for _, tag := range tags {
		err := c.addLabelReq(ctx, config, tag.Key.ValueString(), tag.Value.ValueString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunPrivateZone) updateTags(ctx context.Context, state *CtyunPrivateZoneConfig, plan *CtyunPrivateZoneConfig) error {
	var stateTags []CtyunPrivateZoneTagModel
	var planTags []CtyunPrivateZoneTagModel
	diags := state.Tags.ElementsAs(ctx, &stateTags, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	diags = plan.Tags.ElementsAs(ctx, &planTags, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	// 若plan中有，state中没有，需要新增
	for _, planTag := range planTags {
		found := false
		for _, stateTag := range stateTags {
			if planTag.Key.Equal(stateTag.Key) && planTag.Value.Equal(stateTag.Value) {
				found = true
				break
			}
		}
		if !found {
			err := c.addLabelReq(ctx, state, planTag.Key.ValueString(), planTag.Value.ValueString())
			if err != nil {
				return err
			}
		}
	}
	// 若state中，plan有没有，需要删除
	for _, stateTag := range stateTags {
		found := false
		for _, planTag := range planTags {
			if planTag.Key.Equal(stateTag.Key) && planTag.Value.Equal(stateTag.Value) {
				found = true
				break
			}
		}
		if !found {
			err := c.removeLabelReq(ctx, state, stateTag.TagID.ValueString())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CtyunPrivateZone) addLabelReq(ctx context.Context, state *CtyunPrivateZoneConfig, key string, value string) error {
	params := &ctvpc.CtvpcZoneBindLabelsRequest{
		RegionID:   state.RegionID.ValueString(),
		ZoneID:     state.ID.ValueString(),
		LabelKey:   key,
		LabelValue: value,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcZoneBindLabelsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("内网DNS（id=%s）添加标签失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}

	// 查询列表确认添加成功
	labelID, err := c.getLabelID(ctx, state, key)
	if err != nil {
		return err
	}
	if labelID == "" {
		err = fmt.Errorf("未查询到labelID")
		return err
	}
	return nil
}

func (c *CtyunPrivateZone) removeLabelReq(ctx context.Context, config *CtyunPrivateZoneConfig, labelID string) error {
	params := &ctvpc.CtvpcZoneUnbindLabelsRequest{
		RegionID: config.RegionID.ValueString(),
		ZoneID:   config.ID.ValueString(),
		LabelID:  labelID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcZoneUnbindLabelsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("内网DNS（id=%s）删除标签失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrivateZone) getLabelID(ctx context.Context, config *CtyunPrivateZoneConfig, key string) (string, error) {
	params := &ctvpc.CtvpcListZoneBindLabelsRequest{
		RegionID: config.RegionID.ValueString(),
		ZoneID:   config.ID.ValueString(),
		PageNo:   1,
		PageSize: 50,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListZoneBindLabelsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return "", err
	} else if resp == nil {
		err = fmt.Errorf("获取内网DNS标签列表失败(id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return "", err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return "", err
	} else if resp.ReturnObj == nil || len(resp.ReturnObj.Results) < 1 {
		err = common.InvalidReturnObjError
		return "", err
	}
	for _, result := range resp.ReturnObj.Results {
		if result.LabelKey != nil && *result.LabelKey == key {
			return *result.LabelID, nil
		}
	}
	return "", nil
}

func (c *CtyunPrivateZone) getTags(ctx context.Context, config *CtyunPrivateZoneConfig) ([]*ctvpc.CtvpcListZoneBindLabelsReturnObjResultsResponse, error) {
	params := &ctvpc.CtvpcListZoneBindLabelsRequest{
		RegionID: config.RegionID.ValueString(),
		ZoneID:   config.ID.ValueString(),
		PageNo:   1,
		PageSize: 50,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListZoneBindLabelsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	}
	//if len(resp.ReturnObj.Results) < 1 {
	//	err = common.InvalidReturnObjError
	//	return nil, err
	//}
	return resp.ReturnObj.Results, err
}

type CtyunPrivateZoneTagModel struct {
	TagID types.String `tfsdk:"tag_id"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}
type CtyunPrivateZoneConfig struct {
	RegionID     types.String `tfsdk:"region_id"`
	VpcIDList    types.Set    `tfsdk:"vpc_id_list"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ProxyPattern types.String `tfsdk:"proxy_pattern"`
	TTL          types.Int32  `tfsdk:"ttl"`
	Tags         types.Set    `tfsdk:"tags"`
	ID           types.String `tfsdk:"id"`
	CreateTime   types.String `tfsdk:"create_time"`
	UpdateTime   types.String `tfsdk:"update_time"`
}
