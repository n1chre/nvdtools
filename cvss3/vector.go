// Copyright (c) Facebook, Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cvss3

import (
	"fmt"
	"strings"
)

const (
	prefix          = "CVSS:3.0/"
	partSeparator   = "/"
	metricSeparator = ":"
)

// Vector represents a CVSSv3 vector, holds all metrics inside (base, temporal and environmental)
type Vector struct {
	BaseMetrics
	TemporalMetrics
	EnvironmentalMetrics
}

// String returns this vectors representation as a string
// it shouldn't depend on the order of metrics
func (v Vector) String() string {
	var sb strings.Builder
	fmt.Fprint(&sb, prefix)

	defineables := v.definables()

	first := true
	for _, metric := range order {
		def := defineables[metric]
		if !def.defined() {
			continue
		}
		if !first {
			fmt.Fprint(&sb, partSeparator)
		} else {
			first = false
		}
		fmt.Fprintf(&sb, "%s%s%s", metric, metricSeparator, def)
	}

	return sb.String()
}

// VectorFromString will parse a string into a Vector, or return an error if it can't be parsed
func VectorFromString(str string) (Vector, error) {
	// remove prefix if exists
	str = strings.ToUpper(str)
	if strings.HasPrefix(str, prefix) {
		str = str[len(prefix):]
	}

	var v Vector
	parseables := v.parseables()

	for _, part := range strings.Split(str, partSeparator) {
		tmp := strings.Split(part, metricSeparator)
		if len(tmp) != 2 {
			return v, fmt.Errorf("need two values separated by %s, got %q", metricSeparator, part)
		}

		metric, value := tmp[0], tmp[1]
		if p, ok := parseables[metric]; !ok {
			return v, fmt.Errorf("undefined metric %s with value %s", metric, value)
		} else if err := p.parse(value); err != nil {
			return v, fmt.Errorf("error occurred while parsing metric %s: %v", metric, err)
		}
	}

	return v, nil
}

// Absorb will override only metrics in the current vector from the one given which are defined
// If the other vector specifies only a single metric with all others undefined, the resulting
// vector will contain all metrics it previously did, with only the new one overriden
func (v *Vector) Absorb(other Vector) {
	parseables := v.parseables()
	for metric, defineable := range other.definables() {
		if defineable.defined() {
			parseables[metric].parse(defineable.String())
		}
	}
}

// helpers

var order = []string{"AV", "AC", "PR", "UI", "S", "C", "I", "A", "E", "RL", "RC", "CR", "IR", "AR", "MAV", "MAC", "MPR", "MUI", "MS", "MC", "MI", "MA"}

type defineable interface {
	defined() bool
	String() string
}

type parseable interface {
	parse(string) error
}

func (v *Vector) definables() map[string]defineable {
	return map[string]defineable{
		// base metrics
		"AV": v.BaseMetrics.AttackVector,
		"AC": v.BaseMetrics.AttackComplexity,
		"PR": v.BaseMetrics.PrivilegesRequired,
		"UI": v.BaseMetrics.UserInteraction,
		"S":  v.BaseMetrics.Scope,
		"C":  v.BaseMetrics.Confidentiality,
		"I":  v.BaseMetrics.Integrity,
		"A":  v.BaseMetrics.Availability,
		// temporal metrics
		"E":  v.TemporalMetrics.ExploitCodeMaturity,
		"RL": v.TemporalMetrics.RemediationLevel,
		"RC": v.TemporalMetrics.ReportConfidence,
		// environmental metrics
		"CR":  v.EnvironmentalMetrics.ConfidentialityRequirement,
		"IR":  v.EnvironmentalMetrics.IntegrityRequirement,
		"AR":  v.EnvironmentalMetrics.AvailabilityRequirement,
		"MAV": v.EnvironmentalMetrics.ModifiedAttackVector,
		"MAC": v.EnvironmentalMetrics.ModifiedAttackComplexity,
		"MPR": v.EnvironmentalMetrics.ModifiedPrivilegesRequired,
		"MUI": v.EnvironmentalMetrics.ModifiedUserInteraction,
		"MS":  v.EnvironmentalMetrics.ModifiedScope,
		"MC":  v.EnvironmentalMetrics.ModifiedConfidentiality,
		"MI":  v.EnvironmentalMetrics.ModifiedIntegrity,
		"MA":  v.EnvironmentalMetrics.ModifiedAvailability,
	}
}

func (v *Vector) parseables() map[string]parseable {
	return map[string]parseable{
		// base metrics
		"AV": &v.BaseMetrics.AttackVector,
		"AC": &v.BaseMetrics.AttackComplexity,
		"PR": &v.BaseMetrics.PrivilegesRequired,
		"UI": &v.BaseMetrics.UserInteraction,
		"S":  &v.BaseMetrics.Scope,
		"C":  &v.BaseMetrics.Confidentiality,
		"I":  &v.BaseMetrics.Integrity,
		"A":  &v.BaseMetrics.Availability,
		// temporal metrics
		"E":  &v.TemporalMetrics.ExploitCodeMaturity,
		"RL": &v.TemporalMetrics.RemediationLevel,
		"RC": &v.TemporalMetrics.ReportConfidence,
		// environmental metrics
		"CR":  &v.EnvironmentalMetrics.ConfidentialityRequirement,
		"IR":  &v.EnvironmentalMetrics.IntegrityRequirement,
		"AR":  &v.EnvironmentalMetrics.AvailabilityRequirement,
		"MAV": &v.EnvironmentalMetrics.ModifiedAttackVector,
		"MAC": &v.EnvironmentalMetrics.ModifiedAttackComplexity,
		"MPR": &v.EnvironmentalMetrics.ModifiedPrivilegesRequired,
		"MUI": &v.EnvironmentalMetrics.ModifiedUserInteraction,
		"MS":  &v.EnvironmentalMetrics.ModifiedScope,
		"MC":  &v.EnvironmentalMetrics.ModifiedConfidentiality,
		"MI":  &v.EnvironmentalMetrics.ModifiedIntegrity,
		"MA":  &v.EnvironmentalMetrics.ModifiedAvailability,
	}
}
