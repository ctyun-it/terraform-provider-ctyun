package pgsql

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

const (
	EndpointNamePgSql = "pgsql"
	UrlPordPgSql      = "pgsql-global.ctapi.ctyun.cn"
	UrlTestPgSql      = ""
)

var EndpointPgSqlTest = ctyunsdk.Endpoint{
	EndpointName: EndpointNamePgSql,
	Url:          UrlTestPgSql,
}

var EndPointPgSqlProd = ctyunsdk.Endpoint{
	EndpointName: EndpointNamePgSql,
	Url:          UrlPordPgSql,
}
