package main

import (
	"demo-go/parsing"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cast"
	"log"
	"sync"
	"testing"
)

var varMap = sync.Map{}

func TestExpr(t *testing.T) {

	//接收一个文件流
	inputStream, err := antlr.NewFileStream("./samples/expr_1")
	if err != nil {
		log.Fatal(err)
	}

	//从文件流中创建一个词法解析器
	lexer := parsing.NewExprLexer(inputStream)

	//从词法解析器中创建一个Token Stream
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)

	//从Token Stream中创建一个语法解析
	p := parsing.NewExprParser(tokenStream)

	//添加错误监听器
	p.AddErrorListener(&ErrListener{})

	//开始解析并生成语法树
	tree := p.Prog()

	//遍历树，并设置监听器
	antlr.NewParseTreeWalker().Walk(NewListener(), tree)
}

func NewListener() *TreeListener {
	return &TreeListener{}
}

type TreeListener struct {
	*parsing.BaseExprListener
}

func (l *TreeListener) EnterStat(expr *parsing.StatContext) {
	fmt.Println("enter state", GetPayloadTextFromTree(expr))
	v := compute(expr.GetChildren())
	fmt.Println(v)
}

func compute(trees []antlr.Tree) int {

	if len(trees) == 1 {
		tree := trees[0]
		if tree.GetChildCount() == 1 {
			v := GetPayloadTextFromTree(tree)
			if value, ok := varMap.Load(v); ok {
				return value.(int)
			}
			return cast.ToInt(GetPayloadTextFromTree(tree))
		}
		return compute(tree.GetChildren())
	}
	if len(trees) == 2 {
		return compute([]antlr.Tree{trees[0]})
	}
	if len(trees) == 3 {
		left, mid, right := trees[0], trees[1], trees[2]
		if v, ok := left.GetPayload().(*antlr.CommonToken); ok && v.GetText() == "(" {
			return compute(mid.GetChildren())
		}
		operator := GetPayloadTextFromTree(mid)
		switch operator {
		case "+":
			return compute([]antlr.Tree{left}) + compute([]antlr.Tree{right})
		case "-":
			return compute([]antlr.Tree{left}) - compute([]antlr.Tree{right})
		case "*":
			return compute([]antlr.Tree{left}) * compute([]antlr.Tree{right})
		case "/":
			return compute([]antlr.Tree{left}) / compute([]antlr.Tree{right})
		default:
			break
		}
	}
	if len(trees) == 4 {
		id, expr := trees[0], trees[2]
		varName := GetPayloadTextFromTree(id)
		r := compute([]antlr.Tree{expr})
		varMap.Store(varName, r)
		return r
	}
	panic("syntax error,children is invalid")
}

func GetPayloadTextFromTree(tree antlr.Tree) string {

	payload := tree.GetPayload()

	if v, ok := payload.(*antlr.CommonToken); ok {
		return v.GetText()
	}
	if v, ok := payload.(*antlr.BaseParserRuleContext); ok {
		return v.GetText()
	}
	panic("syntax error,payload is invalid")
}

type ErrListener struct {
}

func (ErrListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	panic(msg)
}

func (ErrListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (ErrListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (ErrListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}
