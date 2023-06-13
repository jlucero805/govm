class Stack:

    def __init__(self):
        self.stack = []

    def push(self, value):
        self.stack.append(value)

    def pop(self):
        return self.stack.pop()

    def __repr__(self) -> str:
        return str(self.stack)

HALT = 0
PUSH = 1
POP = 2
ADD = 3
DUP = 4
LOAD = 5
STORE = 6
JMP = 7
JIF = 8

class Frame:
    def __init__(self):
        self.vars = {}

    def get(self, num):
        return self.vars[num]

    def set(self, var, val):
        self.vars[var] = val

    def __repr__(self):
        return str(self.vars)

class VM:
    def __init__(self, prog):
        self.prog = prog
        self.stack = Stack()
        self.halted = False
        self.ip = 0
        self.frame = Frame()

    def run(self):
        while self.halted == False:
            self.go()
        print(self.stack)
        print(self.frame)

    def go(self):
        self.exec(self.next())

    def exec(self, instr):
        if instr == HALT:
            self.halted = True
        elif instr == PUSH:
            self.stack.push(self.next())
        elif instr == ADD:
            r = self.stack.pop()
            l = self.stack.pop()
            self.stack.push(l + r)
        elif instr == DUP:
            v = self.stack.pop()
            self.stack.push(v)
            self.stack.push(v)
        elif instr == LOAD:
            nex = self.next()
            num = self.frame.get(nex)
            self.stack.push(num)
        elif instr == STORE:
            self.frame.set(self.next(), self.stack.pop())
        elif instr == JMP:
            self.ip = self.next()
        elif instr == JIF:
            addr = self.next()
            if self.stack.pop() > 0:
                self.ip = addr
            
    
    def next(self):
        instr = self.prog[self.ip]
        self.ip += 1
        return instr

VM([
PUSH, 1,
JIF, 5,
HALT,
PUSH, 69,
HALT
]).run()

