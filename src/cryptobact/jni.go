package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>

*/
import "C"
import "log"
import "unsafe"
import "runtime/debug"
import "cryptobact/ui"

// Use JNI_OnLoad to ensure that the go runtime is initialized at a predictable time,
// namely at System.loadLibrary()
//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	return C.JNI_VERSION_1_6
}

//export Java_net_goandroid_cryptobact_Engine_drawFrame
func Java_net_goandroid_cryptobact_Engine_drawFrame(env *C.JNIEnv, clazz C.jclass) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: drawFrame: %v\nstack:\n%s", err, string(debug.Stack()))
		}
	}()
	g.drawFrame()
}

//export Java_net_goandroid_cryptobact_Engine_init
func Java_net_goandroid_cryptobact_Engine_init(env *C.JNIEnv, clazz C.jclass) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: init: %v\nstack:\n%s", err, string(debug.Stack()))
		}
	}()
	g.initGL()
}

//export Java_net_goandroid_cryptobact_Engine_resize
func Java_net_goandroid_cryptobact_Engine_resize(env *C.JNIEnv, clazz C.jclass, width, height C.jint) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: resize: %v\nstack:\n%s", err, string(debug.Stack()))
		}
	}()
	g.resize(int(width), int(height))
}

//export Java_net_goandroid_cryptobact_Engine_onTouch
func Java_net_goandroid_cryptobact_Engine_onTouch(env *C.JNIEnv, clazz C.jclass, action C.jint, x, y C.jfloat) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: touch: %v\nstack:\n%s", err, string(debug.Stack()))
		}
	}()
	actionI := int(action)
	switch action {
	case C.AMOTION_EVENT_ACTION_UP:
		actionI = ui.AMOTION_EVENT_ACTION_UP
	case C.AMOTION_EVENT_ACTION_DOWN:
		actionI = ui.AMOTION_EVENT_ACTION_DOWN
	case C.AMOTION_EVENT_ACTION_MOVE:
		actionI = ui.AMOTION_EVENT_ACTION_MOVE
	}

	g.onTouch(actionI, float32(x), float32(y))
}
