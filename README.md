# pogo
Another parser combinator library for Go

## How do I use it?

 - [Describe your grammar](https://github.com/tgascoigne/pogo/blob/master/json/grammar.go) using the provided combinators, or extend with your own
 - Provide an implementation of both [parser](https://github.com/tgascoigne/pogo/blob/master/json/parser.go) and [lexer](https://github.com/tgascoigne/pogo/blob/master/json/lexer.go) interfaces
 - `pogo.Do` will now produce a parse tree
 - (Optional) Use the provided generator to produce an easy to use parse tree visitor, and [implement it](https://github.com/tgascoigne/pogo/blob/master/json/ast_builder.go) to translate your parse tree to a higher level representation.
 
 TL;DR take a look at the provided [json parser](https://github.com/tgascoigne/pogo/tree/master/json)
 
 ## Should I use it?
 
 Sure, if you want. I can't guarantee maintainance, though!
