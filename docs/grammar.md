### ![](logo/logo.png) Грамматика

```
program = actor list;

actor list = actor, {actor};
actor = "actor", state list, ";";

state list = state, {state};
state = "state", identifier, message list, ";";

message list = message, {message};
message = "message", identifier, command list, ";";

command list = command, {command};
command =
  let command
  | send command
  | set command
  | return command
  | expression;
let command = "let", identifier, "=", expression;
send command = "send", identifier;
set command = "set", identifier;
return command = "return";

expression = list construction;
list construction = disjunction, [":", list construction];
disjunction = conjunction, ["||", disjunction];
conjunction = equality, ["&&", conjunction];
equality = comparison, [("==" | "!="), equality];
comparison = addition, [("<=" | "<" | ">=" | ">"), comparison];
addition = multiplication, [("+" | "-"), addition];
multiplication = unary, [("*" | "/" | "%"), multiplication];
unary = (("-" | "!"), unary) | accessor;

accessor = atom, {list item access};
list item access = "[", expression, "]";

atom =
  number
  | string
  | list definition
  | function call
  | identifier
  | ("(", expression, ")");
number = INTEGER NUMBER | FLOATING-POINT NUMBER;
string =
  SINGLE-QUOTED INTERPRETED STRING
  | DOUBLE-QUOTED INTERPRETED STRING
  | RAW STRING;
list definition = "[", [expression, {",", expression}, [","]], "]";
function call = identifier, "(", [expression, {",", expression}, [","]], ")";
identifier = IDENTIFIER - key words;
key words = "actor" | "state" | "message" | "let" | "send" | "set" | "return";

LINE COMMENT = ? /\/\/.*/ ?;
BLOCK COMMENT = ? /\/\*.*?\*\//s ?;
INTEGER NUMBER = ? /\b((0x[\da-f]+)|(0[0-7]+)|(\d+(e\d+)?)|(\d+e[\+\-]\d+))\b/i ?;
FLOATING-POINT NUMBER = ? /(\.\d+(e[\+\-]\d+)?)\b|\b\d+\.\d*((e[\+\-]\d+)?\b)?/i ?;
SINGLE-QUOTED INTERPRETED STRING = ? /'(\\x[\da-f]{2}|\\.|[^'\n])*?'/i ?;
DOUBLE-QUOTED INTERPRETED STRING = ? /"(\\x[\da-f]{2}|\\.|[^"\n])*?"/i ?;
RAW STRING = ? /`[^`]*?`/ ?
IDENTIFIER = ? /[a-z_]\w*/i ?;
```
