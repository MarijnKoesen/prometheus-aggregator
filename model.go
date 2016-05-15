package main

import (
	"crypto/md5"
	"sort"
)

type sampleKind string

const (
	sampleUnknown sampleKind = ""

	// sampleCounter represents a counter
	sampleCounter sampleKind = "c"

	// sampleGauge represents a gauge
	sampleGauge sampleKind = "g"

	// sampleHistogramLinear represents histogram with linearly spaced buckets.
	// See Prometheus Go client LinearBuckets for details.
	sampleHistogramLinear sampleKind = "hl"
)

// sample represents single sample
type sample struct {
	// name is used to represent sample. It's used as metric name in export to prometheus.
	name string

	// kind of the sample wen mapped to prometheus metric type
	kind sampleKind

	// labels is a set of string pairs mapped to prometheus LabelPairs type
	labels map[string]string

	// value of the sample
	value float64

	// histogramDef is a set of values used in mapping for the histogram types
	histogramDef []string
}

// hash calculates a hash of the sample so it can be recognized.
// Should take all elements other than value under consideration.
func (s *sample) hash() []byte {
	// TODO(szpakas): switch to non-cryptographic hash function like FNV or xxHash
	// TODO(szpakas): hash histogramDef
	hash := md5.New()

	hash.Write([]byte(s.kind))
	hash.Write([]byte("|"))

	hash.Write([]byte(s.name))

	// labels
	if len(s.labels) > 0 {
		hash.Write([]byte("|"))

		// get all keys sorted so hash is repeatable
		var keys []string
		for k := range s.labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			hash.Write([]byte(k))
			hash.Write([]byte("="))
			hash.Write([]byte(s.labels[k]))
			// separator between labels
			if i < len(keys)-1 {
				hash.Write([]byte(";"))
			}
		}
	}

	return hash.Sum([]byte{})
}
