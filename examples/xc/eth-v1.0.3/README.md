# eth-v1.0.3

## XC.sol

 * Add modifier: nonzeroAddress(account)
```
    function setToken(address account) onlyAdmin nonzeroAddress(account) external {
        if (token != account) {
            token = Token(account);
        }
    }
```

 * Add modifier: nonzeroAddress(account)
```
    function setXCPlugin(address account) onlyAdmin nonzeroAddress(account) external {
        if (xcPlugin != account) {
            xcPlugin = XCPlugin(account);
        }
    }
```

## XCPlugin.sol

 * Add modifier: nonzeroAddress(caller)
```
    function addCaller(address caller) onlyAdmin nonzeroAddress(caller) external {
        if (!_existCaller(caller)) {
            callers.push(caller);
        }
    }
```

 * Add modifier: nonzeroAddress(caller)
 * Optimize for loop structure, remove extra var
```
    function deleteCaller(address caller) onlyAdmin nonzeroAddress(caller) external {
        for (uint i = 0; i < callers.length; i++) {
            if (callers[i] == caller) {
                if (i != callers.length - 1 ) {
                    callers[i] = callers[callers.length - 1];
                }
                callers.length--;
                return;
            }
        }
    }
```

 * Optimize for loop structure, remove extra var
```
    function bytes32ToStr(bytes32 b) internal pure returns (string) {
        uint length = b.length;
        for (uint i = 0; i < b.length; i++) {
            if (b[b.length - 1 - i] != "") {
                length -= i;
                break;
            }
        }
        bytes memory bs = new bytes(length);
        for (uint j = 0; j < length; j++) {
            bs[j] = b[j];
        }
        return string(bs);
    }
```

 * Modify var: listOfPublicKey => publicKeys
```
    function _existPublicKey(address publicKey) internal view returns (bool) {
        address[] memory publicKeys = platform.publicKeys;
        for (uint i = 0; i < publicKeys.length; i++) {
            if (publicKeys[i] == publicKey) {
                return true;
            }
        }
        return false;
    }
```

 * Modify var: listOfPublicKey => publicKeys
```
    function addPublicKey(address publicKey) onlyAdmin nonzeroAddress(publicKey) external {
        address[] storage publicKeys = platform.publicKeys;
        for (uint i; i < publicKeys.length; i++) {
            if (publicKey == publicKeys[i]) {
                return;
            }
        }
        publicKeys.push(publicKey);
    }
```

 * The array overbounds problem was fixed
 * Optimize for loop structure, remove extra var
```
    function deletePublicKey(address publicKey) onlyAdmin nonzeroAddress(publicKey) external {
        address[] storage publicKeys = platform.publicKeys;
        for (uint i = 0; i < publicKeys.length; i++) {
            if (publicKeys[i] == publicKey) {
                if (i != publicKeys.length - 1 ) {
                    publicKeys[i] = publicKeys[publicKeys.length - 1];
                }
                publicKeys.length--;
                return;
            }
        }
    }
```

 * Modify var: listOfPublicKey => publicKeys
```
    function _existPublicKey(address publicKey) internal view returns (bool) {
        address[] memory publicKeys = platform.publicKeys;
        for (uint i = 0; i < publicKeys.length; i++) {
            if (publicKeys[i] == publicKey) {
                return true;
            }
        }
        return false;
    }
```

 * Optimize for loop structure, remove extra var
```
    function changeVoters(address publicKey, string txid) internal {
        address[] storage voters = platform.proposals[txid].voters;
        for (uint i = 0; i < voters.length; i++) {
            if (voters[i] == publicKey) {
                return;
            }
        }
        voters.push(publicKey);
    }
```