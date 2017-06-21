package prober

// Prober defines a type that can be used to run a probe.
type Prober interface {
	Probe() error
}

// HealthChecker defines a type that can be used to perform a health check.
type HealthChecker interface {
	HealthCheck() bool
}

// DashboardYielder defines a type that can yield a dashboard.
type DashboardYielder interface {
	GetPort() uint
	HealthChecker
}
