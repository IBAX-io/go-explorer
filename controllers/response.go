/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package controllers

type Response struct {
	Code    int    `json:"code" `
	Data    any    `json:"data" `
	Message string `json:"message" `
}

func (r *Response) ReturnFailureString(str string) {
	r.Code = -1
	r.Message = str
}
func (r *Response) Return(dat any, ct CodeType) {
	r.Code = ct.Code
	r.Message = ct.Message
	r.Data = dat
}
