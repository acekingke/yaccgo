/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

import "strings"

func genTempName(in string) string {
	return "$operator" + in
}

func RemoveTempName(in string) string {

	if len(in) > 9 && in[0:9] == "$operator" {
		return "'" + in[9:] + "' "
	}
	return in
}

func TestPrefix(in string) bool {
	return strings.HasPrefix(in, "$operator")
}
