package controllers

type IController interface {
	Start() error
	Exit()
}
