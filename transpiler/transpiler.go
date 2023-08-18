package transpiler

import (
	"bytes"
	"embed"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/javier-varez/monkey_interpreter/ast"
)

//go:embed runtime/include/* runtime/templates/* runtime/src/* runtime/CMakeLists.txt
var templateFS embed.FS

type astNodeType string

var templates map[astNodeType]*template.Template = map[astNodeType]*template.Template{}
var preamble string

const (
	PROGRAM                     = astNodeType("PROGRAM")
	LET_STATEMENT               = astNodeType("LET_STATEMENT")
	EXPRESSION_STATEMENT        = astNodeType("EXPRESSION_STATEMENT")
	IDENTIFIER_EXPRESSION       = astNodeType("IDENTIFIER_EXPRESSION")
	INTEGER_LITERAL_EXPRESSION  = astNodeType("INTEGER_LITERAL_EXPRESSION")
	CALL_EXPRESSION             = astNodeType("CALL_EXPRESSION")
	FN_LITERAL_EXPRESSION       = astNodeType("FN_LITERAL_EXPRESSION")
	BLOCK_STATEMENT             = astNodeType("BLOCK_STATEMENT")
	STRING_LITERAL_EXPRESSION   = astNodeType("STRING_LITERAL_EXPRESSION")
	PREFIX_EXPRESSION           = astNodeType("PREFIX_EXPRESSION")
	INFIX_EXPRESSION            = astNodeType("INFIX_EXPRESSION")
	BOOL_LITERAL_EXPRESSION     = astNodeType("BOOL_LITERAL_EXPRESSION")
	IF_EXPRESSION               = astNodeType("IF_EXPRESSION")
	RETURN_STATEMENT            = astNodeType("RETURN_STATEMENT")
	ARRAY_LITERAL_EXPRESSION    = astNodeType("ARRAY_LITERAL_EXPRESSION")
	INDEX_OPERATOR_EXPRESSION   = astNodeType("INDEX_OPERATOR_EXPRESSION")
	VAR_ARGS_LITERAL_EXPRESSION = astNodeType("VAR_ARGS_LITERAL_EXPRESSION")
	RANGE_EXPRESSION            = astNodeType("RANGE_EXPRESSION")
)

const runtimeHeaderFile string = "runtime/include/runtime.h"
const runtimeIncludeDir = "runtime/include"
const runtimeSrcDir = "runtime/src"
const runtimeCMakeListsTxt = "runtime/CMakeLists.txt"

var funcs template.FuncMap = map[string]any{
	"Transpile": Transpile,
}

func loadTemplate(nodeType astNodeType, filename string) {
	data, err := templateFS.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read template file: %q", filename)
	}
	templates[nodeType] = template.Must(template.New(filename).Funcs(funcs).Parse(strings.TrimSpace(string(data))))
}

func init() {
	data, err := templateFS.ReadFile(runtimeHeaderFile)
	if err != nil {
		log.Fatalf("Unable to read template file: %q", runtimeHeaderFile)
	}
	preamble = string(data)

	loadTemplate(PROGRAM, "runtime/templates/program.cpp")
	loadTemplate(LET_STATEMENT, "runtime/templates/let_statement.cpp")
	loadTemplate(EXPRESSION_STATEMENT, "runtime/templates/expression_statement.cpp")
	loadTemplate(IDENTIFIER_EXPRESSION, "runtime/templates/identifier_expression.cpp")
	loadTemplate(INTEGER_LITERAL_EXPRESSION, "runtime/templates/integer_literal_expr.cpp")
	loadTemplate(CALL_EXPRESSION, "runtime/templates/call_expr.cpp")
	loadTemplate(FN_LITERAL_EXPRESSION, "runtime/templates/fn_literal_expr.cpp")
	loadTemplate(BLOCK_STATEMENT, "runtime/templates/block_statement.cpp")
	loadTemplate(STRING_LITERAL_EXPRESSION, "runtime/templates/string_literal_expr.cpp")
	loadTemplate(BOOL_LITERAL_EXPRESSION, "runtime/templates/bool_literal_expr.cpp")
	loadTemplate(PREFIX_EXPRESSION, "runtime/templates/prefix_expr.cpp")
	loadTemplate(INFIX_EXPRESSION, "runtime/templates/infix_expr.cpp")
	loadTemplate(IF_EXPRESSION, "runtime/templates/if_expr.cpp")
	loadTemplate(RETURN_STATEMENT, "runtime/templates/return_statement.cpp")
	loadTemplate(ARRAY_LITERAL_EXPRESSION, "runtime/templates/array_literal_expr.cpp")
	loadTemplate(INDEX_OPERATOR_EXPRESSION, "runtime/templates/index_operator_expr.cpp")
	loadTemplate(VAR_ARGS_LITERAL_EXPRESSION, "runtime/templates/var_args_literal_expr.cpp")
	loadTemplate(RANGE_EXPRESSION, "runtime/templates/range_expr.cpp")
}

var indent int = 0
var traceStack []astNodeType

func trace(nodeType astNodeType) {
	log.Println(strings.Repeat("  ", indent), "+", nodeType)
	indent += 1
	traceStack = append(traceStack, nodeType)
}

func untrace(nodeType astNodeType) {
	if traceStack[len(traceStack)-1] != nodeType {
		panic("trace error")
	}

	indent -= 1
	log.Println(strings.Repeat("  ", indent), "-", nodeType)
	traceStack = traceStack[:len(traceStack)-1]
}

func execTemplate(nodeType astNodeType, node ast.Node) string {
	// trace(nodeType)
	// defer untrace(nodeType)

	var buffer bytes.Buffer

	err := templates[nodeType].Execute(&buffer, node)
	if err != nil {
		log.Fatalf("Error parsing template %q", err.Error())
	}

	return buffer.String()
}

func checkErr(err error, fmt string, args ...interface{}) {
	if err != nil {
		log.Fatalf(fmt, args...)
	}
}

func expandBuildEnv() string {
	tmpDir, err := os.MkdirTemp("", "monkey_transpiler_*")
	checkErr(err, "Unable to create temporary dir for monkey: %v", err)

	expandFile := func(srcFilePath, dstFilePath string) {
		srcFile, err := templateFS.Open(srcFilePath)
		checkErr(err, "Unable to read file: %v", err)
		defer srcFile.Close()

		dstFile, err := os.Create(dstFilePath)
		checkErr(err, "Unable to open temp file: %v", err)
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		checkErr(err, "Unable to copy file: %v", err)
	}

	expandDir := func(srcDir, dstDir string) {
		checkErr(os.Mkdir(dstDir, 0755), "Unable to create dir for monkey: %v", err)

		entries, err := templateFS.ReadDir(srcDir)
		checkErr(err, "Unable to read dir: %v", err)
		for _, entry := range entries {
			if !entry.IsDir() {
				srcFilePath := filepath.Join(srcDir, entry.Name())
				dstFilePath := filepath.Join(dstDir, entry.Name())
				expandFile(srcFilePath, dstFilePath)
			}
		}
	}

	expandDir(runtimeIncludeDir, filepath.Join(tmpDir, "include"))
	expandDir(runtimeSrcDir, filepath.Join(tmpDir, "src"))
	expandFile(runtimeCMakeListsTxt, filepath.Join(tmpDir, "CMakeLists.txt"))

	return tmpDir
}

func Compile(program string) string {
	tmpDir := expandBuildEnv()

	file, err := os.Create(filepath.Join(tmpDir, "main.cpp"))
	if err != nil {
		log.Fatalf("Unable to create temporary cpp file: %v", err)
	}

	_, err = file.WriteString(program)
	if err != nil {
		log.Fatalf("Unable to write temporary cpp file: %v", err)
	}

	file.Close()

	buildDir := filepath.Join(tmpDir, "build")
	cmd := exec.Command("cmake", "-S", tmpDir, "-B", buildDir, "-DCMAKE_BUILD_TYPE=release", "-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=true", "-G", "Ninja")
	log.Printf("Running command %q\n", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		log.Fatalf("Error running cmake: %v", err)
	}

	cmd = exec.Command("cmake", "--build", buildDir)
	log.Printf("Running command %q\n", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		log.Fatalf("Error compiling program: %v", err)
	}

	builtProgram := filepath.Join(buildDir, "main")

	var buffer bytes.Buffer
	cmd = exec.Command(builtProgram)
	log.Printf("Running command %q\n", cmd.String())
	cmd.Stderr = &buffer
	cmd.Stdout = &buffer
	if err = cmd.Run(); err != nil {
		log.Println(buffer.String())
		log.Fatalf("Error running compiled program: %v", err)
	}

	return buffer.String()
}

func Transpile(node ast.Node) string {
	switch node := node.(type) {
	case *ast.Program:
		return execTemplate(PROGRAM, node)
	case *ast.LetStatement:
		return execTemplate(LET_STATEMENT, node)
	case *ast.ExpressionStatement:
		return execTemplate(EXPRESSION_STATEMENT, node)
	case *ast.IdentifierExpr:
		return execTemplate(IDENTIFIER_EXPRESSION, node)
	case *ast.IntegerLiteralExpr:
		return execTemplate(INTEGER_LITERAL_EXPRESSION, node)
	case *ast.CallExpr:
		return execTemplate(CALL_EXPRESSION, node)
	case *ast.FnLiteralExpr:
		return execTemplate(FN_LITERAL_EXPRESSION, node)
	case *ast.BlockStatement:
		return execTemplate(BLOCK_STATEMENT, node)
	case *ast.StringLiteralExpr:
		return execTemplate(STRING_LITERAL_EXPRESSION, node)
	case *ast.BoolLiteralExpr:
		return execTemplate(BOOL_LITERAL_EXPRESSION, node)
	case *ast.PrefixExpr:
		return execTemplate(PREFIX_EXPRESSION, node)
	case *ast.InfixExpr:
		return execTemplate(INFIX_EXPRESSION, node)
	case *ast.IfExpr:
		return execTemplate(IF_EXPRESSION, node)
	case *ast.ReturnStatement:
		return execTemplate(RETURN_STATEMENT, node)
	case *ast.ArrayLiteralExpr:
		return execTemplate(ARRAY_LITERAL_EXPRESSION, node)
	case *ast.IndexOperatorExpr:
		return execTemplate(INDEX_OPERATOR_EXPRESSION, node)
	case *ast.VarArgsLiteralExpr:
		return execTemplate(VAR_ARGS_LITERAL_EXPRESSION, node)
	case *ast.RangeExpr:
		return execTemplate(RANGE_EXPRESSION, node)
	default:
		log.Fatalf("Unsupported node type: %T\n", node)
	}

	return ""
}
