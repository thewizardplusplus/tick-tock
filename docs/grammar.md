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
command = send command | set command | out command | exit command;
send command = "send", identifier;
set command = "set", identifier;
out command = "out", string;
exit command = "exit";

string = INTERPRETED STRING | RAW STRING;
identifier = IDENTIFIER - key words;
key words = "actor" | "state" | "message" | "send" | "set" | "out" | "exit";

LINE COMMENT = ? /\/\/.*/ ?;
BLOCK COMMENT = ? /\/\*.*?\*\//s ?;
INTERPRETED STRING = ? /"(\\.|[^"\n])*?"/ ?;
RAW STRING = ? /`[^`]*?`/ ?
IDENTIFIER = ? /[a-z_]\w*/i ?;
```
