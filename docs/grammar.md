### ![](logo/logo.png) Грамматика

```
program = {definition};

definition =
  actor
  | actor class;
actor =
  "actor", identifier, "(", [identifier, {",", identifier}, [","]], ")",
    state, {state},
  ";";
actor class =
  "class", identifier, "(", [identifier, {",", identifier}, [","]], ")",
    state, {state},
  ";";
state =
  "state", identifier, "(", [identifier, {",", identifier}, [","]], ")",
    {message},
  ";";
message =
  "message", identifier, "(", [identifier, {",", identifier}, [","]], ")",
    {command},
  ";";

command =
  let command
  | start command
  | send command
  | set command
  | return command
  | expression;
let command = "let", identifier, "=", expression;
start command =
  "start", (identifier | "[", expression, "]"),
  "(", [expression, {",", expression}, [","]], ")";
send command = "send", identifier, "(", [expression, {",", expression}, [","]], ")";
set command = "set", identifier, "(", [expression, {",", expression}, [","]], ")";
return command = "return";

expression = list construction;
list construction = disjunction, [":", list construction];
disjunction = conjunction, ["||", disjunction];
conjunction = equality, ["&&", conjunction];
equality = comparison, [("==" | "!="), equality];
comparison = bitwise disjunction, [("<=" | "<" | ">=" | ">"), comparison];
bitwise disjunction = bitwise exclusive disjunction, ["|", bitwise disjunction];
bitwise exclusive disjunction = bitwise conjunction, ["^", bitwise exclusive disjunction];
bitwise conjunction = shift, ["&", bitwise conjunction];
shift = addition, [("<<" | ">>>" | ">>"), shift];
addition = multiplication, [("+" | "-"), addition];
multiplication = unary, [("*" | "/" | "%"), multiplication];
unary = (("-" | "~" | "!"), unary) | accessor;

accessor = atom, {accessor key};
accessor key =
  ".", identifier
  | "[", expression, "]";

atom =
  number
  | string
  | list definition
  | hash table definition
  | function call
  | conditional expression
  | identifier
  | "(", expression, ")";
number = INTEGER NUMBER | FLOATING-POINT NUMBER | SYMBOL;
string =
  SINGLE-QUOTED INTERPRETED STRING
  | DOUBLE-QUOTED INTERPRETED STRING
  | RAW STRING;
list definition = "[", [expression, {",", expression}, [","]], "]";
hash table definition = "{", [hash table entry, {",", hash table entry}, [","]], "}";
hash table entry = (identifier | "[", expression, "]"), ":", expression;
function call = identifier, "(", [expression, {",", expression}, [","]], ")";
conditional expression = "when", {conditional case}, ";";
conditional case = "=>", expression, {command};
identifier = IDENTIFIER - key words;
key words =
  "actor"
  | "class"
  | "state"
  | "message"
  | "let"
  | "start"
  | "send"
  | "set"
  | "return"
  | "when";

LINE COMMENT = ? /\/\/.*/ ?;
BLOCK COMMENT = ? /\/\*.*?\*\//s ?;
INTEGER NUMBER = ? /\b((0x[\da-f]+)|(0[0-7]+)|(\d+(e\d+)?)|(\d+e[\+\-]\d+))\b/i ?;
FLOATING-POINT NUMBER = ? /(\.\d+(e[\+\-]\d+)?)\b|\b\d+\.\d*((e[\+\-]\d+)?\b)?/i ?;
SYMBOL = ? /'(\\x[\da-f]{2}|\\.|[^'\n])'/i ?;
SINGLE-QUOTED INTERPRETED STRING = ? /'(\\x[\da-f]{2}|\\.|[^'\n])*?'/i ?;
DOUBLE-QUOTED INTERPRETED STRING = ? /"(\\x[\da-f]{2}|\\.|[^"\n])*?"/i ?;
RAW STRING = ? /`[^`]*?`/ ?
IDENTIFIER = ? /[a-z_]\w*/i ?;
```
