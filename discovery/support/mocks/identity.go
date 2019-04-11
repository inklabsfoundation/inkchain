// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
	"time"

	"github.com/inklabsfoundation/inkchain/msp"
	mspprotos "github.com/inklabsfoundation/inkchain/protos/msp"
)

type Identity struct {
	ExpiresAtStub        func() time.Time
	expiresAtMutex       sync.RWMutex
	expiresAtArgsForCall []struct{}
	expiresAtReturns     struct {
		result1 time.Time
	}
	expiresAtReturnsOnCall map[int]struct {
		result1 time.Time
	}
	GetIdentifierStub        func() *msp.IdentityIdentifier
	getIdentifierMutex       sync.RWMutex
	getIdentifierArgsForCall []struct{}
	getIdentifierReturns     struct {
		result1 *msp.IdentityIdentifier
	}
	getIdentifierReturnsOnCall map[int]struct {
		result1 *msp.IdentityIdentifier
	}
	GetMSPIdentifierStub        func() string
	getMSPIdentifierMutex       sync.RWMutex
	getMSPIdentifierArgsForCall []struct{}
	getMSPIdentifierReturns     struct {
		result1 string
	}
	getMSPIdentifierReturnsOnCall map[int]struct {
		result1 string
	}
	ValidateStub        func() error
	validateMutex       sync.RWMutex
	validateArgsForCall []struct{}
	validateReturns     struct {
		result1 error
	}
	validateReturnsOnCall map[int]struct {
		result1 error
	}
	GetOrganizationalUnitsStub        func() []*msp.OUIdentifier
	getOrganizationalUnitsMutex       sync.RWMutex
	getOrganizationalUnitsArgsForCall []struct{}
	getOrganizationalUnitsReturns     struct {
		result1 []*msp.OUIdentifier
	}
	getOrganizationalUnitsReturnsOnCall map[int]struct {
		result1 []*msp.OUIdentifier
	}
	AnonymousStub        func() bool
	anonymousMutex       sync.RWMutex
	anonymousArgsForCall []struct{}
	anonymousReturns     struct {
		result1 bool
	}
	anonymousReturnsOnCall map[int]struct {
		result1 bool
	}
	VerifyStub        func(msg []byte, sig []byte) error
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		msg []byte
		sig []byte
	}
	verifyReturns struct {
		result1 error
	}
	verifyReturnsOnCall map[int]struct {
		result1 error
	}
	SerializeStub        func() ([]byte, error)
	serializeMutex       sync.RWMutex
	serializeArgsForCall []struct{}
	serializeReturns     struct {
		result1 []byte
		result2 error
	}
	serializeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SatisfiesPrincipalStub        func(principal *mspprotos.MSPPrincipal) error
	satisfiesPrincipalMutex       sync.RWMutex
	satisfiesPrincipalArgsForCall []struct {
		principal *mspprotos.MSPPrincipal
	}
	satisfiesPrincipalReturns struct {
		result1 error
	}
	satisfiesPrincipalReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Identity) ExpiresAt() time.Time {
	fake.expiresAtMutex.Lock()
	ret, specificReturn := fake.expiresAtReturnsOnCall[len(fake.expiresAtArgsForCall)]
	fake.expiresAtArgsForCall = append(fake.expiresAtArgsForCall, struct{}{})
	fake.recordInvocation("ExpiresAt", []interface{}{})
	fake.expiresAtMutex.Unlock()
	if fake.ExpiresAtStub != nil {
		return fake.ExpiresAtStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.expiresAtReturns.result1
}

func (fake *Identity) ExpiresAtCallCount() int {
	fake.expiresAtMutex.RLock()
	defer fake.expiresAtMutex.RUnlock()
	return len(fake.expiresAtArgsForCall)
}

func (fake *Identity) ExpiresAtReturns(result1 time.Time) {
	fake.ExpiresAtStub = nil
	fake.expiresAtReturns = struct {
		result1 time.Time
	}{result1}
}

func (fake *Identity) ExpiresAtReturnsOnCall(i int, result1 time.Time) {
	fake.ExpiresAtStub = nil
	if fake.expiresAtReturnsOnCall == nil {
		fake.expiresAtReturnsOnCall = make(map[int]struct {
			result1 time.Time
		})
	}
	fake.expiresAtReturnsOnCall[i] = struct {
		result1 time.Time
	}{result1}
}

func (fake *Identity) GetIdentifier() *msp.IdentityIdentifier {
	fake.getIdentifierMutex.Lock()
	ret, specificReturn := fake.getIdentifierReturnsOnCall[len(fake.getIdentifierArgsForCall)]
	fake.getIdentifierArgsForCall = append(fake.getIdentifierArgsForCall, struct{}{})
	fake.recordInvocation("GetIdentifier", []interface{}{})
	fake.getIdentifierMutex.Unlock()
	if fake.GetIdentifierStub != nil {
		return fake.GetIdentifierStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getIdentifierReturns.result1
}

func (fake *Identity) GetIdentifierCallCount() int {
	fake.getIdentifierMutex.RLock()
	defer fake.getIdentifierMutex.RUnlock()
	return len(fake.getIdentifierArgsForCall)
}

func (fake *Identity) GetIdentifierReturns(result1 *msp.IdentityIdentifier) {
	fake.GetIdentifierStub = nil
	fake.getIdentifierReturns = struct {
		result1 *msp.IdentityIdentifier
	}{result1}
}

func (fake *Identity) GetIdentifierReturnsOnCall(i int, result1 *msp.IdentityIdentifier) {
	fake.GetIdentifierStub = nil
	if fake.getIdentifierReturnsOnCall == nil {
		fake.getIdentifierReturnsOnCall = make(map[int]struct {
			result1 *msp.IdentityIdentifier
		})
	}
	fake.getIdentifierReturnsOnCall[i] = struct {
		result1 *msp.IdentityIdentifier
	}{result1}
}

func (fake *Identity) GetMSPIdentifier() string {
	fake.getMSPIdentifierMutex.Lock()
	ret, specificReturn := fake.getMSPIdentifierReturnsOnCall[len(fake.getMSPIdentifierArgsForCall)]
	fake.getMSPIdentifierArgsForCall = append(fake.getMSPIdentifierArgsForCall, struct{}{})
	fake.recordInvocation("GetMSPIdentifier", []interface{}{})
	fake.getMSPIdentifierMutex.Unlock()
	if fake.GetMSPIdentifierStub != nil {
		return fake.GetMSPIdentifierStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getMSPIdentifierReturns.result1
}

func (fake *Identity) GetMSPIdentifierCallCount() int {
	fake.getMSPIdentifierMutex.RLock()
	defer fake.getMSPIdentifierMutex.RUnlock()
	return len(fake.getMSPIdentifierArgsForCall)
}

func (fake *Identity) GetMSPIdentifierReturns(result1 string) {
	fake.GetMSPIdentifierStub = nil
	fake.getMSPIdentifierReturns = struct {
		result1 string
	}{result1}
}

func (fake *Identity) GetMSPIdentifierReturnsOnCall(i int, result1 string) {
	fake.GetMSPIdentifierStub = nil
	if fake.getMSPIdentifierReturnsOnCall == nil {
		fake.getMSPIdentifierReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getMSPIdentifierReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *Identity) Validate() error {
	fake.validateMutex.Lock()
	ret, specificReturn := fake.validateReturnsOnCall[len(fake.validateArgsForCall)]
	fake.validateArgsForCall = append(fake.validateArgsForCall, struct{}{})
	fake.recordInvocation("Validate", []interface{}{})
	fake.validateMutex.Unlock()
	if fake.ValidateStub != nil {
		return fake.ValidateStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validateReturns.result1
}

func (fake *Identity) ValidateCallCount() int {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return len(fake.validateArgsForCall)
}

func (fake *Identity) ValidateReturns(result1 error) {
	fake.ValidateStub = nil
	fake.validateReturns = struct {
		result1 error
	}{result1}
}

func (fake *Identity) ValidateReturnsOnCall(i int, result1 error) {
	fake.ValidateStub = nil
	if fake.validateReturnsOnCall == nil {
		fake.validateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Identity) GetOrganizationalUnits() []*msp.OUIdentifier {
	fake.getOrganizationalUnitsMutex.Lock()
	ret, specificReturn := fake.getOrganizationalUnitsReturnsOnCall[len(fake.getOrganizationalUnitsArgsForCall)]
	fake.getOrganizationalUnitsArgsForCall = append(fake.getOrganizationalUnitsArgsForCall, struct{}{})
	fake.recordInvocation("GetOrganizationalUnits", []interface{}{})
	fake.getOrganizationalUnitsMutex.Unlock()
	if fake.GetOrganizationalUnitsStub != nil {
		return fake.GetOrganizationalUnitsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getOrganizationalUnitsReturns.result1
}

func (fake *Identity) GetOrganizationalUnitsCallCount() int {
	fake.getOrganizationalUnitsMutex.RLock()
	defer fake.getOrganizationalUnitsMutex.RUnlock()
	return len(fake.getOrganizationalUnitsArgsForCall)
}

func (fake *Identity) GetOrganizationalUnitsReturns(result1 []*msp.OUIdentifier) {
	fake.GetOrganizationalUnitsStub = nil
	fake.getOrganizationalUnitsReturns = struct {
		result1 []*msp.OUIdentifier
	}{result1}
}

func (fake *Identity) GetOrganizationalUnitsReturnsOnCall(i int, result1 []*msp.OUIdentifier) {
	fake.GetOrganizationalUnitsStub = nil
	if fake.getOrganizationalUnitsReturnsOnCall == nil {
		fake.getOrganizationalUnitsReturnsOnCall = make(map[int]struct {
			result1 []*msp.OUIdentifier
		})
	}
	fake.getOrganizationalUnitsReturnsOnCall[i] = struct {
		result1 []*msp.OUIdentifier
	}{result1}
}

func (fake *Identity) Anonymous() bool {
	fake.anonymousMutex.Lock()
	ret, specificReturn := fake.anonymousReturnsOnCall[len(fake.anonymousArgsForCall)]
	fake.anonymousArgsForCall = append(fake.anonymousArgsForCall, struct{}{})
	fake.recordInvocation("Anonymous", []interface{}{})
	fake.anonymousMutex.Unlock()
	if fake.AnonymousStub != nil {
		return fake.AnonymousStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.anonymousReturns.result1
}

func (fake *Identity) AnonymousCallCount() int {
	fake.anonymousMutex.RLock()
	defer fake.anonymousMutex.RUnlock()
	return len(fake.anonymousArgsForCall)
}

func (fake *Identity) AnonymousReturns(result1 bool) {
	fake.AnonymousStub = nil
	fake.anonymousReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Identity) AnonymousReturnsOnCall(i int, result1 bool) {
	fake.AnonymousStub = nil
	if fake.anonymousReturnsOnCall == nil {
		fake.anonymousReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.anonymousReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Identity) Verify(msg []byte, sig []byte) error {
	var msgCopy []byte
	if msg != nil {
		msgCopy = make([]byte, len(msg))
		copy(msgCopy, msg)
	}
	var sigCopy []byte
	if sig != nil {
		sigCopy = make([]byte, len(sig))
		copy(sigCopy, sig)
	}
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		msg []byte
		sig []byte
	}{msgCopy, sigCopy})
	fake.recordInvocation("Verify", []interface{}{msgCopy, sigCopy})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(msg, sig)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.verifyReturns.result1
}

func (fake *Identity) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *Identity) VerifyArgsForCall(i int) ([]byte, []byte) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].msg, fake.verifyArgsForCall[i].sig
}

func (fake *Identity) VerifyReturns(result1 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Identity) VerifyReturnsOnCall(i int, result1 error) {
	fake.VerifyStub = nil
	if fake.verifyReturnsOnCall == nil {
		fake.verifyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Identity) Serialize() ([]byte, error) {
	fake.serializeMutex.Lock()
	ret, specificReturn := fake.serializeReturnsOnCall[len(fake.serializeArgsForCall)]
	fake.serializeArgsForCall = append(fake.serializeArgsForCall, struct{}{})
	fake.recordInvocation("Serialize", []interface{}{})
	fake.serializeMutex.Unlock()
	if fake.SerializeStub != nil {
		return fake.SerializeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.serializeReturns.result1, fake.serializeReturns.result2
}

func (fake *Identity) SerializeCallCount() int {
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	return len(fake.serializeArgsForCall)
}

func (fake *Identity) SerializeReturns(result1 []byte, result2 error) {
	fake.SerializeStub = nil
	fake.serializeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Identity) SerializeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.SerializeStub = nil
	if fake.serializeReturnsOnCall == nil {
		fake.serializeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.serializeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Identity) SatisfiesPrincipal(principal *mspprotos.MSPPrincipal) error {
	fake.satisfiesPrincipalMutex.Lock()
	ret, specificReturn := fake.satisfiesPrincipalReturnsOnCall[len(fake.satisfiesPrincipalArgsForCall)]
	fake.satisfiesPrincipalArgsForCall = append(fake.satisfiesPrincipalArgsForCall, struct {
		principal *mspprotos.MSPPrincipal
	}{principal})
	fake.recordInvocation("SatisfiesPrincipal", []interface{}{principal})
	fake.satisfiesPrincipalMutex.Unlock()
	if fake.SatisfiesPrincipalStub != nil {
		return fake.SatisfiesPrincipalStub(principal)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.satisfiesPrincipalReturns.result1
}

func (fake *Identity) SatisfiesPrincipalCallCount() int {
	fake.satisfiesPrincipalMutex.RLock()
	defer fake.satisfiesPrincipalMutex.RUnlock()
	return len(fake.satisfiesPrincipalArgsForCall)
}

func (fake *Identity) SatisfiesPrincipalArgsForCall(i int) *mspprotos.MSPPrincipal {
	fake.satisfiesPrincipalMutex.RLock()
	defer fake.satisfiesPrincipalMutex.RUnlock()
	return fake.satisfiesPrincipalArgsForCall[i].principal
}

func (fake *Identity) SatisfiesPrincipalReturns(result1 error) {
	fake.SatisfiesPrincipalStub = nil
	fake.satisfiesPrincipalReturns = struct {
		result1 error
	}{result1}
}

func (fake *Identity) SatisfiesPrincipalReturnsOnCall(i int, result1 error) {
	fake.SatisfiesPrincipalStub = nil
	if fake.satisfiesPrincipalReturnsOnCall == nil {
		fake.satisfiesPrincipalReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.satisfiesPrincipalReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Identity) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.expiresAtMutex.RLock()
	defer fake.expiresAtMutex.RUnlock()
	fake.getIdentifierMutex.RLock()
	defer fake.getIdentifierMutex.RUnlock()
	fake.getMSPIdentifierMutex.RLock()
	defer fake.getMSPIdentifierMutex.RUnlock()
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	fake.getOrganizationalUnitsMutex.RLock()
	defer fake.getOrganizationalUnitsMutex.RUnlock()
	fake.anonymousMutex.RLock()
	defer fake.anonymousMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	fake.satisfiesPrincipalMutex.RLock()
	defer fake.satisfiesPrincipalMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Identity) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ msp.Identity = new(Identity)
