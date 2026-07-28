package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunElbCertificate(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_certificate." + rnd
	resourceFile := "resource_ctyun__elb_certificate.tf"

	datasourceName := "data.ctyun_elb_certificates." + dnd
	datasourceFile := "datasource_ctyun_elb_certificates.tf"

	name := "certificate_" + utils.GenerateRandomString()
	serverCertificateType := "Server"
	caCertificateType := "Ca"
	certificate := "<<EOF\n-----BEGIN CERTIFICATE-----\nMIICsDCCAhkCFG05oHQGnwdAzD7xikVKLaKpMk+OMA0GCSqGSIb3DQEBCwUAMGUx\nCzAJBgNVBAYTAmNuMQswCQYDVQQIDAJiajELMAkGA1UEBwwCYmoxDDAKBgNVBAoM\nA3R5eTEMMAoGA1UECwwDdHl5MQwwCgYDVQQDDAN0eXkxEjAQBgkqhkiG9w0BCQEW\nA3R5eTAeFw0yMjExMjIxMDM5MTZaFw0zMjExMTkxMDM5MTZaMEUxCzAJBgNVBAYT\nAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRn\naXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDG6p7H\nphrTLthG1FO/YXth8qAnmPUeI9FNRv+U/jtFf7pjRg9VbUPrncxeaZfOqeHFkeFJ\n4NHWSbt0GyGEy4VbQwKFdfYiJ8L1wqk8ked4ffS4PGbJOuJO14dwuELNAcsRvDSR\nb9w82Cx6dszIxMLT0K0L+XRURg/zYAP9tALNCsrsncfuD/G6tC5CZudDMO19Ccy+\nucG/JFrbDnvBTTbZkOX0823XxxYeYPePb31/GXL3rpE1mbXh/ywhPPiDin1ILZEY\n72HdSzTPuatj13rGqXfP4zIYMcOHZEetuB2JFhBdENqy9cbkrjlVCmCt1sgrCSQV\nJIUzg45HOb1ckPRVAgMBAAEwDQYJKoZIhvcNAQELBQADgYEA0jiUp0Lop4KA7DgT\nhQNrJag0BmAaMNYMTbTDK8dkHFInVZEKcO5y64IzRLhI1JrDD+aPMqvG/3NSv2Wy\nfJq1fjbOSXtRF6E2XHUFvdFoB148GsrCWdY8y+ZoHkC0VHqVwJlHHgOiMMcfgb62\nj7EnIoiXbueNr9j/Fr4/wXlCSQ0=\n-----END CERTIFICATE-----\nEOF"
	caCertificate := "<<EOF\n-----BEGIN CERTIFICATE-----\nMIICTDCCAbUCFGUlaAW5nRMvi61kRa2ccyCaEY/lMA0GCSqGSIb3DQEBCwUAMGUx\nCzAJBgNVBAYTAmNuMQswCQYDVQQIDAJiajELMAkGA1UEBwwCYmoxDDAKBgNVBAoM\nA3R5eTEMMAoGA1UECwwDdHl5MQwwCgYDVQQDDAN0eXkxEjAQBgkqhkiG9w0BCQEW\nA3R5eTAeFw0yMjExMjIxMDM4MThaFw0zMjExMTkxMDM4MThaMGUxCzAJBgNVBAYT\nAmNuMQswCQYDVQQIDAJiajELMAkGA1UEBwwCYmoxDDAKBgNVBAoMA3R5eTEMMAoG\nA1UECwwDdHl5MQwwCgYDVQQDDAN0eXkxEjAQBgkqhkiG9w0BCQEWA3R5eTCBnzAN\nBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA5QISwOQWeJ/4B7hxC0pxQn76ebBf1Lt2\nKAjl4g3HNK0KXT4pfwQiaWWbaifXF84zkl+t4eu3SNtKBRwm8aOuWGAmKSkEmdbM\n6P6x24SFiY2bluQoZr+pJfhBJPL+h2v0MaR1ksjW4gXhHrs69EpozpW3Q+OQuSYg\n32IV+pOU/RcCAwEAATANBgkqhkiG9w0BAQsFAAOBgQDFGqzklVQSZh+c7yBP5bZM\n4nr75pXzzkkCb2Xv0HfIKO+rj0W642k84WTWhto8Ihkbl71ZMFKi5Ct3GrU19CAk\nGrQhQT61VOeDD5ZsmTsB9FQwraa39zxN3n7RV83/PZE/72PoN2cFAUbXq7RktIvi\nZJP5PzDVNY+6Yx6WsA213g==\n-----END CERTIFICATE-----\nEOF"
	privateKey := "<<EOF\n-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxuqex6Ya0y7YRtRTv2F7YfKgJ5j1HiPRTUb/lP47RX+6Y0YP\nVW1D653MXmmXzqnhxZHhSeDR1km7dBshhMuFW0MChXX2IifC9cKpPJHneH30uDxm\nyTriTteHcLhCzQHLEbw0kW/cPNgsenbMyMTC09CtC/l0VEYP82AD/bQCzQrK7J3H\n7g/xurQuQmbnQzDtfQnMvrnBvyRa2w57wU022ZDl9PNt18cWHmD3j299fxly966R\nNZm14f8sITz4g4p9SC2RGO9h3Us0z7mrY9d6xql3z+MyGDHDh2RHrbgdiRYQXRDa\nsvXG5K45VQpgrdbIKwkkFSSFM4OORzm9XJD0VQIDAQABAoIBACj76FEcYUSHz1nw\nn3y8Gg7ZTbQ66K4YFSTF7x0EsLOmGIIhykEArVDbh2MggH29NN5fKzrsjm+Ha48F\nlAdnY4elK9zRrC/nX10BiZsIONfzo7td/pORhVVXRPmtjV3t86goze/1SzxiEe/9\nkD4BhF7eDPl9oUFH2jt72fao4zbZmAtx1rGFPKl4+b/4KWrx++DHRyO6xBNH3bit\nW92SwLllwGnnnsCpOuItXc4Qh/CyBO4gv79kn2W9ApTeY5eAAau/8icxES7XvL4c\n3SAOQwBx4iQpp3uBLQFMua8KV92BqMesbOCBybHsJZ5hAMPC44qh5e7NH1oYqsSg\nunokMdUCgYEA9zO+LQZjY7v2OdGN3B4YhEg4hbN9B0m+bOlAYFjOpM2KXQnt0Rhs\n7s4sxrvm576Cvl57ILIv1tjJLdz0UZ0l/XkSlwk7qtTd71r2ZroVuOY0vs470iN5\n/L9kOPZSfQqH0Mmhh+IzSrsNFak5BZxQu0n7LPyQZPDNYfg1cO1Bfh8CgYEAzf7y\nfELFq+rhTwbZ46u2vSQmWNO6UJeQHNaULYPb78uOh4898fsNuHQp9A4qpO8YY/Hz\n8vJkHJtumRYQG7T5h98ABPDoqN9YQQDuqQhCRxFLveDDpig8O0rvOPElEv01iPCh\n1o+ki5m8TAY9P0xgz0DYQ2kiH8hQRboMaTKcVwsCgYAaV+FEWxHsZvNuZe6ALpTe\nQ/QCC4afaDRq1tCNc+lRlrXQBGbbiYbSTBZpd0y8FYlJUDg+275NXvzRbmJ68AxE\nXsqkXc+F/PlJsJ/hgqMd+SpVyxSE6FLvpFXB3D4eJSkkDtiv6mMc66IRVN9Gwcm4\nq8GgoamhmCfK8PCBAEeicQKBgEAle+0l/dgjNDYftAoplqYfc7GFfSdLixzv1QS3\nYu2xPZkJCgkoXIVr5wSQxMbHjZjR511oDbS60h3puOpn2KxuzNq9CjZMFndniuoo\nIDtxL1zZeRNsxBTSqNvae+kF4H3cMQlXga3XGcOyza/AYQUo9C9Jtc6f2h9caDD6\nCaUxAoGBAKHQCtnePuAi82Jfz73jQd0oCDFNeTkCf9w9ZlWyk6dozi5kNijctxW1\nLr4z6e5YTDwiERuf2Oic/RfeWaw+tC2qSdb/SzOQrv/6GE/Rn2DyQIVARJZhkABp\nPZRub8Hzjb/NvIha8lfGGQyaR+NfVkqu1ikRPv+0XCnghY7jb/Kk\n-----END RSA PRIVATE KEY-----\nEOF"

	tfPrivateKey := fmt.Sprintf(`private_key=%s`, privateKey)
	updatedName := "certificate_" + utils.GenerateRandomString()
	description := utils.GenerateRandomString()
	tfDescription := fmt.Sprintf(`description="%s"`, description)

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 server证书验证
			// 1.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, serverCertificateType, certificate, tfPrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", serverCertificateType),
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[resourceName].Primary
						createTime := ds.Attributes["create_time"]
						if utils.IsEmptyOrRfc3339(createTime) {
							return nil
						}
						return fmt.Errorf("time format doesn't match")
					},
					//resource.TestCheckResourceAttr(resourceName, "certificate", certificate),
					//resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
			},
			// import state 1
			{ResourceName: resourceName,
				ImportState: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					return fmt.Sprintf("%s", id), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			// import state 2
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionID := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", id, regionID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			// 1.2 Create 更新
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, serverCertificateType, certificate, tfPrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "type", serverCertificateType),
					//resource.TestCheckResourceAttr(resourceName, "certificate", certificate),
					//resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
			},
			// datasource 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, serverCertificateType, certificate, tfPrivateKey) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("ids=%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "certificates.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "certificates.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "certificates.0.type", serverCertificateType),
				),
			},
			// destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, serverCertificateType, certificate, tfPrivateKey),
				Destroy: true,
			},

			// 2 Ca证书验证
			// 2.1 Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, caCertificateType, caCertificate, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", caCertificateType),
					//resource.TestCheckResourceAttr(resourceName, "certificate", certificate),
					//resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
			},
			// 2.2 update
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, caCertificateType, caCertificate, tfDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "type", caCertificateType),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			// 2.3 destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, caCertificateType, caCertificate, tfDescription),
				Destroy: true,
			},
		},
	})
}
