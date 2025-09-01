; fibonacci
        PUSH 0xB00B5
        STORE 0x000FFFFF
        FETCH 0x000FFFFF
        FETCH 0x7FFFFF
        SPAWN SPAWN
:SPAWN  PUSH 1
        DUP

:LOOP   DUP
        DUP 2
        ADD
        DUP
        PUSH -0x7FFFFF
        CMP
        JMPLT LOOP
        KILL