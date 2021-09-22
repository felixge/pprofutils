package utils

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/pprof/profile"
)

type Avg struct {
	Input  []byte
	Output io.Writer
}

func (a *Avg) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(a.Input)
	if err != nil {
		return err
	}

	var (
		countIDX = -1
		delayIDX = -1
	)
	for i, st := range prof.SampleType {
		if st.Type == "contentions" && st.Unit == "count" {
			countIDX = i
		}
		if st.Type == "delay" && st.Unit == "nanoseconds" {
			delayIDX = i
		}
	}
	if countIDX == -1 {
		return errors.New("profile lacks a contentions/count sample type")
	}
	if delayIDX == -1 {
		return errors.New("profile lacks a delay/nanoseconds sample type")
	}

	for i, s := range prof.Sample {
		if countIDX >= len(s.Value) {
			return fmt.Errorf("sample %d has no contentions/count value", i)
		}
		if delayIDX >= len(s.Value) {
			return fmt.Errorf("sample %d has no delay/nanoseconds value", i)
		}
		count, delay := s.Value[countIDX], s.Value[delayIDX]
		s.Value[delayIDX] = delay / count
	}

	return prof.Write(a.Output)
}
