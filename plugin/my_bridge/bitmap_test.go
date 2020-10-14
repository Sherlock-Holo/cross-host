package my_bridge

import (
	"testing"
)

func TestAcquire(t *testing.T) {
	var bitmap [8192]byte

	for i := 0; i < 65535; i++ {
		acquire, ok := Acquire(&bitmap)
		if !ok {
			t.Fatal("acquire failed")
		}

		if acquire != i {
			t.Fatalf("expected acquired %d, got %d", i, acquire)
		}
	}
}

func TestRelease(t *testing.T) {
	var bitmap [8192]byte

	for i := 0; i < 100; i++ {
		acquire, ok := Acquire(&bitmap)
		if !ok {
			t.Fatal("acquire failed")
		}

		if acquire != i {
			t.Fatalf("expected acquired %d, got %d", i, acquire)
		}
	}

	if !Release(&bitmap, 50) {
		t.Fatal("release index 50 failed")
	}

	if Release(&bitmap, 50) {
		t.Fatal("release index 50 twice")
	}

	acquire, ok := Acquire(&bitmap)
	if !ok {
		t.Fatal("acquire failed")
	}

	if acquire != 50 {
		t.Fatalf("expected acquired %d, got %d", 50, acquire)
	}
}
