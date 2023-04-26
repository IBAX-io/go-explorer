/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package buffer

import (
	"github.com/IBAX-io/go-explorer/models"
)

func dataChartServer(reqType *BufferType) error {
	switch reqType {
	case dataRealTime:
	case dataHistory:

	}
	return nil
}

func ecosystemChartServer(reqType *BufferType) error {

	switch reqType {
	case ecoHistory:
		models.GetHistoryEcosystemChartInfo()
	case ecoRealTime:

	}

	return nil
}
