grammar Expr;

prog: stat+ ;

stat: expr NEWLINE
    | ID '=' expr NEWLINE
    | NEWLINE
    ;

expr: expr ('*'| '/') expr
    | expr ('+'| '-') expr
    | INT
    | ID
    | '(' expr ')'
    ;

ID : [a-zA-Z]+ ; //匹配标识符
INT: [0-9]+ ; //匹配整数
NEWLINE: ';' ; //匹配终止符
WS : [ \r\t\n]+ -> skip ; //丢弃空白符