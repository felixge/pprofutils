package legacy

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/pprof/profile"
)

var (
	jemallocHeapHeaderRE = regexp.MustCompile(`^heap_v2/*(\d+)$`)
	jemallocThreadRE     = regexp.MustCompile(`t(\d+): (-?\d+): *(-?\d+) *\[ *(\d+): *(\d+) *\]`)
	jemallocBacktraceRE  = regexp.MustCompile(`@([ x0-9a-f]*)`)
	hexNumberRE          = regexp.MustCompile(`0x[0-9a-f]+`)

	errUnrecognized = fmt.Errorf("unrecognized profile format")
)

type jemallocThreadValue struct {
	value     []int64
	blocksize int64
}

// Jemalloc converts from jemalloc format to protobuf format.
type Jemalloc struct{}

// Convert parses the given text and returns it as protobuf profile.
func (c Jemalloc) Convert(text io.Reader) (*profile.Profile, error) {
	var (
		err       error
		addrs     []uint64
		locs      = make(map[uint64]*profile.Location)
		threadMap = map[string]jemallocThreadValue{}

		p = &profile.Profile{
			SampleType: []*profile.ValueType{
				{Type: "objects", Unit: "count"},
				{Type: "space", Unit: "bytes"},
			},
			PeriodType: &profile.ValueType{Type: "space", Unit: "bytes"},
		}
	)

	s := bufio.NewScanner(text)
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}
		return nil, errUnrecognized
	}

	line := s.Text()
	if header := jemallocHeapHeaderRE.FindStringSubmatch(line); header == nil {
		return nil, errUnrecognized
	}
	p.Period, err = parseJemallocHeapHeader(line)
	if err != nil {
		return nil, err
	}

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if isSpaceOrComment(line) {
			continue
		}

		if isMemoryMapSentinel(line) {
			break
		}

		if sampleData := jemallocThreadRE.FindStringSubmatch(line); sampleData != nil {
			if len(addrs) == 0 {
				// we have not seen the backtrace for this set of allocations yet.
				// This is likely global aggregation, we can ignore it
				continue
			}
			thread, value, blocksize, err := parseJemallocThreadLine(line, p.Period)
			if err != nil {
				return nil, err
			}
			threadMap[thread] = jemallocThreadValue{
				value:     value,
				blocksize: blocksize,
			}
		}

		if sampleData := jemallocBacktraceRE.FindStringSubmatch(line); sampleData != nil {
			// we have next set of addresses, finish the previous one
			if len(addrs) > 0 && len(threadMap) > 0 {
				// complete sample

				var sloc []*profile.Location
				for _, addr := range addrs {
					// Addresses from stack traces point to the next instruction after
					// each call. Adjust by -1 to land somewhere on the actual call.
					addr--
					loc := locs[addr]
					if locs[addr] == nil {
						loc = &profile.Location{
							Address: addr,
						}
						p.Location = append(p.Location, loc)
						locs[addr] = loc
					}
					sloc = append(sloc, loc)
				}
				// iterate the thread map sorted by tid
				tids := make([]string, 0, len(threadMap))
				for tid := range threadMap {
					tids = append(tids, tid)
				}
				sort.Slice(tids, func(i, j int) bool { return tids[i] < tids[j] })
				for _, tid := range tids {
					tv := threadMap[tid]
					p.Sample = append(p.Sample, &profile.Sample{
						Value:    tv.value,
						Location: sloc,
						Label: map[string][]string{
							"thread": {tid},
						},
						NumLabel: map[string][]int64{
							"bytes": {tv.blocksize},
						},
					})
				}

				// reset
				addrs = nil
				threadMap = map[string]jemallocThreadValue{}
			}

			addrs, err = parseHexAddresses(sampleData[1])
			if err != nil {
				return nil, err
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if err := parseAdditionalSections(s, p); err != nil {
		return nil, err
	}
	return p, p.CheckValid()
}

func parseJemallocHeapHeader(line string) (period int64, err error) {
	header := jemallocHeapHeaderRE.FindStringSubmatch(line)
	if header == nil {
		return 0, errUnrecognized
	}

	if len(header[1]) > 0 {
		if period, err = strconv.ParseInt(header[1], 10, 64); err != nil {
			return 0, errUnrecognized
		}
	}

	return period, nil
}

func isSpaceOrComment(line string) bool {
	trimmed := strings.TrimSpace(line)
	return len(trimmed) == 0 || trimmed[0] == '#'
}

var memoryMapSentinels = []string{
	"--- Memory map: ---",
	"MAPPED_LIBRARIES:",
}

// isMemoryMapSentinel returns true if the string contains one of the
// known sentinels for memory map information.
func isMemoryMapSentinel(line string) bool {
	for _, s := range memoryMapSentinels {
		if strings.Contains(line, s) {
			return true
		}
	}
	return false
}

func parseJemallocThreadLine(line string, rate int64) (thread string, value []int64, blocksize int64, err error) {
	sampleData := jemallocThreadRE.FindStringSubmatch(line)
	if len(sampleData) != 6 {
		return "", nil, 0, fmt.Errorf("unexpected number of sample values: got %d, want 6", len(sampleData))
	}

	thread = sampleData[1]

	// This is a local-scoped helper function to avoid needing to pass
	// around rate, sampling and many return parameters.
	addValues := func(countString, sizeString string, label string) error {
		count, err := strconv.ParseInt(countString, 10, 64)
		if err != nil {
			return fmt.Errorf("malformed sample: %s: %v", line, err)
		}
		size, err := strconv.ParseInt(sizeString, 10, 64)
		if err != nil {
			return fmt.Errorf("malformed sample: %s: %v", line, err)
		}
		if count == 0 && size != 0 {
			return fmt.Errorf("%s count was 0 but %s bytes was %d", label, label, size)
		}
		if count != 0 {
			blocksize = size / count
			count, size = scaleHeapSample(count, size, rate)
		}
		value = append(value, count, size)
		return nil
	}

	if err := addValues(sampleData[2], sampleData[3], "inuse"); err != nil {
		return "", nil, 0, fmt.Errorf("malformed sample: %s: %v", line, err)
	}

	// we don't have addrs here yet
	return thread, value, blocksize, nil
}

// scaleHeapSample adjusts the data from a heapz Sample to
// account for its probability of appearing in the collected
// data. heapz profiles are a sampling of the memory allocations
// requests in a program. We estimate the unsampled value by dividing
// each collected sample by its probability of appearing in the
// profile. heapz v2 profiles rely on a poisson process to determine
// which samples to collect, based on the desired average collection
// rate R. The probability of a sample of size S to appear in that
// profile is 1-exp(-S/R).
func scaleHeapSample(count, size, rate int64) (int64, int64) {
	if count == 0 || size == 0 {
		return 0, 0
	}

	if rate <= 1 {
		// if rate==1 all samples were collected so no adjustment is needed.
		// if rate<1 treat as unknown and skip scaling.
		return count, size
	}

	avgSize := float64(size) / float64(count)
	scale := 1 / (1 - math.Exp(-avgSize/float64(rate)))

	return int64(float64(count) * scale), int64(float64(size) * scale)
}

// parseHexAddresses extracts hex numbers from a string, attempts to convert
// each to an unsigned 64-bit number and returns the resulting numbers as a
// slice, or an error if the string contains hex numbers which are too large to
// handle (which means a malformed profile).
func parseHexAddresses(s string) ([]uint64, error) {
	hexStrings := hexNumberRE.FindAllString(s, -1)
	var addrs []uint64
	for _, s := range hexStrings {
		if addr, err := strconv.ParseUint(s, 0, 64); err == nil {
			addrs = append(addrs, addr)
		} else {
			return nil, fmt.Errorf("failed to parse as hex 64-bit number: %s", s)
		}
	}
	return addrs, nil
}

// parseAdditionalSections parses any additional sections in the
// profile, ignoring any unrecognized sections.
func parseAdditionalSections(s *bufio.Scanner, p *profile.Profile) error {
	for !isMemoryMapSentinel(s.Text()) && s.Scan() {
	}
	if err := s.Err(); err != nil {
		return err
	}
	return p.ParseMemoryMapFromScanner(s)
}
