grammar Parallely;

/*
 * Parser Rules
 */
typequantifier : APPROXTYPE | PRECISETYPE;
fulltype : typequantifier INTTYPE | typequantifier FLOATTYPE | typequantifier BOOLTYPE;
processid : INT;
processset : '{' (processid ',')+ '}';

var : VAR;
        
expression : INT # literal
    | var # variable
    | expression MULTIPLY expression # multiply
    | expression DIVISION expression # divide
    | expression ADD expression # add
    | expression MINUS expression # minus
    ;

boolexpression : TRUE # true
    | FALSE # false
    | var # boolvariable
    | expression EQUAL expression # equal
    | expression GREATER expression # greater
    | expression LESS expression # less
    | boolexpression AND boolexpression # and
    | boolexpression OR boolexpression # or
    | NOT boolexpression # not
    ;

declaration : fulltype var # singledeclaration
    | declaration ';' declaration # multipledeclaration
    ;

statement : SKIPSTATEMENT # skipstatement
    | statement ';' statement # seqcomposition
    // | '{' statement '}' # block
    | var ASSIGNMENT expression # expassignment
    | var ASSIGNMENT boolexpression # boolassignment
    | IF boolexpression THEN '{' statement '}' ELSE '{' statement '}' # if
    | SEND '(' processid ',' fulltype ',' var ')' # send
    | var ASSIGNMENT RECEIVE '(' processid ',' fulltype ')' # receive
    ;

parallelprogram : processid ':' '[' declaration ';' statement ']' # singleprogram
    | processset ':' '[' declaration ';' statement ']' # groupedprogram
    ;

program : parallelprogram #single
    | parallelprogram '||' parallelprogram # parcomposition
    ;

        
/*
 * Lexer Rules
 */
fragment A : [aA];
fragment B : [bB];
fragment C : [cC];
fragment D : [dD];
fragment E : [eE];
fragment F : [fF];
fragment G : [gG];
fragment H : [hH];
fragment I : [iI];
fragment J : [jJ];
fragment K : [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];

SKIPSTATEMENT       : S K I P;
IF                  : I F;
THEN                : T H E N;
ELSE                : E L S E;
SEND                : S E N D;
RECEIVE             : R E C E I V E;

TRUE : 'true';
FALSE : 'false';

ASSIGNMENT          : '=';

INT                 : [0-9] +;
INTTYPE            : I N T;
FLOATTYPE          : F L O A T;
BOOLTYPE           : B O O L;
PRECISETYPE        : P R E C I S E;
APPROXTYPE         : A P P R O X;

ADD                 : '+';
MINUS               : '-';
MULTIPLY            : '*';
DIVISION            : '/';

EQUAL               : '==';
GREATER             : '>';
LESS                : '<';
NOT                 : '!';
AND                 : '&';
OR                  : '|';

VAR                 : [a-z] [_0-9A-Za-z]*;

WHITESPACE          : [ \t\r\n\f]+ -> channel(HIDDEN);