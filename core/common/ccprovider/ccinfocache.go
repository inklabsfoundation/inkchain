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

package ccprovider

import (
	"fmt"
	"sync"
)

// ccInfoCacheImpl implements in-memory cache for ChaincodeData
// needed by endorser to verify if the local instantiation policy
// matches the instantiation policy on a channel before honoring
// an invoke
type ccInfoCacheImpl struct {
	sync.RWMutex

	cache        map[string]*ChaincodeData
	cacheSupport CCCacheSupport
}

// NewCCInfoCache returns a new cache on top of the supplied CCInfoProvider instance
func NewCCInfoCache(cs CCCacheSupport) *ccInfoCacheImpl {
	return &ccInfoCacheImpl{
		cache:        make(map[string]*ChaincodeData),
		cacheSupport: cs,
	}
}

func (c *ccInfoCacheImpl) GetChaincodeData(ccname string, ccversion string) (*ChaincodeData, error) {
	// c.cache is guaranteed to be non-nil

	key := ccname + "/" + ccversion

	c.RLock()
	ccdata, in := c.cache[key]
	c.RUnlock()

	if !in {
		var err error

		// the chaincode data is not in the cache
		// try to look it up from the file system
		ccpack, err := c.cacheSupport.GetChaincode(ccname, ccversion)
		if err != nil || ccpack == nil {
			return nil, fmt.Errorf("cannot retrieve package for chaincode %s/%s, error %s", ccname, ccversion, err)
		}

		// we have a non-nil ChaincodeData, put it in the cache
		c.Lock()
		ccdata = ccpack.GetChaincodeData()
		c.cache[key] = ccdata
		c.Unlock()
	}

	return ccdata, nil
}
