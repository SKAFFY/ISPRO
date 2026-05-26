package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EntriesCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "entries_created_total",
			Help: "Total number of entries created",
		},
	)

	LinksCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "links_created_total",
			Help: "Total number of links created",
		},
	)

	EntriesDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "entries_deleted_total",
			Help: "Total number of entries deleted",
		},
	)

	LinksDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "links_deleted_total",
			Help: "Total number of links deleted",
		},
	)
)
