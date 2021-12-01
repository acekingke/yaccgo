/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package symbol

import "testing"

func TestNewSymbol(t *testing.T) {
	//t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	sym := NewSymbol(1, "start")
	if sym.Name != "start" {
		t.Errorf("value name %s fail", sym.Name)
	}
}
