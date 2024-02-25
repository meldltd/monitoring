package checks

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"monitoring/spec"
	"strconv"
	"time"
)

func checkTLS(url string, checkSpec *spec.CheckSpec) error {
	if nil == checkSpec.DSNParams {
		return errors.New("DSNParams user and password must be set for SSH check!")
	}

	conn, err := tls.Dial("tcp", url, nil)
	if nil != err {
		return errors.New("TLS connection failed")
	}
	defer conn.Close()

	return handleTLSCheckMethods(checkSpec, conn.ConnectionState().PeerCertificates[0])
}

func handleTLSCheckMethods(checkSpec *spec.CheckSpec, cert *x509.Certificate) error {
	switch checkSpec.Method {
	case spec.CONNECT:
		return nil
	case spec.EXPIRES:
		days := (*checkSpec.DSNParams)["daysBefore"].(string)
		if "" == days {
			return errors.New("daysBefore must be defined for TLS Expiry Check")
		}
		daysBefore, err := strconv.Atoi(days)
		if err != nil {
			return errors.New("daysBefore is not integer")
		}

		if cert.NotAfter.Before(time.Now().Add(time.Hour * time.Duration(24*daysBefore))) {
			return spec.WillExpire
		}

		return nil
	case spec.QUERY:
		return spec.NotImplemented
	case spec.CONTAINS:
		return spec.NotImplemented
	case spec.STATUS:
		return spec.NotImplemented
	}

	return spec.NoCheckPerformed
}

func (c *CheckHandler) CheckTLS(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkTLS(spec.DSN, spec)
}
