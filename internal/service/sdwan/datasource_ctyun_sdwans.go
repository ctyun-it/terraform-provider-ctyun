package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &CtyunSdwans{}
)

func NewCtyunSdwans() datasource.DataSource {
	return &CtyunSdwans{}
}

type CtyunSdwans struct {
	meta *common.CtyunMetadata
}

type CtyunSdwansConfig struct {
	ID        types.String      `tfsdk:"id"`
	RegionID  types.String      `tfsdk:"region_id"`
	ProjectID types.String      `tfsdk:"project_id"`
	SdwanID   types.String      `tfsdk:"sdwan_id"`
	Sdwans    []CtyunSdwanInfos `tfsdk:"sdwans"`
}

type CtyunSdwanInfos struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Desc        types.String `tfsdk:"description"`
	CreatedTime types.String `tfsdk:"created_time"`
}

func (c *CtyunSdwans) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwans"
}

func (c *CtyunSdwans) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10035786/10035852`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
			},
			"sdwan_id": schema.StringAttribute{
				Optional:    true,
				Description: "SD-WAN ID，用于精确查询",
			},
			"sdwans": schema.ListNestedAttribute{
				Computed:    true,
				Description: "SD-WAN列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "SD-WAN ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "SD-WAN名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "SD-WAN描述",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSdwans) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwans) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSdwansConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &sdwan.SdwanListSdwanRequest{}
	if !state.SdwanID.IsNull() {
		request.SdwanID = state.SdwanID.ValueStringPointer()
	}

	tflog.Info(ctx, "查询SD-WAN列表", map[string]interface{}{
		"region_id": state.RegionID.ValueString(),
	})

	response, err := c.meta.Apis.SdkSdwanApis.SdwanListSdwanApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		return
	} else if response.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *response.Message)
		return
	} else if response.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析结果
	var sdwanList []CtyunSdwanInfos
	for _, sdwanItem := range response.ReturnObj.Result {
		sdwanInfo := CtyunSdwanInfos{
			ID:          types.StringValue(*sdwanItem.SdwanID),
			Name:        types.StringValue(*sdwanItem.SdwanName),
			CreatedTime: types.StringValue(*sdwanItem.CreateTime),
		}

		if sdwanItem.Description != nil {
			sdwanInfo.Desc = types.StringValue(*sdwanItem.Description)
		}

		sdwanList = append(sdwanList, sdwanInfo)
	}

	state.Sdwans = sdwanList
	state.ID = types.StringValue(fmt.Sprintf("%s,%s", state.RegionID.ValueString(), state.SdwanID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
