/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	builder "github.com/acekingke/yaccgo/Builder"
	parser "github.com/acekingke/yaccgo/Parser"
	utils "github.com/acekingke/yaccgo/Utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yaccgo",
	Short: "yaccgo is a yacc generator",
	Long:  "Understandable yacc generator , it can generate go/js/rust code",
	Run:   nil,
}

func init() {
	c := &cobra.Command{
		Use:   "generate [flags]",
		Short: "generate filetype input.y output.go",
		Long: `
		"generate filetype input.y output.go
		filetype can be go/js/rust:
			go : generate go code",
			typescript : generate typescript code",
			rust : generate rust code",
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				cmd.Help()
			}
			cmdGenerate(args)
		},
	}
	rootCmd.AddCommand(c)
	debugCmd := &cobra.Command{
		Use:   "debug input.y",
		Short: "open debug mode",
		Long: `	`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
			}
			RunDebugCmd(args)
		},
	}
	rootCmd.AddCommand(debugCmd)

}

type genfun func(input string, file string) error

func cmdGenerate(args []string) {
	switch args[0] {
	case "go":
		genCommonFunc(args[1], args[2], builder.GoGenFromString)
	case "typescript":
		genCommonFunc(args[1], args[2], builder.TsGenFromString)
	case "rust":
		fmt.Println("not support rust yet")
	}
}

func genCommonFunc(in string, out string, gen genfun) {
	// read all input file to string
	if f, err := os.Open(in); err != nil {
		panic(err)
	} else {
		defer f.Close()
		if inputbytes, err := ioutil.ReadAll(f); err != nil {
			panic(err)
		} else {
			if err = gen(string(inputbytes), out); err != nil {
				panic(err)
			}
		}

	}
}

func RunDebugCmd(args []string) {
	utils.DebugFlags = true
	genCommonFunc(args[0], "", func(in string, n string) error {
		_ = n
		_, err := parser.ParseAndBuild(in)
		return err
	})
}
