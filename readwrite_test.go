package proto9

import (
	"bytes"
	"reflect"
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
			niters := 10

			for j := 0; j < niters; j++ {
				fuzzer.Fuzz(fc)
				var buf bytes.Buffer
				msize := uint32(2 * 1024 * 1024)
				err := WriteFcall(fc, msize, &buf)
				if err != nil {
					continue
				}

				if uint64(buf.Len()-5) != fc.EncodedSize() {
					t.Fatalf("EncodedSize %d did not match actual size %d for %#v", fc.EncodedSize(), buf.Len()-5, fc)
				}

				fc2, err := ReadFcall(msize, &buf)
				if err != nil {
					t.Fatalf("encoding %v failed with error %s", fc, err)
				}
				if !reflect.DeepEqual(fc, fc2) {
					t.Fatalf("%#v\n should equal\n %#v", fc2, fc)
				}
			}
		}(i)
	}

}
