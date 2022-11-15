package observer

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	DefaultNamespace = "default"
	DefaultSubsystem = "activities"
	LabelCategory    = "category"
)

type ActivitiesObserver interface {
	Count(category string)
	DurationOf(category string, startedAt time.Time)
}

type activitiesObserver struct {
	counter  *prometheus.CounterVec
	duration *prometheus.SummaryVec
}

func NewActivitiesObserver() ActivitiesObserver {

	var (
		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: DefaultNamespace,
			Subsystem: DefaultSubsystem,
			Name:      "counter",
			Help:      "Counter of activities by category",
		}, []string{LabelCategory})

		duration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:  DefaultNamespace,
			Subsystem:  DefaultSubsystem,
			Name:       "duration",
			Help:       "Duration of activities by category",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{LabelCategory})
	)

	prometheus.MustRegister(counter)
	prometheus.MustRegister(duration)

	return &activitiesObserver{
		counter:  counter,
		duration: duration,
	}
}

func (o *activitiesObserver) Count(category string) {
	o.counter.WithLabelValues(category).Inc()
}

func (o *activitiesObserver) DurationOf(category string, startedAt time.Time) {
	o.duration.WithLabelValues(category).Observe(ms(startedAt))
}

func ms(startAt time.Time) float64 {
	return float64(time.Since(startAt).Nanoseconds()) / 1e6
}
