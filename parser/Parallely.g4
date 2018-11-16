grammar Parallely;

/*
 * Parser Rules
 */
typequantifier : APPROXTYPE | PRECISETYPE;
fulltype : typequantifier INTTYPE | typequantifier FLOATTYPE | typequantifier BOOLTYPE;
processid : INT | VAR;
// processset : VAR '=' {' processid (',' processid)+ '}';
probability : FLOAT # floatprob
    | VAR # varprob
    ;

var : VAR;
        
expression : INT # literal
    | var # variable
    | expression MULTIPLY expression # multiply
    | expression DIVISION expression # divide
    | expression ADD expression # add
    | expression MINUS expression # minus
    | expression GREATER expression # greater
    // | expression '[' probability ']' expression # prob
    ;

declaration : fulltype var # singledeclaration
    | fulltype '[' INT ']' var # arraydec
    | declaration ';' declaration # multipledeclaration
    ;

globaldec : GLOBALVAR '=' '{' processid (',' processid)+ '}' # singleglobaldec
    | globaldec ';' globaldec # multipleglobaldec
    ;

statement : SKIPSTATEMENT # skipstatement
    | statement ';' statement (';')? # seqcomposition
    | var '[' expression ']' ASSIGNMENT expression # arrayassignment
    | var ASSIGNMENT var '[' expression ']' # arrayload
    | var ASSIGNMENT '(' fulltype ')' var # cast
    | var ASSIGNMENT expression # expassignment
    | var ASSIGNMENT expression '[' probability ']' expression # probassignment
    | IF var THEN '{' statement '}' ELSE '{' statement '}' # if
    | SEND '(' processid ',' fulltype ',' var ')' # send
    | CONDSEND '(' var ',' processid ',' fulltype ',' var ')' # condsend
    | var ASSIGNMENT RECEIVE '(' processid ',' fulltype ')' # receive
    | var ',' var ASSIGNMENT CONDRECEIVE '(' processid ',' fulltype ')' # condreceive
    | FOR VAR IN GLOBALVAR  DO '{' statement '}' # forloop
    | REPEAT INT '{' statement (';')? '}' # repeat
    | var ASSIGNMENT VAR '(' (expression)? ')' # func
    ;

parallelprogram : processid ':' '[' declaration ';' statement (';')? ']' # singleprogram
    | VAR IN GLOBALVAR ':' '[' declaration ';' statement ']' # groupedprogram
    | parallelprogram '||' parallelprogram # parcomposition    
    ;

program : (globaldec ';')? parallelprogram #single
    ;

sequentialprogram : (globaldec ';')? declaration ';' statement #sequential
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
CONDSEND            : C O N D S E N D;
RECEIVE             : R E C E I V E;
CONDRECEIVE         : C O N D R E C E I V E;
FOR                 : F O R;
IN                  : I N;
DO                  : D O;
REPEAT              : R E P E A T;

TRUE : 'true';
FALSE : 'false';

ASSIGNMENT          : '=';

INT                 : [0-9] +;
FLOAT               : [0-9]+ '.' [0-9]+;

INTTYPE            : I N T;
INTTHIRTYTWOTYPE   : I N T '3' '2';
INTSIXTYPE         : I N T '6' '4';
FLOATTYPE          : F L O A T;
FLOATTYPETWO       : F L O A T '6' '4';
FLOATTYPETHREE     : F L O A T '3' '2';
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
GLOBALVAR           : [A-Z] [_0-9A-Za-z]*;

WHITESPACE          : [ \t\r\n\f]+ -> channel(HIDDEN);
