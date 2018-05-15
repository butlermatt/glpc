## Syntax Grammar

```glpc
program        → import
               | declaration* EOF ;
```

### Imports

```glpc
import         → STRING ";" ;
```

### Declarations

```glpc
declaration    → classDecl
               | fnDecl
               | varDecl ;

classDecl      → "class" IDENTIFIER ( ":" IDENTIFIER )?
                 "{" function* "}" ;
fnDecl         → "fn" function ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
```

### Utility Rules

```glpc
function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
arguments      → expression ( "," expression )* ;
```

### Statements

```glpc
statement      → block
               | doWhileStmt
               | forStmt
               | ifStmt
               | printStmt
               | returnStmt
               | whileStmt
               | exprStmt ;

block          → "{" ( declaration | statement )* "}" ;
exprStmt       → expression ";" ;
doWhileStmt    → "do" statement "while" "(" expression ")" ";" ;
forStmt        → "for" "(" ( varDecl | exprStmt )
                 expression? ";" expression? ")" statement ;
ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
printStmt      → "print" expression ";" ;
returnStmt     → "return" expression? ";" ;
whileStmt      → "while" "(" expression ")" statement ;
```

### Expressions

```glpc
expression     → assignment ;

assignment     → ( call "." )? IDENTIFIER 
                 ("=" | "+=" | "-=" | "*=" | "/=" | "%=" | "~/=") 
                 assignment 
               | logic_or;

logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "*" | "/" | "~/" | "%" ) unary )* ;

unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
primary        → "true" | "false" | "null" | "this"
               | NUMBER | STRING | IDENTIFIER | "(" expression ")"
               | "super" "." IDENTIFIER ;
```

## Lexical Grammar

The lexical grammar is used by the scanner to group characters into tokens.
Where the syntax is [context free][], the lexical grammar is [regular][] -- note
that there are no recursive rules.

[context free]: https://en.wikipedia.org/wiki/Context-free_grammar
[regular]: https://en.wikipedia.org/wiki/Regular_grammar

```glpc
NUMBER         → DIGIT+ ( "." DIGIT+ )? ;
STRING         → '"' <any char except '"' or '\n'>* '"' 
               | "'" <any character except "'" or "\n">* "'"
               | "`" <any character except "`">* "`" 
IDENTIFIER     → ALPHA ( ALPHA | DIGIT )* ;
ALPHA          → 'a' ... 'z' | 'A' ... 'Z' | '_' ;
DIGIT          → '0' ... '9' ;
```