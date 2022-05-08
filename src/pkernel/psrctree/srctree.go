package psrctree

import (
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/pkernel/pcoverage"
	"auiauto/putils"
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// 源码树, 用于覆盖率可视化以及差异分析
type SourceTree struct {
	// 根节点
	MRoot *SourceNode `json:"root"`
	// 正确用例总数
	MTotalPassed int `json:"-"`
	// 错误用例总数
	MTotalFailed int `json:"-"`
	// 关联类与源文件(一个源文件中可能包含多个类, 需要维护好这种关系才能做好覆盖率可视化)
	MClsSrcMap map[string]string `json:"-"`
}

// 前端点击源文件后向后端发起请求, 拿到代码以及覆盖信息, 因此后端需要维护一颗全局的source tree
// 每次调用ReadSourceTree时更新
// 这里也是并发时一个很头疼的点, 除了前端访问, 我们不能对GSourceTree做任何读操作, 并对写加锁. 这样在没有前端访问时脏数据就不会传播出来了.
var gSourceTreeMutex sync.Mutex
var gSourceTree *SourceTree

func GetGSourceTreeJS() *SourceTree {
	gSourceTreeMutex.Lock()
	defer gSourceTreeMutex.Unlock()
	return gSourceTree
}

func SetGSourceTree(sourceTree *SourceTree) {
	gSourceTreeMutex.Lock()
	defer gSourceTreeMutex.Unlock()
	gSourceTree = sourceTree
}

func (sourceTree *SourceTree) Foreach(f func(cur *SourceNode)) {
	sourceTree.foreachDFS(sourceTree.MRoot, f)
}

func (sourceTree *SourceTree)foreachDFS(cur *SourceNode, f func(cur *SourceNode)) {
	f(cur)
	for _, child := range cur.MChildren {
		sourceTree.foreachDFS(child, f)
	}
}

// 根据dotClassPath查询节点, 并转换成CodesAndCoverLines形式返回, 用于可视化显示
func (sourceTree *SourceTree) GetCodesAndCoverLines(dotClassPath string) (*CodesAndCoverLines, *perrorx.ErrorX) {
	sp := strings.Split(dotClassPath, ".")

	cur := sourceTree.MRoot
	for i := 0; i < len(sp); i++ {
		if nex, exist := cur.MChildrenMap[sp[i]]; exist {
			cur = nex
		} else {
			return nil, perrorx.NewErrorXGetCodesAndCoverLines("class " + dotClassPath + " does not exist", nil)
		}
	}

	codes := make([]string, cur.MTotalNum)
	codesType := make([]int, cur.MTotalNum)
	for i := 0; i < cur.MTotalNum; i++ {
		codes[i] = cur.MCodes[i]
		if cur.MEventIds[i] != nil && len(cur.MEventIds[i]) != 0 {
			// 每行代码后接一条注释"// 1, 2, ...", 用于表示这行代码被哪些动作执行到过
			codes[i] += " //"
			for _, id := range cur.MEventIds[i] {
				codes[i] += " " + strconv.Itoa(id)
			}
		}
		// -1:未覆盖
		// 0: 覆盖,  lightgreen
		// 1: 可疑度rank1, DarkRed
		// 2: 可疑度rank10, Red
		// 3: 可疑度rank100, LightCoral
		// 4: 可疑度rank other, Yellow
		codesType[i] = -1
		if cur.MPassed[i] + cur.MFailed[i] > 0 {
			if cur.MRanks[i] == 0 {
				codesType[i] = 0
			} else if 0 < cur.MRanks[i] && cur.MRanks[i] <= 10 {
				codesType[i] = 1
			} else if 10 < cur.MRanks[i] && cur.MRanks[i] <= 50 {
				codesType[i] = 2
			} else if 50 < cur.MRanks[i] && cur.MRanks[i] <= 100 {
				codesType[i] = 3
			} else {
				codesType[i] = 4
			}
		}
	}
	return &CodesAndCoverLines{codes, codesType}, nil
}

// 根据dotClassPath(此时已经去掉了$)寻找其在src tree上对应的节点
// 对应class和源文件不是建简单的事情, 为此我们做了一些匹配规则:
// 1. 直接看搜索最后一步的节点是否和dotClassPath相对应
// 2. 若最后一步没找到, 则去掉dotClassPath的Kt后缀再找一次(kotlin一些源文件编译成class后会加Kt)
// 3. 若去掉Kt找不到, 则(此时不要去Kt)去clsSrcMap找, clsSrcMap是我们通过简单语法分析得到的class与源文件的映射关系
func (sourceTree *SourceTree) FindSrcNode(dotClassPath string) *SourceNode {
	sp := strings.Split(dotClassPath, ".")
	cur := sourceTree.MRoot
	for i := 0; i < len(sp); i++ {
		if nex, exist := cur.MChildrenMap[sp[i]]; exist {
			cur = nex
		} else {
			if i == len(sp)-1 {
				if strings.HasSuffix(sp[i], "Kt") {
					if nex, exist := cur.MChildrenMap[sp[i][0:len(sp[i])-2]]; exist {
						cur = nex
						break
					}
				}
				if sourceName, exist := sourceTree.MClsSrcMap[dotClassPath]; exist {
					if nex, exist := cur.MChildrenMap[sourceName]; exist {
						cur = nex
						break
					}
				}
			} else if 'a' <= sp[i][0] && sp[i][0] <= 'z' {
				fst := make([]byte, 1)
				fst[0] = sp[i][0]-'a'+'A'
				// 包名首字母大写在类文件中会自动转为小写, 特殊处理这种情况
				upFirst := (string)(fst)+sp[i][1:]
				if nex, exist := cur.MChildrenMap[upFirst]; exist {
					cur = nex
					continue
				}
			}
			return nil
		}
	}
	return cur
}

// 为srctree添加覆盖, 具体而言是根据dotClassPath找到对应节点, 根据coverLines和eventIds为节点中相应的代码行做标记,
// 通过status判断当前是被正确用例标记还是错误用例标记
// eventId表示只统计被这个动作覆盖的行, 例如3的话表示统计被3覆盖的语句, 3!的话表示统计只被3覆盖的语句, 在覆盖显示时才会生效, 分析时默认为空
func (sourceTree *SourceTree) addCover(dotClassPath string, coverLines []int, eventIds []string, status string,
	eventId string) *perrorx.ErrorX {
	cur := sourceTree.FindSrcNode(dotClassPath)
	if cur == nil {
		return perrorx.NewErrorXAddCover("class " + dotClassPath + " does not exist", nil)
	}
	if eventIds != nil {
		for i := 0; i < len(coverLines); i++ {
			// kt伴生对象可能导致class的代码行数大于源码行数, 特殊处理这种情况
			if coverLines[i]-1 >= cur.MTotalNum {
				continue
			}
			ids, err := strconv.ParseInt(eventIds[i], 10, 64)
			if err != nil {
				return perrorx.NewErrorXParseInt(eventIds[i], nil)
			}
			temp := make([]int, 0)
			one := int64(1)
			for id := 0; id <= 62; id++ {
				if (ids & (one<<id)) != 0 {
					temp = append(temp, id)
				}
			}
			cur.MEventIds[coverLines[i]-1] = temp
		}
	}
	if status == "true" {
		for i := 0; i < len(coverLines); i++ {
			// kt伴生对象可能导致class的代码行数大于源码行数, 特殊处理这种情况
			if coverLines[i]-1 >= cur.MTotalNum {
				continue
			}
			if eventIds != nil && eventId != "" {
				ok := false
				for _, id := range cur.MEventIds[coverLines[i]-1] {
					if strconv.Itoa(id) == eventId || strconv.Itoa(id)+"!" == eventId {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				if ok && strings.HasSuffix(eventId, "!") &&
					len(cur.MEventIds[coverLines[i]-1]) != 1 {
					continue
				}
			}
			if cur.MPassed[coverLines[i]-1] == 0 {
				cur.MCoverNum += 1
			}
			cur.MPassed[coverLines[i]-1] += 1
		}
	} else {
		for i := 0; i < len(coverLines); i++ {
			// kt伴生对象可能导致class的代码行数大于源码行数, 特殊处理这种情况
			if coverLines[i]-1 >= cur.MTotalNum {
				continue
			}
			if eventIds != nil && eventId != "" {
				ok := false
				for _, id := range cur.MEventIds[coverLines[i]-1] {
					if strconv.Itoa(id) == eventId || strconv.Itoa(id)+"!" == eventId {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				if ok && strings.HasSuffix(eventId, "!") &&
					len(cur.MEventIds[coverLines[i]-1]) != 1 {
					continue
				}
			}
			if cur.MFailed[coverLines[i]-1] == 0 {
				cur.MCoverNum += 1
			}
			cur.MFailed[coverLines[i]-1] += 1
		}
	}
	return nil
}

// 读取projectId下的clssrcmap
func readClsSrcMap(projectId string) (map[string]string, *perrorx.ErrorX) {
	mapPath := pdba.DBURLProjectIdClssrcmap(projectId)
	if file, err := os.Open(mapPath); err != nil {
		return nil, perrorx.NewErrorXOpen(mapPath, nil)
	} else {
		ans := make(map[string]string)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, "@") {
				continue
			}
			sp := strings.Split(line, "@")
			if len(sp) != 2 {
				continue
			}
			ans[sp[0]] = sp[1]
		}
		return ans, nil
	}

}

// 递归生成源码树
func readSourceTreeDFS(filePath string, name string, fullName string, isSource bool) (*SourceNode, *perrorx.ErrorX) {
	cur := &SourceNode{MName: name, MFullName: fullName, MTotalNum: 0, MCoverNum: 0, MChildren: make([]*SourceNode, 0),
		MChildrenMap: make(map[string]*SourceNode), MCodes: make([]string, 0), MEventIds: nil, MPassed: nil, MFailed: nil}
	if isSource {
		if file, err := os.Open(filePath); err != nil{
			return nil, perrorx.NewErrorXOpen(filePath, nil)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				cur.MCodes = append(cur.MCodes, scanner.Text())
			}
			cur.MTotalNum = len(cur.MCodes)
			cur.MEventIds = make([][]int, cur.MTotalNum)
			cur.MPassed = make([]int, cur.MTotalNum)
			cur.MFailed = make([]int, cur.MTotalNum)
			cur.MRanks = make([]int, cur.MTotalNum)
		}
		return cur, nil
	}
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil, perrorx.NewErrorXReadDir(filePath, err.Error(), nil)
	}
	for _, file := range files {
		if file.IsDir() {
			childName := file.Name()
			child, err := readSourceTreeDFS(path.Join(filePath, childName), childName,
				fullName+childName+".", false)
			if err != nil {
				return nil, perrorx.TransErrorX(err)
			}
			cur.MChildren = append(cur.MChildren, child)
			cur.MChildrenMap[childName] = child
			cur.MTotalNum += child.MTotalNum
		} else {
			childName := file.Name()
			// src文件下可能会放置一些开源协议, 没有.后缀, 需要特殊判断防止数组越界
			if strings.Contains(childName, ".") {
				childName = childName[0:strings.Index(childName, ".")]
				child, err := readSourceTreeDFS(path.Join(filePath, file.Name()), childName,
					fullName+childName, true)
				if err != nil {
					return nil, perrorx.TransErrorX(err)
				}
				cur.MChildren = append(cur.MChildren, child)
				cur.MChildrenMap[childName] = child
				cur.MTotalNum += child.MTotalNum
			}
		}
	}
	return cur, nil
}

// 根据projectId/src目录下的源文件构造源码树
func ReadSourceTree(projectId string) (*SourceTree, *perrorx.ErrorX) {
	srcPath := pdba.DBURLProjectIdSRC(projectId)
	if !putils.FileExist(srcPath) {
		return nil, perrorx.NewErrorXFileNotFound(srcPath, nil)
	}
	root, err := readSourceTreeDFS(srcPath, "", "", false)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	root.MFullName = projectId
	root.MName = projectId
	clsSrcMap, err := readClsSrcMap(projectId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}

	// 创建新的GSourceTree
	sourceTree := &SourceTree{root, 0, 0, clsSrcMap}
	SetGSourceTree(sourceTree)
	return sourceTree, nil
}

// 根据projectId/testcases/caseName/cover.json的内容进行GSourceTree覆盖
// 注意eventId只在覆盖显示时生效, 分析时需要为空
func CoverSourceTreeStd(sourceTree *SourceTree, projectId string, caseName string, eventId string) *perrorx.ErrorX {
	coverPath := pdba.DBURLProjectIdTestcaseCover(projectId, caseName)
	err := CoverSourceTree(sourceTree, coverPath, eventId)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	return nil
}

// 根据path的内容进行GSourceTree覆盖
// 注意eventId只在覆盖显示时生效, 分析时需要为空
func CoverSourceTree(sourceTree *SourceTree, path string, eventId string) *perrorx.ErrorX {
	coverage, err := pcoverage.ReadCoverage(path)
	if err != nil {
		return perrorx.TransErrorX(err)
	}
	if coverage.MStatus == "true" {
		sourceTree.MTotalPassed += 1
	} else {
		sourceTree.MTotalFailed += 1
	}
	for _, coverClass := range coverage.MCoverClasses {
		err := sourceTree.addCover(coverClass.MClassName, coverClass.MLines, coverClass.MEventIds, coverage.MStatus, eventId)
		if err != nil {
			return perrorx.TransErrorX(err)
		}
	}
	return nil
}

// 递归计算SourceNode的MCoverNum
// 这个函数不整合到CoverSourceTree是考虑到效率问题, CoverSourceTree可能会被多次调用, 但只需要计算一次CalCoverNumDFS
func CalCoverNumDFS(cur *SourceNode) int {
	if len(cur.MChildren) == 0 {
		return cur.MCoverNum
	}
	cur.MCoverNum = 0
	for _, nex := range cur.MChildren {
		cur.MCoverNum += CalCoverNumDFS(nex)
	}
	return cur.MCoverNum
}

// 根据projectId/src目录下的源文件构造源码树, 并根据database/projectId/testcases/caseName/cover.json的内容进行源码树覆盖
func ReadSourceTreeAndCover(projectId string, caseName string, eventId string) (*SourceTree, *perrorx.ErrorX) {
	sourceTree, err := ReadSourceTree(projectId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	err = CoverSourceTreeStd(sourceTree, projectId, caseName, eventId)
	if err != nil {
		return nil, perrorx.TransErrorX(err)
	}
	CalCoverNumDFS(sourceTree.MRoot)
	return sourceTree, nil
}