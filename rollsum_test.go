// rollsum_test.go, tests the rollsum.go library
// Copyright (C) 2013  Greg Baker
// GPLv2, see LICENSE.txt

package rollsum

import "testing"
import "hash/adler32"
import "crypto/rand"


// Tests by writing up to 4k of data into the checksum
func TestEquivalence(t *testing.T) {
	testarr := make([]byte, 4*1024);
	sum := 0
	for sum != len(testarr) {
		num, _ := rand.Read(testarr[sum:])
		sum += num
	}
	for i := 0; i <= len(testarr); i++ {
		 var rolling = New(255);
		 rolling.Write(testarr[:i])
		 start := 0;
		 if i > 255 {
			 start = i-255;
		 }

		expected := adler32.Checksum(testarr[start:i]);
		actual := rolling.Sum32();
		if expected != actual {
			t.Fatalf("%d: expected %x, got %x, %x -> %x",
				i, expected, actual,
				testarr[start - 1], testarr[i - 1]);
		}
	}
}

// Writes up to 4MB of random data into the rolling checksum
func TestEquivalenceLong(t *testing.T) {
	testarr := make([]byte, 4*1024*1024)
	sum := 0
	for sum != len(testarr) {
		num, _ := rand.Read(testarr[sum:])
		sum += num
	}

	for window := 32; window <= 40960; window*=2 {
		var rolling = New(uint32(window));
		rolling.Write(testarr);
		start := len(testarr) - window
		expected := adler32.Checksum(testarr[start:]);
		actual := rolling.Sum32();
		if expected != actual {
			t.Fatalf("%d: expected %x, got %x",
				window, expected, actual);
		}
	}
}

func TestEmpty(t *testing.T) {
	rolling := New(32);
	expected := uint32(1);
	actual := rolling.Sum32();
	if expected != actual {
		t.Fail();
	}
}

// Test from the Wikipedia page on Adler32
func TestWikipedia(t *testing.T) {
	expected := uint32(0x11E60398);
	teststr := "Wikipedia";
	rolling := New(32);
	rolling.Write([]byte(teststr));
	if expected != rolling.Sum32() {
		t.Fatalf("expected %x, got %x", expected, rolling.Sum32());
	}
}

