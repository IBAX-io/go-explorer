/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package get_data

import "github.com/IBAX-io/go-explorer/models"

type HonorNode struct {
	Signal chan bool
}

func (p *HonorNode) SendSignal() {
	select {
	case p.Signal <- true:
	default:
		//If there is still unprocessed content in the channel, not continue to send
	}
}

func (p *HonorNode) ReceiveSignal() {
	if p.Signal == nil {
		p.Signal = make(chan bool)
	}
	for {
		select {
		case <-p.Signal:
			models.InsertHonorNodeInfo()
		}
	}
}
