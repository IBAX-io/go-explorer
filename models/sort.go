/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"github.com/IBAX-io/go-explorer/storage"
)

type LeaderboardSlice []storage.HonorNodeModel // Sort by LeaderboardSlice.number from largest to smallest

func (a LeaderboardSlice) Len() int { // Override the Len() method
	return len(a)
}
func (a LeaderboardSlice) Swap(i, j int) { // Override the Swap() method
	a[i], a[j] = a[j], a[i]
}

func (a LeaderboardSlice) Less(i, j int) bool { // Rewrite the Less() method, sort from largest to smallest
	//return a[j].PkgAccountedFor.LessThan(a[i].PkgAccountedFor)

	return a[j].NodeBlock < a[i].NodeBlock
}
