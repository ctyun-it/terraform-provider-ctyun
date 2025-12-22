package mongodb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var (
	_ datasource.DataSource              = &CtyunMongodbWhiteLists{}
	_ datasource.DataSourceWithConfigure = &CtyunMongodbWhiteLists{}
)

func NewCtyunMongodbWhiteLists() datasource.DataSource {
	return &CtyunMongodbWhiteLists{}
}

type CtyunMongodbWhiteLists struct {
	meta *common.CtyunMetadata
}

func (c *CtyunMongodbWhiteLists) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongodb_white_lists"
}

func (c *CtyunMongodbWhiteLists) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10089535`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "MongoDB实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"white_lists": schema.ListNestedAttribute{
				Computed:    true,
				Description: "白名单分组列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_whitelist_name": schema.StringAttribute{
							Computed:    true,
							Description: "白名单分组名称",
						},
						"ip_list": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "IP列表",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunMongodbWhiteLists) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunMongodbWhiteLists) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MongodbWhiteListsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeReq := &mongodb.MongodbDescribeIpWhitelistRequest{
		ProdInstId: data.InstanceID.ValueString(),
	}
	regionId := c.meta.GetExtraIfEmpty(data.RegionID.ValueString(), common.ExtraRegionId)

	headers := &mongodb.MongodbDescribeIpWhitelistRequestHeaders{
		RegionID: regionId,
	}
	if !data.ProjectID.IsNull() {
		headers.ProjectID = data.ProjectID.ValueStringPointer()
	}

	tflog.Info(ctx, "查询MongoDB白名单列表", map[string]interface{}{
		"instance_id": data.InstanceID.ValueString(),
	})

	describeResp, err := c.meta.Apis.SdkMongodbApis.MongodbDescribeIpWhitelistApi.Do(ctx, c.meta.Credential, describeReq, headers)
	if err != nil {
		resp.Diagnostics.AddError("查询MongoDB白名单列表失败", err.Error())
		return
	}

	if describeResp.StatusCode != 800 {
		if describeResp.Message != nil {
			resp.Diagnostics.AddError("查询MongoDB白名单列表失败", fmt.Sprintf("API返回错误: %s", *describeResp.Message))
		} else {
			resp.Diagnostics.AddError("查询MongoDB白名单列表失败", fmt.Sprintf("API返回错误，状态码: %d", describeResp.StatusCode))
		}
		return
	}

	if describeResp.ReturnObj == nil {
		resp.Diagnostics.AddError("查询MongoDB白名单列表失败", "API返回空结果")
		return
	}

	// 转换API响应数据到Terraform模型
	var whiteLists []MongodbWhiteListModel
	for _, group := range describeResp.ReturnObj.WhitelistGroup {
		// 处理IP列表，API返回的是逗号分隔的字符串
		var ipListSlice []string
		if group.IpList != "" {
			ipListSlice = strings.Split(group.IpList, ",")
		}

		ipList, diags := types.ListValueFrom(ctx, types.StringType, ipListSlice)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		whiteLists = append(whiteLists, MongodbWhiteListModel{
			IpWhitelistName: types.StringValue(group.GroupName),
			IpList:          ipList,
		})
	}

	data.WhiteLists = whiteLists
	data.ID = types.StringValue(data.InstanceID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type MongodbWhiteListModel struct {
	IpWhitelistName types.String `tfsdk:"ip_whitelist_name"`
	IpList          types.List   `tfsdk:"ip_list"`
}

type MongodbWhiteListsDataSourceModel struct {
	ID         types.String            `tfsdk:"id"`
	InstanceID types.String            `tfsdk:"instance_id"`
	RegionID   types.String            `tfsdk:"region_id"`
	ProjectID  types.String            `tfsdk:"project_id"`
	WhiteLists []MongodbWhiteListModel `tfsdk:"white_lists"`
}
