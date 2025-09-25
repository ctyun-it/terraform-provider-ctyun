package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	YamlError = "不满足YAML格式"
)

type validatorYaml struct {
	mustExistKey []string
}

func Yaml(key ...string) validator.String {
	return &validatorYaml{mustExistKey: key}
}

func (v validatorYaml) Description(_ context.Context) string {
	return YamlError
}

func (v validatorYaml) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorYaml) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	value := request.ConfigValue.ValueString()
	// 检查空字符串
	if strings.TrimSpace(value) == "" {
		response.Diagnostics.AddError(YamlError, "不能为空")
		return
	}

	// 尝试解析 YAML
	var parsed map[string]interface{}
	decoder := yaml.NewDecoder(strings.NewReader(value))
	err := decoder.Decode(&parsed)
	if err != nil {
		response.Diagnostics.AddError(YamlError, fmt.Sprintf("语法错误: %v", err))
		return
	}

	for _, key := range v.mustExistKey {
		if _, ok := parsed[key]; !ok {
			response.Diagnostics.AddError(YamlError, fmt.Sprintf("必须存在key: %s", key))
			return
		}
	}

	// 检查是否存在未解析的内容（多文档场景）
	if decoder.Decode(&struct{}{}) == nil {
		{
			response.Diagnostics.AddError(YamlError, "存在多个 YAML 文档，只允许单个文档")
			return
		}
	}
	return
}
