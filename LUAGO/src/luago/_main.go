package main

import (
	"fmt"
	"io/ioutil"
	"luago/binchunk"
	. "luago/vm"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		proto := binchunk.Undump(data)
		list(proto)
	}
}

// 打印函数原型信息
func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p) // 递归打印子函数
	}
}

//
func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d,%d> (%d instructions\n)",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))
	fmt.Printf("%d locals, %d constants, %d function\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

// 打印指令的序号，行号，十六进制表示
func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != OpArgN {
			if b > 0xFF {
				// 如果操作数B或C的最高位是1就认为它表示常量表索引，按负数输出
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
		
	case IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == OpArgU {
			fmt.Printf(" %d", bx)
		}
		
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
		
	case IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
		
	}
}

// 打印常量表，局部变量表，Upvalue表
func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants(%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

// 把常量表里的常量转化成字符串
func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "? "
	}
}

// 根据Upvalue索引从调试信息里找出Upvalue名字
func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}
