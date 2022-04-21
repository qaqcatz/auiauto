package panalyze

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/psrctree"
	"auiauto/pkernel/psuscode"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"strconv"
	"sync"
	"time"
)

// 随机采样分析, 在指定目录下生成分析文件rdanalyze_casePrefix_analyzeFile_factor.txt
// 格式:
// 第一行四个部分, 组合总数n(拆分成正确用例数和错误用例数), 每个组合下的随机次数m, 聚焦的语句文件(不能为空), 聚焦语句数量s
// 接下来n个组合, 每个组合:
// 	第一行两个数, 正确用例数, 错误用例数
// 	接下来m次随机, 每次随机:
//	  s个数, 代表聚焦语句文件中每个聚焦语句的排名

// 随机采样分析, 首先获取projectId/test/tester下以casePrefix为前缀的正确用例和错误用例, 接着根据passNum和crashNum进行有放回采样,
// 每组采样重复randNum次, 每次调用rdAnalyze进行分析, 获取susCodeSlice, 按顺序输出排名(
// 切片中的语句顺序和analyzeFile的顺序是一致的, 有需要的话可以自己按顺序还原到各个类的各个语句上)
// 最后在testPath下生成分析文件rdanalyze_casePrefix_analyzeFile_factor.txt
func RDAnalyzeStd(projectId string, analyzeFile string, factor string,
	tester string, casePrefix string, randNum int) *perrorx.ErrorX {
	testPath := pdba.DBURLProjectIdTester(projectId, tester)
	// 获取dirPath下的正确用例和错误用例(name->path)
	passPaths, crashPaths, err := readPassCrash0(testPath, casePrefix)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	passNum := len(passPaths)
	crashNum := len(crashPaths)
	// 读取analyzeFile, 获取聚焦语句总数
	analyzeClasses, err := psuscode.ReadAnalyzeFile(projectId, analyzeFile)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	totalLineNum := analyzeClasses.GetTotalLineNum()

	// 处理\r
	defer func () {
		fmt.Println()
	} ()

	ans := ""
	// 第一行四个部分, 组合总数n(拆分成正确用例数和错误用例数), 每个组合下的随机次数m, 聚焦的语句文件(不能为空), 聚焦语句数量
	ans += strconv.Itoa(passNum) + " " + strconv.Itoa(crashNum) + " " + strconv.Itoa(randNum) + " " +
		analyzeFile + " " + strconv.Itoa(totalLineNum) + "\n"
	// 枚举组合情况
	cur := 0
	total := passNum * crashNum
	for x := 1; x <= passNum; x++ {
		for y := 1; y <= crashNum; y++ {
			cur++
			fmt.Print("\r" + strconv.Itoa(cur) + "/" + strconv.Itoa(total))
			// 每个组合第一行两个部分, 正确用例数, 错误用例数
			ans += strconv.Itoa(x) + " " + strconv.Itoa(y) + "\n"
			// 并发跑每个随机, 用err_记录并发过程中有没有错误
			var wg sync.WaitGroup
			var mutex sync.Mutex
			var err_ *perrorx.ErrorX = nil
			// 重复randNum次
			for r := 0; r < randNum; r++ {
				wg.Add(1)
				go func () {
					// 根据x, y有放回采样
					_, susCodeSlice, err := rdSampleAndGenSusCodeSlice(projectId, analyzeFile, factor, passPaths, crashPaths, x, y)
					mutex.Lock()
					if err_ == nil && err != nil {
						err_ = perrorx.TransErrorX(err)
					}
					if err == nil {
						// 注意susCodeSlice中的语句顺序和analyzeFile的顺序是一致的, 不用担心乱掉
						for _, susCode := range susCodeSlice {
							ans += strconv.Itoa(susCode.MRank) + " "
						}
						ans += "\n"
					}
					mutex.Unlock()
					wg.Done()
				} ()
			}
			wg.Wait()
			if err_ != nil {
				return perrorx.TransErrorX(err_)
			}
		}
	}
	// 最后在testPath下生成分析文件rdanalyze_casePrefix_analyzeFile_factor.txt
	rdAnalyzePath := path.Join(testPath, "rdanalyze_"+casePrefix+"_"+analyzeFile+"_"+factor+".txt")
	err_ := ioutil.WriteFile(rdAnalyzePath, []byte(ans), 0777)
	if err_ != nil {
		return perrorx.NewErrorXWriteFile(rdAnalyzePath, err_.Error(), nil)
	}
	return nil
}

// passPaths中随机有放回采样x个, crashPaths中有放回采样y个进行分析, 返回src tree和sus code slice
// 注意拼接上cover.json
func rdSampleAndGenSusCodeSlice(projectId string, analyzeFile string, factor string, passPaths []string, crashPaths []string, x int, y int) (*psrctree.SourceTree, psuscode.SusCodes, *perrorx.ErrorX) {
	// 随机种子
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	samplePaths := make([]string, 0)
	// passNames中采样x个
	for i := 0; i < x; i++ {
		id := rd.Intn(len(passPaths))
		samplePaths = append(samplePaths, path.Join(passPaths[id], "cover.json"))
	}
	// crashNames中采样y个
	for i := 0; i < y; i++ {
		id := rd.Intn(len(crashPaths))
		samplePaths = append(samplePaths, path.Join(crashPaths[id], "cover.json"))
	}
	srcTree, susCodeSlice, err := readSourceTreeAndAnalyze(projectId, samplePaths, analyzeFile, factor)
	if err != nil {
		return nil, nil, perrorx.TransErrorX(err)
	}
	return srcTree, susCodeSlice, nil
}
