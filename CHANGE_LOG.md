# Change Log

## [v1.4](https://github.com/thewizardplusplus/tick-tock/tree/v1.4) (2020-06-24)

- support state parameters;
- support message parameters:
  - fix texts of errors on message processing;
- examples:
  - improve the example with dining philosophers;
  - add the example with the random counter.

## [v1.3](https://github.com/thewizardplusplus/tick-tock/tree/v1.3) (2020-06-23)

- support message parameters;
- support expression evaluation:
  - support operations:
    - optimize evaluation of conditional expressions;
- examples:
  - improve the example with dining philosophers;
  - improve the example with the guessing game;
  - improve the example with pi evaluation;
  - improve the example with ping-pong.

## [v1.2](https://github.com/thewizardplusplus/tick-tock/tree/v1.2) (2020-06-15)

- support commands:
  - return command (early exit from message processing);
- support expression evaluation:
  - support types:
    - symbol constants (based on real numbers);
  - support operations:
    - conditional expression (`when...;`);
- examples:
  - improve the example with the guessing game;
  - add the new example with the [maze](https://10print.org/).

## [v1.2-alpha](https://github.com/thewizardplusplus/tick-tock/tree/v1.2-alpha) (2020-05-19)

- support expression evaluation:
  - support types:
    - boolean (based on other types);
  - support operations:
    - lazy disjunction (`||`);
    - lazy conjunction (`&&`);
    - equality:
      - equal (`==`);
      - not equal (`!=`);
    - comparison:
      - less (`<`);
      - less or equal (`<=`);
      - greater (`>`);
      - greater or equal (`>=`);
    - logical negation (`!`);
- рантайм:
  - логические константы:
    - `false` — число 0;
    - `true` — число 1;
  - функции:
    - функции для работы с логическими значениями:
      - `bool(value: any): bool` — преобразует значение в логический тип: возвращает строго 0 или 1;
    - математические функции:
      - `is_nan(x: num): bool`;
    - функции для работы со строками:
      - `strb(value: any): str` — преобразует переданное значение в строку, как логическое: если `value` истинно, возвращает `"true"`, иначе — `"false"`;
    - функции для ввода/вывода:
      - исправлена функция `inln(count: num): nil|str`:
        - исправлено описание функции (отсутствие символа `'\n'` при отрицательном `count` на самом деле трактуется, как ошибка);
        - исправлена работа функции (символ `'\n'` больше не добавляется в результирующую строку);
- add the new example with the guessing game.

## [v1.1](https://github.com/thewizardplusplus/tick-tock/tree/v1.1) (2020-04-24)

- commands:
  - support:
    - variable definition;
    - expression evaluation;
  - remove:
    - `out` command;
    - `sleep` command;
    - `exit` command;
- support expression evaluation:
  - support types:
    - nil;
    - real numbers;
    - LISP-like lists;
    - strings (based on lists);
  - support operations:
    - for numbers:
      - addition (`+`);
      - subtraction (`-`);
      - negation (unary `-`);
      - multiplication (`*`);
      - division (`/`);
      - modulo (`%`);
    - for lists:
      - in-place definition (`[..., ..., ...]`);
      - construction from a head and a tail (`... : ...`);
      - concatenation (`+`);
      - indexing (`...[...]`);
    - function call;
- рантайм:
  - константы:
    - общие константы:
      - `nil: nil` — значение нулевого типа;
    - математические константы:
      - `nan: num`;
      - `inf: num` — положительная бесконечность;
      - `pi: num`;
      - `e: num`;
  - функции:
    - общие функции:
      - `type(value: any): str` — возвращает имя типа значения `value`;
    - математические функции:
      - `floor(x: num): num`;
      - `ceil(x: num): num`;
      - `trunc(x: num): num`;
      - `round(x: num): num`;
      - `sin(x: num): num`;
      - `cos(x: num): num`;
      - `tn(x: num): num`;
      - `arcsin(x: num): num`;
      - `arccos(x: num): num`;
      - `arctn(x: num): num`;
      - `angle(x: num, y: num): num` — atan2;
      - `pow(base: num, exponent: num): num`;
      - `sqrt(x: num): num`;
      - `exp(x: num): num`;
      - `ln(x: num): num`;
      - `lg(x: num): num`;
      - `abs(x: num): num`;
      - генерация псевдослучайных чисел:
        - `seed(seed: num): nil` — устанавливает начальное состояние генератора псевдослучайных чисел;
        - `random(): num` — возвращает псевдослучайное число в диапазоне [0; 1);
    - функции для работы со списками:
      - `size(list: list<any>): num` — возвращает размер (длину) списка `list`;
      - `head(list: list<any>): any` — возвращает голову списка `list`; список не должен быть пустым;
      - `tail(list: list<any>): list<any>` — возвращает хвост списка `list`; список не должен быть пустым;
    - функции для работы со строками:
      - `num(text: str): nil|num` — парсит число из строки `text`; при ошибке парсинга будет возвращён `nil`;
      - `str(value: any): str` — преобразует значение `value` в строку;
      - `strs(text: str): str` — преобразует строку `text` в другую строку, экранируя её символы и окружая всю строку кавычками;
      - `strl(list: list<str>): str` — преобразует список строк `list` в строку, отображая при этом строки как строки;
    - системные функции:
      - `env(name: str): nil|str` — возвращает значение указанной переменной окружения, если она установлена; в противном случае возвращается `nil`;
      - `time(): num` — возвращает текущее UNIX-время по UTC в секундах;
      - `sleep(duration: num): nil` — останавливает выполнение скрипта на `duration` секунд; `duration` может быть вещественным числом;
      - `exit(exitCode: num): nil` — завершает выполнение скрипта; код возврата будет равен `exitCode`;
      - функции для ввода/вывода:
        - `in(count: num): nil|str` — считывает `count` символов из stdin и возвращает их в виде строки; если количество будет отрицательным, то будут считаны все доступные символы; при ошибке чтения будет возвращён `nil`;
        - `inln(count: num): nil|str` — считывает `count` символов из stdin и возвращает их в виде строки; если количество будет отрицательным, то будут считаны все символы до символа `'\n'` или, если он отсутствует, все доступные символы; при ошибке чтения будет возвращён `nil`;
        - `out(text: str): nil` — выводит строку `text` в stdout;
        - `outln(text: str): nil` — выводит строку `text` в stdout и переводит строку (добавляет символ `'\n'`);
        - `err(text: str): nil` — выводит строку `text` в stderr;
        - `errln(text: str): nil` — выводит строку `text` в stderr и переводит строку (добавляет символ `'\n'`);
- add the new example with pi evaluation.

## [v1.0](https://github.com/thewizardplusplus/tick-tock/tree/v1.0) (2019-01-14)
