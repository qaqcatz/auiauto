package pstatistic

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/psrctree"
	"auiauto/putils"
	"io/ioutil"
	"path"
	"strings"
)

// 获取dir下的所有以casePrefix为前缀的pass用例和crash(返回json路径), 用名字是否包含pass|crash来判断
func readPassCrashCover(dirPath string, casePrefix string) ([]string, []string, *perrorx.ErrorX) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, perrorx.NewErrorXReadDir(dirPath, err.Error(), nil)
	}
	passPaths := make([]string, 0)
	crashPaths := make([]string, 0)
	for _, fi := range dir {
		if fi.IsDir() && strings.HasPrefix(fi.Name(), casePrefix) {
			if strings.Contains(fi.Name(), "pass") && putils.FileExist(path.Join(dirPath, fi.Name(), "cover.json")) {
				passPaths = append(passPaths, path.Join(dirPath, fi.Name(), "cover.json"))
			} else if strings.Contains(fi.Name(), "crash") && putils.FileExist(path.Join(dirPath, fi.Name(), "cover.json")) {
				crashPaths = append(crashPaths, path.Join(dirPath, fi.Name(), "cover.json"))
			}
		}
	}
	return passPaths, crashPaths, nil
}

// 根据projectId/src创建源码树, 读取cases列表的路径进行源码树覆盖, 通过factor进行差异分析, 根据analyzeFile进行分析结果过滤
func readSourceTreeAndCoverEach(projectId string, cases []string) ([]int, *perrorx.ErrorX) {
	sourceTree, err := psrctree.ReadSourceTree(projectId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	ans := make([]int, 0)
	for _, casePath := range cases {
		err = psrctree.CoverSourceTree(sourceTree, casePath, "")
		if err != nil {
			return nil, perrorx.TransErrorX(err)
		}
		psrctree.CalCoverNumDFS(sourceTree.MRoot)
		ans = append(ans, sourceTree.MRoot.MCoverNum)
	}
	return ans, nil
}

// 统计随机Testing结果
func RDTestingAnalyzeStd(projectId string, tester string, casePrefix string) ([]int, *perrorx.ErrorX) {
	testPath := pdba.DBURLProjectIdTester(projectId, tester)
	// 获取dirPath下的正确用例和错误用例(name->path)
	passPaths, crashPaths, err := readPassCrashCover(testPath, casePrefix)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	cases := make([]string, 0)
	for i := 0; i < len(passPaths); i++ {
		cases = append(cases, passPaths[i])
		if i < len(crashPaths) {
			cases = append(cases, crashPaths[i])
		}
	}
	if len(passPaths) < len(crashPaths) {
		for i := len(passPaths); i < len(crashPaths); i++ {
			cases = append(cases, crashPaths[i])
		}
	}

	ans, err := readSourceTreeAndCoverEach(projectId, cases)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	return ans, nil
}
