package logic

// type Logic interface {
// 	Init()
// 	Action() (delta_X, delta_Y float64)
// 	Update()
// }

type Follow struct {
}

// Nothing to update
func (f *Follow) Init() {}

func Action() (delta_X, delta_Y float64) {
	return 0.0, 0.0
}

func Update() {}
