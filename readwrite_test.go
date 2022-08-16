package proto9

import (
	"bytes"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/google/gofuzz"
)

func TestReadWrite(t *testing.T) {
	// Property test random values in parallel, just ensure round trip.
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i <= 0xff; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fc, err := FcallFromKind(byte(i))
			if err != nil {
				return
			}

			fuzzer := fuzz.New().NilChance(0.0).NumElements(0, 32)
			niters := runtime.NumCPU() * 3

			for j := 0; j < niters; j++ {
				fuzzer.Fuzz(fc)
				var buf bytes.Buffer
				msize := uint32(2 * 1024 * 1024)
				err := WriteFcall(fc, msize, &buf)
				if err != nil {
					continue
				}
				fc2, err := ReadFcall(msize, &buf)
				if err != nil {
					t.Fatalf("encoding %v failed with error %s", fc, err)
				}
				if !reflect.DeepEqual(fc, fc2) {
					t.Fatalf("%#v\n should equal\n %#v", fc, fc2)
				}
			}
		}(i)
	}

}