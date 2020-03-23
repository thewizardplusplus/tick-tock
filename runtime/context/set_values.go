package context

// ValueHolder ...
//go:generate mockery -name=ValueHolder -case=underscore
type ValueHolder interface {
	SetValue(name string, value interface{})
}

// ValueGroup ...
type ValueGroup map[string]interface{}

// SetValues ...
func SetValues(holder ValueHolder, values ValueGroup) {
	for name, value := range values {
		holder.SetValue(name, value)
	}
}
