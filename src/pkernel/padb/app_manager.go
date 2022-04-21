package padb

import (
	"auiauto/pconfig"
	"auiauto/pdba"
	"auiauto/perrorx"
	"auiauto/putils"
)

// 安装app
func AdbInstall(avd string, apkPath string) (string, *perrorx.ErrorX) {
	_, err_ := AdbReconnectOffline()
	if err_ != nil {
		return "", perrorx.TransErrorX(err_)
	}
	_, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " install -t " + apkPath)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return "install " + apkPath, nil
}

// 启动应用进程
func AdbStart(avd string, applicationId string, mainActivity string) (string, *perrorx.ErrorX) {
	ans, err := AdbShell(avd, "am start "+applicationId+"/"+mainActivity)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 强制杀死应用进程
func AdbKill(avd string, applicationId string) (string, *perrorx.ErrorX) {
	ans, err := AdbShell(avd, " am force-stop "+applicationId)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

// 卸载app
func AdbUnInstall(avd string, applicationId string) (string, *perrorx.ErrorX) {
	_, err_ := AdbReconnectOffline()
	if err_ != nil {
		return "", perrorx.TransErrorX(err_)
	}
	_, err := putils.Sh(pconfig.GConfig.MAdb + " -s " + avd + " uninstall " + applicationId)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return "uninstall " + applicationId, nil
}

// 安装启动apk, 根据project的值分类操作:
// 1. xmu.wrxlab.antrance, 首先执行kill卸载操作(不捕获异常), 然后从kernel目录安装, setting启动
// 2. 数据库的其他应用, 首先读取app的applicaitonId和mainActivity, 执行kill卸载操作(不捕获异常), 之后重新安装启动
// 3. 未知应用, 报错
func AdbInstallAndStart(avd string, projectId string) (string, *perrorx.ErrorX) {
	if projectId == "xmu.wrxlab.antrance" {
		_, _ = AdbKill(avd, projectId)
		_, _ = AdbUnInstall(avd, projectId)
		_, err := AdbInstall(avd, pdba.DBURLKernelAntranceAPK())
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		// enable
		ans, err := AdbShell(avd, "settings put secure enabled_accessibility_services xmu.wrxlab.antrance/.Antrance")
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		return ans, nil
	} else if putils.FileExist(pdba.DBURLProjectId(projectId)) {
		applicationId, mainActivity, err := GetAppConfig(projectId)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		_, _ = AdbKill(avd, applicationId)
		_, _ = AdbUnInstall(avd, applicationId)
		_, err = AdbInstall(avd, pdba.DBURLProjectIdAPKAPP(projectId))
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		ans, err := AdbStart(avd, applicationId, mainActivity)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		return ans, nil
	} else {
		return "", perrorx.NewErrorXFileNotFound(projectId+"apk/app.apk", nil)
	}
}

// 只适用于数据库中存在的应用, 只进行重新安装操作
func AdbOnlyInstall(avd string, projectId string) (string, *perrorx.ErrorX) {
	if putils.FileExist(pdba.DBURLProjectId(projectId)) {
		applicationId, _, err := GetAppConfig(projectId)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		_, _ = AdbKill(avd, applicationId)
		_, _ = AdbUnInstall(avd, applicationId)
		ans, err := AdbInstall(avd, pdba.DBURLProjectIdAPKAPP(projectId))
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		return ans, nil
	} else {
		return "", perrorx.NewErrorXFileNotFound(projectId+"apk/app.apk", nil)
	}
}

// 只适用于数据库中存在的应用, 只进行重新启动操作
func AdbOnlyStart(avd string, projectId string) (string, *perrorx.ErrorX) {
	if putils.FileExist(pdba.DBURLProjectId(projectId)) {
		applicationId, mainActivity, err := GetAppConfig(projectId)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		_, _ = AdbKill(avd, applicationId)
		ans, err := AdbStart(avd, applicationId, mainActivity)
		if err != nil {
			return "", perrorx.TransErrorX(err)
		}
		return ans, nil
	} else {
		return "", perrorx.NewErrorXFileNotFound(projectId+"apk/app.apk", nil)
	}
}

// 获取当前顶层activity信息, adb shell dumpsys window | grep  mCurrentFocus
func AdbGetTopActivity(avd string) (string, *perrorx.ErrorX) {
	ans, err := AdbShell(avd, "dumpsys window | grep  mCurrentFocus")
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
