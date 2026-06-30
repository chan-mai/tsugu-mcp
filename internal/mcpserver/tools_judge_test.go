package mcpserver

import (
	"context"
	"strings"
	"testing"
)

func TestHandleShares(t *testing.T) {
	res, out, err := handleShares(context.Background(), nil, sharesToolInput{
		SpouseName: "花子",
		Children:   []shareHeir{{Name: "一郎", Alive: true}, {Name: "二郎", Alive: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %+v", res.Content)
	}
	if out.Result.Sum != "1" {
		t.Errorf("sum = %q, want 1", out.Result.Sum)
	}
	frac := map[string]string{}
	for _, s := range out.Result.Shares {
		frac[s.Name] = s.Fraction
	}
	if frac["花子"] != "1/2" || frac["一郎"] != "1/4" || frac["二郎"] != "1/4" {
		t.Errorf("shares = %v", frac)
	}
}

func TestHandleShares_Substitution(t *testing.T) {
	// 子B先死亡、孫D・Eが代襲(2階層スキーマの確認)
	res, out, err := handleShares(context.Background(), nil, sharesToolInput{
		SpouseName: "花子",
		Children: []shareHeir{
			{Name: "A", Alive: true},
			{Name: "B", Substitutes: []shareSub{{Name: "D", Alive: true}, {Name: "E", Alive: true}}},
		},
	})
	if err != nil || res.IsError {
		t.Fatalf("unexpected: %v %+v", err, res.Content)
	}
	frac := map[string]string{}
	for _, s := range out.Result.Shares {
		frac[s.Name] = s.Fraction
	}
	if frac["D"] != "1/8" || frac["E"] != "1/8" || frac["A"] != "1/4" {
		t.Errorf("代襲 shares = %v", frac)
	}
}

func TestHandlePattern(t *testing.T) {
	res, out, err := handlePattern(context.Background(), nil, patternToolInput{Method: "agreement", Renunciation: true})
	if err != nil || res.IsError {
		t.Fatalf("unexpected: %v %+v", err, res.Content)
	}
	if out.Result.Primary.Key != "B" {
		t.Errorf("key = %q, want B", out.Result.Primary.Key)
	}
	if len(out.Result.Modifiers) == 0 || !strings.Contains(strings.Join(out.Result.Modifiers, " "), "相続放棄") {
		t.Errorf("missing renunciation modifier: %v", out.Result.Modifiers)
	}
}

func TestHandleNotify(t *testing.T) {
	res, out, err := handleNotify(context.Background(), nil, notifyToolInput{DeathDate: "2025-06-15"})
	if err != nil || res.IsError {
		t.Fatalf("unexpected: %v %+v", err, res.Content)
	}
	if !strings.HasPrefix(out.Result.Deadline, "2028-06-15") {
		t.Errorf("deadline = %q, want 2028-06-15", out.Result.Deadline)
	}
}

func TestHandleNotify_BadDate(t *testing.T) {
	res, _, _ := handleNotify(context.Background(), nil, notifyToolInput{DeathDate: "not-a-date"})
	if !res.IsError {
		t.Error("malformed date should error")
	}
}
