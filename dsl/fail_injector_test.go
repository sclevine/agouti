package dsl

func InjectFail(failFunc func(message string, callerSkip ...int)) {
    failDSL = failFunc
}
