//go:build nodebug

package framework

//
// This implementation of debug output drops all debug output.
//
// Use the build tag 'nodebug' to explicitly drop debug output
// (the default is to output debug messages to the TinyGo serial
// output)
//

const DebugEnabled = false

func Debug(msg string) {
	_ = msg
}
