// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls_test

import (
	"os"
	"testing"

	"github.com/searKing/golang/go/crypto/tls"
	testing_ "github.com/searKing/golang/go/testing"
)

// This code is borrowed from https://github.com/ory/x/blob/master/tlsx/cert_test.go
var (
	certFixture = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVFRENDQXZpZ0F3SUJBZ0lKQU5mK0lUMU1HaHhCTUEwR0NTcUdTSWI` +
		`zRFFFQkN3VUFNSUdaTVFzd0NRWUQKVlFRR0V3SlZVekVMTUFrR0ExVUVDQXdDUTBFeEVqQVFCZ05WQkFjTUNWQmhiRzhnUVd4MGJ6RWlNQ0FHQ` +
		`TFVRQpDZ3daVDI1bFEyOXVZMlZ5YmlCYmRHVnpkQ0J3ZFhKd2IzTmxYVEVjTUJvR0ExVUVBd3dUYjI1bFkyOXVZMlZ5CmJpMTBaWE4wTG1OdmJ` +
		`URW5NQ1VHQ1NxR1NJYjNEUUVKQVJZWVpuSmxaR1Z5YVdOQVkyOXVaV052Ym1ObGNtNHUKWTI5dE1CNFhEVEU0TURnd016RTJNakUwT0ZvWERUR` +
		`TVNVEl4TmpFMk1qRTBPRm93Z1lReEN6QUpCZ05WQkFZVApBbFZUTVFzd0NRWURWUVFJREFKRFFURVNNQkFHQTFVRUJ3d0pVR0ZzYnlCQmJIUnZ` +
		`NU0l3SUFZRFZRUUxEQmxQCmJtVkRiMjVqWlhKdUlGdDBaWE4wSUhCMWNuQnZjMlZkTVRBd0xnWURWUVFERENkaGNHa3RjMlZ5ZG1salpTMXcKY` +
		`205NGFXVmtMbTl1WldOdmJtTmxjbTR0ZEdWemRDNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQgpEd0F3Z2dFS0FvSUJBUURXVzF` +
		`KQnZweC9vZkYwei80QnkrYmdBcCtoYnlxblVsQ2FnYmlneE9QTHY3aUg4TSt1CjNENkRlSVkzQzdkV0thTjRnYXZHd1MvN3I0UWxXSWdvK09NR` +
		`HQ1M25OZDVvakwvNWY5R1E0ZGRObW53b25EeEYKVThrd1lMWURMTkJIQzJqMzFBNVNueHo0S1NkVE03Rmc0OFBJeTNBaWFGMkhEcURZVlJpWkV` +
		`ackl4U3JTSmFKZgp1WGVCSUVBcFBpUG1IOURObGw2VVo3ODZvZitJWWVLV2VuY0MvbGpPaGlJSnJWL3NEZTc2QVFjdXY5T29XaUdiCklGVFMyW` +
		`ExSRGF0YzByQXhWdlFiTnMzeWlFYjh3UzBaR0F4cTBuZk9pMGZkYVBIODdFc25MdkpqWk5PcXIvTVMKSW5BYmN2ZmlwckxxaEdLQTVIN2hKVGZ` +
		`EcFJ6WWxBcm5maTJMQWdNQkFBR2piakJzTUFrR0ExVWRFd1FDTUFBdwpDd1lEVlIwUEJBUURBZ1hnTUZJR0ExVWRFUVJMTUVtQ0htOWhkR2hyW` +
		`ldWd1pYSXViMjVsWTI5dVkyVnliaTEwClpYTjBMbU52YllJbllYQnBMWE5sY25acFkyVXRjSEp2ZUdsbFpDNXZibVZqYjI1alpYSnVMWFJsYzN` +
		`RdVkyOXQKTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFCMVBibCtSbW50RW9jbHlqWXpzeWtLb2lYczNwYTgzQ2dEWjZwQwpncnY0TFF4U29FZ` +
		`kowNGY4YkQ0SUlZRkdDWmZWTkcwVnBFWHJObGs2VWJzVmRUQUJ0cUNndUpUV3dER1VBaDZYCjNiRmhyWm5QZXhzLy9Rd2dEQWRxSWYwRWd3Y0R` +
		`VRzc2R0lkZms3MGUxWnV4Y2h4ZDhVQkNwQUlkZVUwOHZWa3kKNFBXdjJLNGFENEZqQ2hLeENONWtoTjUwRk1QY2FJK3hWZ2Q0N3RQaFZOOWxRa` +
		`W9HRENoc1Q1dkFSazdiYS9jZQowUTlOV2RpTWZMRWdMZGNCb2JaS0Z0RnJsS3R5ek9nRGpMdlh2TFFzL3MybWVyU0k5Zmt3b09CRVArN2o3Wm5` +
		`zCkFqeTlNZmh3cWJUcFc3S3BDU0ZhMFZULzJ1OTVaUmNQdnJYbGRLUnlnQjRXdUFScgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`
	keyFixture = `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBMWx0U1FiNmNmNkh4ZE0vK0Fjdm00QUtm` +
		`b1c4cXAxSlFtb0c0b01Uank3KzRoL0RQCnJ0dytnM2lHTnd1M1ZpbWplSUdyeHNFdis2K0VKVmlJS1BqakE3ZWQ1elhlYUl5LytYL1JrT0hYVF` +
		`pwOEtKdzgKUlZQSk1HQzJBeXpRUnd0bzk5UU9VcDhjK0NrblV6T3hZT1BEeU10d0ltaGRodzZnMkZVWW1SR2F5TVVxMGlXaQpYN2wzZ1NCQUtU` +
		`NGo1aC9RelpaZWxHZS9PcUgvaUdIaWxucDNBdjVZem9ZaUNhMWY3QTN1K2dFSExyL1RxRm9oCm15QlUwdGx5MFEyclhOS3dNVmIwR3piTjhvaE` +
		`cvTUV0R1JnTWF0SjN6b3RIM1dqeC9PeExKeTd5WTJUVHFxL3oKRWlKd0czTDM0cWF5Nm9SaWdPUis0U1UzdzZVYzJKUUs1MzR0aXdJREFRQUJB` +
		`b0lCQVFET2xyRE9RQ0NnT2JsMQo5VWMrLy84QkFrWksxZExyODc5UFNacGhCNkRycTFqeld6a3RzNEprUHZKTGR2VTVDMlJMTGQ0WjdmS0t4UH` +
		`U4CjZuZy8xSzhsMC85UTZHL3puME1kK1B4R2dBSjYvbHFPNFJTTlZGVGdWVFRXRm9pZEQvZ1ljYjFrRDRsaCtuZTIKRG1uemtWQU40MU90Tlp4` +
		`K0g3RVJEZUpwRTdoenFSOEhodnhxZU82Z25CMXJkZ3JRSE9MV1lSdmM1cGd2QS9BTwpYcTBRVXIrQWlUcTR0UW5oYjhDbDhJK2lLRmF5ZzZvY0` +
		`FnQXVCZkZBMnVBd29CL25LajZXTHlJVHV0NWE1VDBQCmxpbVJaYllGUTFyeHBJaVpUMmFja0NxUjN1Yk9qdVBGOCtJZHVWSmNXN05WcTFRSlls` +
		`RkFrSnVhTnpaRDlNMGkKUCs3WTgvTGhBb0dCQVBEYTg2cU9pazZpamNaajJtKzFub3dycnJINjdCRzhqRzdIYzJCZzU1M2VXWHZnQ3Z6RQppMk` +
		`xYU3J6VVV6SGN2aHFQRVZqV2RPbk1rVHkxK2VoZDRnV3FTZW9iUlFqcHAxYU40clA5dVcvOStZaHVoTlZWCnJ2QUh3ZHBTaTRlelovNEVERmxl` +
		`YUd5dXNWSkcvU1lJM096bnVQU051NW1lcysxN05Hb2pBZWtaQW9HQkFPUFYKMG5oRy9rNitQLzdlRXlqL2tjU3lPeUE5MzYvV05yVUU3bDF4b2` +
		`YyK3laSVVhUitOcE1manpmcVJqaitRWmZIZwpJS0kvYmJGWGtlWm9nWG5seHk0T1YvSmtKZy9oTHo2alJUQjhYTW9kbEhwVnFOaEZYcWJhV1Bj` +
		`a0h3WkhaVFU0CkNsQWg0QWZrZ2hpVWVrS2lhcTFNMWNyOE5CTWlyeTR2WWhKVXVReERBb0dCQUpyTG5aOFlUVHVNcmFHN3V6L2cKY2kyVVJZcU` +
		`53ZnNFT3gxWGdvZUd3RlZ0K2dUclVTUnpEVUpSSysrQVpwZTlUMUN5Y211dUtTVzZHLzN3MXRUSQp3ZUx5TnQ4Rzk2OXF1K21jOXY3SEtzOFhZ` +
		`N0NUbHp1ay9mRzJpcGhPUk83S0Z5UGlaaTFweDZOU0F4VG1HdnkrCjVYNDh6MW9kWFZ5MTZ0M09PVG1kbGpUQkFvR0FTYk5SY2pjRTdOUCtQNl` +
		`AyN3J3OW16Tk1qUkYyMnBxZzk4MncKamVuRVRTRDZjNWJHcXI1WEg1SkJmMXkyZHpsdXdOK1BydXgxdjNoa2FmUkViZm8yaEY5L2M1bVI5bkVS` +
		`cDJHSgpjRFhLamxjalFLK1UvdUR4eldlMGY3M2ZpMWh0Rk5vYisrLzVXSlJDd1ZER2UrZXVPb0V3WjRsT0R5S1pLSWVMCllnS21HYUVDZ1lBMF` +
		`prd3k5ejFXczRBTmpHK1lsYVV4cEtMY0pGZHlDSEtkRnI2NVdZc21HcU5rSmZHU0dlQjYKUkhNWk5Nb0RUUmhtaFFoajhNN04rRk10WkFVT01k` +
		`ZFovMWN2UkV0Rlc3KzY2dytYWnZqOUNRL3VlY3RwL3FiKwo2ZG5PYnJkbUxpWitVL056R0xLbUZnSlRjOVg3ZndtMTFQU2xpWkswV3JkblhLbn` +
		`praDlPaFE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`

	certFileContent = `-----BEGIN CERTIFICATE-----
MIIEEDCCAvigAwIBAgIJANf+IT1MGhxBMA0GCSqGSIb3DQEBCwUAMIGZMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCVBhbG8gQWx0bzEiMCAGA1UE
CgwZT25lQ29uY2VybiBbdGVzdCBwdXJwb3NlXTEcMBoGA1UEAwwTb25lY29uY2Vy
bi10ZXN0LmNvbTEnMCUGCSqGSIb3DQEJARYYZnJlZGVyaWNAY29uZWNvbmNlcm4u
Y29tMB4XDTE4MDgwMzE2MjE0OFoXDTE5MTIxNjE2MjE0OFowgYQxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJUGFsbyBBbHRvMSIwIAYDVQQLDBlP
bmVDb25jZXJuIFt0ZXN0IHB1cnBvc2VdMTAwLgYDVQQDDCdhcGktc2VydmljZS1w
cm94aWVkLm9uZWNvbmNlcm4tdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDWW1JBvpx/ofF0z/4By+bgAp+hbyqnUlCagbigxOPLv7iH8M+u
3D6DeIY3C7dWKaN4gavGwS/7r4QlWIgo+OMDt53nNd5ojL/5f9GQ4ddNmnwonDxF
U8kwYLYDLNBHC2j31A5Snxz4KSdTM7Fg48PIy3AiaF2HDqDYVRiZEZrIxSrSJaJf
uXeBIEApPiPmH9DNll6UZ786of+IYeKWencC/ljOhiIJrV/sDe76AQcuv9OoWiGb
IFTS2XLRDatc0rAxVvQbNs3yiEb8wS0ZGAxq0nfOi0fdaPH87EsnLvJjZNOqr/MS
InAbcvfiprLqhGKA5H7hJTfDpRzYlArnfi2LAgMBAAGjbjBsMAkGA1UdEwQCMAAw
CwYDVR0PBAQDAgXgMFIGA1UdEQRLMEmCHm9hdGhrZWVwZXIub25lY29uY2Vybi10
ZXN0LmNvbYInYXBpLXNlcnZpY2UtcHJveGllZC5vbmVjb25jZXJuLXRlc3QuY29t
MA0GCSqGSIb3DQEBCwUAA4IBAQB1Pbl+RmntEoclyjYzsykKoiXs3pa83CgDZ6pC
grv4LQxSoEfJ04f8bD4IIYFGCZfVNG0VpEXrNlk6UbsVdTABtqCguJTWwDGUAh6X
3bFhrZnPexs//QwgDAdqIf0EgwcDUG76GIdfk70e1Zuxchxd8UBCpAIdeU08vVky
4PWv2K4aD4FjChKxCN5khN50FMPcaI+xVgd47tPhVN9lQioGDChsT5vARk7ba/ce
0Q9NWdiMfLEgLdcBobZKFtFrlKtyzOgDjLvXvLQs/s2merSI9fkwoOBEP+7j7Zns
Ajy9MfhwqbTpW7KpCSFa0VT/2u95ZRcPvrXldKRygB4WuARr
-----END CERTIFICATE-----`
	keyFileContent = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1ltSQb6cf6HxdM/+Acvm4AKfoW8qp1JQmoG4oMTjy7+4h/DP
rtw+g3iGNwu3VimjeIGrxsEv+6+EJViIKPjjA7ed5zXeaIy/+X/RkOHXTZp8KJw8
RVPJMGC2AyzQRwto99QOUp8c+CknUzOxYOPDyMtwImhdhw6g2FUYmRGayMUq0iWi
X7l3gSBAKT4j5h/QzZZelGe/OqH/iGHilnp3Av5YzoYiCa1f7A3u+gEHLr/TqFoh
myBU0tly0Q2rXNKwMVb0GzbN8ohG/MEtGRgMatJ3zotH3Wjx/OxLJy7yY2TTqq/z
EiJwG3L34qay6oRigOR+4SU3w6Uc2JQK534tiwIDAQABAoIBAQDOlrDOQCCgObl1
9Uc+//8BAkZK1dLr879PSZphB6Drq1jzWzkts4JkPvJLdvU5C2RLLd4Z7fKKxPu8
6ng/1K8l0/9Q6G/zn0Md+PxGgAJ6/lqO4RSNVFTgVTTWFoidD/gYcb1kD4lh+ne2
DmnzkVAN41OtNZx+H7ERDeJpE7hzqR8HhvxqeO6gnB1rdgrQHOLWYRvc5pgvA/AO
Xq0QUr+AiTq4tQnhb8Cl8I+iKFayg6ocAgAuBfFA2uAwoB/nKj6WLyITut5a5T0P
limRZbYFQ1rxpIiZT2ackCqR3ubOjuPF8+IduVJcW7NVq1QJYlFAkJuaNzZD9M0i
P+7Y8/LhAoGBAPDa86qOik6ijcZj2m+1nowrrrH67BG8jG7Hc2Bg553eWXvgCvzE
i2LXSrzUUzHcvhqPEVjWdOnMkTy1+ehd4gWqSeobRQjpp1aN4rP9uW/9+YhuhNVV
rvAHwdpSi4ezZ/4EDFleaGyusVJG/SYI3OznuPSNu5mes+17NGojAekZAoGBAOPV
0nhG/k6+P/7eEyj/kcSyOyA936/WNrUE7l1xof2+yZIUaR+NpMfjzfqRjj+QZfHg
IKI/bbFXkeZogXnlxy4OV/JkJg/hLz6jRTB8XModlHpVqNhFXqbaWPckHwZHZTU4
ClAh4AfkghiUekKiaq1M1cr8NBMiry4vYhJUuQxDAoGBAJrLnZ8YTTuMraG7uz/g
ci2URYqNwfsEOx1XgoeGwFVt+gTrUSRzDUJRK++AZpe9T1CycmuuKSW6G/3w1tTI
weLyNt8G969qu+mc9v7HKs8XY7CTlzuk/fG2iphORO7KFyPiZi1px6NSAxTmGvy+
5X48z1odXVy16t3OOTmdljTBAoGASbNRcjcE7NP+P6P27rw9mzNMjRF22pqg982w
jenETSD6c5bGqr5XH5JBf1y2dzluwN+Prux1v3hkafREbfo2hF9/c5mR9nERp2GJ
cDXKjlcjQK+U/uDxzWe0f73fi1htFNob++/5WJRCwVDGe+euOoEwZ4lODyKZKIeL
YgKmGaECgYA0Zkwy9z1Ws4ANjG+YlaUxpKLcJFdyCHKdFr65WYsmGqNkJfGSGeB6
RHMZNMoDTRhmhQhj8M7N+FMtZAUOMddZ/1cvREtFW7+66w+XZvj9CQ/uectp/qb+
6dnObrdmLiZ+U/NzGLKmFgJTc9X7fwm11PSliZK0WrdnXKnzkh9OhQ==
-----END RSA PRIVATE KEY-----`
)

func TestHTTPSCertificate(t *testing.T) {
	tmpCertFile, _ := os.CreateTemp("", "test-cert")
	tmpCert := tmpCertFile.Name()
	tmpKeyFile, _ := os.CreateTemp("", "test-key")
	tmpKey := tmpKeyFile.Name()
	defer func() {
		_ = os.Remove(tmpCert)
		_ = os.Remove(tmpKey)
	}()
	_ = os.WriteFile(tmpCert, []byte(certFileContent), 0600)
	_ = os.WriteFile(tmpKey, []byte(keyFileContent), 0600)

	// 1. no TLS
	cert, err := tls.LoadCertificates("", "", "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.EqualError(err, tls.ErrNoCertificatesConfigured); !ok {
		t.Error(msg)
	}

	// 2. inconsistent TLS (i): warning only
	cert, err = tls.LoadCertificates("", "", "", "x")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.EqualError(err, tls.ErrInvalidCertificateConfiguration); !ok {
		t.Error(msg)
	}

	// 3. inconsistent TLS (ii): warning only
	cert, err = tls.LoadCertificates("x", "", "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.EqualError(err, tls.ErrInvalidCertificateConfiguration); !ok {
		t.Error(msg)
	}

	// 4. invalid TLS file
	cert, err = tls.LoadCertificates("", "", tmpCert, "x")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 5. invalid TLS string (i)
	cert, err = tls.LoadCertificates(certFixture, "{}", "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}
	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 6. invalid TLS string (ii)
	cert, err = tls.LoadCertificates("{}", keyFixture, "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 7. valid TLS files
	cert, err = tls.LoadCertificates("", "", tmpCert, tmpKey)
	if ok, msg := testing_.NonNil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 8. valid TLS strings
	cert, err = tls.LoadCertificates(certFixture, keyFixture, "", "")
	testing_.NonZero(t, cert)

	if ok, msg := testing_.NonError(err); !ok {
		t.Error(msg)
	}

	// 9. invalid TLS file content
	cert, err = tls.LoadCertificates("", "", certFixture, keyFixture)

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 10. invalid TLS string content
	cert, err = tls.LoadCertificates(certFileContent, keyFileContent, "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 11. mismatched TLS file content
	cert, err = tls.LoadCertificates("", "", keyFileContent, certFileContent)

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

	// 12. mismatched TLS string content
	cert, err = tls.LoadCertificates(keyFixture, certFixture, "", "")

	if ok, msg := testing_.Nil(cert); !ok {
		t.Error(msg)
	}

	if ok, msg := testing_.Error(err); !ok {
		t.Error(msg)
	}

}
