package pgsql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	explanmodifier "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/planmodifier"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &CtyunPgsqlWhiteList{}
	_ resource.ResourceWithConfigure = &CtyunPgsqlWhiteList{}
)

type CtyunPgsqlWhiteList struct {
	meta *common.CtyunMetadata
	name string
}

func (c *CtyunPgsqlWhiteList) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_white_list"
	c.name = response.TypeName
}
func NewCtyunPgsqlWhiteList() resource.Resource {
	return &CtyunPgsqlWhiteList{}
}

func (c *CtyunPgsqlWhiteList) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPgsqlWhiteList) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理PostgreSQL实例的白名单", "关系数据库PostgreSQL版", "https://www.ctyun.cn/document/10034019/10161484"),
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
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MySQL实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: "废弃字段，请不要指定",
				Description:        "企业项目ID",
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Description: "修改模式。创建时必填，导入时可不填写。 cover(覆盖) ， append(追加) ， delete(删除),若分组下的ip被全部删除，则会将该分组也删除，默认分组(default)则会被设置成只允许本机访问，即只有127.0.0.1这个白名单ip)",
				Validators: []validator.String{
					stringvalidator.OneOf("cover", "append", "delete"),
				},
				PlanModifiers: []planmodifier.String{
					explanmodifier.NullIgnoreString(),
				},
			},
			"ip_list": schema.SetAttribute{
				Optional:    true,
				Description: "ip列表,数量限制：1-1000。创建、更新时必填，导入时可不填写",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1000),
				},
				PlanModifiers: []planmodifier.Set{
					explanmodifier.NullIgnoreSet(),
				},
			},
			"ip_list_result": schema.SetAttribute{
				Computed:    true,
				Description: "变更后最终的ip列表,数量限制：1-1000",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1000),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "白名单id",
			},
		},
	}
}

func (c *CtyunPgsqlWhiteList) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlWhiteListConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 开始创建备份集
	err = c.updateWhiteListRequest(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMergePostgresqlWhiteList(ctx, &plan)
	if err != nil {
		return
	}
	//plan.ID = types.StringValue(plan.BackupName.ValueString())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPgsqlWhiteList) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlWhiteListConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMergePostgresqlWhiteList(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
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

func (c *CtyunPgsqlWhiteList) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan CtyunPostgresqlWhiteListConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunPostgresqlWhiteListConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Mode.IsUnknown() && !plan.Mode.IsNull() && state.Mode.IsNull() {
		state.Mode = plan.Mode
		response.Diagnostics.AddWarning("mode的更新仅写入状态文件", "在import时，状态文件中mode为null，允许用模板中的值进行一次更新，该更新不触发远程调用")
	}
	if !plan.IpList.IsUnknown() && !plan.IpList.IsNull() && state.IpList.IsNull() {
		state.IpList = plan.IpList
		response.Diagnostics.AddWarning("ip_list的更新仅写入状态文件", "在import时，状态文件中ip_list为null，允许用模板中的值进行一次更新，该更新不触发远程调用")
	}

	state.ProjectID = plan.ProjectID
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunPgsqlWhiteList) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	return
}

func (c *CtyunPgsqlWhiteList) updateWhiteListRequest(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) error {
	if config.Mode.IsNull() || config.Mode.IsUnknown() {
		err := errors.New("创建、更新阶段mode必填！")
		return err
	}
	if config.IpList.IsNull() || config.IpList.IsUnknown() {
		err := errors.New("创建、更新阶段ip_list必填！")
		return err
	}
	var ips []string
	diags := config.IpList.ElementsAs(ctx, &ips, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	params := &pgsql.PgsqlUpdateWhiteListRequest{
		ProdInstId: config.InstID.ValueString(),
		Mode:       config.Mode.ValueString(),
		IpList:     ips,
	}
	header := &pgsql.PgsqlUpdateWhiteListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlUpdateWhiteListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例添加白名单ip失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return err
	}
	return nil
}

func (c *CtyunPgsqlWhiteList) getAndMergePostgresqlWhiteList(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) error {
	resp, err := c.getWhiteIpList(ctx, config)
	if err != nil {
		return err
	}
	ips, diags := types.SetValueFrom(ctx, types.StringType, resp.ReturnObj)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.IpListResult = ips
	config.ID = types.StringValue(fmt.Sprintf("%s_white_list", config.InstID.ValueString()))
	if config.IpList.IsNull() || config.IpList.IsUnknown() {
		config.IpList = types.SetNull(types.StringType)
	}
	if config.Mode.IsNull() || config.Mode.IsUnknown() {
		config.Mode = types.StringNull()
	}
	return nil
}

func (c *CtyunPgsqlWhiteList) getWhiteIpList(ctx context.Context, config *CtyunPostgresqlWhiteListConfig) (*pgsql.PgsqlGetWhiteListResponse, error) {
	params := &pgsql.PgsqlGetWhiteListRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	header := &pgsql.PgsqlGetWhiteListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlGetWhiteListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("postgresql实例获取白名单ip失败，接口返回nil，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(*resp.Error, "PG_2001") || strings.Contains(resp.Message, "未找到实例") {
			err = common.ResourceNotExistError
			return nil, err
		}
		err = fmt.Errorf(" API return error. Message: %s Error: %s", resp.Message, *resp.Error)
		return nil, err
	} else if resp.ReturnObj == nil || len(resp.ReturnObj) == 0 {
		err = common.ResourceNotExistError
		return nil, err
	}

	return resp, nil
}

type CtyunPostgresqlWhiteListConfig struct {
	InstID       types.String `tfsdk:"instance_id"`
	RegionID     types.String `tfsdk:"region_id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Mode         types.String `tfsdk:"mode"`
	IpList       types.Set    `tfsdk:"ip_list"`
	IpListResult types.Set    `tfsdk:"ip_list_result"`
	ID           types.String `tfsdk:"id"`
}
