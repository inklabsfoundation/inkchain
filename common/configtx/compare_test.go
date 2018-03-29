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

package configtx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cb "github.com/inklabsfoundation/inkchain/protos/common"
)

func TestCompareConfigValue(t *testing.T) {
	// Normal equality
	assert.True(t, comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}.equals(comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}), "Should have found identical config values to be identical")

	// Different Mod Policy
	assert.False(t, comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}.equals(comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "bar",
			Value:     []byte("bar"),
		}}), "Should have detected different mod policy")

	// Different Value
	assert.False(t, comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}.equals(comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("foo"),
		}}), "Should have detected different value")

	// Different Version
	assert.False(t, comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}.equals(comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   1,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}), "Should have detected different version")

	// One nil value
	assert.False(t, comparable{
		ConfigValue: &cb.ConfigValue{
			Version:   0,
			ModPolicy: "foo",
			Value:     []byte("bar"),
		}}.equals(comparable{}), "Should have detected nil other value")

}

func TestCompareConfigPolicy(t *testing.T) {
	// Normal equality
	assert.True(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}), "Should have found identical config policies to be identical")

	// Different mod policy
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "bar",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}), "Should have detected different mod policy")

	// Different version
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   1,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}), "Should have detected different version")

	// Different policy type
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  2,
				Value: []byte("foo"),
			},
		}}), "Should have detected different policy type")

	// Different policy value
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("bar"),
			},
		}}), "Should have detected different policy value")

	// One nil value
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{}), "Should have detected one nil value")

	// One nil policy
	assert.False(t, comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type:  1,
				Value: []byte("foo"),
			},
		}}.equals(comparable{
		ConfigPolicy: &cb.ConfigPolicy{
			Version:   0,
			ModPolicy: "foo",
			Policy: &cb.Policy{
				Type: 1,
			},
		}}), "Should have detected one nil policy")
}

func TestCompareConfigGroup(t *testing.T) {
	// Normal equality
	assert.True(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}), "Should have found identical config groups to be identical")

	// Different mod policy
	assert.False(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "bar",
		}}), "Should have detected different mod policy")

	// Different version
	assert.False(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   1,
			ModPolicy: "foo",
		}}), "Should have detected different version")

	// Different groups
	assert.False(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}), "Should have detected different groups entries")

	// Different values
	assert.False(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}), "Should have detected fifferent values entries")

	// Different policies
	assert.False(t, comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar3": nil},
		}}.equals(comparable{
		ConfigGroup: &cb.ConfigGroup{
			Version:   0,
			ModPolicy: "foo",
			Groups:    map[string]*cb.ConfigGroup{"Foo1": nil, "Bar1": nil},
			Values:    map[string]*cb.ConfigValue{"Foo2": nil, "Bar2": nil},
			Policies:  map[string]*cb.ConfigPolicy{"Foo3": nil, "Bar4": nil},
		}}), "Should have detected fifferent policies entries")
}
