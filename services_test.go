package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPopulateTLSConfigCertificate(t *testing.T) {
	assert := assert.New(t)

	cert_data := "-----BEGIN CERTIFICATE-----\nMIIC+jCCAeKgAwIBAgIRALoWCCvCs3+Ssdt6AX05O6kwDQYJKoZIhvcNAQELBQAw\nEjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xNzA1MTEwNDI0MjJaFw0xODA1MTEwNDI0\nMjJaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQDtPLAe953YqenqauNx2c8gEfY0dgUeVdzl6dnn1mMz4qTAmIA/ZsER\n4irM1isotYlDT6CtdFUZFahQ2Ttqorash2jruUbqANfk/yOHRdXAf5gIiULkVenQ\nPUzqSaTTkBPemxS83JYseseOJQT0KhCkma4lF7pPlU87BkDgex3sshbC0LpzdbQh\nZfsBmFuFqcTekGe10jecZRUyxgG6GhSaAhdRmcJFDhooKo5CO2yA9HmfUnyWfjYJ\n6XqZc+2YTwdV+eQ1wwdU3GPqmNQoFTDLatFfIWSVsC1J4RNU3OLgu9xsLucWZ++6\n+bOX/CsACpEGeSD6yr476dvN3ltZMRq/AgMBAAGjSzBJMA4GA1UdDwEB/wQEAwIF\noDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBQGA1UdEQQNMAuC\nCWxvY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEAVmShb3GbDknAG5Ncsya+qp6P\na95vj5nm/NwxjNorh8mX8vUK5NBGCUKORdND6jEZFq5MN/P2iSxd87X8iy1Xmfbr\nLtcuvktHtJxf2j5KA0cfCbFA+g+Zp210uPm6C05BI1js4gHoi03SMvMAeXQZH2Mg\ndSGw4o3rIDTxAIOv3jrbniUi5kY1bwGQcrmpX7u87OSnXJEbXMl0N36k/QA64Ihg\nVhdw2MHRVtqI9PMpAyXsoUXsUcuXnwJLUoCdKsqcLJWw6uClFA7W1qHPQZ0Z3np6\np9j+vmMxOHBylK/mnvvFu2ke7UpdYDPyA5MOhXnMqY00tRMsov0/4DlBp6bFUg==\n-----END CERTIFICATE-----\n"
	key_data := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7TywHved2Knp6mrjcdnPIBH2NHYFHlXc5enZ59ZjM+KkwJiA\nP2bBEeIqzNYrKLWJQ0+grXRVGRWoUNk7aqK2rIdo67lG6gDX5P8jh0XVwH+YCIlC\n5FXp0D1M6kmk05AT3psUvNyWLHrHjiUE9CoQpJmuJRe6T5VPOwZA4Hsd7LIWwtC6\nc3W0IWX7AZhbhanE3pBntdI3nGUVMsYBuhoUmgIXUZnCRQ4aKCqOQjtsgPR5n1J8\nln42Cel6mXPtmE8HVfnkNcMHVNxj6pjUKBUwy2rRXyFklbAtSeETVNzi4LvcbC7n\nFmfvuvmzl/wrAAqRBnkg+sq+O+nbzd5bWTEavwIDAQABAoIBAQDNNGFjZ+wxLUgY\nbLywDicLJn5AgpWK19prRQwnbVoB11mK/l5weQEn5un+pIJQhDZm8smZP7ccK4+b\n30t1wakhMz+eJnUgk/orKkYhDFcIW1W2jIQ/3dCRP3T4cxsPpOCK/LnDY0rCzrEu\nUvcl+/zJY+UuUIfnrs8Jzy7u/Y+02+w/8knxwX8X+RzGRArFY5eRrAJlxmPuCTab\nwOqUrx3DjaFdPWNBcGtn1kInF1pMGgnn7MBQky+qWJHkXAF7Xto3mk0durDU1zQv\nxXLEfsp5K6lm+/anokO/UR36qOrAcihlA8/fbYaMQhUFEtlo96jpg/G5GPnlS8HG\nsfI5RBCBAoGBAPxM3Iybr/Q60ZPl6wcUduHC8v9gePWeb/oStdbhLgwPgrqMOi6P\nQISmxKTNku4fu+JhuEzrOQvmWwYDkSzPWWJ4L8qd5mr3XGWgDNfEEdCyYg4fnPCI\nIMkFOdIjtjXd51U47WhsSVzLxwOmVi7hoBXx7SdObyCI+5e9hcrw5vHNAoGBAPC3\nR3HREsUPcGX0cFeBrisM29FP48yhMdoAPKXfYZTZNERfNc+v7w1pZuDeipLid0f3\nxhSDXpjqzIohov9lm3HAuPRF9h+VLCL3HvEAcES/1/qVZ4/a3nWsrcUTo8wcWaeJ\nwyuzc7xdWR8KQHkjnQ0HxZWKfoBxg6hXyhvhpmK7AoGBANDMV1XyXnrT1q/8fjY8\nxGnwKaQZVeGHvooJw/1SHAaVK45xEJGJsk5VqbXt/6QcFSSz1I+rt2lWuYvPleys\nqP+qEXswlAmALzJXc2l5dXjut+GSXhJdxiw2q/Rx45DO1W0dELTzsP8gEdK+bOKd\njRu8PJTj/2nAk96vVTNvjOEdAoGAQYtWS9qHBtt2WnydnmY0O4qrzCm5uH1n6plH\n6k6R7oraHYfjSiL4r6k0lyRhjS9XFWSVLf67Yl4ExdP04yASnH3Cntjx9JWUyAyM\nA0mASGgIcjX+VgBdtKMJSfFYF3rcuq7bUunxKKguXTJYbePRnruwBFEKswS1ub/1\ny3O874MCgYBk83x3WHCg3BiRN8lm0Rrs8RIHRLmqC7LNQGM/miIDHhvyFvWjaM3L\nHAyK+aTmoYCM0/b/nGEzg+gufN2nSn8zDwH/Hao/GzROT5pcX800FGOdcFp+o0WP\npxgqoZs+r1SiFbncU4zlfHd2zgo2M5/dKUu+CR5LmZVO4Xwx21JQ2g==\n-----END RSA PRIVATE KEY-----\n"

	tls_config := TLSConfig{
		"CERT": cert_data,
		"KEY":  key_data,
	}
	config, err := populateTLSConfig(tls_config)
	assert.Nil(err)
	loaded_cert := config.Certificates[0].Certificate[0]
	loaded_key := config.Certificates[0].PrivateKey

	var cert_buffer bytes.Buffer
	err = pem.Encode(&cert_buffer, &pem.Block{Type: "CERTIFICATE", Bytes: loaded_cert})
	assert.Nil(err)
	loaded_cert_pem := cert_buffer.Bytes()

	var key_buffer bytes.Buffer
	err = pem.Encode(&key_buffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(loaded_key.(*rsa.PrivateKey))})
	assert.Nil(err)
	loaded_key_pem := key_buffer.Bytes()

	assert.Equal(cert_data, string(loaded_cert_pem))
	assert.Equal(key_data, string(loaded_key_pem))
}

func TestPopulateTLSConfigInsecureSkipVerify(t *testing.T) {
	assert := assert.New(t)

	tls_config := TLSConfig{
		"InsecureSkipVerify": "true",
	}
	config, err := populateTLSConfig(tls_config)
	assert.Nil(err)
	assert.Equal(true, config.InsecureSkipVerify)
}

func TestPopulateTLSConfigNextProtos(t *testing.T) {
	assert := assert.New(t)

	tls_config := TLSConfig{
		"NextProtos": "foo,bar",
	}
	config, err := populateTLSConfig(tls_config)
	assert.Nil(err)
	assert.Equal([]string{"foo", "bar"}, config.NextProtos)
}
