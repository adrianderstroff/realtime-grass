package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
)

func CheckGLError() {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		errorType := string(err)
		switch err {
		case gl.INVALID_OPERATION:
			errorType = "INVALID OPERATION"
		case gl.INVALID_ENUM:
			errorType = "INVALID ENUM"
		case gl.INVALID_VALUE:
			errorType = "INVALID VALUE"
		case gl.OUT_OF_MEMORY:
			errorType = "OUT OF MEMORY"
		}
		panic("OpenGL error: " + errorType)
	} else {
		fmt.Println("No error")
	}
}
