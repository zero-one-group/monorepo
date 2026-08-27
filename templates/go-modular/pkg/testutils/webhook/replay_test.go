package webhook

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeReceiver is the smallest receiver that satisfies the harness: it dedupes on
// event id, tracks a per-order status with a terminal set, and "verifies" a header.
// It exists only to prove the harness itself behaves; real receivers live in modules/.
type fakeReceiver struct {
	seen     map[string]bool
	status   map[string]string
	terminal map[string]bool
}

type fakeEvent struct {
	ID      string `json:"id"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func newFake() *fakeReceiver {
	return &fakeReceiver{
		seen:     map[string]bool{},
		status:   map[string]string{},
		terminal: map[string]bool{"COMPLETED": true, "CANCELLED": true},
	}
}

func (f *fakeReceiver) Handle(_ context.Context, headers map[string]string, payload []byte) (Result, error) {
	var ev fakeEvent
	if err := json.Unmarshal(payload, &ev); err != nil || headers["x-signature"] != "ok" || len(payload) != len(jsonOf(ev)) {
		return Result{ProviderEventID: ev.ID, Outcome: OutcomeRejected}, nil
	}
	if f.seen[ev.ID] {
		return Result{ProviderEventID: ev.ID, Outcome: OutcomeIgnored}, nil
	}
	f.seen[ev.ID] = true
	if f.terminal[f.status[ev.OrderID]] {
		return Result{ProviderEventID: ev.ID, Outcome: OutcomeIgnored}, nil
	}
	f.status[ev.OrderID] = ev.Status
	return Result{ProviderEventID: ev.ID, Outcome: OutcomeProcessed}, nil
}

func jsonOf(v any) []byte { b, _ := json.Marshal(v); return b }

func fx(name, id, order, status string) Fixture {
	return Fixture{
		Name:    name,
		Headers: map[string]string{"x-signature": "ok"},
		Payload: jsonOf(fakeEvent{ID: id, OrderID: order, Status: status}),
	}
}

func TestHarness_IdempotentReceiverPasses(t *testing.T) {
	AssertIdempotent(t, newFake(), fx("paid", "evt-1", "o-1", "PAID"))
}

func TestHarness_LateEventIsIgnored(t *testing.T) {
	r := newFake()
	AssertLateEventIgnored(t, r,
		[]Fixture{
			fx("accepted", "evt-1", "o-1", "PREPARING"),
			fx("completed", "evt-2", "o-1", "COMPLETED"),
		},
		fx("late-accept", "evt-3", "o-1", "PREPARING"),
	)
	require.Equal(t, "COMPLETED", r.status["o-1"])
}

func TestHarness_BadSignatureRejected(t *testing.T) {
	AssertRejectsBadSignature(t, newFake(), fx("paid", "evt-1", "o-1", "PAID"))
}

func TestLoad_ReadsFixtureAndOptionalHeaders(t *testing.T) {
	dir := filepath.Join(moduleRoot(t), "testdata", "webhooks", "_selftest")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	require.NoError(t, os.WriteFile(filepath.Join(dir, "paid.json"), []byte(`{"id":"evt-1"}`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "paid.headers.json"), []byte(`{"x-signature":"ok"}`), 0o644))

	fx := Load(t, "_selftest", "paid")
	require.Equal(t, "paid", fx.Name)
	require.JSONEq(t, `{"id":"evt-1"}`, string(fx.Payload))
	require.Equal(t, "ok", fx.Headers["x-signature"])
}

func TestModuleRootPointsAtApi(t *testing.T) {
	root := moduleRoot(t)
	require.FileExists(t, root+"/go.mod")
}
