package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	IdpTypeVirtual = "virtual"
	IdpTypeIam     = "iam"

	IdpProtocolSaml = "saml"
	IdpProtocolOidc = "oidc"
)

var IdpTypes = []string{
	IdpTypeVirtual,
	IdpTypeIam,
}

var IdpProtocols = []string{
	IdpProtocolSaml,
	IdpProtocolOidc,
}

const (
	IdpTypeMapScene1 = iota
)

const (
	IdpProtocolMapScene1 = iota
)

var IdpTypeMap = utils.Must(
	[]any{
		IdpTypeVirtual,
		IdpTypeIam,
	},
	map[utils.Scene][]any{
		IdpTypeMapScene1: {
			0,
			1,
		},
	},
)

var IdpProtocolMap = utils.Must(
	[]any{
		IdpProtocolSaml,
		IdpProtocolOidc,
	},
	map[utils.Scene][]any{
		IdpProtocolMapScene1: {
			0,
			1,
		},
	},
)
