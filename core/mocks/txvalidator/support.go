/*
Copyright Ziggurat Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package support

import (
	mockpolicies "github.com/inklabsfoundation/inkchain/common/mocks/policies"
	"github.com/inklabsfoundation/inkchain/common/policies"
	"github.com/inklabsfoundation/inkchain/core/ledger"
	"github.com/inklabsfoundation/inkchain/msp"
	"github.com/inklabsfoundation/inkchain/protos/common"
)

type Support struct {
	LedgerVal     ledger.PeerLedger
	MSPManagerVal msp.MSPManager
	ApplyVal      error
}

// Ledger returns LedgerVal
func (ms *Support) Ledger() ledger.PeerLedger {
	return ms.LedgerVal
}

// MSPManager returns MSPManagerVal
func (ms *Support) MSPManager() msp.MSPManager {
	return ms.MSPManagerVal
}

// Apply returns ApplyVal
func (ms *Support) Apply(configtx *common.ConfigEnvelope) error {
	return ms.ApplyVal
}

func (ms *Support) PolicyManager() policies.Manager {
	return &mockpolicies.Manager{}
}

func (cs *Support) GetMSPIDs(cid string) []string {
	return []string{"DEFAULT"}
}
