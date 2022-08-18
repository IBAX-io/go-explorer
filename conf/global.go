/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package conf

import (
	"fmt"
	"time"
)

type serverModel struct {
	Mode                 string        `yaml:"mode"`                    // run mode
	Host                 string        `yaml:"host"`                    // server host
	Port                 int           `yaml:"port"`                    // server port
	EnableHttps          bool          `yaml:"enable_https"`            // enable https
	CertFile             string        `yaml:"cert_file"`               // cert file path
	KeyFile              string        `yaml:"key_file"`                // key file path
	JwtPubKeyPath        string        `yaml:"jwt_public_key_path"`     // jwt public key path
	JwtPriKeyPath        string        `yaml:"jwt_private_key_path"`    // jwt private key path
	TokenExpireSecond    time.Duration `yaml:"token_expire_second"`     // token expire second
	SystemStaticFilePath string        `yaml:"system_static_file_path"` // system static file path
	DocsApi              string        `yaml:"docs_api"`                // api docs request address
}

type UrlModel struct {
	Base string `yaml:"base_url"`
}

type LogConfig struct {
	LogTo     string
	LogLevel  string
	LogFormat string
}

func (r *serverModel) Str() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
