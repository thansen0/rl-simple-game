package logic

type Logic interface {
	Init()
	Action() (delta_X, delta_Y float64)
	Update()
}

// NOTE
// None of the game logic is in any of these interface implementations yet,
// they're a work in progress

// func (l *Logic) Init() {
// 	fmt.Println("This function should be re-implemented")
// }

// // takes in state and returns movements
// func (l *Logic) Action() (delta_X, delta_Y float64) {
// 	fmt.Println("This function should be re-implemented")
// 	return 0, 0
// }

// // update Q matrix
// func (l *Logic) Update() {
// 	fmt.Println("This function should be re-implemented")
// }
