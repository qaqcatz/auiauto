package perrorx

import (
	"strconv"
)

var GErrorOpen = "OpenError"
func NewErrorXOpen(filePath string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorOpen, "open " + filePath + " error", cause)
}

var GErrorFileExist = "FileExistError"
func NewErrorXFileExist(filePath string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorFileExist, "file " + filePath + " exists", cause)
}

var GErrorFileNotFound = "FileNotFoundError"
func NewErrorXFileNotFound(filePath string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorFileNotFound, "file " + filePath + " not found", cause)
}

var GErrorReadFile = "ReadFileError"
func NewErrorXReadFile(filePath string, originDesc string, cause *ErrorX) *ErrorX {
	if originDesc != "" {
		originDesc = "(" + originDesc + ")"
	}
	return NewErrorX(GErrorReadFile, "read " + filePath + " error" + originDesc, cause)
}

var GErrorReadDir = "ReadDirError"
func NewErrorXReadDir(filePath string, originDesc string, cause *ErrorX) *ErrorX {
	if originDesc != "" {
		originDesc = "(" + originDesc + ")"
	}
	return NewErrorX(GErrorReadDir, "read " + filePath + " error" + originDesc, cause)
}

var GErrorReadAll = "ReadAllError"
func NewErrorXReadAll(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorReadAll, desc, cause)
}

var GErrorWriteFile = "WriteFileError"
func NewErrorXWriteFile(filePath string, originDesc string, cause *ErrorX) *ErrorX {
	if originDesc != "" {
		originDesc = "(" + originDesc + ")"
	}
	return NewErrorX(GErrorWriteFile, "write " + filePath + " error" + originDesc, cause)
}

var GErrorRename = "RenameError"
func NewErrorXRename(filePath string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorRename, "rename " + filePath + " error", cause)
}

var GErrorFileCopy = "FileCopyError"
func NewErrorXFileCopy(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorFileCopy, desc, cause)
}

var GErrorMarshal = "MarshalError"
func NewErrorXMarshal(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorMarshal, desc, cause)
}

var GErrorUnmarshal = "UnmarshalError"
func NewErrorXUnmarshal(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorUnmarshal, desc, cause)
}

var GErrorBindJson = "BindJsonError"
func NewErrorXBindJson(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorBindJson, desc, cause)
}

var GErrorShellStart = "ShellStartError"
func NewErrorXShellStart(cmd string, desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorShellStart, cmd + " error: " + desc, cause)
}

var GErrorShellExecute = "ShellExecuteError"
func NewErrorXShellExecute(cmd string, desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorShellExecute, cmd + " error: " + desc, cause)
}

var GErrorADBShellBlock = "ADBShellBlockError"
func NewErrorXADBShellBlock(cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorADBShellBlock, "adb shell blocked", cause)
}

var GErrorGetAppConfig = "GetAppConfigError"
func NewErrorXGetAppConfig(appConfigPath string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGetAppConfig, "load " + appConfigPath + " error", cause)
}

var GErrorConnectAvd = "ConnectAvdError"
func NewErrorXConnectAvd(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorConnectAvd, desc, cause)
}

var GErrorInvalidAddress = "InvalidAddressError"
func NewErrorXInvalidAddress(address string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorInvalidAddress, "invalid address: " + address, cause)
}

var GErrorAntranceRquest = "AntranceRequestError"
func NewErrorXAntranceRquest(url string, desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorAntranceRquest, url + " : " + desc, cause)
}

var GErrorGetStmtLog = "GetStmtLogError"
func NewErrorXGetStmtLog(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGetStmtLog, desc, cause)
}

var GErrorGetUITree = "GetUITreeError"
func NewErrorXGetUITree(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGetUITree, desc, cause)
}

var GErrorAtoI = "AtoIError"
func NewErrorXAtoI(number string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorAtoI, "parse " + number + " error", cause)
}

var GErrorParseInt = "ParseIntError"
func NewErrorXParseInt(number string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorParseInt, "parse " + number + " error", cause)
}

var GErrorParseFloat = "ParseFloatError"
func NewErrorXParseFloat(number string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorParseFloat, "parse " + number + " error", cause)
}


var GErrorSplitN = "SplitNError"
func NewErrorXSplitN(sz int, n int, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorSplitN, "split size " + strconv.Itoa(sz) + " != " + strconv.Itoa(n), cause)
}

var GErrorStmtlogParse = "StmtlogParseError"
func NewErrorXStmtlogParse(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorStmtlogParse, desc, cause)
}

var GErrorCalLineCoverage = "CalLineCoverageError"
func NewErrorXCalLineCoverage(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorCalLineCoverage, desc, cause)
}

var GErrorLineCoverageDFS = "LineCoverageDFSError"
func NewErrorXLineCoverageDFS(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorLineCoverageDFS, desc, cause)
}

var GErrorGetCodesAndCoverLines = "GetCodesAndCoverLinesError"
func NewErrorXGetCodesAndCoverLines(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGetCodesAndCoverLines, desc, cause)
}

var GErrorAddCover = "AddCoverError"
func NewErrorXAddCover(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorAddCover, desc, cause)
}

var GErrorWaitAvaUITree = "WaitAvaUITreeError"
func NewErrorXWaitAvaUITree(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorWaitAvaUITree, desc, cause)
}

var GErrorPerform = "PerformError"
func NewErrorXPerform(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorPerform, desc, cause)
}

var GErrorPerforms = "PerformsError"
func NewErrorXPerforms(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorPerforms, desc, cause)
}

var GErrorReadSourceTreeAndAnalyze = "ReadSourceTreeAndAnalyzeError"
func NewErrorXReadSourceTreeAndAnalyze(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorReadSourceTreeAndAnalyze, desc, cause)
}

var GErrorReadSourceTreeAndCombineAnalyze = "ReadSourceTreeAndCombineAnalyzeError"
func NewErrorXReadSourceTreeAndCombineAnalyze(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorReadSourceTreeAndCombineAnalyze, desc, cause)
}

var GErrorCCL = "CCLError"
func NewErrorXCCL(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorCCL, desc, cause)
}

var GErrorRStmtlogNow = "RStmtlogNowError"
func NewErrorXRStmtlogNow(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorRStmtlogNow, desc, cause)
}

var GErrorGenInitSnapBeforeTest = "GenInitSnapBeforeTestError"
func NewErrorXGenInitSnapBeforeTest(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGenInitSnapBeforeTest, desc, cause)
}

var GErrorSaveStmtLogAndLineCoverage = "SaveStmtLogAndLineCoverageError"
func NewErrorXSaveStmtLogAndLineCoverage(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorSaveStmtLogAndLineCoverage, desc, cause)
}

var GErrorTesting = "TestingError"
func NewErrorXTesting(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorTesting, desc, cause)
}

var GErrorRTesting = "RTestingError"
func NewErrorXRTesting(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorRTesting, desc, cause)
}

var GErrorTimeout = "TimeoutError"
func NewErrorXTimeout(cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorTimeout, "", cause)
}

var GErrorGenEvent = "GenEventError"
func NewErrorXGenEvent(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorGenEvent, desc, cause)
}

var GErrorRDAnalyze = "RDAnalyzeError"
func NewErrorXRDAnalyze(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorRDAnalyze, desc, cause)
}

var GErrorSaveCharts = "SaveChartsError"
func NewErrorXSaveCharts(desc string, cause *ErrorX) *ErrorX {
	return NewErrorX(GErrorSaveCharts, desc, cause)
}