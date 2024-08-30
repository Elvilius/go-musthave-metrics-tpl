package models

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync/atomic"

	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/shirou/gopsutil/v4/mem"
)

var memStats runtime.MemStats

// Create new Metrics
func NewAllocMetric() *AllocMetric {
	return &AllocMetric{metric: &Metrics{ID: "Alloc", MType: Gauge}}
}

func NewBuckHashSysMetric() *BuckHashSysMetric {
	return &BuckHashSysMetric{metric: &Metrics{ID: "BuckHashSys", MType: Gauge}}
}

func NewFreesMetric() *FreesMetric {
	return &FreesMetric{metric: &Metrics{ID: "Frees", MType: Gauge}}
}

func NewGCCPUFractionMetric() *GCCPUFractionMetric {
	return &GCCPUFractionMetric{metric: &Metrics{ID: "GCCPUFraction", MType: Gauge}}
}

func NewGCSysMetric() *GCSysMetric {
	return &GCSysMetric{metric: &Metrics{ID: "GCSys", MType: Gauge}}
}

func NewHeapAllocMetric() *HeapAllocMetric {
	return &HeapAllocMetric{metric: &Metrics{ID: "HeapAlloc", MType: Gauge}}
}

func NewHeapIdleMetric() *HeapIdleMetric {
	return &HeapIdleMetric{metric: &Metrics{ID: "HeapIdle", MType: Gauge}}
}

func NewHeapInuseMetric() *HeapInuseMetric {
	return &HeapInuseMetric{metric: &Metrics{ID: "HeapInuse", MType: Gauge}}
}

func NewMCacheSys() *MCacheSysMetric {
	return &MCacheSysMetric{metric: &Metrics{ID: "MCacheSys", MType: Gauge}}
}

func NewHeapObjectsMetric() *HeapObjectsMetric {
	return &HeapObjectsMetric{metric: &Metrics{ID: "HeapObjects", MType: Gauge}}
}

func NewHeapReleasedMetric() *HeapReleasedMetric {
	return &HeapReleasedMetric{metric: &Metrics{ID: "HeapReleased", MType: Gauge}}
}

func NewHeapSysMetric() *HeapSysMetric {
	return &HeapSysMetric{metric: &Metrics{ID: "HeapSys", MType: Gauge}}
}

func NewLastGCMetric() *LastGCMetric {
	return &LastGCMetric{metric: &Metrics{ID: "LastGC", MType: Gauge}}
}

func NewLookupsMetric() *LookupsMetric {
	return &LookupsMetric{metric: &Metrics{ID: "Lookups", MType: Gauge}}
}

func NewMCacheInuseMetric() *MCacheInuseMetric {
	return &MCacheInuseMetric{metric: &Metrics{ID: "MCacheInuse", MType: Gauge}}
}

func NewMSpanSysMetric() *MSpanSysMetric {
	return &MSpanSysMetric{metric: &Metrics{ID: "MSpanSys", MType: Gauge}}
}

func NewMallocsMetric() *MallocsMetric {
	return &MallocsMetric{metric: &Metrics{ID: "Mallocs", MType: Gauge}}
}

func NewNextGCMetric() *NextGCMetric {
	return &NextGCMetric{metric: &Metrics{ID: "NextGC", MType: Gauge}}
}

func NewNumForcedGC() *NumForcedGCMetric {
	return &NumForcedGCMetric{metric: &Metrics{ID: "NumForcedGC", MType: Gauge}}
}

func NewNumGCMetric() *NumGCMetric {
	return &NumGCMetric{metric: &Metrics{ID: "NumGC", MType: Gauge}}
}

func NewOtherSysMetric() *OtherSysMetric {
	return &OtherSysMetric{metric: &Metrics{ID: "OtherSys", MType: Gauge}}
}

func NewPauseTotalNsMetric() *PauseTotalNsMetric {
	return &PauseTotalNsMetric{metric: &Metrics{ID: "PauseTotalNs", MType: Gauge}}
}

func NewStackInuseMetric() *StackInuseMetric {
	return &StackInuseMetric{metric: &Metrics{ID: "StackInuse", MType: Gauge}}
}

func NewStackSysMetric() *StackSysMetric {
	return &StackSysMetric{metric: &Metrics{ID: "StackSys", MType: Gauge}}
}

func NewTotalAlloc() *TotalAllocMetric {
	return &TotalAllocMetric{metric: &Metrics{ID: "TotalAlloc", MType: Gauge}}
}

func NewSysMetric() *SysMetric {
	return &SysMetric{metric: &Metrics{ID: "Sys", MType: Gauge}}
}

func NewMSpanInuseMetric() *MSpanInuseMetric {
	return &MSpanInuseMetric{metric: &Metrics{ID: "MSpanInuse", MType: Gauge}}
}

func NewRandomValueMetric() *RandomValueMetric {
	return &RandomValueMetric{metric: &Metrics{ID: "RandomValue", MType: Gauge}}
}

func NewPollCountMetric() *PollCountMetric {
	return &PollCountMetric{metric: &Metrics{ID: "PollCount", MType: Counter}}
}

func NewTotalMemoryMetric() *TotalMemoryMetric {
	return &TotalMemoryMetric{metric: &Metrics{ID: "TotalMemory", MType: Counter}}
}

func NewFreeMemoryMetric() *FreeMemoryMetric {
	return &FreeMemoryMetric{metric: &Metrics{ID: "FreeMemory", MType: Counter}}
}

func NewCPUutilization1Metric() *CPUutilization1Metric {
	return &CPUutilization1Metric{metric: &Metrics{ID: "CPUutilization1Metric", MType: Counter}}
}

// Update Metrics
func (m *AllocMetric) Update() {
	runtime.ReadMemStats(&memStats)

	alloc := float64(memStats.Alloc)
	m.metric = &Metrics{ID: "Alloc", MType: Gauge, Value: &alloc}
}

func (m *BuckHashSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	buckHashSys := float64(memStats.BuckHashSys)
	m.metric = &Metrics{ID: "BuckHashSys", MType: Gauge, Value: &buckHashSys}
}

func (m *FreesMetric) Update() {
	runtime.ReadMemStats(&memStats)

	frees := float64(memStats.Frees)
	m.metric = &Metrics{ID: "Frees", MType: Gauge, Value: &frees}
}

func (m *GCCPUFractionMetric) Update() {
	runtime.ReadMemStats(&memStats)

	gCCPUFraction := memStats.GCCPUFraction
	m.metric = &Metrics{ID: "GCCPUFraction", MType: Gauge, Value: &gCCPUFraction}
}

func (m *GCSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	gCSys := float64(memStats.GCSys)
	m.metric = &Metrics{ID: "GCSys", MType: Gauge, Value: &gCSys}
}

func (m *HeapAllocMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapAlloc := float64(memStats.HeapAlloc)
	m.metric = &Metrics{ID: "HeapAlloc", MType: Gauge, Value: &heapAlloc}
}

func (m *HeapIdleMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapIdle := float64(memStats.HeapIdle)
	m.metric = &Metrics{ID: "HeapIdle", MType: Gauge, Value: &heapIdle}
}

func (m *HeapInuseMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapInuse := float64(memStats.HeapInuse)
	m.metric = &Metrics{ID: "HeapInuse", MType: Gauge, Value: &heapInuse}
}

func (m *HeapObjectsMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapObjects := float64(memStats.HeapObjects)
	m.metric = &Metrics{ID: "HeapObjects", MType: Gauge, Value: &heapObjects}
}

func (m *HeapReleasedMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapReleased := float64(memStats.HeapReleased)
	m.metric = &Metrics{ID: "HeapReleased", MType: Gauge, Value: &heapReleased}
}

func (m *HeapSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	heapSys := float64(memStats.HeapSys)
	m.metric = &Metrics{ID: "HeapSys", MType: Gauge, Value: &heapSys}
}

func (m *LastGCMetric) Update() {
	runtime.ReadMemStats(&memStats)

	lastGC := float64(memStats.LastGC)
	m.metric = &Metrics{ID: "LastGC", MType: Gauge, Value: &lastGC}
}

func (m *LookupsMetric) Update() {
	runtime.ReadMemStats(&memStats)

	lookups := float64(memStats.Lookups)
	m.metric = &Metrics{ID: "Lookups", MType: Gauge, Value: &lookups}
}

func (m *MCacheInuseMetric) Update() {
	runtime.ReadMemStats(&memStats)

	mCacheInuse := float64(memStats.MCacheInuse)
	m.metric = &Metrics{ID: "MCacheInuse", MType: Gauge, Value: &mCacheInuse}
}

func (m *MCacheSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	mCacheSys := float64(memStats.MCacheSys)
	m.metric = &Metrics{ID: "MCacheSys", MType: Gauge, Value: &mCacheSys}
}

func (m *MSpanInuseMetric) Update() {
	runtime.ReadMemStats(&memStats)

	mSpanInuse := float64(memStats.MSpanInuse)
	m.metric = &Metrics{ID: "MSpanInuse", MType: Gauge, Value: &mSpanInuse}
}

func (m *MSpanSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	mSpanSys := float64(memStats.MSpanSys)
	m.metric = &Metrics{ID: "MSpanSys", MType: Gauge, Value: &mSpanSys}
}

func (m *MallocsMetric) Update() {
	runtime.ReadMemStats(&memStats)

	mallocs := float64(memStats.Mallocs)
	m.metric = &Metrics{ID: "Mallocs", MType: Gauge, Value: &mallocs}
}

func (m *NextGCMetric) Update() {
	runtime.ReadMemStats(&memStats)

	nextGC := float64(memStats.NextGC)
	m.metric = &Metrics{ID: "NextGC", MType: Gauge, Value: &nextGC}
}

func (m *NumForcedGCMetric) Update() {
	runtime.ReadMemStats(&memStats)

	numForcedGC := float64(memStats.NumForcedGC)
	m.metric = &Metrics{ID: "NumForcedGC", MType: Gauge, Value: &numForcedGC}
}

func (m *NumGCMetric) Update() {
	runtime.ReadMemStats(&memStats)

	numGC := float64(memStats.NumGC)
	m.metric = &Metrics{ID: "NumGC", MType: Gauge, Value: &numGC}
}

func (m *OtherSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	otherSys := float64(memStats.OtherSys)
	m.metric = &Metrics{ID: "OtherSys", MType: Gauge, Value: &otherSys}
}

func (m *PauseTotalNsMetric) Update() {
	runtime.ReadMemStats(&memStats)

	pauseTotalNs := float64(memStats.PauseTotalNs)
	m.metric = &Metrics{ID: "PauseTotalNs", MType: Gauge, Value: &pauseTotalNs}
}

func (m *StackInuseMetric) Update() {
	runtime.ReadMemStats(&memStats)

	stackInuse := float64(memStats.StackInuse)
	m.metric = &Metrics{ID: "StackInuse", MType: Gauge, Value: &stackInuse}
}

func (m *StackSysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	stackSys := float64(memStats.StackSys)
	m.metric = &Metrics{ID: "StackSys", MType: Gauge, Value: &stackSys}
}

func (m *SysMetric) Update() {
	runtime.ReadMemStats(&memStats)

	sys := float64(memStats.Sys)
	m.metric = &Metrics{ID: "Sys", MType: Gauge, Value: &sys}
}

func (m *TotalAllocMetric) Update() {
	runtime.ReadMemStats(&memStats)

	totalAlloc := float64(memStats.TotalAlloc)
	m.metric = &Metrics{ID: "TotalAlloc", MType: Gauge, Value: &totalAlloc}
}

func (m *RandomValueMetric) Update() {
	randomValue := rand.Float64()
	m.metric = &Metrics{ID: "RandomValue", MType: Gauge, Value: &randomValue}
}

func (m *PollCountMetric) Update() {
	fmt.Println(m.count)
	atomic.AddInt64(&m.count, 1)

	m.metric = &Metrics{ID: "PollCount", MType: Counter, Delta: &m.count}
}

func (m *TotalMemoryMetric) Update() {
	v, _ := mem.VirtualMemory()
	totalMemory := float64(v.Total)
	m.metric = &Metrics{ID: "TotalMemory", MType: Gauge, Value: &totalMemory}
}

func (m *FreeMemoryMetric) Update() {
	v, _ := mem.VirtualMemory()
	free := float64(v.Free)
	m.metric = &Metrics{ID: "FreeMemory", MType: Gauge, Value: &free}
}

func (m *CPUutilization1Metric) Update() {
	v, _ := cpu.Percent(0, true)
	percentages := v[1]
	m.metric = &Metrics{ID: "CPUutilization1", MType: Gauge, Value: &percentages}
}

// Get metrics
func (m *AllocMetric) Get() Metrics {
	return *m.metric
}

func (m *BuckHashSysMetric) Get() Metrics {
	return *m.metric
}

func (m *FreesMetric) Get() Metrics {
	return *m.metric
}

func (m *GCCPUFractionMetric) Get() Metrics {
	return *m.metric
}

func (m *GCSysMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapAllocMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapIdleMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapInuseMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapObjectsMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapReleasedMetric) Get() Metrics {
	return *m.metric
}

func (m *HeapSysMetric) Get() Metrics {
	return *m.metric
}

func (m *LastGCMetric) Get() Metrics {
	return *m.metric
}

func (m *LookupsMetric) Get() Metrics {
	return *m.metric
}

func (m *MCacheInuseMetric) Get() Metrics {
	return *m.metric
}

func (m *MCacheSysMetric) Get() Metrics {
	return *m.metric
}

func (m *MSpanInuseMetric) Get() Metrics {
	return *m.metric
}

func (m *MSpanSysMetric) Get() Metrics {
	return *m.metric
}

func (m *MallocsMetric) Get() Metrics {
	return *m.metric
}

func (m *NextGCMetric) Get() Metrics {
	return *m.metric
}

func (m *NumForcedGCMetric) Get() Metrics {
	return *m.metric
}

func (m *NumGCMetric) Get() Metrics {
	return *m.metric
}

func (m *OtherSysMetric) Get() Metrics {
	return *m.metric
}

func (m *PauseTotalNsMetric) Get() Metrics {
	return *m.metric
}

func (m *StackInuseMetric) Get() Metrics {
	return *m.metric
}

func (m *StackSysMetric) Get() Metrics {
	return *m.metric
}

func (m *SysMetric) Get() Metrics {
	return *m.metric
}

func (m *TotalAllocMetric) Get() Metrics {
	return *m.metric
}

func (m *RandomValueMetric) Get() Metrics {
	return *m.metric
}

func (m *PollCountMetric) Get() Metrics {
	return *m.metric
}

func (m *TotalMemoryMetric) Get() Metrics {
	return *m.metric
}

func (m *FreeMemoryMetric) Get() Metrics {
	return *m.metric
}

func (m *CPUutilization1Metric) Get() Metrics {
	return *m.metric
}
