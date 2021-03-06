package parser

import (
	"DiniSQL/MiniSQL/src/Interpreter/lexer"
	"DiniSQL/MiniSQL/src/Interpreter/types"
	"log"
)

type lexerWrapper struct {
	impl        *lexer.LexerImpl
	channelSend chan<- types.DStatements
	lastLiteral string
	err         error
}

func newLexerWrapper(li *lexer.LexerImpl, channel chan<- types.DStatements) *lexerWrapper {
	return &lexerWrapper{
		impl:        li,
		channelSend: channel,
	}
}

func (l *lexerWrapper) Lex(lval *yySymType) int {
	r, err := l.impl.Lex(lval.LastToken)
	if err != nil {
		log.Fatal(err)
	}
	l.lastLiteral = r.Literal

	tokVal := r.Token
	lval.str = r.Literal
	lval.LastToken = tokVal
	return tokVal
}

func (l *lexerWrapper) Error(errStr string) {
	l.err = wrapParseError(l.lastLiteral, errStr)
}
