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
command = let command | send command | set command | expression;
let command = "let", identifier, "=", expression;
send command = "send", identifier;
set command = "set", identifier;

expression = list construction;
list construction = addition, [":", list construction];
addition = multiplication, [("+" | "-"), addition];
multiplication = unary, [("*" | "/" | "%"), multiplication];
unary = ("-", unary) | accessor;

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
string = INTERPRETED STRING | RAW STRING;
list definition = "[", [expression, {",", expression}, [","]], "]";
function call = identifier, "(", [expression, {",", expression}, [","]], ")";
identifier = IDENTIFIER - key words;
key words = "actor" | "state" | "message" | "let" | "send" | "set";

LINE COMMENT = ? /\/\/.*/ ?;
BLOCK COMMENT = ? /\/\*.*?\*\//s ?;
INTEGER NUMBER = ? /\b((0x[\da-f]+)|(0[0-7]+)|(\d+(e\d+)?)|(\d+e[\+\-]\d+))\b/i ?;
FLOATING-POINT NUMBER = ? /(\.\d+(e[\+\-]\d+)?)\b|\b\d+\.\d*((e[\+\-]\d+)?\b)?/i ?;
INTERPRETED STRING = ? /"(\\x[\da-f]{2}|\\.|[^"\n])*?"/i ?;
RAW STRING = ? /`[^`]*?`/ ?
IDENTIFIER = ? /[a-z_]\w*/i ?;
```
