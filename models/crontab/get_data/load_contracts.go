/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package get_data

import (
	"github.com/IBAX-io/go-ibax/packages/smart"
	log "github.com/sirupsen/logrus"
)

type LoadContracts struct {
	Signal chan bool
}

func (p *LoadContracts) SendSignal() {
	select {
	case p.Signal <- true:
	default:
		//If there is still unprocessed content in the channel, not continue to send
	}
}

func (p *LoadContracts) ReceiveSignal() {
	if p.Signal == nil {
		p.Signal = make(chan bool)
	}
	for {
		select {
		case <-p.Signal:
			if err := smart.LoadContracts(); err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Load Contracts Failed")
			}
		}
	}
}
