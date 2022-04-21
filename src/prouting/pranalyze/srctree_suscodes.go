package pranalyze

import (
	"auiauto/pkernel/psrctree"
	"auiauto/pkernel/psuscode"
)

// 将返回给前端的结果简单打个包
type SourceTreeAndSusCodes struct {
	// 需要显示的源码树
	MSourceTree    *psrctree.SourceTree `json:"sourceTree"`
	// 需要展示的可以语句切片
	MSusCodesSlice psuscode.SusCodes    `json:"susCodesSlice"`
}
