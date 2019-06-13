package pg_test

import (
	"bytes"
	"fmt"
	"github.com/hetianyi/gox/pg"
	"math"
	"os"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	a := 'p'

	fmt.Fprintf(os.Stdout, "%c", a)
}

func Test11(t *testing.T) {
	var buffer bytes.Buffer
	finish := int(math.Floor(float64(15) / float64(100) * float64(50-2)))
	for i := 0; i < finish; i++ {
		buffer.WriteRune('-')
	}
	fs := buffer.String()
	buffer.Reset()
	for i := 0; i < 50-2-finish-1; i++ {
		buffer.WriteRune(' ')
	}
	bs := buffer.String()
	fmt.Fprintf(os.Stdout, "\r%c%s%c%s%c", '[', fs, '>', bs, ']')
}

func Test2(t *testing.T) {
	fmt.Println()
	p := pg.New(100, 50)
	for i := 0; i < 100; i++ {
		p.Increase()
		time.Sleep(time.Millisecond * 100)
	}
}
