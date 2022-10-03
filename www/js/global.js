// http请求前缀
let Prefix = "http://localhost:8082/"

// 页面同一时间只能有一个操作
let OpLock = false;

// 加锁
function acquireLock() {
    if (OpLock) {
        alert("请等待上次操作执行完毕");
        return false;
    }
    OpLock = true;
    document.getElementById("head").innerHTML = "Waiting...";
    return true;
}

// 解锁
function releaseLock() {
    OpLock = false;
    document.getElementById("head").innerHTML = "AUIAuto";
}

// 通用的http get请求, 请求参数自动携带avd, status ok时回调func, 失败时弹窗提示, 用户可以调用await等待请求处理完毕
// @param url: 请求url, 如uitree
// @param param: 请求参数, 如avd=emulator-5554
// @param func: status ok时回调func
async function httpGet(url, param, func) {
    return new Promise(function(resolve, reject) {
        let avd = "";
        if (document.getElementById("avd") != null) {
            avd = document.getElementById("avd").value;
        }
        console.log("GET", Prefix+url+"?avd="+avd+"&"+param)
        let xhr = new XMLHttpRequest();
        xhr.open('GET', Prefix+url+"?avd="+avd+"&"+param, true);
        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                let res = xhr.responseText;
                if(xhr.status === 200) {
                    func(res);
                } else {
                    alert("GET "+ url + " " + param + " error: " + res);
                }
                resolve(0);
            } // end of if (xhr.readyState === 4)
        } // end of xhr.onreadystatechange
        xhr.send();
    });
}

// 通用的http post请求, 请求参数自动携带avd, status ok时回调func, 失败时弹窗提示, 用户可以调用await等待请求处理完毕
// @param url: 请求url, 如edit
// @param param: 请求参数, 如avd=emulator-5554
// @param json: json字符串
// @param func: status ok时回调func
async function httpPost(url, param, json, func) {
    // console.log(json)
    return new Promise(function(resolve, reject) {
        let avd = "";
        if (document.getElementById("avd") != null) {
            avd = document.getElementById("avd").value;
        }
        console.log("POST", Prefix+url+"?avd="+avd+"&"+param+"\nbody\n"+json)
        let xhr = new XMLHttpRequest();
        xhr.open("POST", Prefix+url+"?avd="+avd+"&"+param, true);
        xhr.setRequestHeader('content-type', 'application/json');
        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                let res = xhr.responseText;
                if(xhr.status === 200) {
                    func(res);
                } else {
                    alert("POST "+ url + " " + param + " error: " + res);
                }
                resolve(0);
            } // end of if (xhr.readyState === 4)
        }
        xhr.send(json);
    });
}