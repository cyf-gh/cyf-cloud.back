package ccDoc

import (
	"../cc/err"
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type (
	 Model struct {
			File   string       `json:"file"`
			Desc   string       `json:"desc"`
			Path   string       `json:"path"`
			Childs []ChildModel `json:"childs"`
	}
	MatchFunc struct {
		Pattern *regexp.Regexp
		Proc func( params[]string, i int ) error
	}
	ChildModel struct {
		Path  string        `json:"path"`
		Brief string        `json:"brief"`
		Note  string        `json:"note"`
		Args  []Arg         `json:"args"`
		Return string       `json:"return"`
		Type string         `json:"type"`
		NeedValidation bool `json:"needValidation"`
		IsDeprecated bool 	`json:"isDeprecated"`
	}
	Arg struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
)

func ResetCurrentDesc() {
	currentDesc = nil
}

func InChildDoc() {
	inChildDoc = true
}

func FinishChildDoc() {
	currentDocModel.Childs = append( currentDocModel.Childs, *currentChildModel )
	prevChildModel = &currentDocModel.Childs[len(currentDocModel.Childs)-1]
	currentChildModel = &ChildModel{}
	inChildDoc = false
}

var (
	Matchfuncs      []MatchFunc
	currentDocModel *Model
	currentChildModel, prevChildModel *ChildModel
	inChildDoc bool
	currentDesc *string
	)

func init() {
	inChildDoc = false
	Matchfuncs = []MatchFunc {
		// 匹配主路由
		{
			Pattern: regexp.MustCompile(`^(\s*)cc\.AddActionGroup\((\s*)\"(.+)\"(.+)$`),
			Proc:    func( params[]string, i int ) error {
				currentDocModel.Path = params[3]
				return nil
			},
		},
		// 注释
		// 非cc格式注释
		{
			Pattern: regexp.MustCompile(`(\s*)// [^\\](.+)$`),
			Proc:    func( params[]string, i int ) error {
				// 首行为该文件的描述
				if i == 0 {
					currentDocModel.Desc = params[2]
				}
				// 如果有注释跟随之前的描述，则添加
				if currentDesc != nil {
					*currentDesc += params[2]
				}
				return nil
			},
		},
		// 匹配子路由brief
		{
			Pattern: regexp.MustCompile(`^(\s*)// \\brief(\s*)(.+)$`),
			Proc:    func( params[]string, i int ) error {
				InChildDoc()
				currentChildModel.Brief = params[3]
				currentDesc = &currentChildModel.Brief
				return nil
			},
		},
		// 匹配子路由note
		{
			Pattern: regexp.MustCompile(`^(\s*)// \\note(\s*)(.+)$`),
			Proc:    func( params[]string, i int ) error {
				InChildDoc()
				currentChildModel.Note = params[3]
				currentDesc = &currentChildModel.Note
				return nil
			},
		},
		// 匹配子路由 return
		{
			Pattern: regexp.MustCompile(`^(\s*)// \\return(\s*)(.+)$`),
			Proc:    func( params[]string, i int ) error {
				InChildDoc()
				currentChildModel.Return = params[3]
				currentDesc = &currentChildModel.Return
				return nil
			},
		},
		// 匹配子路由arg
		{
			Pattern: regexp.MustCompile(`^(\s*)// \\arg\[(.+)\](\s*)(.+)$`),
			Proc:    func( params[]string, i int ) error {
				InChildDoc()
				arg := Arg{
					Name: params[2],
					Desc: params[4],
				}
				currentDesc = &arg.Name
				currentChildModel.Args = append( currentChildModel.Args, arg )
				return nil
			},
		},
		// 匹配子路由其余属性
		{
			Pattern: regexp.MustCompile(`^(\s*)// \\(.+)(\s*)(.+)$`),
			Proc:    func( params[]string, i int ) error {
				InChildDoc()
				return nil
			},
		},
		// 匹配子路由路由
		{
			Pattern: regexp.MustCompile(`^(\s*)a.(GET|POST|WS)(.*)\((\s*)\"(.*)\"(.+)$`),
			Proc:    func( params[]string, i int ) error {
				currentChildModel.Path = currentDocModel.Path + params[5]
				currentChildModel.Type = params[2]+params[3]
				FinishChildDoc()
				return nil
			},
		},
		// 匹配路由是否需要验证
		// 一定正确
		{
			Pattern: regexp.MustCompile(`^(.*)cc.HerDeprecated(.*)$`),
			Proc: func( params[]string, i int ) error {
				prevChildModel.IsDeprecated = true
				return nil
			},
		},
		// 匹配该函数是否已经弃用
		{
			Pattern: regexp.MustCompile(`^(.*)(ByAtk|DM1CheckPermission|MakeClipboardKey)(.*)$`),
			Proc:    func( params[]string, i int ) error {
				prevChildModel.NeedValidation = true
				return nil
			},
		},
		// 注释
		// 所有注释
		{
			Pattern: regexp.MustCompile(`(\s*)//(.+)$`),
			Proc:    func( params[]string, i int ) error {
				// 如果有注释跟随之前的描述，则添加
				if currentDesc != nil {
					*currentDesc += params[2]
				}
				return nil
			},
		},
		// 代码
		{
			Pattern: regexp.MustCompile(`^(.+)$`),
			Proc:    func( params[]string, i int ) error {
				if inChildDoc {
					FinishChildDoc()
				}
				ResetCurrentDesc()
				return nil
			},
		},
	}
	currentDocModel = &Model{}
	currentChildModel = &ChildModel{}
	prevChildModel = &ChildModel{}
}

func main() {

}

func GenerateDocJson( inputDir, outputDir string ) string {
	httpGoPath := inputDir

	httpDir, e := ioutil.ReadDir( httpGoPath ); err.Assert( e )

	PthSep := string(os.PathSeparator)

	var GoFiles []string
	var GoFileNames []string
	for _, fi := range httpDir {
		if !fi.IsDir() {
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				GoFiles = append( GoFiles, httpGoPath + PthSep + fi.Name() )
				GoFileNames = append(GoFileNames, fi.Name())
			}
		}
	}

	var docModels []Model

	for i, f := range GoFiles {
		fi, e := os.Open( f ); err.Assert( e )
		br := bufio.NewReader(fi)
		currentDocModel.File = GoFileNames[i]
		for {
			l, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}

			for i, m := range Matchfuncs {
				// 匹配到了之后立即退出
				if match( m.Pattern, string(l), i, m.Proc ) {
					break
				}
			}
		}
		if currentDocModel.Path != "" {
			// 忽略没有路由的文件
			docModels = append(docModels, *currentDocModel)
		}
		currentDocModel = &Model{}
		currentChildModel = &ChildModel{}
		fi.Close()
	}
	bss, e := json.Marshal( docModels )
	if outputDir != "" {
		ioutil.WriteFile( os.Args[2], bss, 0777 )
	}
	return string(bss)
}

func match( pattern *regexp.Regexp, str string, strIndex int, proc func( params[]string, strIndex int ) error ) bool {
	params := pattern.FindStringSubmatch(string(str))
	if len(params) > 0 {
		proc( params, strIndex )
		return true
	} else {
		return false
	}
}