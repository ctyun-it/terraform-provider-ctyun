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
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunPrivateZoneRecord{}
	_ resource.ResourceWithConfigure   = &CtyunPrivateZoneRecord{}
	_ resource.ResourceWithImportState = &CtyunPrivateZoneRecord{}
)

type CtyunPrivateZoneRecord struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunPrivateZoneRecord) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_zone_record"
}

func (c *CtyunPrivateZoneRecord) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunPrivateZoneRecord() resource.Resource {
	return &CtyunPrivateZoneRecord{}
}

func (c *CtyunPrivateZoneRecord) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPrivateZoneRecordConfig
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

func (c *CtyunPrivateZoneRecord) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026757/10224466",
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
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "内网DNS id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: "内网DNS记录类型，支持: A / CNAME / MX / AAAA / TXT, 大小写不敏感。" +
					"A-将域名指向一个IPv4地址；" +
					"CNAME-将域名指向另一个域名；MX-邮件交换记录，用于指定接收电子邮件的服务器；" +
					"MX-邮件交换记录，用于指定接收电子邮件的服务器；" +
					"AAAA-将域名指向一个IPv6地址；" +
					"TXT-文本记录，可以包含任意文本信息；",
				//"SRV-记录提供特定服务的服务器",
				Validators: []validator.String{
					stringvalidator.OneOf("A", "CNAME", "MX", "AAAA", "TXT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value_list": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "支持更新。当type=A，value_list必须是 IPv4 地址；" +
					"当type=CNAME，value_list填写您要指向的别名，只能写一个域名；" +
					"当type=MX，value_list 填写邮箱服务器地址，最多可以输入8个不重复地址；" +
					"当type=AAAA，valueList 填写IPv6地址，最多可以输入8个不重复地址；" +
					"当type= TXT 时，valueList 填写文本记录值(合法字符包含大小写字母、数字、空格，文本记录)",
				//"当type=SRV时，valueList填写指定服务的服务器地址，最多可以输入8个不重复地址。 数组中元素格式要求如下：_服务名._协议名.主机名:端口号；",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 8),
				},
			},
			"ttl": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "zone ttl，支持更新。TTL指解析记录在本地DNS服务器的缓存时间。如果您的服务地址经常更换，建议TTL值设置相对小些，反之，建议设置相对大些。",
				Default:     int32default.StaticInt32(300),
				Validators: []validator.Int32{
					int32validator.Between(300, 2147483647),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "DNS记录集的 name",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "DNS记录集描述，支持更新",
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "内网DNS记录id",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启解析记录，默认启用，支持更新。",
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (c *CtyunPrivateZoneRecord) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunPrivateZoneRecordConfig
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

func (c *CtyunPrivateZoneRecord) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPrivateZoneRecordConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPrivateZoneRecord) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPrivateZoneRecordConfig
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

func (c *CtyunPrivateZoneRecord) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunPrivateZoneRecordConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunPrivateZoneRecord) create(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcCreatePrivateZoneRecordRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		ZoneID:      config.ZoneID.ValueString(),
		RawType:     config.Type.ValueString(),
		TTL:         config.TTL.ValueInt32(),
		Name:        config.Name.ValueStringPointer(),
	}
	var values []string
	diags := config.ValueList.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params.ValueList = values
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.Description = config.Description.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreatePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建内网DNS记录失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.ZoneRecordID)
	// 是否需要关闭解析
	err = c.controlEnable(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunPrivateZoneRecord) getAndMerge(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcShowPrivateZoneRecordRequest{
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowPrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取内网DNS记录(id=%s)详情失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.Name = types.StringValue(*resp.ReturnObj.Name)
	config.Description = types.StringValue(*resp.ReturnObj.Description)
	config.ZoneID = types.StringValue(*resp.ReturnObj.ZoneID)
	config.Type = types.StringValue(*resp.ReturnObj.RawType)
	config.TTL = types.Int32Value(resp.ReturnObj.TTL)
	config.CreatedTime = types.StringValue(*resp.ReturnObj.CreatedAt)
	config.UpdatedTime = types.StringValue(*resp.ReturnObj.UpdatedAt)
	// 处理value
	var values []string
	for _, value := range resp.ReturnObj.Value {
		if value == nil {
			continue
		}
		if config.Type.ValueString() == "TXT" {
			result := strings.Trim(*value, `"`)
			values = append(values, result)
		} else {
			values = append(values, *value)
		}

	}
	valueTmp, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.ValueList = valueTmp
	return nil
}

func (c *CtyunPrivateZoneRecord) update(ctx context.Context, state *CtyunPrivateZoneRecordConfig, plan *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcUpdatePrivateZoneRecordAttributeRequest{
		RegionID:     state.RegionID.ValueString(),
		ZoneRecordID: state.ID.ValueString(),
		TTL:          plan.TTL.ValueInt32(),
	}
	var values []string
	diags := plan.ValueList.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params.ValueList = values

	if !plan.Description.Equal(state.Description) {
		params.Description = plan.Description.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdatePrivateZoneRecordAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新内网DNS记录(id=%s)属性失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	plan.ID = state.ID
	err = c.controlEnable(ctx, plan)
	if err != nil {
		return err
	}
	state.Enabled = plan.Enabled
	return nil
}

func (c *CtyunPrivateZoneRecord) delete(ctx context.Context, config CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcDeletePrivateZoneRecordRequest{
		ClientToken:  uuid.NewString(),
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeletePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除内网DNS记录(id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrivateZoneRecord) controlEnable(ctx context.Context, config *CtyunPrivateZoneRecordConfig) (err error) {
	if config.Enabled.ValueBool() {
		err = c.enableRecord(ctx, config)
	} else {
		err = c.disableRecord(ctx, config)
	}
	return
}

func (c *CtyunPrivateZoneRecord) enableRecord(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcEnablePrivateZoneRecordRequest{
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcEnablePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("开启内网DNS记录(id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunPrivateZoneRecord) disableRecord(ctx context.Context, config *CtyunPrivateZoneRecordConfig) error {
	params := &ctvpc.CtvpcDisablePrivateZoneRecordRequest{
		RegionID:     config.RegionID.ValueString(),
		ZoneRecordID: config.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDisablePrivateZoneRecordApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("关闭内网DNS记录(id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunPrivateZoneRecordConfig struct {
	RegionID    types.String `tfsdk:"region_id"`
	ZoneID      types.String `tfsdk:"zone_id"`
	Type        types.String `tfsdk:"type"`
	ValueList   types.Set    `tfsdk:"value_list"`
	TTL         types.Int32  `tfsdk:"ttl"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	ID          types.String `tfsdk:"id"`
	CreatedTime types.String `tfsdk:"create_time"`
	UpdatedTime types.String `tfsdk:"update_time"`
}
