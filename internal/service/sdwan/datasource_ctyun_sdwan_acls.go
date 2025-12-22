package sdwan

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sdwan"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &CtyunSdwanAcls{}
)

func NewCtyunSdwanAcls() datasource.DataSource {
	return &CtyunSdwanAcls{}
}

type CtyunSdwanAcls struct {
	meta *common.CtyunMetadata
}

type CtyunSdwanAclsConfig struct {
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Acls      []AclInfo    `tfsdk:"acls"`
	ID        types.String `tfsdk:"id"`
}

type AclInfo struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (c *CtyunSdwanAcls) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdwan_acls"
}

func (c *CtyunSdwanAcls) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10035786/10035852`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "访问控制名称（模糊查询）",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"acls": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ACL ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "ACL名称",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunSdwanAcls) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunSdwanAcls) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSdwanAclsConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &sdwan.SdwanGetSdwanAclRequest{
		PageNo:   1,
		PageSize: 1000,
	}

	if !plan.ProjectID.IsNull() {
		request.ProjectID = plan.ProjectID.ValueStringPointer()
	}

	if !plan.Name.IsNull() {
		request.AclName = plan.Name.ValueStringPointer()
	}

	response, err := c.meta.Apis.SdkSdwanApis.SdwanGetSdwanAclApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		return
	} else if response.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *response.Message)
		return
	} else if response.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var acls []AclInfo
	for _, acl := range response.ReturnObj.Result {
		acls = append(acls, AclInfo{
			ID:   types.StringValue(*acl.AclID),
			Name: types.StringValue(*acl.Name),
		})
	}

	plan.Acls = acls
	plan.ID = types.StringValue("sdwan_acls")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
