package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbCertificate{}
	_ resource.ResourceWithConfigure   = &CtyunElbCertificate{}
	_ resource.ResourceWithImportState = &CtyunElbCertificate{}
)

type CtyunElbCertificate struct {
	meta *common.CtyunMetadata
}

func NewCtyunElbCertificate() resource.Resource {
	return &CtyunElbCertificate{}
}

func (c *CtyunElbCertificate) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbCertificate) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunElbCertificate) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_certificate"
}

func (c *CtyunElbCertificate) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10155416**`,
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
				Description: "支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
					validator2.Desc(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "证书类型。取值范围：Server（服务器证书）、Ca（Ca证书）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.CertificateTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Description: "服务器证书私钥，type=Server服务器证书此字段必填",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.CertificateTypeServer),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("type"),
						types.StringValue(business.CertificateTypeCA),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate": schema.StringAttribute{
				Required:    true,
				Description: "type为Server 该字段表示服务器证书公钥Pem内容;type为Ca 该字段表示Ca证书Pem内容",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "证书ID",
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
		},
	}
}

func (c *CtyunElbCertificate) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunElbCertificateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,检查证书有效性
	isValid, err := c.checkBeforeCreateCertificate(ctx, plan)
	if !isValid {
		err = fmt.Errorf("服务器证书/Ca证书无效，证书创建失败")
		return
	}
	err = c.createCertificate(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeCertificate(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbCertificate) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbCertificateConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeCertificate(ctx, &state)
	if err != nil {
		// 有待确定
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunElbCertificate) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取tf文件中配置
	var plan CtyunElbCertificateConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbCertificateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}
	//更新证书名字和描述
	err = c.updateElbCertificate(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步本地信息
	err = c.getAndMergeCertificate(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbCertificate) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbCertificateConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	//调用证书删除接口
	params := &ctelb.CtelbDeleteCertificateRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		CertificateID: state.ID.ValueString(),
	}

	// 同步接口，无需轮询
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteCertificateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
}

func (c *CtyunElbCertificate) checkBeforeCreateCertificate(ctx context.Context, plan CtyunElbCertificateConfig) (isValid bool, err error) {

	if plan.RegionID.ValueString() == "" {
		err = fmt.Errorf("region id不能为空！")
		return
	}

	// 若type=server类型，调用检查server证书接口
	if plan.Type.ValueString() == business.CertificateTypeServer {

		// 判断privateKey和certificate是否有值
		if plan.Certificate.ValueString() == "" || plan.PrivateKey.ValueString() == "" {
			err = fmt.Errorf("当证书类型为Server(服务器证书)时，服务器证书私钥和证书不能为空!")
			return
		}

		params := &ctelb.CtelbCheckServerCertRequest{
			Certificate: plan.Certificate.ValueString(),
			PrivateKey:  plan.PrivateKey.ValueString(),
		}
		resp, err2 := c.meta.Apis.SdkCtElbApis.CtelbCheckServerCertApi.Do(ctx, c.meta.SdkCredential, params)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}
		isValid = *resp.ReturnObj.IsValid
		return isValid, nil
	} else if plan.Type.ValueString() == business.CertificateTypeCA {
		// 若type=ca 类型， 调用检查ca证书接口

		if plan.Certificate.ValueString() == "" {
			err = fmt.Errorf("当证书类型为Ca(Ca证书)时，证书不能为空")
			return
		}
		params := &ctelb.CtelbCheckCaCertRequest{
			Certificate: plan.Certificate.ValueString(),
		}
		resp, err2 := c.meta.Apis.SdkCtElbApis.CtelbCheckCaCertApi.Do(ctx, c.meta.SdkCredential, params)
		if err2 != nil {
			return false, err2
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}
		isValid = *resp.ReturnObj.IsValid
		return isValid, nil
	} else {
		err = fmt.Errorf("证书类型有误！，当前证书类型为：" + plan.Type.ValueString())
		return false, err
	}
}

func (c *CtyunElbCertificate) createCertificate(ctx context.Context, config *CtyunElbCertificateConfig) (err error) {
	params := &ctelb.CtelbCreateCertificateRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		Name:        config.Name.ValueString(),
		RawType:     config.Type.ValueString(),
		Certificate: config.Certificate.ValueString(),
	}
	if config.Description.ValueString() != "" {
		params.Description = config.Description.ValueString()
	}
	if config.Type.ValueString() == business.CertificateTypeServer {
		params.PrivateKey = config.PrivateKey.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateCertificateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	if resp.ReturnObj.ID != "" {
		config.ID = types.StringValue(resp.ReturnObj.ID)
	} else {
		err = fmt.Errorf("id为空")
		return
	}
	return
}

func (c *CtyunElbCertificate) getAndMergeCertificate(ctx context.Context, config *CtyunElbCertificateConfig) (err error) {
	params := &ctelb.CtelbShowCertificateRequest{
		RegionID:      config.RegionID.ValueString(),
		CertificateID: config.ID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowCertificateApi.Do(ctx, c.meta.SdkCredential, params)
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
	// 解析证书详情
	config.Status = types.StringValue(returnObj.Status)
	config.CreatedTime = types.StringValue(returnObj.CreatedTime)
	config.UpdatedTime = types.StringValue(returnObj.UpdatedTime)
	config.Name = types.StringValue(returnObj.Name)
	config.Description = types.StringValue(returnObj.Description)
	return
}

func (c *CtyunElbCertificate) updateElbCertificate(ctx context.Context, state *CtyunElbCertificateConfig, plan *CtyunElbCertificateConfig) (err error) {
	params := &ctelb.CtelbUpdateCertificateRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		CertificateID: state.ID.ValueString(),
	}
	if plan.ProjectID.ValueString() != "" && plan.ProjectID.ValueString() != state.ProjectID.ValueString() {
		params.ProjectID = plan.ProjectID.ValueString()
	}
	if plan.Name.ValueString() != "" && plan.Name.ValueString() != state.Name.ValueString() {
		params.Name = plan.Name.ValueString()
	}
	if plan.Description.ValueString() != "" && plan.Description.ValueString() != state.Description.ValueString() {
		params.Description = plan.Description.ValueString()
	}
	// 若projectID, 证书名称和证书描述为空的话，不必更新直接返回
	if params.ProjectID == "" && params.Name == "" && params.Description == "" {
		return
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateCertificateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

type CtyunElbCertificateConfig struct {
	RegionID    types.String `tfsdk:"region_id"`    //资源池ID
	Name        types.String `tfsdk:"name"`         //	唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description types.String `tfsdk:"description"`  //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·！@#￥%……&*（） —— -+={}\|《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128
	Type        types.String `tfsdk:"type"`         //证书类型。取值范围：Server（服务器证书）、Ca（Ca证书）
	PrivateKey  types.String `tfsdk:"private_key"`  //服务器证书私钥，服务器证书此字段必填
	Certificate types.String `tfsdk:"certificate"`  //type为Server该字段表示服务器证书公钥Pem内容;type为Ca该字段表示Ca证书Pem内容
	ID          types.String `tfsdk:"id"`           //证书ID
	Status      types.String `tfsdk:"status"`       //状态: ACTIVE / INACTIVE
	CreatedTime types.String `tfsdk:"created_time"` //创建时间，为UTC格式
	UpdatedTime types.String `tfsdk:"updated_time"` //更新时间，为UTC格式
	AzName      types.String `tfsdk:"az_name"`      //可用区名称
	ProjectID   types.String `tfsdk:"project_id"`   //项目ID
}
