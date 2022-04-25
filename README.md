# HProg (Hel programming language)

Dynamically-typed programming language.

## Sample VM Instructions

```
hprog> print(1+1)

== test ==
0000    1 CONSTANT         0 '1 (VT_INT)'
0002    1 CONSTANT         1 '1 (VT_INT)'
0004    1 INSTRUC_ADDITION
0005    1 INSTRUC_PRINT
0006    1 INSTRUC_RETURN

2 (VT_INT)

```

# Samples

## Variables

- Variable declaration

```
decl a = 10
decl b[]
decl c = 10
```

## Functions

```
fn (bool) testVariable() = return a < c ? true : false

fn (array) populateArray(startRange, endRange) = {
    ...
}
```
