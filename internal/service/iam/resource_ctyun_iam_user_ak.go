package iam

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctiam"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CtyunIamUserAk{}
	_ resource.ResourceWithConfigure   = &CtyunIamUserAk{}
	_ resource.ResourceWithImportState = &CtyunIamUserAk{}
)

func NewCtyunIamUserAk() resource.Resource {
	return &CtyunIamUserAk{}
}

type CtyunIamUserAk struct {
	meta       *common.CtyunMetadata
	iamService *business.IamService
}

func (c *CtyunIamUserAk) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam_user_ak"
}

func (c *CtyunIamUserAk) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10345725/10355289`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "用户ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"ak": schema.StringAttribute{
				Computed:    true,
				Description: "用户AK",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sk": schema.StringAttribute{
				Computed:    true,
				Description: "用户SK",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "密钥状态",
				Default:     booldefault.StaticBool(true),
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunIamUserAk) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunIamUserAkConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	ak, sk, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.AK = types.StringValue(ak)
	plan.SK = types.StringValue(sk)

	// 设置状态
	if !plan.Enabled.ValueBool() {
		err = c.disable(ctx, plan)
		if err != nil {
			return
		}
	}

	// 查询
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *CtyunIamUserAk) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIamUserAkConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunIamUserAk) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIamUserAkConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan CtyunIamUserAkConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (c *CtyunIamUserAk) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunIamUserAkConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *CtyunIamUserAk) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ak],[userID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunIamUserAkConfig
	var ak, userID string
	err = terraform_extend.Split(request.ID, &ak, &userID)
	if err != nil {
		return
	}

	cfg.AK = types.StringValue(ak)
	cfg.UserID = types.StringValue(userID)

	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *CtyunIamUserAk) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.iamService = business.NewIamService(meta)
}

// create 创建AK
func (c *CtyunIamUserAk) create(ctx context.Context, plan CtyunIamUserAkConfig) (ak, sk string, err error) {
	params := &ctiam.CtiamCreateAkRequest{
		UserId: plan.UserID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamCreateAkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ak, sk = utils.SecString(resp.ReturnObj.AppId), utils.SecString(resp.ReturnObj.AppKey)
	return
}

// update 更新AK
func (c *CtyunIamUserAk) update(ctx context.Context, plan, state CtyunIamUserAkConfig) (err error) {
	if plan.Enabled.Equal(state.Enabled) {
		return
	}
	if plan.Enabled.ValueBool() {
		err = c.enable(ctx, plan)
	} else {
		err = c.disable(ctx, plan)
	}
	return
}

// delete 删除AK
func (c *CtyunIamUserAk) delete(ctx context.Context, state CtyunIamUserAkConfig) (err error) {
	if state.Enabled.ValueBool() {
		err = c.disable(ctx, state)
		if err != nil {
			return
		}
	}
	err = c.remove(ctx, state)
	if err != nil {
		return
	}
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	return
}

// remove 移入回收站
func (c *CtyunIamUserAk) remove(ctx context.Context, plan CtyunIamUserAkConfig) (err error) {
	params := &ctiam.CtiamAkToRecycleBinRequest{
		UserId: plan.UserID.ValueString(),
		Ak:     plan.AK.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtIamApis.CtiamAkToRecycleBinApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// destroy 从回收站删除
func (c *CtyunIamUserAk) destroy(ctx context.Context, plan CtyunIamUserAkConfig) (err error) {
	params := &ctiam.CtiamDeleteAkRequest{
		UserId: plan.UserID.ValueString(),
		Ak:     plan.AK.ValueString(),
	}

	resp, err := c.meta.Apis.SdkCtIamApis.CtiamDeleteAkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// enable 启用AK
func (c *CtyunIamUserAk) enable(ctx context.Context, plan CtyunIamUserAkConfig) (err error) {
	params := &ctiam.CtiamOnlineAkRequest{
		UserId: plan.UserID.ValueString(),
		Ak:     plan.AK.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamOnlineAkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// disable 禁用AK
func (c *CtyunIamUserAk) disable(ctx context.Context, plan CtyunIamUserAkConfig) (err error) {
	params := &ctiam.CtiamPauseAkRequest{
		UserId: plan.UserID.ValueString(),
		Ak:     plan.AK.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtIamApis.CtiamPauseAkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if *resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	}
	return
}

// query 查询ak
func (c *CtyunIamUserAk) query(ctx context.Context, plan CtyunIamUserAkConfig) (res *ctiam.CtiamQueryAkReturnObjAccessKeyUserListAccessKeyListResponse, err error) {
	aks, err := c.iamService.QueryAkList(ctx, plan.UserID.ValueString())
	if err != nil {
		return
	}
	for _, a := range aks {
		if utils.SecString(a.AccessKey) == plan.AK.ValueString() {
			return a, nil
		}
	}
	err = common.ResourceNotExistError
	return
}

// getAndMerge 查询并merge
func (c *CtyunIamUserAk) getAndMerge(ctx context.Context, cfg *CtyunIamUserAkConfig) (err error) {
	resp, err := c.query(ctx, *cfg)
	if err != nil {
		return
	}
	cfg.CreateTime = types.StringValue(utils.FromUnixToUTC(resp.CreatedTime))
	cfg.Enabled = types.BoolValue(map[string]bool{business.AkEnabled: true, business.AkDisabled: false}[utils.SecString(resp.Status)])
	sk, err := c.iamService.DecryptSK(utils.SecString(resp.SecretKey), c.meta.SdkCredential.GetAccessKey())
	if err != nil {
		return err
	}
	cfg.SK = types.StringValue(sk)
	cfg.ID = types.StringValue(fmt.Sprintf("%s,%s", cfg.AK.ValueString(), cfg.UserID.ValueString()))
	return
}

type CtyunIamUserAkConfig struct {
	ID         types.String `tfsdk:"id"`
	UserID     types.String `tfsdk:"user_id"`
	AK         types.String `tfsdk:"ak"`
	SK         types.String `tfsdk:"sk"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	CreateTime types.String `tfsdk:"create_time"`
}

// decrypt 解密SK
func (c *CtyunIamUserAk) decrypt(secretSK, ak string) (decrypted string, err error) {
	decodedTxt, err := hex.DecodeString(secretSK)
	if err != nil {
		return
	}

	decodedKey, err := hex.DecodeString(ak)
	if err != nil {
		return
	}

	decrypted, err = utils.Decrypt(decodedTxt, decodedKey)
	if err != nil {
		return
	}
	return
}
