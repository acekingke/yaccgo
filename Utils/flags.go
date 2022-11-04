/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package utils

import "strings"

var DebugFlags bool = false
var PackFlags bool = true
var HttpDebug bool = false
var DebugPackTab bool = false
var GenDotGraph bool = false
var GenDotPath string = "./yaccgo.png"
var ObjectMode bool = false

func RemoveTempName(in string) string {

	if len(in) > 9 && in[0:9] == "$operator" {
		return "'" + in[9:] + "' "
	}
	return in
}

func EscapeDotGraph(in string) string {
	res := strings.ReplaceAll(in, "<", "\\<")
	res = strings.ReplaceAll(res, ">", "\\>")
	return res
}
