package lcs_test

import (
	"encoding/hex"
	"fmt"

	"github.com/the729/lcs"
)

func ExampleMarshal_struct() {
	type MyStruct struct {
		Boolean    bool
		Bytes      []byte
		Label      string
		unexported uint32
	}
	type Wrapper struct {
		Inner *MyStruct `lcs:"optional"`
		Name  string
	}

	bytes, err := lcs.Marshal(&Wrapper{
		Name: "test",
		Inner: &MyStruct{
			Bytes: []byte{0x01, 0x02, 0x03, 0x04},
			Label: "hello",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x", bytes)
	// Output: 010004000000010203040500000068656c6c6f0400000074657374
}

func ExampleUnmarshal_struct() {
	type MyStruct struct {
		Boolean    bool
		Bytes      []byte
		Label      string
		unexported uint32
	}
	type Wrapper struct {
		Inner *MyStruct `lcs:"optional"`
		Name  string
	}

	bytes, _ := hex.DecodeString("010004000000010203040500000068656c6c6f0400000074657374")
	out := &Wrapper{}
	err := lcs.Unmarshal(bytes, out)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s, Label: %s\n", out.Name, out.Inner.Label)
	// Output: Name: test, Label: hello
}

type TransactionArgument interface {
	isTransactionArg()
}
type TxnArgU64 uint64
type TxnArgAddress [32]byte
type TxnArgString string

func (TxnArgU64) isTransactionArg()     {}
func (TxnArgAddress) isTransactionArg() {}
func (TxnArgString) isTransactionArg()  {}

type Program struct {
	Code    []byte
	Args    []TransactionArgument `lcs:"enum=txn_arg"`
	Modules [][]byte
}

func (*Program) EnumTypes() []lcs.EnumVariant {
	return []lcs.EnumVariant{
		{"txn_arg", 0, TxnArgU64(0)},
		{"txn_arg", 1, TxnArgAddress([32]byte{})},
		{"txn_arg", 2, TxnArgString("")},
	}
}

func ExampleMarshal_libra_program() {
	prog := &Program{
		Code: []byte("move"),
		Args: []TransactionArgument{
			TxnArgString("CAFE D00D"),
			TxnArgString("cafe d00d"),
		},
		Modules: [][]byte{{0xca}, {0xfe, 0xd0}, {0x0d}},
	}

	bytes, err := lcs.Marshal(prog)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%X\n", bytes)
	// Output:
	// 040000006D6F766502000000020000000900000043414645204430304402000000090000006361666520643030640300000001000000CA02000000FED0010000000D
}

func ExampleUnmarshal_libra_program() {
	bytes, _ := hex.DecodeString("040000006D6F766502000000020000000900000043414645204430304402000000090000006361666520643030640300000001000000CA02000000FED0010000000D")
	out := &Program{}
	err := lcs.Unmarshal(bytes, out)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", out)
	// Output:
	// &{Code:[109 111 118 101] Args:[CAFE D00D cafe d00d] Modules:[[202] [254 208] [13]]}
}
