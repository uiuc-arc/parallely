grammar Parallely;

/*
 * Parser Rules
 */
typequantifier : APPROXTYPE | PRECISETYPE | DYNTYPE;

basictype : typequantifier INTTYPE
    | typequantifier FLOATTYPE
    | typequantifier BOOLTYPE
    | typequantifier INTTHIRTYTWOTYPE
    | typequantifier INTSIXTYPE
    | typequantifier FLOATTYPETWO
    | typequantifier FLOATTYPETHREE
    ;

fulltype : basictype #singletype
    | basictype '[' ']' #arraytype
    ;

probability : FLOAT # floatprob
    | VAR # varprob
    ;

var : VAR # localvariable
    | GLOBALVAR # globalvariable
    | VAR '$' processid # threadvariable
    ;

fvar : VAR;

processid : INT # namedp
    | VAR # variablep
    | VAR IN GLOBALVAR # groupedp;

expression : INT # literal
    | FLOAT # fliteral
    | var # variable
    // | GLOBALVAR # globalvariable
    // | var ('[' expression ']')+ #arrayvar
    | expression MULTIPLY expression # multiply
    | expression DIVISION expression # divide
    | expression ADD expression # add
    | expression MINUS expression # minus
    | expression GREATER expression # greater
    | expression LESS expression # less
    | expression GEQ expression # geq
    | expression LEQ expression # leq
    | expression EQUAL expression # eq
    | '(' expression ')' # select
    | expression AND expression #andexp
    | expression OR expression #orexp    
    ;

declaration : basictype var # singledeclaration
    | basictype ('[' (INT)? ']')+ var # arraydeclaration
    | basictype ('[' (GLOBALVAR)? ']')+ var # dynarraydeclaration
    | declaration '@' INT (',' INT)* # annotateddec
    // | declaration ';' declaration # multipledeclaration
    ;

globaldec : GLOBALVAR '=' '{' processid (',' processid)* '}' # singleglobaldec
    | basictype GLOBALVAR # globalconst
    | basictype '[' (INT)? ']' GLOBALVAR # globalarray
    | EXTERN basictype GLOBALVAR # globalexternal
    | globaldec '@' INT (',' INT)* # annotatedgdec        
    // | globaldec ';' globaldec # multipleglobaldec
    ;

statement : SKIPSTATEMENT # skipstatement
    | var ('[' expression ']')+ ASSIGNMENT expression # arrayassignment
    | var ASSIGNMENT var ('[' expression ']')+ # arrayload
    | var ASSIGNMENT '(' fulltype ')' var # cast
    | var ASSIGNMENT expression # expassignment
    | GLOBALVAR ASSIGNMENT expression # gexpassignment
    | var ASSIGNMENT precise=expression '[' probability ']' approx=expression # probassignment
    | APPROXIMATE '(' var ',' expression ')' # approximate
    | var ASSIGNMENT condition=var '?' ifvar=var elsevar=var # condassignment
    | assigned=var ASSIGNMENT '(' lvar=var GEQ rvar=var ')' '?' ifvar=var elsevar=var # dyncondassignmentgeq
    | IF var THEN '{' (ifs+=statement ';')+ '}' # ifonly
    | IF var THEN '{' (ifs+=statement ';')+ '}' ELSE '{' (elses+=statement ';')+ '}' # if
    | SEND '(' processid ',' fulltype ',' var ')' # send
    | CONDSEND '(' var ',' processid ',' fulltype ',' var ')' # condsend
    | DYNSEND '(' processid ',' fulltype ',' var ')' # dynsend
    | DYNCONDSEND '(' var ',' processid ',' fulltype ',' var ')' # dyncondsend
    | var ASSIGNMENT RECEIVE '(' processid ',' fulltype ')' # receive
    | var ',' var ASSIGNMENT CONDRECEIVE '(' processid ',' fulltype ')' # condreceive
    | var ASSIGNMENT DYNRECEIVE '(' processid ',' fulltype ')' # dynreceive
    | var ',' var ASSIGNMENT DYNCONDRECEIVE '(' processid ',' fulltype ')' # dyncondreceive
    | FOR VAR IN GLOBALVAR  DO '{' (statement ';')+ '}' # forloop
    | REPEAT INT '{' (statement ';')+ '}' # repeat
    | REPEAT var '{' (statement ';')+ '}' # repeatlvar
    | REPEAT GLOBALVAR '{' (statement ';')+ '}' # repeatvar
    | WHILE '(' cond=expression ')' '{' (body+=statement ';')+ '}' # while
    | avar+=var(','avar+=var)* ASSIGNMENT fname=fvar '('(invar+=var)?(',' invar+=var)*')' # func
    // | var ASSIGNMENT TRACK '(' var ',' probability ')' # track
    | var ASSIGNMENT TRACK '(' var ',' eps=FLOAT ',' delta=probability ')' # track
    | var ASSIGNMENT TRACK '(' var ',' eps=var ',' delta=var ')' # trackvar
    // | var ASSIGNMENT CHECK '(' var ',' probability ')' # check
    | CHECK '(' var ',' eps=FLOAT ',' delta=probability ')' # speccheck
    | assigned=var ASSIGNMENT CHECK '(' checkedvar=var ',' eps=FLOAT ',' delta=probability ')' # speccheckwithresult
    | CHECKARRAY '(' var ',' eps=FLOAT ',' delta=probability ')' # speccheckarray
    | code=COMMENT # instrument
    | '<' DUMMY INT '>' # dummy
    | TRY '{' (trys+=statement ';')+ '}'
        CHECK '{' check=expression '}'
        RECOVER '{' (recovers+=statement ';')+ '}' # recover
    | TRY '{' (trys+=statement ';')+ '}'
        CHECK '{' check=expression '}'
        RECOVERWITH(processid*) '{' (recovers+=statement ';')+ '}' # recoverwith
    | TRY '{' (trys+=statement ';')+ '}'
        RECOVERFROM(processid) '{' (recovers+=statement ';')+ '}' # recoverfrom        
    | statement '@' annotation+=INT (',' annotation+=INT)* # annotated
    ;

program : processid ':' '[' (declaration ';')*  (statement ';')+ ']' # single
    ;

parallelprogram : (funcspec ';')* (globaldec ';')* program ('||' program)* # parcomposition
    ;

sequentialprogram : (funcspec ';')* (globaldec ';')* (declaration ';')* (statement ';')+ # sequential
    ;

singlerelyspec : FLOAT LEQ (FLOAT '*')? 'R' '(' (VAR | GLOBALVAR) (',' (VAR | GLOBALVAR))* ')'
    ;

relyspec : singlerelyspec (AND singlerelyspec)*
    ;

interval : '[' FLOAT ',' FLOAT ']'
    ;

varchiselspec : var IN interval
    ;

funcchiselspec : var IN '<' INT ',' INT ',' FLOAT (',' FLOAT ',' interval)+ '>'
    ;

checkchiselspec : var ENSURES '<' expression GEQ var (',' expression GEQ var)* '>'
    ;

singlechiselspec : FLOAT LEQ 'R' '(' (FLOAT GEQ 'd' '(' var ')') (',' (FLOAT GEQ 'd' '(' var ')'))* ')'
    ;

chiselspec : singlechiselspec (AND singlechiselspec)* (varchiselspec | funcchiselspec | checkchiselspec)*
    ;

funcspec : SPEC funcname=var '(' funcargs+=var? (',' funcargs+=var)* ')' '=' '(' ACC ':' accexps+=expression? (',' accexps+=expression)* ',' REL ':' relspecs+=expression? (',' relspecs+=expression)* ')'
        // REL ':'  expression? (', ' expression)* ';' ')'
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
fragment Y: [yY];

SKIPSTATEMENT       : S K I P;
IF                  : I F;
THEN                : T H E N;
ELSE                : E L S E;
SEND                : S E N D;
CONDSEND            : C O N D S E N D;
DYNSEND             : D Y N S E N D;
DYNCONDSEND         : D Y N C O N D S E N D;
RECEIVE             : R E C E I V E;
CONDRECEIVE         : C O N D R E C E I V E;
DYNRECEIVE          : D Y N R E C E I V E;
DYNCONDRECEIVE      : D Y N C O N D R E C E I V E;
FOR                 : F O R;
IN                  : I N;
DO                  : D O;
REPEAT              : R E P E A T;
WHILE               : W H I L E;
TRACK               : T R A C K;
CHECK               : C H E C K;
CHECKARRAY          : C H E C K A R R A Y;
TRY                 : T R Y;
RECOVER             : R E C O V E R;
RECOVERWITH         : R E C O V E R W I T H;
RECOVERFROM         : R E C O V E R F R O M;
DUMMY               : D U M M Y;
APPROXIMATE         : A P P R O X I M A T E;
ENSURES             : E N S U R E S;
SPEC                : S P E C;
ACC                 : A C C;
REL                 : R E L;

TRUE : 'true';
FALSE : 'false';

ASSIGNMENT          : '=';

INT                 : ('-')?[0-9] +;
FLOAT               : ('-')?[0-9]+ '.' [0-9]+;

INTTYPE            : I N T;
INTTHIRTYTWOTYPE   : I N T '3' '2';
INTSIXTYPE         : I N T '6' '4';
FLOATTYPE          : F L O A T;
FLOATTYPETWO       : F L O A T '6' '4';
FLOATTYPETHREE     : F L O A T '3' '2';
BOOLTYPE           : B O O L;
PRECISETYPE        : P R E C I S E;
APPROXTYPE         : A P P R O X;
DYNTYPE            : D Y N A M I C;
EXTERN              : E X T E R N;

ADD                 : '+';
MINUS               : '-';
MULTIPLY            : '*';
DIVISION            : '/';

EQUAL               : '==';
GREATER             : '>';
LESS                : '<';
GEQ                 : '>=';
LEQ                 : '<=';
NOT                 : '!';
AND                 : '&&';
OR                  : '||';

VAR                 : [a-z] [._0-9A-Za-z]*;
GLOBALVAR           : [A-Z] [_0-9A-Za-z]*;

WHITESPACE          : [ \t\r\n\f]+ -> channel(HIDDEN);

COMMENT             : '##' ~('\r' | '\n')* '##';
