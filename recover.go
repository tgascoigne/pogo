package pogo

type RecoverFunc func(start ParseState, end ParseState) (Parsed, ParseState)

// Recover parses p. If any errors are raised while parsing p, recoverfn is called and its
// results become the result of Recover. Errors are sent to the error handler.
func Recover(p interface{}, recoverfn RecoverFunc) Parser {
	return doRecover(p, recoverfn, true)
}

// TryRecover parses p. If any errors are raised while parsing p, recoverfn is called and its
// results become the result of Recover. Errors are not sent to the error handler.
func TryRecover(p interface{}, recoverfn RecoverFunc) Parser {
	return doRecover(p, recoverfn, false)
}

// Try parses p. If any errors are raised while parsing p, NilParsed is returned.
// Errors are not sent to the error handler.
func Try(p interface{}) Parser {
	return TryRecover(p, func(start ParseState, end ParseState) (Parsed, ParseState) {
		return NilParsed, end
	})
}

// Root parses p. If any errors are raised while parsing p, ErrParsed is returned.
// Errors are sent to the error handler.
func Root(p interface{}) Parser {
	return Recover(p, func(start ParseState, end ParseState) (Parsed, ParseState) {
		return ErrParsed, end
	})
}

// RecoverTo parses p. If any errors are raised while parsing p, the lexer is fast-fowarded
// until after the next instance of token. Errors are sent to the error handler.
func RecoverTo(token TokenType, p interface{}) Parser {
	return Recover(p, func(start ParseState, end ParseState) (Parsed, ParseState) {
		var item Item
		ps := end
		for {
			item, ps = nextTok(ps)
			if item.Token == token || item.Token == TOK_EOF {
				break
			}
		}
		return NilParsed, ps
	})
}

func doRecover(p interface{}, recoverfn RecoverFunc, reportError bool) Parser {
	return func(ps ParseState) (Parsed, ParseState) {
		var newState ParseState
		var result Parsed
		result, newState = Do(p, ps)
		// if the parser failed, pass the start state and the new state to recoverfn
		// and allow it to give us a new result and state
		if result == ErrParsed {
			result, newState = recoverfn(ps, newState)
		}
		return result, newState
	}
}
