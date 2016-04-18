## Example Trees

#
### INVALID
#
+10-57-+30


```
+
10  -
    57  -
        +
        30

10 + (57 - (-(+(30)))
```

#
### VALID
#
*21-30+10/2*9,12

```
*
21  -
    30  +
        10  /
            2   *
                9   12

21 * (30 - (10 + (2 / (9 * 12))))
```

