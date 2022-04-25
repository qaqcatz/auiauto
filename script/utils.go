package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
)

var Projects = []string{
	"ActivityDiary-118",        // 0,  1,     191
	"AlarmClock_389",           // 1,  1,     234
	"AmazeFileManager-1796",    // 2,  4(1),  255
	"AmazeFileManager-1837",    // 3   5(1),  276
	"and-bible-261",            // 4   1,     219
	"and-bible-375",            // 5   12(1), 119
	"and-bible-480",            // 6   10(1), 148
	"Anki-Android-4586",        // 7   78(1), 107
	"Anki-Android-4589",        // 8   58(1), 58
	"Anki-Android-4977",        // 9   2(2),  307
	"Anki-Android-5638",        // 10  1,     276
	"Anki-Android-5756",        // 11  1,     146
	"Anki-Android-6145",        // 12  1,     114
	"AntennaPod-3138",          // 13  8(3),  362
	"AntennaPod_4645",          // 14  1,     311
	"AnyMemo_422",              // 15  1,     307
	"AnyMemo_440",              // 16  1,     241
	"APhotoManager_139",        // 17  2(1),  212
	"collect-3222",             // 18  1,     207
	"commons-1391",             // 19  1,     171
	"commons-1581",             // 20  1,     280
	"commons-2123",             // 21  8(8),  357
	"dagger-46",                // 22  1,     352
	"Easy_xkcd_134",            // 23  24(1), 191
	"FirefoxLite-4881",         // 24  1,     174
	"FirefoxLite-4942",         // 25  1,     352
	"FirefoxLite-5085",         // 26  1,     253
	"geohashdroid-73",          // 27  1,     275
	"Images-to-PDF_585",        // 28  52(18) 299
	"Images-to-PDF_771",        // 29  9(9)   362
	"k-9_3255",                 // 30  1,     281
	"markor_194",               // 31  1,     265
	"nextcloud-1918",           // 32  3(1),  281
	"nextcloud-4026",           // 33  1,     170
	"nextcloud-4792",           // 34  1,     186
	"nextcloud-5173",           // 35  1,     233
	"open-event-attendee-2198", // 36  1,     325
	"openlauncher-67",          // 37  1,     209
	"opentasks_629",            // 38  1,     251
	"osmeditor4android-729",    // 39  1,     172
	"Scarlet-Notes-114",        // 40  1,     201
	"screenrecorder-32",        // 41  1,     252
	"Simple-Music-Player_128",  // 42  1,     128
	"Simple-Music-Player_204",  // 43  1,     174
	"ToGoZip_10",               // 44  1,     209
	"ToGoZip_20",               // 45  1,     195
}

var Projects2 = []string{
	"ActivityDiary-118",        // 0,  1,     191
	"AlarmClock_389",           // 1,  1,     234
	"AmazeFileManager-1796",    // 2,  4(1),  255
	"AmazeFileManager-1837",    // 3   5(1),  276
	"and-bible-261",            // 4   1,     219
	"and-bible-375",            // 5   12(1), 119
	"and-bible-480",            // 6   10(1), 148
	"Anki-Android-4586",        // 7   78(1), 107
	"Anki-Android-4589",        // 8   58(1), 58
	"Anki-Android-4977",        // 9   2(2),  307
	"Anki-Android-5638",        // 10  1,     276
	"Anki-Android-5756",        // 11  1,     146
	"Anki-Android-6145",        // 12  1,     114
	"AntennaPod-3138",          // 13  8(3),  362
	"AntennaPod_4645",          // 14  1,     311
	"AnyMemo_422",              // 15  1,     307
	"AnyMemo_440",              // 16  1,     241
	"APhotoManager_139",        // 17  2(1),  212
	"collect-3222",             // 18  1,     207
	"commons-1391",             // 19  1,     171
	"commons-1581",             // 20  1,     280
	"commons-2123",             // 21  8(8),  357
	"dagger-46",                // 22  1,     352
	"Easy_xkcd_134",            // 23  24(1), 191
	"FirefoxLite-4881",         // 24  1,     174
	"FirefoxLite-4942",         // 25  1,     352
	"FirefoxLite-5085",         // 26  1,     253
	"geohashdroid-73",          // 27  1,     275
	"Images-to-PDF_585",        // 28  52(18) 299
	"Images-to-PDF_771",        // 29  9(9)   362
	"k-9_3255",                 // 30  1,     281
	"markor_194",               // 31  1,     265
	"nextcloud-1918",           // 32  3(1),  281
	"nextcloud-4026",           // 33  1,     170
	"nextcloud-4792",           // 34  1,     186
	"nextcloud-5173",           // 35  1,     233
	"open-event-attendee-2198", // 36  1,     325
	"openlauncher-67",          // 37  1,     209
	"opentasks_629",            // 38  1,     251
	"osmeditor4android-729",    // 39  1,     172
	"Scarlet-Notes-114",        // 40  1,     201
	"screenrecorder-32",        // 41  1,     252
	"Simple-Music-Player_128",  // 42  1,     128
	"Simple-Music-Player_204",  // 43  1,     174
	"ToGoZip_10",               // 44  1,     209
	"ToGoZip_20",               // 45  1,     195
}

// 调用/bin/bash -c执行myCmdStr, 出现异常时默认合并输出流和错误流到error中
func Sh(myCmdStr string) (string, error) {
	exeCmd := exec.Command("/bin/bash", "-c", myCmdStr)

	stdout, _ := exeCmd.StdoutPipe()
	stderr, _ := exeCmd.StderrPipe()
	defer stdout.Close()
	defer stderr.Close()

	if err := exeCmd.Start(); err != nil {
		return "", err
	}
	// 输出流
	scanner := bufio.NewScanner(stdout)
	var outBuf bytes.Buffer
	for scanner.Scan() {
		outBuf.WriteString(scanner.Text() + "\n")
	}
	myOut := outBuf.String()
	// 错误流
	scanner = bufio.NewScanner(stderr)
	var errBuf bytes.Buffer
	for scanner.Scan() {
		errBuf.WriteString(scanner.Text() + "\n")
	}
	myErr := errBuf.String()

	if myErr != "" {
		// 异常时默认合并错误流和输出流
		return "", errors.New(myErr + "\n" + myOut)
	}
	return myOut, nil
}

// 判断文件是否存在
func FileExist(myPath string) bool {
	_, err := os.Stat(myPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 将文件src拷贝到文件dst, dst不存在会新建, 但文件夹要自己创建
func FileCopy(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return errors.New(src + " status error!")
	}
	if !sourceFileStat.Mode().IsRegular() {
		return errors.New(src + " is not a regular file")
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	return nil
}

// cp -r from to
func CopyR(from, to string) error {
	_, err := Sh("cp -r " + from + " " + to)
	if err != nil {
		return err
	}
	return nil
}
