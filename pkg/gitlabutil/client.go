package gitlabutil

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

// httpClientWithTimeout 创建一个带有超时设置的 HTTP 客户端
func httpClientWithTimeout(timeout time.Duration, skipSSLVerify bool) *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLVerify,
		},
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
	return client
}

// NewClient 创建一个新的 GitLab 客户端
// token: GitLab 用户令牌
// apiURL: GitLab API 地址
// skipSSLVerify: 是否跳过 SSL 验证
func NewClient(token string, apiURL string, skipSSLVerify bool) (*gitlab.Client, error) {
	if len(token) == 0 {
		return nil, errors.New("Gitlab user token not set")
	}

	opts := []gitlab.ClientOptionFunc{
		gitlab.WithHTTPClient(httpClientWithTimeout(time.Second*90, skipSSLVerify)),
	}

	if len(apiURL) > 0 {
		if !strings.HasSuffix(apiURL, "/") {
			apiURL += "/"
		}
		if !strings.HasSuffix(apiURL, "api/v4/") {
			apiURL += "api/v4/"
		}
		opts = append(opts, gitlab.WithBaseURL(apiURL))
	}

	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
