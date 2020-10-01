package cli

import (
	"bufio"
	"fmt"
	"github.com/kpango/glg"
	"os"
	"strings"
)

type (
	CliFunc func( []string ) error
	CliFuncPack struct {
		F CliFunc
		Desc string
		Group string
	}
)

var (
	CliFuncs map[string] *CliFuncPack
	groups map[string] bool
	ir * bufio.Reader
)

func init() {
	CliFuncs = make(map[string] *CliFuncPack)
	groups = make(map[string]bool)
	ir = bufio.NewReader(os.Stdin)
	initClis()
}

func proc() {
	defer func() {
		if r := recover(); r != nil {
			glg.Error(r)
		}
	}()
	print("cyf-cloud> ")
	input, e := ir.ReadString('\n')

	for input[0:1] == " " {
		input = strings.TrimPrefix( input," ")
	}
	args := strings.Split( input, " ")

	fa := strings.Replace( args[0], "\n", "", -1 )
	fa = strings.Replace( fa, "\r", "", -1 )
	f := CliFuncs[fa]
	// 是否有该命令
	if f == nil {
		fmt.Println("command \"" + fa + "\" does not exist")
		return
	}

	// 移除首个元素
	args = append(args[:0], args[1:]...)
	e = f.F( args )
	if e != nil {
		panic( e )
	}
}

func Run() {
	print("cyf-cloud> ")
	for {
		proc()
	}
}

func Register( name string, f *CliFuncPack ) {
	if CliFuncs[name] != nil {
		glg.Warn("cli: "+name+ " overwrote")
	}
	groups[f.Group] = true
	CliFuncs[name] = f
}