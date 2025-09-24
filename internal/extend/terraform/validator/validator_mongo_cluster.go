package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

const (
	MongodbClusterError = "只有mongodb类型为集群版，该字段才可以填写"
)

type validatorMongodbClusterField struct {
}

func (v validatorMongodbClusterField) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// 获取prod_id，确认Cluster字段已经包含其中，如果字段为空则返回报错
	config := request.Config
	prodID, err := getStringAttribute(config, path.Root("prod_id"))
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if strings.Contains(prodID.ValueString(), "Cluster") {

	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	if value == "" {
		return
	}
	pattern := `^eip-[a-z0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		response.Diagnostics.AddError(IpError, IpError)
		return
	}
}

func MongodbClusterFieldValidate() validator.String {
	return &validatorMongodbClusterField{}
}

func (v validatorMongodbClusterField) Description(ctx context.Context) string {
	return EipError
}

func (v validatorMongodbClusterField) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func getStringAttribute(config tfsdk.Config, attrPath path.Path) (types.String, error) {
	var value types.String
	diags := config.GetAttribute(context.Background(), attrPath, &value)
	if diags.HasError() {
		return types.StringValue(""), fmt.Errorf("failed to get attribute: %v", diags.Errors())
	}
	return value, nil
}
