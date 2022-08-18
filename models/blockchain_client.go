/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

func sendPostFormRequest(reqUrl string, reqBody url.Values) ([]byte, error) {
	resp, err := http.PostForm(reqUrl, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("sendPostFormRequest err:")
		return nil, err
	}
	return data, nil
}
