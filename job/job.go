package job

// Resource is essentially just the data that we'll be using and doing work with.
type Resource string

// ProcessFunc processes resources and may or may not return an error.
type ProcessFunc func(r Resource) (uint32, error)

// ProcessResultFunc processes a result and may or may not return an error.
type ProcessResultFunc func(r Result) Result

// Job describes information associated with a specific job
type Job struct {
	ID       uint
	Resource Resource
}

// Result describes the output after processing a job.
type Result struct {
	Job Job
	// This is simply the resultant value after processing.
	Value uint32
	// Determines whether or not the resultant outcome is an error.
	Err error
}
