/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var centrifugoTimeout = time.Second * 5

const (
	CryptoError = "Crypto"
)

type CentJWT struct {
	Sub string
	jwt.StandardClaims
}

type CentJWTToken struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

func GetJWTCentToken(userID, expire int64) (*CentJWTToken, error) {
	if conf.GetCentrifugoConn().Enable {
		var ret CentJWTToken
		centJWT := CentJWT{
			Sub: strconv.FormatInt(userID, 10),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Second * time.Duration(expire)).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, centJWT)
		result, err := token.SignedString([]byte(conf.GetCentrifugoConn().Secret))

		if err != nil {
			log.WithFields(log.Fields{"type": CryptoError, "error": err}).Error("JWT centrifugo error")
			return &ret, err
		}
		ret.Token = result
		ret.Url = conf.GetCentrifugoConn().Socket
		return &ret, nil
	} else {
		var ret CentJWTToken
		return &ret, errors.New("centrifugo not enable")
	}
}

func WriteChannelByte(channel string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), centrifugoTimeout)
	defer cancel()
	return conf.GetCentrifugoConn().Conn().Publish(ctx, channel, data)
}
