package gw_errors

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/getsentry/sentry-go"

	//"reflect"
	"runtime"

	"github.com/morikuni/failure/v2"
)

type ErrorCode string

const (
	GenericError        ErrorCode = "GenericError" //一般的なエラー。デフォルトでこれを指定
	InternalServerError ErrorCode = "InternalServerError"
	BadRequest          ErrorCode = "BadRequest"
	NotFound            ErrorCode = "NotFound"
	Forbidden           ErrorCode = "Forbidden"
	UnknownError        ErrorCode = "UnknownError" // どれにも当てはまらないエラー
)

func errorLog(err error, objList ...interface{}) {
	//error message, relational data
	errorMessageList := []string{"err: " + err.Error()}
	for ind, obj := range objList {
		relationalStr := fmt.Sprintf("error relational data %v: %v\n", ind, obj)
		errorMessageList = append(errorMessageList, relationalStr)
	}
	//stacktrace
	stackTraceStr, ok := CallStackOf(err)
	if !ok {
		pc, fileName, line, ok := runtime.Caller(2)

		if ok {
			stackTraceStr = fmt.Sprintf("memory address: %v, file: %v, line: %v \n", pc, fileName, line)
		} else {
			stackTraceStr = fmt.Sprintf("can't get line data \n")
		}
	}
	errorMessageList = append(errorMessageList, "", stackTraceStr)

	//logging & send to sentry server
	log.Println(strings.Join(errorMessageList, "\n"))
	sentry.CaptureMessage(strings.Join(errorMessageList, "\n"))
}

func New(errStr string) error {
	return failure.New(GenericError, failure.Message(errStr))
}
func Errorf(format string, a ...interface{}) error {
	return Wrap(fmt.Errorf(format, a...))
}
func HasStack(err error) (failure.CallStack, bool) {
	stack := failure.CallStackOf(err)
	if stack == nil || len(stack.Frames()) == 0 {
		return nil, false
	} else {
		return stack, true
	}
}
func CallStackOf(err error) (stackTrace string, ok bool) {
	if stack, ok := HasStack(err); ok {
		out := &bytes.Buffer{}
		for _, f := range stack.Frames() {
			p := f.Path()
			fmt.Fprintf(out, "%s:%d [%s.%s]\n", p, f.Line(), f.Pkg(), f.Func())
		}
		return out.String(), ok
	} else {
		return "", false
	}
}
func Wrap(err error, objList ...interface{}) error {
	const KEY_WRAP_COUNT = "wrapCnt"
	if err == nil {
		return err
	}

	// originalStack := failure.CallStackOf(err)
	// originalMessage := failure.MessageOf(err)
	originalCode := failure.CodeOf(err)

	//errに、すでにwrapされた回数があれば、それを取得して、+1する
	var failureCtx *failure.Context
	objStrList := []string{}
	for _, obj := range objList {
		// relationalStr := fmt.Sprintf("error relational data %v: %v\n", ind, obj)
		objStr := "nil"
		if obj != nil {
			objStr = fmt.Sprintf("%v", obj)
		}
		objStrList = append(objStrList, objStr)
	}
	if objStrList != nil && len(objStrList) > 0 {
		failureCtx = &failure.Context{"params": fmt.Sprintf("%v", objStrList)}
	}

	if originalCode != nil {
		var err error
		if failureCtx != nil {
			err = failure.New(err, failureCtx)
		} else {
			err = failure.New(err)
		}

		stack := failure.CallStackOf(err)
		if stack != nil {
			log.Println("stacktrace---------------------------------------------------------------------------:")
			for _, f := range stack.Frames() {
				p := f.Path()
				log.Printf("%s:%d [%s.%s]\n", p, f.Line(), f.Pkg(), f.Func())
			}
		} else {
			log.Println("stacktrace is not found!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!.")

		}
		return failure.New(err, failureCtx)
	} else {
		if failureCtx != nil {
			return failure.Wrap(err, failureCtx)
		} else {
			return failure.Wrap(err)
		}
	}

	//stacktrace
	// pc, fileName, line, ok := runtime.Caller(1)
	// stackTraceStr := ""
	// if ok {
	// 	stackTraceStr = fmt.Sprintf("memory address: %v, file: %v, line: %v \n", pc, fileName, line)
	// } else {
	// 	stackTraceStr = fmt.Sprintf("can't get line data \n")
	// }
	// errorMessageList = append(errorMessageList, "", stackTraceStr)
	// return pkg_errors.Wrap(err, strings.Join(errorMessageList, "\n"))
}

func ReturnError(err error, objList ...interface{}) error {
	if err != nil {
		errorLog(err, objList)
	}
	return Wrap(err, objList...)
}
func ReturnErrorStr(errStr string) error {
	if errStr != "" {
		err := New(errStr)
		errorLog(err)
		return err
	}
	return nil
}

func PanicError(err error, objList ...interface{}) {
	if err != nil {
		err = Wrap(err, objList...)
		errorLog(err, objList)
		panic(err)
	}
}
func PanicErrorStr(errStr string, objList ...interface{}) {
	if errStr != "" {
		err := New(errStr)
		errorLog(err, objList...)
		panic(err)
	}
}
func PanicErrorWithFunc(err error, f func(), objList ...interface{}) {
	if err != nil {
		err = Wrap(err, objList...)
		errorLog(err, objList)
		//c.Status(status)
		f()
		panic(err)
	}
}

func PrintError(err error, objList ...interface{}) {
	if err != nil {
		err = Wrap(err, objList...)
		errorLog(err, objList)
	}
}
func PrintErrorStr(errStr string, objList ...interface{}) {
	if errStr != "" {
		err := New(errStr)
		errorLog(err, objList)
	}
}
