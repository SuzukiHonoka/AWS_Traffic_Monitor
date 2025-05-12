package instance

type List []*Instance

// Check checks all instances in the list
func (s List) Check() {
	for _, v := range s {
		v.Check()
	}
}
