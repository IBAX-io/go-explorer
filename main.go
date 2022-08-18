/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"runtime"

	"github.com/IBAX-io/go-explorer/cmd"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @contact.name   IBXC Official Website
// @contact.url    https://ibax.io
// @contact.email  support@ibax.io

// @license.name  Apache 2.0
// @license.url https://github.com/IBAX-io/go-ibax/blob/main/LICENSE
func main() {
	runtime.LockOSThread()
	cmd.Execute()
}
