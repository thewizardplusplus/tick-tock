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
sleep command = "sleep", number, ",", number;
exit command = "exit";

number = INTEGER NUMBER | FLOATING-POINT NUMBER;
string = INTERPRETED STRING | RAW STRING;
identifier = IDENTIFIER - key words;
key words = "actor" | "state" | "message" | "send" | "set" | "out" | "exit";

LINE COMMENT = ? /\/\/.*/ ?;
BLOCK COMMENT = ? /\/\*.*?\*\//s ?;
INTEGER NUMBER = ? /\b((0x[\da-f]+)|(0[0-7]+)|(\d+(e\d+)?)|(\d+e[\+\-]\d+))\b/i ?;
FLOATING-POINT NUMBER = ? /(\.\d+(e[\+\-]\d+)?)\b|\b\d+\.\d*((e[\+\-]\d+)?\b)?/i ?;
INTERPRETED STRING = ? /"(\\x[\da-f]{2}|\\.|[^"\n])*?"/i ?;
RAW STRING = ? /`[^`]*?`/ ?
IDENTIFIER = ? /[a-z_]\w*/i ?;
```
