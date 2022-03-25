# Hel programming language

Dynamically-typed programming language with the scope of providing functionallity
of processing data by using the internal goroutines.

## Principles:

- What is written, will be changed
- “Simple” is a process, not a goal
- Production ready from the start

## Top goals/questions
- How to split a structure in batches so it can have operations done concurrently.
- How do you achieve map-reduce on large datasets with lazy loading using concurrency

Current VM instructions for 1+2

```
== test ==
0000    1 CONSTANT         0 '{%!g(value.VALUE_TYPE=3) {%!g(bool=false) %!g(int=0) 1}}'
0002    1 CONSTANT         1 '{%!g(value.VALUE_TYPE=3) {%!g(bool=false) %!g(int=0) 2}}'
0004    1 INSTRUC_ADDITION
0005    1 INSTRUC_RETURN
RETURN
STACK POP:, value.Value{VT:3, _V:value.V{b:false, i:0, f64:3}}%
```

VARIABLES

- Variable declaration

```
decl a = 10
decl b[]
decl c = 10
```

- Variable type declaration

https://go101.org/article/type-system-overview.html

```
decl int8 a = 10
decl string b[]
decl int8 c = 10
```

# .............................................

CONTROL FLOW EXPRESSIONS

---


FUNCTIONS

```
decl fn testVariable() {
    return a < c ? true : false
} (bool)

decl fn testVariable() =  {
    a < c ? true : false
} (bool)

decl cncr fn testVariable() =  {
    a < c ? true : false
} (bool)

b = testVariable()

decl func populateArray(startRange: int, endRange: int) =  {
    ...
} (array)

b = populateArray(1, 100)
```
