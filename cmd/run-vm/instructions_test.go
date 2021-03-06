/*
 * Copyright 2019, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/offchainlabs/arb-avm/code"
	"github.com/offchainlabs/arb-avm/vm"
	"github.com/offchainlabs/arb-util/protocol"
	"github.com/offchainlabs/arb-util/value"
)

// This is to test that a machine can be built and run
// It creates a macine with 4 steps and runs it
// There is no automated test check so pass/fail must be verified visually
func TestMachineAdd(t *testing.T) {
	insns := make([]value.Operation, 4)
	fmt.Println("Setting up insns")
	i := 0
	insns[i] = value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(2)}
	i++
	insns[1] = value.ImmediateOperation{Op: code.ADD, Val: value.NewInt64Value(4)}
	i++
	insns[i] = value.BasicOperation{Op: code.LOG}
	i++
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	assertion := m.ExecuteAssertion(80000, protocol.NewTimeBounds(0, 100000))
	fmt.Println(assertion.NumSteps)
}

// base operation tests for one, two, or three operands
// Push the required number of operands (from passed in values)
// Run the given instruction
// Push the expected result to the stack of a second machine
// Compare the two machines
func unaryIntOpTest(x, expected *big.Int, oper value.Opcode) (bool, string) {
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewIntValue(x))

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: value.Opcode(oper)}); err != nil {
		tmp := "RunInstruction error:"
		tmp += err.Error()
		return false, tmp
	}
	knownMachine.Stack().Push(value.NewIntValue(expected))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		tmp := "machines not equal: "
		tmp += err
		return false, tmp
	}

	return true, ""
}

func binaryIntOpTest(x, y, expected *big.Int, oper value.Opcode) (bool, string) {
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewIntValue(y))
	m.Stack().Push(value.NewIntValue(x))

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: value.Opcode(oper)}); err != nil {
		tmp := "RunInstruction error:"
		tmp += err.Error()
		return false, tmp
	}
	knownMachine.Stack().Push(value.NewIntValue(expected))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		tmp := "machines not equal: "
		tmp += err
		return false, tmp
	}

	return true, ""
}

func binaryValueOpTest(x, y value.Value, expected *big.Int, oper value.Opcode) (bool, string) {
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(y)
	m.Stack().Push(x)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: value.Opcode(oper)}); err != nil {
		tmp := "RunInstruction error:"
		tmp += err.Error()
		return false, tmp
	}
	knownMachine.Stack().Push(value.NewIntValue(expected))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		tmp := "machines not equal: "
		tmp += err
		return false, tmp
	}

	return true, ""
}

func tertiaryIntOpTest(x, y, z, expected *big.Int, oper value.Opcode) (bool, string) {
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewIntValue(z))
	m.Stack().Push(value.NewIntValue(y))
	m.Stack().Push(value.NewIntValue(x))

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: value.Opcode(oper)}); err != nil {
		tmp := "RunInstruction error:"
		tmp += err.Error()
		return false, tmp
	}
	knownMachine.Stack().Push(value.NewIntValue(expected))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		tmp := "machines not equal: "
		tmp += err
		return false, tmp
	}

	return true, ""
}

// This test is to test an operation missing the second value
func TestAddMissingValue(t *testing.T) {
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.ADD}); err == nil {
		t.Error("tried to pop empty stack expected")
	}
	knownMachine.Stack().Push(value.NewInt64Value(2))
	if ok, _ := vm.Equal(knownMachine, m); ok {
		tmp := "machines equal expected different"
		t.Error(tmp)
	}
}

//*************************************
// These are the tests for each opcode
//*************************************
func TestAdd(t *testing.T) {
	// test 3+4=7
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(4), big.NewInt(7), code.ADD)
	if !res {
		t.Error(err)
	}
	// test 0+0=0
	res, err = binaryIntOpTest(big.NewInt(0), big.NewInt(0), big.NewInt(0), code.ADD)
	if !res {
		t.Error(err)
	}
	// test (2**256-1)+4=3
	res, err = binaryIntOpTest(math.U256(big.NewInt(-1)), big.NewInt(4), big.NewInt(3), code.ADD)
	if !res {
		t.Error(err)
	}
	// test (2**256-2)+1=0xffff...ff
	res, err = binaryIntOpTest(math.U256(big.NewInt(-2)), big.NewInt(1), math.U256(big.NewInt(-1)), code.ADD)
	if !res {
		t.Error(err)
	}
	// test 3 + tuple = error
	tup := value.NewEmptyTuple()
	res, err = binaryValueOpTest(value.NewInt64Value(3), tup, big.NewInt(7), code.ADD)
	if res {
		t.Error("expected error")
	} else {
		fmt.Println(err)
	}
}

func TestMul(t *testing.T) {
	// test 3*4=12
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(4), big.NewInt(12), code.MUL)
	if !res {
		t.Error(err)
	}
	// test 3*0=0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(0), big.NewInt(0), code.MUL)
	if !res {
		t.Error(err)
	}
	// test (2**256-1)*1=(2**256-1)
	res, err = binaryIntOpTest(math.U256(big.NewInt(-1)), big.NewInt(1), math.U256(big.NewInt(-1)), code.MUL)
	if !res {
		t.Error(err)
	}
	// test (2**256-1)*2=(2**256-2)
	res, err = binaryIntOpTest(math.U256(big.NewInt(-1)), big.NewInt(2), math.U256(big.NewInt(-2)), code.MUL)
	if !res {
		t.Error(err)
	}
}

func TestSub(t *testing.T) {
	// test 4-3=1
	res, err := binaryIntOpTest(big.NewInt(4), big.NewInt(3), big.NewInt(1), code.SUB)
	if !res {
		t.Error(err)
	}
	// test 0-1=0xffff...ff
	res, err = binaryIntOpTest(big.NewInt(0), big.NewInt(1), math.U256(big.NewInt(-1)), code.SUB)
	if !res {
		t.Error(err)
	}
}

func TestDiv(t *testing.T) {
	// test 6/2=3
	res, err := binaryIntOpTest(big.NewInt(6), big.NewInt(2), big.NewInt(3), code.DIV)
	if !res {
		t.Error(err)
	}
	// test -6/2!=-3 (should be unsigned division)
	res, err = binaryIntOpTest(math.U256(big.NewInt(-6)), big.NewInt(2), math.U256(big.NewInt(-3)), code.DIV)
	if res {
		t.Error("should not be -3")
	}
	// test 6/0=0
	res, err = binaryIntOpTest(big.NewInt(6), big.NewInt(0), big.NewInt(0), code.DIV)
	if res {
		t.Error("Divide by 0 expected")
	}
}

func TestSDiv(t *testing.T) {
	// test -6/3=-2
	res, err := binaryIntOpTest(math.U256(big.NewInt(-6)), math.U256(big.NewInt(3)), math.U256(big.NewInt(-2)), code.SDIV)
	if !res {
		t.Error(err)
	}
	// test 6/-3=-2
	res, err = binaryIntOpTest(math.U256(big.NewInt(6)), math.U256(big.NewInt(-3)), math.U256(big.NewInt(-2)), code.SDIV)
	if !res {
		t.Error(err)
	}
	// test -6/-3=2
	res, err = binaryIntOpTest(math.U256(big.NewInt(-6)), math.U256(big.NewInt(-3)), math.U256(big.NewInt(2)), code.SDIV)
	if !res {
		t.Error(err)
	}
	// test 6/3=2
	res, err = binaryIntOpTest(math.U256(big.NewInt(6)), math.U256(big.NewInt(3)), math.U256(big.NewInt(2)), code.SDIV)
	if !res {
		t.Error(err)
	}
	// test 6/0=0
	res, err = binaryIntOpTest(big.NewInt(6), big.NewInt(0), big.NewInt(0), code.SDIV)
	if res {
		t.Error("Divide by 0 expected")
	}
}

func TestMod(t *testing.T) {
	// test 8%3=2
	res, err := binaryIntOpTest(big.NewInt(8), big.NewInt(3), big.NewInt(2), code.MOD)
	if !res {
		t.Error(err)
	}
	// test 0%3=0
	res, err = binaryIntOpTest(big.NewInt(0), big.NewInt(3), big.NewInt(0), code.MOD)
	if !res {
		t.Error(err)
	}
	// test 3%8=0 - divide by zero error
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(0), big.NewInt(0), code.MOD)
	if res {
		t.Error("Divide by 0 expected")
	}
}

func TestSMod(t *testing.T) {
	// test -8%3=-2
	res, err := binaryIntOpTest(math.U256(big.NewInt(-8)), math.U256(big.NewInt(3)), math.U256(big.NewInt(-2)), code.SMOD)
	if !res {
		t.Error(err)
	}
	// test -8%-3=-2 (per spec modulo by a negative number takes the sign of the dividend
	res, err = binaryIntOpTest(math.U256(big.NewInt(-8)), math.U256(big.NewInt(-3)), math.U256(big.NewInt(-2)), code.SMOD)
	if !res {
		t.Error(err)
	}
	// test 8%3=2
	res, err = binaryIntOpTest(math.U256(big.NewInt(8)), math.U256(big.NewInt(3)), math.U256(big.NewInt(2)), code.SMOD)
	if !res {
		t.Error(err)
	}
	// test -8%0=0
	res, err = binaryIntOpTest(math.U256(big.NewInt(-8)), math.U256(big.NewInt(0)), math.U256(big.NewInt(0)), code.SMOD)
	if res {
		t.Error("Divide by 0 expected")
	}
}

func TestAddMod(t *testing.T) {
	// test (8+5)%3=1
	res, err := tertiaryIntOpTest(big.NewInt(8), big.NewInt(5), big.NewInt(3), big.NewInt(1), code.ADDMOD)
	if !res {
		t.Error(err)
	}
	// test ((2**256-1)+1)%7=2 - shows that internal addition does not 256 bit truncate
	res, err = tertiaryIntOpTest(math.U256(big.NewInt(-1)), big.NewInt(1), big.NewInt(7), big.NewInt(2), code.ADDMOD)
	if !res {
		t.Error(err)
	}
	// test (0+0)%7=0
	res, err = tertiaryIntOpTest(math.U256(big.NewInt(0)), big.NewInt(0), big.NewInt(7), big.NewInt(0), code.ADDMOD)
	if !res {
		t.Error(err)
	}
}

func TestMulMod(t *testing.T) {
	// test (8*2)%3=1
	res, err := tertiaryIntOpTest(big.NewInt(8), big.NewInt(2), big.NewInt(3), big.NewInt(1), code.MULMOD)
	if !res {
		t.Error(err)
	}
	// test ((2**256-1)*2)%7=2 - shows that internal addition does not 256 bit truncate
	res, err = tertiaryIntOpTest(math.U256(big.NewInt(-1)), big.NewInt(2), big.NewInt(7), big.NewInt(2), code.MULMOD)
	if !res {
		t.Error(err)
	}
}

func TestExp(t *testing.T) {
	// test 3^2=9
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(2), big.NewInt(9), code.EXP)
	if !res {
		t.Error(err)
	}
	// test 2 exp 256 = 0 - test wrap
	res, err = binaryIntOpTest(big.NewInt(2), big.NewInt(256), big.NewInt(0), code.EXP)
	if !res {
		t.Error(err)
	}
}

func TestSignextend(t *testing.T) {
	// test
	res, err := binaryIntOpTest(big.NewInt(-1), big.NewInt(0), math.U256(big.NewInt(-1)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(1), big.NewInt(0), math.U256(big.NewInt(1)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(127), big.NewInt(0), math.U256(big.NewInt(127)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(128), big.NewInt(0), math.U256(big.NewInt(-128)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(254), big.NewInt(0), math.U256(big.NewInt(-2)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(257), big.NewInt(0), math.U256(big.NewInt(1)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(65534), big.NewInt(1), math.U256(big.NewInt(-2)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(65537), big.NewInt(1), math.U256(big.NewInt(1)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
	// test
	res, err = binaryIntOpTest(big.NewInt(65537), big.NewInt(2), math.U256(big.NewInt(65537)), code.SIGNEXTEND)
	if !res {
		t.Error(err)
	}
}

func TestLt(t *testing.T) {
	// test 3<9 res 1
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(1), code.LT)
	if !res {
		t.Error(err)
	}
	// test 9<3 res 0
	res, err = binaryIntOpTest(big.NewInt(9), big.NewInt(3), big.NewInt(0), code.LT)
	if !res {
		t.Error(err)
	}
	// test 3<3 res 0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.LT)
	if !res {
		t.Error(err)
	}
	// test 0xfffffffffffffffffffffffffffffffc((2**256)-4)<9 res 0
	res, err = binaryIntOpTest(math.U256(big.NewInt(-4)), big.NewInt(9), big.NewInt(0), code.LT)
	if !res {
		t.Error(err)
	}
	// test 9< tuple res 0
	res, err = binaryValueOpTest(value.NewInt64Value(9), value.NewEmptyTuple(), big.NewInt(0), code.LT)
	if res {
		t.Error("expected error")
	}
}

func TestGt(t *testing.T) {
	// test 3>9 res 0
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(0), code.GT)
	if !res {
		t.Error(err)
	}
	// test 9>3 res 1
	res, err = binaryIntOpTest(big.NewInt(9), big.NewInt(3), big.NewInt(1), code.GT)
	if !res {
		t.Error(err)
	}
	// test 3>3 res 0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.GT)
	if !res {
		t.Error(err)
	}
	// test 0xfffffffffffffffffffffffffffffffc(-4)>9 res 1
	res, err = binaryIntOpTest(math.U256(big.NewInt(-4)), big.NewInt(9), big.NewInt(1), code.GT)
	if !res {
		t.Error(err)
	}
}

func TestSlt(t *testing.T) {
	// test 3 < 9 = 1
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(1), code.SLT)
	if !res {
		t.Error(err)
	}
	// test 9 < 3 = 0
	res, err = binaryIntOpTest(big.NewInt(9), big.NewInt(3), big.NewInt(0), code.SLT)
	if !res {
		t.Error(err)
	}
	// test 3 < 3 = 0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.SLT)
	if !res {
		t.Error(err)
	}
	// test -3 < 3 = 1
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), big.NewInt(3), big.NewInt(1), code.SLT)
	if !res {
		t.Error(err)
	}
	// test -3 < -4 = 0
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), math.U256(big.NewInt(-4)), big.NewInt(0), code.SLT)
	if !res {
		t.Error(err)
	}
	// test -3 < -2 = 1
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), math.U256(big.NewInt(-2)), big.NewInt(1), code.SLT)
	if !res {
		t.Error(err)
	}
}

func TestSgt(t *testing.T) {
	// test 3 > 9 = 0
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(0), code.SGT)
	if !res {
		t.Error(err)
	}
	// test 9 > 3 = 1
	res, err = binaryIntOpTest(big.NewInt(9), big.NewInt(3), big.NewInt(1), code.SGT)
	if !res {
		t.Error(err)
	}
	// test 3 > 3 = 0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.SGT)
	if !res {
		t.Error(err)
	}
	// test -3 > 3 = 0
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), big.NewInt(3), big.NewInt(0), code.SGT)
	if !res {
		t.Error(err)
	}
	// test -3 > -4 = 1
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), math.U256(big.NewInt(-4)), big.NewInt(1), code.SGT)
	if !res {
		t.Error(err)
	}
	// test -3 > -2 = 0
	res, err = binaryIntOpTest(math.U256(big.NewInt(-3)), math.U256(big.NewInt(-2)), big.NewInt(0), code.SGT)
	if !res {
		t.Error(err)
	}
}

func TestEq(t *testing.T) {
	// test 3==9 = 0
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(0), code.EQ)
	if !res {
		t.Error(err)
	}
	// test 3==3 = 1
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(1), code.EQ)
	if !res {
		t.Error(err)
	}

	var vals [8]value.Value
	vals[0] = value.NewInt64Value(1)
	vals[1] = value.NewInt64Value(2)
	vals[2] = value.NewInt64Value(3)
	vals[3] = value.NewInt64Value(4)
	tup, _ := value.NewTupleOfSizeWithContents(vals, 4)

	vals[0] = value.NewInt64Value(1)
	vals[1] = value.NewInt64Value(2)
	vals[2] = value.NewInt64Value(3)
	vals[3] = value.NewInt64Value(4)
	tup2, _ := value.NewTupleOfSizeWithContents(vals, 4)

	vals[0] = value.NewInt64Value(1)
	vals[1] = value.NewInt64Value(2)
	vals[2] = value.NewInt64Value(7) // one value is different
	vals[3] = value.NewInt64Value(4)
	tup3, _ := value.NewTupleOfSizeWithContents(vals, 4)
	// test matching tuples
	res, err = binaryValueOpTest(tup, tup2, big.NewInt(1), code.EQ)
	if !res {
		t.Error(err)
	}
	// test different tuples
	res, err = binaryValueOpTest(tup, tup3, big.NewInt(0), code.EQ)
	if !res {
		t.Error(err)
	}
	// test different types
	res, err = binaryValueOpTest(tup, value.NewInt64Value(1), big.NewInt(0), code.EQ)
	if !res {
		t.Error("expected fail")
	}
}

func TestIszero(t *testing.T) {
	// test 0 isZero = 1
	res, err := unaryIntOpTest(big.NewInt(0), big.NewInt(1), code.ISZERO)
	if !res {
		t.Error(err)
	}
	// test 2 isZero = 0
	res, err = unaryIntOpTest(big.NewInt(3), big.NewInt(0), code.ISZERO)
	if !res {
		t.Error(err)
	}
}

func TestAnd(t *testing.T) {
	// test 0x03 and 0x09 = 0x01
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(1), code.AND)
	if !res {
		t.Error(err)
	}
	// test 0x03 and 0x03 = 0x03
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(3), code.AND)
	if !res {
		t.Error(err)
	}
}

func TestOr(t *testing.T) {
	// test 0x03 or 0x09 = 0x0b
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(11), code.OR)
	if !res {
		t.Error(err)
	}
	// test 0x03 or 0x03 = 0x03
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(3), code.OR)
	if !res {
		t.Error(err)
	}
}

func TestXor(t *testing.T) {
	// test 0x03 xor 0x09 = 0x0a
	res, err := binaryIntOpTest(big.NewInt(3), big.NewInt(9), big.NewInt(10), code.XOR)
	if !res {
		t.Error(err)
	}
	// test 0x03 xor 0x03 = 0x00
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.XOR)
	if !res {
		t.Error(err)
	}
}

func TestNot(t *testing.T) {
	// test !0x00 = 0xffffffffffffffffffffffffffffffff(-1)
	res, err := unaryIntOpTest(big.NewInt(0), math.U256(big.NewInt(-1)), code.NOT)
	if !res {
		t.Error(err)
	}
	// test !0x03 = 0xfffffffffffffffffffffffffffffffc(-4)
	res, err = unaryIntOpTest(big.NewInt(3), math.U256(big.NewInt(-4)), code.NOT)
	if !res {
		t.Error(err)
	}
	// test !0xfffffffffffffffffffffffffffffffc(-4) = 0x03(3)
	res, err = unaryIntOpTest(math.U256(big.NewInt(-4)), math.U256(big.NewInt(3)), code.NOT)
	if !res {
		t.Error(err)
	}
}

func TestByte(t *testing.T) {
	// test 31st byte of 16 = 16
	res, err := binaryIntOpTest(big.NewInt(16), big.NewInt(31), big.NewInt(16), code.BYTE)
	if !res {
		t.Error(err)
	}
	// test 3rd byte of 3 = 0
	res, err = binaryIntOpTest(big.NewInt(3), big.NewInt(3), big.NewInt(0), code.BYTE)
	if !res {
		t.Error(err)
	}
}

func TestSha3(t *testing.T) {
	// test
	hash, _ := new(big.Int).SetString("80084422859880547211683076133703299733277748156566366325829078699459944778998", 10)
	res, err := unaryIntOpTest(big.NewInt(1), hash, code.SHA3)
	if !res {
		t.Error(err)
	}
}

func TestPop(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))
	a := m.Stack().Count()
	if a != 1 {
		tmp := "PUSH failed stack size = "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.POP}); err != nil {
		tmp := "POP failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	a = m.Stack().Count()
	if a != 0 {
		tmp := "POP stack size check failed"
		t.Error(tmp)
	}
}

func TestSpush(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.SPUSH}); err != nil {
		tmp := "SPUSH failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	a := m.Stack().Count()
	if a != 1 {
		tmp := "SPUSH stack size check failed"
		t.Error(tmp)
	}
}

func TestRpush(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.RPUSH}); err != nil {
		tmp := "RPUSH failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	a := m.Stack().Count()
	if a != 1 {
		tmp := "RPUSH stack size check failed"
		t.Error(tmp)
	}
}

func TestRset(t *testing.T) {
	//test:
	// 1. push value
	// 2. run RSET
	// 3. verify machines different
	// 4. push value to known
	// 5. run RSET on known
	// 6. verify machines match
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(5))

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.RSET}); err != nil {
		tmp := "RSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	a := m.Stack().Count()
	if a != 0 {
		tmp := "RSET stack size check failed"
		t.Error(tmp)
	}
	if ok, _ := vm.Equal(knownMachine, m); ok {
		t.Error("machines equal expected different")
	}

	knownMachine.Stack().Push(value.NewInt64Value(5))
	if _, err := vm.RunInstruction(knownMachine, value.BasicOperation{Op: code.RSET}); err != nil {
		tmp := "RSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestInbox(t *testing.T) {
	//test:
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knowninbox := protocol.NewEmptyInbox()

	var tok protocol.TokenType
	tok[0] = 15
	tok[20] = 1

	dest := [32]byte{}
	copy(dest[:], math.U256(big.NewInt(7)).Bytes())

	m.SendOnchainMessage(protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(3), dest))
	m.DeliverOnchainMessage()

	knowninbox.SendMessage(protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(3), dest))
	knowninbox.DeliverMessages()

	vm.NewMachineAssertionContext(m, [2]uint64{0, 10000})

	var tokint big.Int
	var bigtok [32]byte
	bigtok[0] = 15
	bigtok[20] = 1
	tokint.SetBytes(bigtok[:])
	var vals [8]value.Value
	vals[0] = value.NewInt64Value(1)
	vals[1] = value.NewIntValue(&tokint)
	vals[2] = value.NewInt64Value(3)
	vals[3] = value.NewInt64Value(4)
	tup, _ := value.NewTupleOfSizeWithContents(vals, 4)

	m.Stack().Push(tup)
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.INBOX}); err != nil {
		tmp := "INBOX failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	knownMachine.Stack().Push(knowninbox.Receive())
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestJump(t *testing.T) {
	//test:
	insns := make([]value.Operation, 5)
	i := 0 // insn 0
	insns[i] = value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(1)}
	i++ // insn 1
	insns[i] = value.ImmediateOperation{Op: code.ADD, Val: value.NewInt64Value(4)}
	i++ // insn 2
	insns[i] = value.ImmediateOperation{Op: code.SUB, Val: value.NewInt64Value(5)}
	i++ // insn 3
	insns[i] = value.BasicOperation{Op: code.LOG}
	i++ // insn 4
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// run NOP to push value 1
	vm.RunInstruction(m, m.GetOperation())
	// push 2 to set jump point
	var nextHash [32]byte
	codept := value.CodePointValue{2, value.BasicOperation{Op: code.SUB}, nextHash}
	m.Stack().Push(codept)
	// JUMP
	vm.RunInstruction(m, value.BasicOperation{Op: code.JUMP})
	// PC should now be 2 - immediate operation that pushes 5 and subtracts
	vm.RunInstruction(m, m.GetOperation())
	// verify sub was executed
	knownMachine.Stack().Push(value.NewInt64Value(4))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestCJump(t *testing.T) {
	//test:
	insns := make([]value.Operation, 5)
	i := 0 // insn 0
	insns[i] = value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(1)}
	i++ // insn 1
	insns[i] = value.ImmediateOperation{Op: code.ADD, Val: value.NewInt64Value(4)}
	i++ // insn 2
	insns[i] = value.ImmediateOperation{Op: code.SUB, Val: value.NewInt64Value(5)}
	i++ // insn 3
	insns[i] = value.BasicOperation{Op: code.LOG}
	i++ // insn 4
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	saveMachine := m.Clone().(*vm.Machine)
	saveKnownMachine := knownMachine.Clone().(*vm.Machine)

	// run NOP to push value 1
	vm.RunInstruction(m, m.GetOperation())
	// push 0 for conditional
	m.Stack().Push(value.NewInt64Value(0))
	// push 2 to set jump point
	var nextHash [32]byte
	codept := value.CodePointValue{2, value.BasicOperation{Op: code.SUB}, nextHash}
	m.Stack().Push(codept)
	// CJUMP
	vm.RunInstruction(m, value.BasicOperation{Op: code.CJUMP})
	// PC should now be 2 - immediate operation that pushes 5 and subtracts
	vm.RunInstruction(m, m.GetOperation())
	// verify sub was executed
	knownMachine.Stack().Push(value.NewInt64Value(4))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}

	// repeat test with conditional set to 1
	m = saveMachine
	knownMachine = saveKnownMachine
	// run NOP to push value 1
	vm.RunInstruction(m, m.GetOperation())
	// push 1 for conditional
	m.Stack().Push(value.NewInt64Value(1))
	// push 2 to set jump point
	codept = value.CodePointValue{2, value.BasicOperation{Op: code.SUB}, nextHash}
	m.Stack().Push(codept)
	// CJUMP
	vm.RunInstruction(m, value.BasicOperation{Op: code.CJUMP})
	// PC should now be 2 - immediate operation that pushes 5 and subtracts
	vm.RunInstruction(m, m.GetOperation())
	// verify sub was executed
	knownMachine.Stack().Push(value.NewInt64Value(4))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestStackempty(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.STACKEMPTY}); err != nil {
		tmp := "STACKEMPTY failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// stack should have value 1 - stack was empty
	a := m.Stack().Count()
	if a != 1 {
		tmp := "STACKEMPTY stack size check failed expected 3 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 1
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.STACKEMPTY}); err != nil {
		tmp := "STACKEMPTY failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	a = m.Stack().Count()
	if a != 2 {
		tmp := "STACKEMPTY stack size check failed expected 7 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// push 0 to knownMachine as result of second STACKEMPTY call
	knownMachine.Stack().Push(value.NewInt64Value(0))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestPcpush(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.PCPUSH}); err != nil {
		tmp := "PCPUSH failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// stack should have one item - current codepoint
	a := m.Stack().Count()
	if a != 1 {
		tmp := "PCPUSH stack size check failed expected 3 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 1
	var nextHash [32]byte
	codept := value.CodePointValue{0, value.BasicOperation{Op: code.HALT}, nextHash}
	knownMachine.Stack().Push(codept)
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestAuxpush(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(4))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.AUXPUSH}); err != nil {
		tmp := "AUXPUSH failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// auxstack should have one item - value popped from stack
	a := m.AuxStack().Count()
	if a != 1 {
		tmp := "AUXPUSH stack size check failed expected 1 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// stack should be empty
	a = m.Stack().Count()
	if a != 0 {
		tmp := "AUXPUSH stack size check failed expected 3 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 4
	knownMachine.AuxStack().Push(value.NewInt64Value(4))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestAuxpop(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.AuxStack().Push(value.NewInt64Value(5))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.AUXPOP}); err != nil {
		tmp := "AUXPOP failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// auxstack should be empty
	a := m.AuxStack().Count()
	if a != 0 {
		tmp := "AUXPOP stack size check failed expected 1 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// stack should have one item - value popped from auxstack
	a = m.Stack().Count()
	if a != 1 {
		tmp := "AUXPOP stack size check failed expected 3 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 5
	knownMachine.Stack().Push(value.NewInt64Value(5))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestAuxstckempty(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// auxstack should be empty
	a := m.AuxStack().Count()
	if a != 0 {
		tmp := "AUXPOP stack size check failed expected 1 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	// check aux stack empty and push results on data stack
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.AUXSTACKEMPTY}); err != nil {
		tmp := "AUXSTACKEMPTY failed"
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 1
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}

	m.AuxStack().Push(value.NewInt64Value(5))
	// auxstack should not be empty
	a = m.AuxStack().Count()
	if a != 1 {
		tmp := "AUXSTACKEMPTY stack size check failed expected 3 found "
		tmp += strconv.FormatInt(a, 10)
		t.Error(tmp)
	}
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.AUXSTACKEMPTY}); err != nil {
		tmp := "AUXSTACKEMPTY failed"
		tmp += err.Error()
		t.Error(err)
	}
	// verify known and unknown match
	knownMachine.AuxStack().Push(value.NewInt64Value(5))
	knownMachine.Stack().Push(value.NewInt64Value(0))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestNop(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// check NOP does nothing
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.NOP}); err != nil {
		tmp := "NOP failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}

	// check NOP does nothing
	// immediate operation pushes value then does nothing
	if _, err := vm.RunInstruction(m, value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(1)}); err != nil {
		tmp := "NOP failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestErrpush(t *testing.T) {
	// test
	insns := make([]value.Operation, 5)
	i := 0 // insn 0
	insns[i] = value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(1)}
	i++ // insn 1
	insns[i] = value.ImmediateOperation{Op: code.ADD, Val: value.NewInt64Value(4)}
	i++ // insn 2
	insns[i] = value.ImmediateOperation{Op: code.SUB, Val: value.NewInt64Value(5)}
	i++ // insn 3
	insns[i] = value.BasicOperation{Op: code.LOG}
	i++ // insn 4
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// push codepoint onto stack
	var nextHash [32]byte
	codept := value.CodePointValue{4, value.BasicOperation{Op: code.HALT}, nextHash}
	m.Stack().Push(codept)
	knownMachine.Stack().Push(codept)
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// run errset to set the error handler
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.ERRSET}); err != nil {
		tmp := "ERRSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown different
	if ok, _ := vm.Equal(knownMachine, m); ok {
		tmp := "machines equal expected different"
		t.Error(tmp)
	}
	// set known to match
	if _, err := vm.RunInstruction(knownMachine, value.BasicOperation{Op: code.ERRSET}); err != nil {
		tmp := "ERRSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// run errpush to push error handler to data stack
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.ERRPUSH}); err != nil {
		tmp := "ERRPUSH failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown different
	if ok, _ := vm.Equal(knownMachine, m); ok {
		tmp := "machines equal expected different"
		t.Error(tmp)
	}
	// push error handler code point to known data stack
	knownMachine.Stack().Push(codept)
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestErrset(t *testing.T) {
	// test
	insns := make([]value.Operation, 5)
	i := 0 // insn 0
	insns[i] = value.ImmediateOperation{Op: code.NOP, Val: value.NewInt64Value(1)}
	i++ // insn 1
	insns[i] = value.ImmediateOperation{Op: code.ADD, Val: value.NewInt64Value(4)}
	i++ // insn 2
	insns[i] = value.ImmediateOperation{Op: code.SUB, Val: value.NewInt64Value(5)}
	i++ // insn 3
	insns[i] = value.BasicOperation{Op: code.LOG}
	i++ // insn 4
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// push codepoint onto stack
	var nextHash [32]byte
	codept := value.CodePointValue{4, value.BasicOperation{Op: code.HALT}, nextHash}
	m.Stack().Push(codept)
	knownMachine.Stack().Push(codept)
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// run errset to set the error handler
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.ERRSET}); err != nil {
		tmp := "ERRSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown different
	if ok, _ := vm.Equal(knownMachine, m); ok {
		tmp := "machines equal expected different"
		t.Error(tmp)
	}
	// set known to match
	if _, err := vm.RunInstruction(knownMachine, value.BasicOperation{Op: code.ERRSET}); err != nil {
		tmp := "ERRSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// check NOP does nothing
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.ERROR}); err == nil {
		tmp := "ERROR failed - should have generated error"
		t.Error(tmp)
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestDup0(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.DUP0}); err != nil {
		tmp := "DUP0 failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestDup1(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))
	m.Stack().Push(value.NewInt64Value(2))
	knownMachine.Stack().Push(value.NewInt64Value(1))
	knownMachine.Stack().Push(value.NewInt64Value(2))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.DUP1}); err != nil {
		tmp := "DUP1 failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestDup2(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))
	m.Stack().Push(value.NewInt64Value(2))
	m.Stack().Push(value.NewInt64Value(3))
	knownMachine.Stack().Push(value.NewInt64Value(1))
	knownMachine.Stack().Push(value.NewInt64Value(2))
	knownMachine.Stack().Push(value.NewInt64Value(3))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.DUP2}); err != nil {
		tmp := "DUP2 failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestSwap2(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(1))
	m.Stack().Push(value.NewInt64Value(2))
	m.Stack().Push(value.NewInt64Value(3))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.SWAP2}); err != nil {
		tmp := "SWAP2 failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(3))
	knownMachine.Stack().Push(value.NewInt64Value(2))
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestTget(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	tup := value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2))

	m.Stack().Push(tup)
	m.Stack().Push(value.NewInt64Value(1))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TGET}); err != nil {
		tmp := "TGET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match one item value = 1
	knownMachine.Stack().Push(value.NewInt64Value(2))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// test with only int on stack
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TGET}); err == nil {
		tmp := "TGET expected fail"
		t.Error(tmp)
	}
	// verify known and unknown match expect empty stack
	knownMachine.Stack().Pop()
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// test A out of range
	m.Stack().Push(value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2)))
	m.Stack().Push(value.NewInt64Value(3))
	var nextHash [32]byte
	codept := value.CodePointValue{0, value.BasicOperation{Op: code.HALT}, nextHash}
	m.SetPC(codept)
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TGET}); err == nil {
		tmp := "TGET expected fail"
		t.Error(tmp)
	}
	// verify known and unknown match expect empty stack
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestTset(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewInt64Value(3))
	m.Stack().Push(value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2)))
	m.Stack().Push(value.NewInt64Value(1))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TSET}); err != nil {
		tmp := "TSET failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(3)))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// test with only tuple on stack
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TSET}); err == nil {
		tmp := "TSET expected fail"
		t.Error(tmp)
	}
	// verify known and unknown match expect empty stack
	knownMachine.Stack().Pop()
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// test incorrect A
	m.Stack().Push(value.NewInt64Value(3))
	m.Stack().Push(value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2)))
	m.Stack().Push(value.NewInt64Value(4))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TSET}); err == nil {
		tmp := "TSET expected fail"
		t.Error(tmp)
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestTlen(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.Stack().Push(value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2)))
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TLEN}); err != nil {
		tmp := "TLEN failed - "
		tmp += err.Error()
		t.Error(tmp)
	}
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(2))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// test A not a tuple
	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TLEN}); err == nil {
		t.Error("TLEN expected fail")
	}
	// verify known and unknown match expect empty stack
	knownMachine.Stack().Pop()
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestType(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	testValues := []value.Value{
		value.NewEmptyTuple(),
		value.NewTuple2(value.NewInt64Value(1), value.NewInt64Value(2)),
		value.ErrorCodePoint,
		value.NewInt64Value(100),
	}

	resultValues := []value.Value{
		value.NewInt64Value(3),
		value.NewInt64Value(3),
		value.NewInt64Value(1),
		value.NewInt64Value(0),
	}

	for i := range testValues {
		m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
		knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

		m.Stack().Push(testValues[i])
		if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.TYPE}); err != nil {
			tmp := "TYPE failed - "
			tmp += err.Error()
			t.Error(tmp)
		}

		knownMachine.Stack().Push(resultValues[i])

		if ok, err := vm.Equal(knownMachine, m); !ok {
			t.Error(err)
		}
	}
}

func TestBreakpoint(t *testing.T) {
	// test
	insns := make([]value.Operation, 1)
	i := 0
	insns[i] = value.BasicOperation{Op: code.HALT}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	if _, err := vm.RunInstruction(m, value.BasicOperation{Op: code.BREAKPOINT}); err == nil {
		t.Error("Breakpoint didn't block")
	}
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}

func TestLog(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.LOG},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	m.Stack().Push(value.NewInt64Value(5))
	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	// verify out message
	logs := ad.Logs
	if len(logs) != 1 {
		t.Error("No log value generated")
	}
	if !logs[0].Equal(value.NewInt64Value(5)) {
		t.Error("log value incorrect")
	}
}

func TestSendFungible(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.SEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// fungible value=10
	var tok protocol.TokenType
	tok[0] = 15
	tok[20] = 0
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(7),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(10), [32]byte{}))

	// send token 15 value=7 to dest 4
	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 1 {
		t.Error("No out message generated")
	}

	dest := [32]byte{}
	dest[31] = 4
	knownmessage := protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(7), dest)
	if !msgs[0].Equals(knownmessage) {
		t.Error("Out message incorrect")
	}
}

func TestSendNonFungible(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.SEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// test send of non fungible
	var tok protocol.TokenType
	tok[0] = 16
	tok[20] = 1
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(7),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(7), [32]byte{}))

	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 1 {
		t.Error("No out message generated")
	}

	dest := [32]byte{}
	dest[31] = 4
	knownmessage := protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(7), dest)
	if !msgs[0].Equals(knownmessage) {
		t.Error("Out message incorrect")
	}
}

func TestSendLowBalance(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.SEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// test send with insufficient funds
	var tok protocol.TokenType
	tok[0] = 17
	tok[20] = 0
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(17),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(10), [32]byte{}))

	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	knownMachine.Stack().Push(tup)
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 0 {
		t.Error("No out message generated")
	}
}

func TestNbsend1(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.NBSEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	var tok protocol.TokenType
	tok[0] = 15
	tok[20] = 1
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(10),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})

	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(10), [32]byte{}))

	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))

	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}

	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 1 {
		t.Error("No out message generated")
	}
}

func TestNBSendFungible(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.NBSEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// fungible value=10
	var tok protocol.TokenType
	tok[0] = 15
	tok[20] = 0
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(7),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(10), [32]byte{}))

	// send token 15 value=7 to dest 4
	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 1 {
		t.Error("No out message generated")
	}

	dest := [32]byte{}
	dest[31] = 4
	knownmessage := protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(7), dest)
	if !msgs[0].Equals(knownmessage) {
		t.Error("Out message incorrect")
	}
}

func TestNBSendNonFungible(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.NBSEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// test send of non fungible
	var tok protocol.TokenType
	tok[0] = 16
	tok[20] = 1
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(7),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(7), [32]byte{}))

	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(1))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 1 {
		t.Error("No out message generated")
	}

	dest := [32]byte{}
	dest[31] = 4
	knownmessage := protocol.NewMessage(value.NewInt64Value(1), tok, big.NewInt(7), dest)
	if !msgs[0].Equals(knownmessage) {
		t.Error("Out message incorrect")
	}
}

func TestNBSendLowBalance(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.NBSEND},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	// test send with insufficient funds
	var tok protocol.TokenType
	tok[0] = 17
	tok[20] = 0
	var bigtok [32]byte
	copy(bigtok[:], tok[:])
	tup, _ := value.NewTupleFromSlice([]value.Value{
		value.NewInt64Value(1),
		value.NewInt64Value(4),
		value.NewInt64Value(17),
		value.NewIntValue(new(big.Int).SetBytes(bigtok[:])),
	})
	m.Stack().Push(tup)

	// add tokens to balanceTracker
	m.SendOnchainMessage(protocol.NewMessage(value.NewEmptyTuple(), tok, big.NewInt(10), [32]byte{}))

	ad := m.ExecuteAssertion(10, protocol.NewTimeBounds(0, 1000))
	// verify known and unknown match
	knownMachine.Stack().Push(value.NewInt64Value(0))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
	msgs := ad.OutMsgs
	// verify out message
	if len(msgs) != 0 {
		t.Error("No out message generated")
	}
}

func TestGettime(t *testing.T) {
	// test
	insns := []value.Operation{
		value.BasicOperation{Op: code.GETTIME},
		value.BasicOperation{Op: code.HALT},
	}

	m := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)
	knownMachine := vm.NewMachine(insns, value.NewInt64Value(1), false, 100)

	m.ExecuteAssertion(10, [2]uint64{5, 10})

	// verify known and unknown match
	knownMachine.Stack().Push(value.NewTuple2(value.NewInt64Value(5), value.NewInt64Value(10)))
	if ok, err := vm.Equal(knownMachine, m); !ok {
		t.Error(err)
	}
}
