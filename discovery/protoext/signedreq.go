/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import (
	"github.com/gogo/protobuf/proto"
	"github.com/inklabsfoundation/inkchain/protos/discovery"
)

// SignedRequestToRequest deserializes this SignedRequest's payload
// and returns the serialized Request in its object form.
// Returns an error in case the operation fails.
func SignedRequestToRequest(sr *discovery.SignedRequest) (*discovery.Request, error) {
	req := &discovery.Request{}
	return req, proto.Unmarshal(sr.Payload, req)
}
