/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package rules

import (
	"reflect"
	"testing"

	symbol "github.com/acekingke/yaccgo/Symbol"
)

func TestNewProductoinRule(t *testing.T) {
	type args struct {
		LeftPart *symbol.Symbol
		RighPart []*symbol.Symbol
	}
	tests := []struct {
		name string
		args args
		want *ProductoinRule
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProductoinRule(tt.args.LeftPart, tt.args.RighPart); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProductoinRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
