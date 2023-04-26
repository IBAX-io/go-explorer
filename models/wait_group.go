/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import "sync"

var ChartWG = &sync.WaitGroup{}
var HistoryWG = &sync.WaitGroup{}
var RealtimeWG = &sync.WaitGroup{}
