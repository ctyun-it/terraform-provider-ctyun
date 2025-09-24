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
	certificate := "<<EOF\n-----BEGIN CERTIFICATE-----\nMIIDEzCCAfsCFBS3/LXUQVm1sPtyKN06Fiac89I6MA0GCSqGSIb3DQEBCwUAMEYx\nCzAJBgNVBAYTAmNuMQswCQYDVQQIDAJqczELMAkGA1UEBwwCc3oxDjAMBgNVBAoM\nBWN0eXVuMQ0wCwYDVQQLDARpYWFzMB4XDTI1MDUxNDEzMTIwMloXDTI2MDUxNDEz\nMTIwMlowRjELMAkGA1UEBhMCY24xCzAJBgNVBAgMAmpzMQswCQYDVQQHDAJwZDEO\nMAwGA1UECgwFY3R5dW4xDTALBgNVBAsMBGlhYXMwggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQDYg+VEnKCezCRHNScxGPprYN4rsvJh0NtMczQTkb4xZp7z\nV69ua9yU2tFiGXaP1yzQu6nNcAfhdbO6kAvuQthwIQGy1+tVPJsC2SMsRhOQrvvg\nbZgCcajXL/+L1+4KTohKnLLyke3Bx2z9CqLLJ1xggx8r7hXSev4DFO2oq1NUxt/f\n354BK3aTdFVsH1QgiOi+RwK0VXv0Wr7QgRckCqq5iCSg2O2/1gYXMIQDKhcyvXCs\nq3MwgqWI27YioTOgD68xxdQs2vhHzqHu9KavfWyoXadfAsX5aFlBM6aIhreQiEWb\nFfh+lDrDuU5u9AMDzo4WOvglhq1GHt2EBKYcfW2nAgMBAAEwDQYJKoZIhvcNAQEL\nBQADggEBAECLpCZg60zVb5bJLMS+OPiBX0IkBpzU8zb/T2eNChwvFC2kCvqJUdVC\nnoD6No2w9ZG8dzDxM+8mhY9H5N6InBr9wtzGME1B9+cBN/kjCLO5blwVj3EgDIC6\nzlYlAdZYzaOBbgp+9YnE4M2i0lkPGnxiSN1g3m9/9LD4D1hOIpBQ4J/4iHw9RzlN\nvAcY1HkS3JAYDg+4/VAL8WjaB829YpamPuxpRRxCSfkl/syZfgSaTHxFAUaCqaR9\nQsdkrXjO98LcFBQ9TtgjWXYhy21e5j7ObRijWPtmLHDPGQdUvTzNnpJ83YzjoRR2\nCpofzsknFWuwRsfefk7rcEZXPDUwc1Y=\n-----END CERTIFICATE-----\nEOF"
	caCertificate := "<<EOF\n-----BEGIN CERTIFICATE-----\nMIIDbTCCAlWgAwIBAgIUSqUBZkCUWOTP8DbL72F2B8HVixIwDQYJKoZIhvcNAQEL\nBQAwRjELMAkGA1UEBhMCY24xCzAJBgNVBAgMAmpzMQswCQYDVQQHDAJzejEOMAwG\nA1UECgwFY3R5dW4xDTALBgNVBAsMBGlhYXMwHhcNMjUwNTE0MTMwNzI1WhcNMjYw\nNTE0MTMwNzI1WjBGMQswCQYDVQQGEwJjbjELMAkGA1UECAwCanMxCzAJBgNVBAcM\nAnN6MQ4wDAYDVQQKDAVjdHl1bjENMAsGA1UECwwEaWFhczCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAKtWVxAhtLVKZZ8sKGJiW+vA0mVLZOqO0HkUrKi9\nHrLV2vuyr+yLblaltNhW7IZEcW1UKR9FcfQ0kdW+cfZOGkca6b3YI4gyZ0so83uG\n4mxs/FPs372sUsA9/VxmbN0cuDQXEloctMeERMIXXM4wGhD7X9IiRDqhQK7nW4/b\nsfpJU9PPss++55GaUitS6cRjMlPHxnz3lApXJf41D/1pI0FWIQMEi1pwXQVfZEHe\ny6Li+JChXOwPIsTzbqH3j1NbVJljjRCHQbRt6PjwOh+BNoPRuJdKlcKCLFECHHjg\nN3BYkbSTSbD4GKDBmT1M0kMf5E0+VdO/rB4/IAqURyV784ECAwEAAaNTMFEwHQYD\nVR0OBBYEFCmOin4Dc9aEummONDnRCgTsvMFLMB8GA1UdIwQYMBaAFCmOin4Dc9aE\nummONDnRCgTsvMFLMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB\nAI6kCWpofIAxJFrDZpkUmFJpBbWWEGT4ssPe6W977U6+kJ2uaaD7dB+9UwocqOAv\nCAImbSe2NkyStDFmxFzO+jpZ02WCO+GNnuRMEvuidkm/YC0ZN2Sxd4BEMt8OYfI2\nrD9AyMfSRmH6oTwKBAFqsdQ18etCEqp2UqjiiQFQ1/hPZQLhYxd6uzKObkyFFZII\nODOww2jt0wvcLUKkeja3HfafMUocDDsIW6WQthxH7MXyJuuxAsIcMljW3/M3KXQV\nbhBV0KWnOIDOo9X6Eg2QQWq6C26mYHgdtkYWK+zWQqTF9BGbZwZvkczJq9hMCaVz\nEYJghKoiYZor/LEES5R/4nI=\n-----END CERTIFICATE-----\nEOF"
	privateKey := "<<EOF\n-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA2IPlRJygnswkRzUnMRj6a2DeK7LyYdDbTHM0E5G+MWae81ev\nbmvclNrRYhl2j9cs0LupzXAH4XWzupAL7kLYcCEBstfrVTybAtkjLEYTkK774G2Y\nAnGo1y//i9fuCk6ISpyy8pHtwcds/QqiyydcYIMfK+4V0nr+AxTtqKtTVMbf39+e\nASt2k3RVbB9UIIjovkcCtFV79Fq+0IEXJAqquYgkoNjtv9YGFzCEAyoXMr1wrKtz\nMIKliNu2IqEzoA+vMcXULNr4R86h7vSmr31sqF2nXwLF+WhZQTOmiIa3kIhFmxX4\nfpQ6w7lObvQDA86OFjr4JYatRh7dhASmHH1tpwIDAQABAoIBACxbCPL/a2BczenQ\nl9DVHyg0Vg95v3IOiX6l8zs3FlGkhlev1P64Fh5gnZiNQt7A5Ct25phxpQupQ0Kr\nE+ahxxlhiamL1mhF3DO/LBd727I+Rtd5XDTT+BTtflq1x8xQvlcatwY+owiZxFgq\nf6NC0wfqlM45HjtaewNCnFV7k7zACLwTK0Eg0m9H6grh+d/wj/toMv+1fX5bPk3I\nW7dD0kWqAWoT+siWKYxQBtOfj7QmzR9OEU8M9Ofn7B/Btg/PgDgEkzgugvN+Kg+j\nlcQ8THL1mfpm8QuE00oP4yXK0gR3p6nxbyiSl7FzdA6d5OSDDGO83kvYi+Hnpj4/\nWMuhQIECgYEA9bQjy1Wq+74OJ7j4or0oJodf5+obckBuJuaP0vdJds6Tkkvyo8ni\nN7dzE+G0ADDc6AhjvnSUSqumCWiqPgHy8ggPZBK6uGZyCPUML6MVpcs0zQP9U0TK\n52Fl6snatpUSRmbgIXfwoIVhLcoLBVZ+VRi6k2zrM5i2leOR+qMkgGECgYEA4Zag\nuMVHlz3/JpdtOmih20UMID64Pd/XdDzzGzmdcOfI7LW8XLZdBJfEszd+8zQ+ZdKN\nokFl/DJZaBbw3H7tG5gi6ISQ1ylgnDivWaZa38e6DZIb2AV2ry2nHiOdoW38iNc3\nYenQVGrm88vPB8ZoRMJKxmjGWBwpAtVPkNZRywcCgYAjgASP82/B8cLf0GT0Nnpf\nnr2np3aRumdv7W4oYBxYMx61S43iRmAOBs50t+6m/Gheih/HO9bVNxnFUD8QuH1e\npPV8UwxxajPdwXIiS+Dr1IFrgcEPT+g1KChulP2p1b4PN/v4OfklaEATOpb/DXEx\nYGJMvLb6/HydCYVk/j9e4QKBgF8fuAgR0B1bXCL2yW1Ov6mJVRKnv+L/Z4exR8Xc\nJqJ8aoI+I/5oO5L7OEtmBFcjWYhwH58iupF83ayBt/ESaIMxRM8NT1IbNYzKbk58\nUmd4feDJEoqlyyUVOBoCZK60hT9imQzlnh5qsZZuA4/Avyj3ULaXacOvpC94qUNm\nzDizAoGBAOUg9tErHeNLOL+foJCy+ib8bQlwFpaizNMNKDwcdVihtS0ROUUp87FG\nMKVvp+2lXRi6av0R3rnpETmYx/JVictJmt2eDK6LBK/iNW4KeitJHYGwrx8OjUOM\noVTVU5+Z0VYc7QPyGGCqzcuySZZGaI0MwKOe7jshCzTtgwFzfLHT\n-----END RSA PRIVATE KEY-----\nEOF"

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
					//resource.TestCheckResourceAttr(resourceName, "certificate", certificate),
					//resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
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
