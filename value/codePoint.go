package value

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/offchainlabs/arb-avm/code"
	"golang.org/x/crypto/sha3"
	"io"
)

type Operation interface {
	GetOp() code.Opcode
	TypeCode() uint8
	Marshal(wr io.Writer) error
}

type BasicOperation struct {
	Op code.Opcode
}

type ImmediateOperation struct {
	Op  code.Opcode
	Val Value
}

func NewBasicOperationFromReader(rd io.Reader) (BasicOperation, error) {
	op, err := code.NewOpcodeFromReader(rd)
	return BasicOperation{op}, err
}

func NewImmediateOperationFromReader(rd io.Reader) (ImmediateOperation, error) {
	op, err := code.NewOpcodeFromReader(rd)
	if err != nil {
		return ImmediateOperation{}, err
	}
	val, err := UnmarshalValue(rd)
	return ImmediateOperation{op, val}, err
}

func (op BasicOperation) Marshal(wr io.Writer) error {
	return op.Op.Marshal(wr)
}

func (op ImmediateOperation) Marshal(wr io.Writer) error {
	if err := op.Op.Marshal(wr); err != nil {
		return err
	}
	return MarshalValue(op.Val, wr)
}

func (op BasicOperation) TypeCode() uint8 {
	return 0
}

func (op ImmediateOperation) TypeCode() uint8 {
	return 1
}

func (op BasicOperation) GetOp() code.Opcode {
	return op.Op
}

func (op BasicOperation) String() string {
	return fmt.Sprintf("Basic(%v)", code.InstructionNames[op.Op])
}

func (op ImmediateOperation) String() string {
	return fmt.Sprintf("Immediate(%v, %v)", code.InstructionNames[op.Op], op.Val)
}

func (op ImmediateOperation) GetOp() code.Opcode {
	return op.Op
}

func NewOperationFromReader(rd io.Reader) (Operation, error) {
	var immediateCount uint8
	err := binary.Read(rd, binary.BigEndian, &immediateCount)
	if err != nil {
		return nil, err
	}
	if immediateCount == 0 {
		return NewBasicOperationFromReader(rd)
	} else if immediateCount == 1 {
		return NewImmediateOperationFromReader(rd)
	} else {
		return nil, errors.New("immediate count must be 0 or 1")
	}
}

func MarshalOperation(op Operation, wr io.Writer) error {
	typ := op.TypeCode()
	if err := binary.Write(wr, binary.BigEndian, &typ); err != nil {
		return err
	}
	return op.Marshal(wr)
}

const CodePointCode = 1

func NewCodePointForProofFromReader(rd io.Reader) (CodePointValue, error) {
	var op Operation
	op, err := NewOperationFromReader(rd)
	if err != nil {
		return CodePointValue{}, err
	}
	var nextHash [32]byte
	_, err = io.ReadFull(rd, nextHash[:])
	return CodePointValue{0, op, nextHash}, err
}

type CodePointValue struct {
	InsnNum  int64
	Op       Operation
	NextHash [32]byte
}

//func NewCodePointValue(point CodePointValue) CodePointValue {
//	return CodePointValue{point}
//}

func NewCodePointValueFromReader(rd io.Reader) (CodePointValue, error) {
	var insnNum int64
	if err := binary.Read(rd, binary.BigEndian, &insnNum); err != nil {
		return CodePointValue{}, err
	}
	var op Operation
	op, err := NewOperationFromReader(rd)
	if err != nil {
		return CodePointValue{}, err
	}
	var nextHash [32]byte
	_, err = io.ReadFull(rd, nextHash[:])
	return CodePointValue{insnNum, op, nextHash}, err
}

func (cv CodePointValue) TypeCode() uint8 {
	return TypeCodeCodePoint
}

func (cv CodePointValue) InternalTypeCode() uint8 {
	return TypeCodeCodePoint
}

func (cv CodePointValue) Clone() Value {
	return CodePointValue{cv.InsnNum, cv.Op, cv.NextHash}
}

func (cv CodePointValue) CloneShallow() Value {
	return CodePointValue{cv.InsnNum, cv.Op, cv.NextHash}
}

func (cv CodePointValue) Equal(val Value) bool {
	if val.TypeCode() == TypeCodeHashOnly {
		return cv.Hash() == val.Hash()
	} else if val.TypeCode() != TypeCodeCodePoint {
		return false
	} else {
		if cv.InsnNum != val.(CodePointValue).InsnNum {
			return false
		}
		// for now only check InsnNum
		//if cv.Op != val.(CodePointValue).Op {
		//	return false
		//}
		//if cv.NextHash != val.(CodePointValue).NextHash {
		//	return false
		//}
		return true
	}
}

func (cv CodePointValue) Size() int64 {
	return 1
}

var ErrorCodePointHash [32]byte
var HaltCodePointHash [32]byte

var ErrorCodePoint CodePointValue
var HaltCodePoint CodePointValue

func init() {
	ErrorCodePointHash = sha256.Sum256([]byte("ErrorCodePointHash"))
	HaltCodePointHash = sha256.Sum256([]byte("HaltCodePointHash"))

	HaltCodePoint = CodePointValue{-1, BasicOperation{code.NOP}, [32]byte{}}
	ErrorCodePoint = CodePointValue{-2, BasicOperation{code.NOP}, [32]byte{}}
}

func (cv CodePointValue) Hash() [32]byte {
	if cv.InsnNum == -1 {
		return HaltCodePointHash
	} else if cv.InsnNum == -2 {
		return ErrorCodePointHash
	}

	switch op := cv.Op.(type) {
	case ImmediateOperation:
		var codePointData [66]byte
		codePointData[0] = CodePointCode
		codePointData[1] = byte(op.Op)
		valHash := op.Val.Hash()
		copy(codePointData[2:], valHash[:])
		copy(codePointData[34:], cv.NextHash[:])
		d := sha3.NewLegacyKeccak256()
		d.Write(codePointData[:])
		ret := [32]byte{}
		d.Sum(ret[:0])
		return ret
	case BasicOperation:
		var codePointData [34]byte
		codePointData[0] = CodePointCode
		codePointData[1] = byte(op.Op)
		copy(codePointData[2:], cv.NextHash[:])
		d := sha3.NewLegacyKeccak256()
		d.Write(codePointData[:])
		ret := [32]byte{}
		d.Sum(ret[:0])
		return ret
	default:
		panic(fmt.Sprintf("Bad operation type: %T in with pc %d", op, cv.InsnNum))
	}
}

func (cv CodePointValue) Marshal(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, &cv.InsnNum); err != nil {
		return err
	}
	if err := cv.Op.Marshal(w); err != nil {
		return err
	}
	_, err := w.Write(cv.NextHash[:])
	return err
}

func (cv CodePointValue) MarshalForProof(w io.Writer) error {
	if err := cv.Op.Marshal(w); err != nil {
		return err
	}
	_, err := w.Write(cv.NextHash[:])
	return err
}

func (cv CodePointValue) String() string {
	return fmt.Sprintf("CodePoint(%v, %v)", cv.InsnNum, cv.Op)
}
