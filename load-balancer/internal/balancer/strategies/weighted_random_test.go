package strategies

import (
	"math/rand"
	"testing"
	"time"
)

func TestWeightedRandom(t *testing.T) {
	rand.Seed(1) // фиксируем seed для воспроизводимости теста

	backends := []*Backend{
		{
			Host:  "server1",
			Alive: true,
			Metrics: Metrics{
				CPU: 0.2, MemUsage: 0.3,
				Latency: 20 * time.Millisecond, ErrorRate: 0.01,
			},
		},
		{
			Host:  "server2",
			Alive: true,
			Metrics: Metrics{
				CPU: 0.8, MemUsage: 0.9,
				Latency: 100 * time.Millisecond, ErrorRate: 0.05,
			},
		},
		{
			Host:  "server3",
			Alive: true,
			Metrics: Metrics{
				CPU: 0.5, MemUsage: 0.2,
				Latency: 50 * time.Millisecond, ErrorRate: 0,
			},
		},
	}

	// Предположим, хотим подчеркнуть важность latency (alphaLatency=1.0), а CPU чуть меньше (alphaCPU=0.5)
	wr := NewWeightedRandom(backends,
		0.5, // alphaCPU
		0.2, // alphaMem
		1.0, // alphaLatency
		2.0, // alphaErrorRate
	)

	// Пересчитаем weights (обычно уже вызывается в конструкторе)
	wr.RecalcScoresAndWeights()

	// Проверим, что веса не нулевые для Alive
	for _, b := range backends {
		if !b.Alive {
			continue
		}
		b.mux.RLock()
		if b.Weight <= 0 {
			t.Errorf("Backend %s has non-positive weight: %.4f (score=%.4f)",
				b.Host, b.Weight, b.Score)
		}
		b.mux.RUnlock()
	}

	// Cделаем несколько выборов и посмотрим, распределяется ли чаще на "лучшем" бэкенде
	chosenCount := map[string]int{
		"server1": 0,
		"server2": 0,
		"server3": 0,
	}
	total := 10000
	for i := 0; i < total; i++ {
		host, err := wr.GetBackend()
		if err != nil {
			t.Fatalf("GetBackend error: %v", err)
		}
		chosenCount[host]++
	}

	// Выведем статистику
	t.Logf("Chosen results (of %d): %v", total, chosenCount)

	// В идеале, "наименее загруженный" (по score) должен иметь более высокую вероятность
	// Допустим, сервер1: (CPU=0.2, Mem=0.3, Lat=20ms, Err=0.01) — ожидается, что score будет меньше server2.
	// Проверим, что server2 действительно выбирается реже (так как score больше => weight меньше).
	if chosenCount["server2"] > chosenCount["server1"] {
		t.Errorf("Expected server2 to be chosen less often than server1, but got server2=%d, server1=%d",
			chosenCount["server2"], chosenCount["server1"])
	}
}
