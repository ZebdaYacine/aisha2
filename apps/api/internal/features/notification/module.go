package notification

// Module is the composition boundary for this feature. Business behavior is
// intentionally added through the documented domain, application, data, and
// server layers as the feature is implemented.
type Module struct{}

func NewModule() Module { return Module{} }
