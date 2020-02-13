package service

import "fmt"

func ControllerRecover() {
	if err := recover(); err != nil {
		if err = fmt.Errorf("error: %v", err); err != nil {
			panic(err)
		}
	}
}
