package Error

import log "github.com/sirupsen/logrus"

const ERRORMSG = "GET error != nil"

// HandleError is a function to handle error using Errorln, do nothing if error is nil
func HandleError(e error) {
	if e != nil {
		log.Errorln(ERRORMSG, e)
	}
}

// HandlePanic is a function to handle error using Panicln, do nothing if error is nil
func HandlePanic(e error) {
	if e != nil {
		log.Panicln(ERRORMSG, e)
	}
}

// HandleFatal is a function to handle error using Fatalln, do nothing if error is nil
func HandleFatal(e error) {
	if e != nil {
		log.Fatalln(ERRORMSG, e)
	}
}
