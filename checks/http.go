package checks

import (
	"crypto/tls"
	"io"
	"monitoring/spec"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func checkHttps(url string, checkSpec *spec.CheckSpec) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: spec.IHTTPS == checkSpec.Type},
	}

	client := http.Client{
		Timeout:   time.Duration(checkSpec.Timeout) * time.Second,
		Transport: tr,
	}

	resp, err := client.Get(url)
	if nil != err {
		return err
	}

	err = handleCheckMethods(checkSpec, resp)
	return err
}

func checkHttp(url string, checkSpec *spec.CheckSpec) error {
	client := http.Client{
		Timeout: time.Duration(checkSpec.Timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if nil != err {
		return spec.ConnectionFailed
	}

	err = handleCheckMethods(checkSpec, resp)
	return err
}

func handleCheckMethods(checkSpec *spec.CheckSpec, resp *http.Response) error {
	switch checkSpec.Method {
	case spec.CONNECT:
		return nil

	case spec.CONTAINS:
		content, err := io.ReadAll(resp.Body)
		if nil == checkSpec.Expect {
			return spec.ExpectUndefined
		}

		if nil != err {
			return err
		}

		if strings.Contains(string(content), *checkSpec.Expect) {
			return nil
		}

		return spec.ExpectFailed

	case spec.STATUS:
		if nil == checkSpec.Expect {
			return spec.ExpectUndefined
		}

		if strconv.Itoa(resp.StatusCode) != *checkSpec.Expect {
			return spec.ExpectFailed
		}

		return nil

	case spec.QUERY:
		return spec.NotImplemented

	}

	return spec.NoCheckPerformed
}

func (c *CheckHandler) CheckHTTPS(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkHttps(spec.DSN, spec)
}

func (c *CheckHandler) CheckIHTTPS(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkHttps(spec.DSN, spec)
}

func (c *CheckHandler) CheckHTTP(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkHttp(spec.DSN, spec)
}
