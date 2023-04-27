package mobile

const (
	CStateStop = 0
	CStateRun  = 1
)

const (
	CAndroidRootPath = "/sdcard/Documents"
	CAndroidFullPath = CAndroidRootPath + "/hidden_lake"
)

func SwitchState(state int) int {
	if state == CStateStop {
		return CStateRun
	}
	return CStateStop
}

func ButtonTextFromState(pState int, pServiceName string) string {
	if pState == CStateStop {
		return "Run " + pServiceName
	}
	return "Stop " + pServiceName
}
