package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunElbCertificates{}
	_ datasource.DataSourceWithConfigure = &ctyunElbCertificates{}
)

type ctyunElbCertificates struct {
	meta *common.CtyunMetadata
}

func (c *ctyunElbCertificates) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_certificates"
}

func (c *ctyunElbCertificates) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026756/10155416`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "证书ID列表，以,分隔",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "证书名称，以,分隔，必须与ID顺序严格对应",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "证书类型。Ca或Server，以,分隔，必须与ID和name的顺序严格对应",
				Validators: []validator.String{
					stringvalidator.OneOf(business.CertificateTypes...),
				},
			},
			"certificates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "证书列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "证书ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "证书类型: server/ca",
						},
						"private_key": schema.StringAttribute{
							Computed:    true,
							Description: "服务器证书私钥",
						},
						"certificate": schema.StringAttribute{
							Computed:    true,
							Description: "type为Server该字段表示服务器证书公钥Pem内容;type为Ca该字段表示Ca证书Pem内容",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态: ACTIVE / INACTIVE",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunElbCertificates) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunElbCertificatesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	params := &ctelb.CtelbListCertificateRequest{
		RegionID: regionId,
	}
	if config.IDs.ValueString() != "" {
		params.IDs = config.IDs.ValueString()
	}
	if config.Name.ValueString() != "" {
		params.Name = config.Name.ValueString()
	}
	if config.Type.ValueString() != "" {
		params.RawType = config.Type.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListCertificateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 将数据写入到certificate中
	var certificates []CtyunElbCertificatesInfoModel
	for _, certificateItem := range resp.ReturnObj {
		var certificate CtyunElbCertificatesInfoModel
		certificate.RegionID = types.StringValue(certificateItem.RegionID)
		certificate.AzName = types.StringValue(certificateItem.AzName)
		certificate.ProjectID = types.StringValue(certificateItem.ProjectID)
		certificate.ID = types.StringValue(certificateItem.ID)
		certificate.Name = types.StringValue(certificateItem.Name)
		certificate.Description = types.StringValue(certificateItem.Description)
		certificate.Type = types.StringValue(certificateItem.RawType)
		certificate.PrivateKey = types.StringValue(certificateItem.PrivateKey)
		certificate.Certificate = types.StringValue(certificateItem.Certificate)
		certificate.Status = types.StringValue(certificateItem.Status)
		certificate.CreatedTime = types.StringValue(certificateItem.CreatedTime)
		certificate.UpdatedTime = types.StringValue(certificateItem.UpdatedTime)
		certificates = append(certificates, certificate)
	}
	config.Certificates = certificates
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunElbCertificates) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewElbCertificates() datasource.DataSource {
	return &ctyunElbCertificates{}
}

type CtyunElbCertificatesConfig struct {
	RegionID     types.String                    `tfsdk:"region_id"` //资源池ID
	IDs          types.String                    `tfsdk:"ids"`       //证书ID列表，以,分隔
	Name         types.String                    `tfsdk:"name"`      //证书名称，以,分隔，必须与ID顺序严格对应
	Type         types.String                    `tfsdk:"type"`      //证书类型。Ca或Server，以,分隔，必须与ID和name的顺序严格对应
	Certificates []CtyunElbCertificatesInfoModel `tfsdk:"certificates"`
}

type CtyunElbCertificatesInfoModel struct {
	RegionID    types.String `tfsdk:"region_id"`    //资源池ID
	AzName      types.String `tfsdk:"az_name"`      //可用区名称
	ProjectID   types.String `tfsdk:"project_id"`   //项目ID
	ID          types.String `tfsdk:"id"`           //证书ID
	Name        types.String `tfsdk:"name"`         //名称
	Description types.String `tfsdk:"description"`  //描述
	Type        types.String `tfsdk:"type"`         //证书类型: server / ca
	PrivateKey  types.String `tfsdk:"private_key"`  //服务器证书私钥
	Certificate types.String `tfsdk:"certificate"`  //type为Server该字段表示服务器证书公钥Pem内容;type为Ca该字段表示Ca证书Pem内容
	Status      types.String `tfsdk:"status"`       //状态: ACTIVE / INACTIVE
	CreatedTime types.String `tfsdk:"created_time"` //创建时间，为UTC格式
	UpdatedTime types.String `tfsdk:"updated_time"` //更新时间，为UTC格式
}
