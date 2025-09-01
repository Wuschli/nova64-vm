# NOVA64 CPU SPECIFICATION  
**Confidential – Circa 1991 – TechCorp Systems Division**

---

## Overview

The Nova64 is a 32-bit stack-oriented processor designed for modular computing platforms. It supports multi-tasking, dynamic I/O, and flexible memory addressing using a reserved stack alias `TOP` (0xFFFFFFFF). Its architecture is ideal for embedded cyberpunk computing devices, industrial controllers, and experimental programmable cartridges.

- **Word size:** 32-bit  
- **Stack depth:** 1024 entries recommended  
- **Registers:** IP (Instruction Pointer), SP (Stack Pointer), X (Compare Register), FLAGS (status flags)  
- **Ports:** Generic LINK, Displays, Input devices, Cartridges, Audio/Video, Status/Clock  

---
## Instruction Set
because all `<val>`,`<label>`,`<offset>`,`<addr>` and `<port>` values are encoded into the instruction word, all those can only use 24bit values (-1, because of `TOP` keyword)
The 4 bytes of an instruction are stored in memory in the following way: `<OPCODE>` `<VALUE>`  `<VALUE>` `<VALUE>` where the 3 bytes of `<VALUE>` are in little endian order

| Mnemonic        | Opcode (hex) | Description                                                                                              |
| --------------- | ------------ | -------------------------------------------------------------------------------------------------------- |
| `NOOP`          | 0x00         | Does nothing                                                                                             |
| `PUSH <val>`    | 0x01         | Push 24-bit value onto stack                                                                             |
| `DROP`          | 0x02         | Remove TOS                                                                                               |
| `DUP <offset>`  | 0x03         | Duplicate \[offset] of stack (default offset = 0) offset is encoded in upper 24bit or from stack (`TOP`) |
| `SWAP`          | 0x04         | Swap top two stack values                                                                                |
| `ADD`           | 0x10         | Pop A, Pop B, push A+B                                                                                   |
| `SUB`           | 0x11         | Pop A, Pop B, push A-B                                                                                   |
| `MUL`           | 0x12         | Pop A, Pop B, push A*B                                                                                   |
| `DIV`           | 0x13         | Pop A, Pop B, push A/B                                                                                   |
| `MOD`           | 0x14         | Pop A, Pop B, push A%B                                                                                   |
| `AND`           | 0x20         | Pop A, Pop B, push A&B                                                                                   |
| `OR`            | 0x21         | Pop A, Pop B, push A\|B                                                                                  |
| `XOR`           | 0x22         | Pop A, Pop B, push A^B                                                                                   |
| `NOT`           | 0x23         | Logical NOT of top value                                                                                 |
| `CMP`           | 0x30         | Compare top two values, set X: -1 <, 0 ==, 1 >                                                           |
| `JMP <label>`   | 0x40         | Jump to label or address from stack (`TOP`)                                                              |
| `CALL <label>`  | 0x41         | Push the current IP to stack and jump to label or address from stack (`TOP`)                             |
| `RET`           | 0x42         | Pop the return address from stack and set IP to it                                                       |
| `JMPZ <label>`  | 0x43         | Jump if X == 0 (TOP supported)                                                                           |
| `JMPLT <label>` | 0x44         | Jump if X < 0 (TOP supported)                                                                            |
| `JMPGT <label>` | 0x45         | Jump if X > 0 (TOP supported)                                                                            |
| `FETCH <addr>`  | 0x50         | Push value from memory (TOP = stack)                                                                     |
| `STORE <addr>`  | 0x51         | Store top of stack to memory (TOP = stack)                                                               |
| `IN <port>`     | 0x60         | Read from port (TOP = port from stack)                                                                   |
| `OUT <port>`    | 0x61         | Write to port (TOP = port from stack)                                                                    |
| `SPAWN <label>` | 0x70         | Spawn new task at label or stack address (`TOP`)                                                         |
| `YIELD`         | 0x71         | Yield CPU cycle to next task                                                                             |
| `WAIT <port>`   | 0x72         | Wait for interrupt on port (`TOP` = port from stack)                                                     |
| `KILL`          | 0xFF         | Kill the active task. CPU is halted when all tasks are killed                                            |

---

## Ports

| Port  | Function                                                |
| ----- | ------------------------------------------------------- |
| 1–4   | LINK0–LINK3, generic inter-CPU communications           |
| 5–8   | Display0–Display3, framebuffer and commands             |
| 9     | Keyboard                                                |
| 10    | Mouse / Pointer                                         |
| 11–14 | Cartridge0–Cartridge3                                   |
| 15    | Status: RAM size, task count, CPU version, display info |
| 16    | Clock / Timer                                           |
| 17    | Audio IN/OUT                                            |
| 18    | Video IN/OUT                                            |
| 19-32 | Reserved for future expansion                           |

---

## Notes on TOP

`TOP` (0xFFFFFF) is a **reserved stack alias**. When used in any instruction expecting an address or port, the CPU will pop the top value from the stack to use as the operand:

- `JMP TOP` → Jump to address on stack  
- `CALL TOP` → Call subroutine at address on stack  
- `FETCH TOP` / `STORE TOP` → Memory indirect addressing  
- `IN TOP` / `OUT TOP` → Dynamic I/O  
- `WAIT TOP` → Wait on port specified on stack  
- `SPAWN TOP` → Spawn task at stack address  

This provides a fully flexible, stack-driven mechanism for dynamic code execution and I/O.

---

**End of Document**  
© 1991 TechCorp Systems Division – All rights reserved.
