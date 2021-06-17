package JGScript

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
    "bufio"
)

func Run(name, code string) {
	ctx := MakeContext(nil, &BPos{0, 0, 0, &name, &code})
	res := EvalCode(name, "__main__", code, ctx)
	if res.has() {
		fmt.Println(ToStr(res.e))
		os.Exit(-3)
	}
}

func RunCompiled(name string) {
	src, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var ast AST
	decoder := gob.NewDecoder(src)
	decoder.Decode(&ast)
	ctx := MakeContext(nil, ast.S)
	res := ExecAst(&ast, ctx)
	if res.has() {
		fmt.Println(ToStr(res.e))
		os.Exit(-3)
	}
}

func RunFile(file string) {
	if strings.HasSuffix(file, ".jgs") {
		src, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		} else {
			Run(file, string(src))
		}
	} else {
		RunCompiled(file)
	}
}

func CompileFile(file string, out string) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	code := string(src)
	lexer := NewLexer(&file, &code)
	parser := MakeParser(lexer)
	ast := parser.Parse()
	if out == "" {
		out = file
		if strings.HasSuffix(out, ".jgs") {
			out = strings.Replace(out, ".jgs", ".jgc", -1)
		} else {
			out += ".jgc"
		}
	}
	f, _ := os.Create(out)
	defer f.Close()
	encoder := gob.NewEncoder(f)
	encoder.Encode(*ast)
}

func Compile(file, out string, rec bool) {
	files := []string{}
	if rec {
		getDirs(file, &files)
	} else {
		files = append(files, file)
	}
	pathout := ""
	if len(files) == 1 && out != "" {
		pathout = out
	}
	for _, f := range files {
		fmt.Println("Compiling: " + f + " ...")
		CompileFile(f, pathout)
	}
	fmt.Println("Complete.")
}

func getDirs(dir string, file *[]string) {
	fns, e := ioutil.ReadDir(dir)
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
	for _, fn := range fns {
		if fn.IsDir() {
			getDirs(path.Join(dir, fn.Name()), file)
		} else {
			if strings.HasSuffix(fn.Name(), ".jgs") {
				*file = append(*file, path.Join(dir, fn.Name()))
			}
		}
	}
}

func RunInterpreter() {
    name := "MainInterpreter"
	ctx := MakeContext(nil, &BPos{0, 0, 0, &name, nil})
    code := ""
    bc := 0
    for {
        if bc == 0 {
            fmt.Print(">>> ")
        } else {
            fmt.Print(getIndent(bc))
        }
        reader := bufio.NewReader(os.Stdin)
		in, _ := reader.ReadString('\n')
		input := strings.TrimRight(in, "\r\n")
        if input == "exit" {
            break
        } else {
            if bc > 0 {
                code += "\n" + input
            } else {
                code = input
            }
            for _, c := range input {
                if c == '{' || c == '[' || c == '(' {
                    bc ++
                } else if c == '}' || c == ']' || c == ')' {
                    bc --
                }
            }
            if bc == 0 {
                ctx.s.Code = &code
                res := EvalCode(name, name, code, ctx)
	            if res.has() {
		            fmt.Println(ToStr(res.e))
	            } else if (*res.t).GetT() != "null" {
                    if (*res.t).GetT() == "list" {
                        printList(res.t)
                        fmt.Println()
                    } else {
                        fmt.Println((*res.t).str())
                    }
                }
            }
        }
    }
}

func getIndent(bc int) string {
    s := "... "
    for i := 0; i < bc; i++ {
        s += "    "
    }
    return s
}
