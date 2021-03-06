package jwt

import (
	"context"
	"testing"
	"time"
)

func TestBlocklist(t *testing.T) {
	b := NewBlocklist(0)
	c := Map{"username": "kataras", "age": 27}
	sc := Claims{Expiry: Clock().Add(2 * time.Minute).Unix()}
	token, err := Sign(testAlg, testSecret, Merge(c, sc))
	if err != nil {
		t.Fatal(err)
	}

	b.InvalidateToken(token, sc.Expiry)
	if !b.Has(token) {
		t.Fatalf("expected token to be in the list")
	}

	if b.Count() != 1 {
		t.Fatalf("expected list to contain a single token entry")
	}

	if err = b.ValidateToken(token, Claims{}, nil); err != ErrBlocked {
		t.Fatalf("expected error: ErrBlock but got: %v", err)
	}

	if err = b.ValidateToken(token, Claims{}, ErrExpired); err != ErrExpired {
		t.Fatalf("expected error: ErrExpired as it respects the previous one but got: %v", err)
	}

	if b.Has(token) {
		t.Fatalf("expected token to be removed as the validate token's error was ErrExpired")
	}

	b.InvalidateToken(token, sc.Expiry)
	if removed := b.GC(); removed != 0 {
		t.Fatalf("expected nothing to be removed because the expiration is before current time but got: %d", removed)
	}

	b.Del(token)

	if count := b.Count(); count != 0 {
		t.Fatalf("expected count to be zero but got: %d", count)
	}

	if err = b.ValidateToken(token, Claims{}, nil); err != nil {
		t.Fatalf("expected no error as this token is now not blocked")
	}

	b.InvalidateToken([]byte{}, 1)
	if got := b.Count(); got != 0 {
		t.Fatalf("expected zero entries as the token was empty but got: %d", got)
	}

	if b.Has([]byte{}) {
		t.Fatalf("expected Has to always return false as the given token was empty")
	}

	// Test GC expired.
	b.InvalidateToken([]byte("expired one"), 1)
	if got := b.Count(); got != 1 {
		t.Fatalf("expected upsert not append")
	}
	if removed := b.GC(); removed != 1 {
		t.Fatalf("expected one token to be removed as it's expired")
	}

	// test automatic gc
	ctx, cancel := context.WithCancel(context.Background())
	b = NewBlocklistContext(ctx, 500*time.Millisecond)
	for i := 0; i < 10; i++ {
		b.InvalidateToken(MustGenerateRandom(92), Clock().Add(time.Second).Unix())
	}
	time.Sleep(2 * time.Second)
	cancel()

	if got := b.Count(); got != 0 {
		t.Fatalf("expected all entries to be removed but: %d", got)
	}
}
