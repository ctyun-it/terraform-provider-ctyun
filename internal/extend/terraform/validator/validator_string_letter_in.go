package validator

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringLetterIn(letters ...int32) validator.String {
	str := ""
	mapping := make(map[int32]struct{})
	for _, char := range letters {
		mapping[char] = struct{}{}
		str = str + string(char)
	}
	return &validatorStringLetterIn{
		letter: mapping,
		str:    str,
	}
}

type validatorStringLetterIn struct {
	str    string
	letter map[int32]struct{}
}

func (v validatorStringLetterIn) Description(_ context.Context) string {
	return v.getError()
}

func (v validatorStringLetterIn) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validatorStringLetterIn) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	for _, c := range request.ConfigValue.ValueString() {
		_, ok := v.letter[c]
		if !ok {
			getError := v.getError()
			response.Diagnostics.AddError(getError, getError)
			return
		}
	}
}

func (v validatorStringLetterIn) getError() string {
	return "字符串字符不在所在范围内：" + v.str
}
