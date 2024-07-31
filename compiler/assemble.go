package compile

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"strata/parser"
	"strconv"
	"strings"
)

const (
	PUSH int64 = iota
	LOAD
	STORE
	ENVNEW
	ENVEND
	FUNDEF
	FUNEND
	CALL
	RET
	ADD
	MUL
	SUB
	DIV
	GT
	LT
	GTEQ
	LTEQ
	JMPR
	JMPRIF
	EQ
)

func stringToNestedList(input string) [][]string {
	split1 := strings.Split(input, "\n")
	lists := [][]string{}

	// furthur split
	for _, line := range split1 {
		lists = append(lists, strings.Split(line, " "))
	}
	return lists
}

func instrCode(input string) int64 {
	switch input {
	case "push":
		return PUSH
	case "load":
		return LOAD
	case "store":
		return STORE
	case "envnew":
		return ENVNEW
	case "envend":
		return ENVEND
	case "fundef":
		return FUNDEF
	case "funend":
		return FUNEND
	case "call":
		return CALL
	case "ret":
		return RET
	case "add":
		return ADD
	case "sub":
		return SUB
	case "mul":
		return MUL
	case "div":
		return DIV
	case "lt":
		return LT
	case "eq":
		return EQ
	case "jmpr":
		return JMPR
	case "jmprif":
		return JMPRIF
	default:
		panic("instr doesn't exist")
	}
}

func symbolsToInts(inputs [][]string) [][]int64 {
	var counter int64 = 0
	sym := make(map[string]int64)
	lists := [][]int64{}

	for _, input := range inputs {
		if len(input) == 2 {
			if input[0] == "push" || input[0] == "jmpr" || input[0] == "jmprif" {
				integer, _ := strconv.Atoi(input[1])
				lists = append(lists, []int64{instrCode(input[0]), int64(integer)})
			} else if val, ok := sym[input[1]]; ok {
				lists = append(lists, []int64{instrCode(input[0]), val})
			} else {
				counter += 1
				sym[input[1]] = counter
				lists = append(lists, []int64{instrCode(input[0]), sym[input[1]]})
			}
		} else {
			lists = append(lists, []int64{instrCode(input[0])})
		}
	}
	return lists
}

func flat(input [][]int64) []int64 {
	bytes := []int64{}

	for _, input := range input {
		for _, num := range input {
			bytes = append(bytes, num)
		}
	}

	return bytes
}

func Assemble(input string) []int64 {
	ast := parser.Parse(input)
	instrs := Compile(ast)
	nestedList := stringToNestedList(instrs)
	nestedInts := symbolsToInts(nestedList)
	bytes := flat(nestedInts)
	return bytes
}

func AssembleAll(input string) []int64 {
	asts := parser.ParseAll(input)
	code := []string{}
	for _, ast := range asts {
		code = append(code, Compile(ast))
	}
	finalCode := strings.Join(code, "\n")
	nestedList := stringToNestedList(finalCode)
	nestedInts := symbolsToInts(nestedList)
	bytes := flat(nestedInts)
	return bytes
}

func writeInt64ToFile(fileName string, nums []int64) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	for _, num := range nums {
		err = binary.Write(file, binary.LittleEndian, num)
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}
}
func writeToFile(filename string, numbers []int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, num := range numbers {
		err := binary.Write(file, binary.LittleEndian, num)
		if err != nil {
			return err
		}
	}

	return nil
}

func Write(input string, fileName string) {
	bytecode := AssembleAll(input)
	writeToFile(fileName, bytecode)
}

func readFromFile(filename string) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []int64
	for {
		var num int64
		err := binary.Read(file, binary.LittleEndian, &num)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}

func GetBytesFromFile(input string) []int64 {
	nums, _ := readFromFile(input)
	return nums
}
