package main

import (
	"flag"
	"fmt"
	"github.com/gitdxj/errgen"
	"github.com/tallstoat/pbparser"
)

func main(){
	flag.Parse()
	filename := "./resource/student.proto"
	pf, err := pbparser.ParseFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	errgen.GenerateErrFile(pf)
}
