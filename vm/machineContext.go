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

package vm

import (
	"github.com/offchainlabs/arb-util/machine"
	"github.com/offchainlabs/arb-util/protocol"
	"github.com/offchainlabs/arb-util/value"
)

type MachineAssertionContext struct {
	machine    *Machine
	timeBounds protocol.TimeBounds
	numSteps   uint32
	outMsgs    []protocol.Message
	logs       []value.Value
}

func NewMachineAssertionContext(m *Machine, timeBounds protocol.TimeBounds) *MachineAssertionContext {
	outMsgs := make([]protocol.Message, 0)
	ret := &MachineAssertionContext{
		m,
		timeBounds,
		0,
		outMsgs,
		make([]value.Value, 0),
	}
	ret.machine.SetContext(ret)
	return ret
}

func (ac *MachineAssertionContext) LoggedValue(data value.Value) error {
	ac.logs = append(ac.logs, data)
	return nil
}

func (ac *MachineAssertionContext) Send(data value.Value, tokenType value.IntValue, currency value.IntValue, dest value.IntValue) error {
	tokType := [21]byte{}
	tokBytes := tokenType.ToBytes()
	copy(tokType[:], tokBytes[:])
	newMsg := protocol.NewMessage(data, tokType, currency.BigInt(), dest.ToBytes())
	ac.outMsgs = append(ac.outMsgs, newMsg)
	return nil
}

func (ac *MachineAssertionContext) OutMessageCount() int {
	return len(ac.outMsgs)
}

func (ac *MachineAssertionContext) GetTimeBounds() value.Value {
	return ac.timeBounds.AsValue()
}

func (ac *MachineAssertionContext) NotifyStep() {
	ac.numSteps++
}

func (ac *MachineAssertionContext) Finalize(m *Machine) *protocol.Assertion {
	ac.machine.SetContext(&machine.MachineNoContext{})
	return protocol.NewAssertion(ac.machine.Hash(), ac.numSteps, ac.outMsgs, ac.logs)
}

func (ac *MachineAssertionContext) EndContext() {
	ac.machine.SetContext(&machine.MachineNoContext{})
}
