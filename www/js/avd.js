// 记录当前选中的ui node
let SelectedNode = null;

// 选中ui tree上的节点, 并做相应的矩形显示和ui tree上的跳转, 以及显示节点细节
function setSelectedNode(node, jump) {
    if (SelectedNode !== null) {
        SelectedNode.htmlNode.setAttribute("class", "treeLi"); // 清除之前的选择状态
    }
    SelectedNode = node;
    if (node != null) syncScreen();
    if (node != null) syncUiTree(jump);
    syncTable();
}

// 设置SelectedNode的同时在画布上显示矩形框
function syncScreen() {
    let uiNode = SelectedNode;
    let cvr = document.getElementById("screenshotRect");
    let cxtr = cvr.getContext("2d");
    cxtr.clearRect(0, 0, cvr.width, cvr.height)
    cxtr.lineWidth=10;
    cxtr.strokeStyle="#000000";
    cxtr.strokeRect(uiNode.left,uiNode.top,uiNode.right-uiNode.left,uiNode.bottom-uiNode.top);
}

// 设置SelectedNode的同时在ui tree做跳转
function syncUiTree(jump) {
    SelectedNode.htmlNode.setAttribute("class", "treeLiSelected")
    if (jump) {
        document.getElementById('uiTree').scrollTop=
            $(Root.htmlNode).height() * ((SelectedNode.nodeId) / Root.nodeNum);
    }
}

// 设置SelectedNode的同时显示节点细节
function syncTable() {
    let uiNode = SelectedNode;
    document.getElementById('myDepth').innerHTML = uiNode === null ? "" : uiNode.dp;
    document.getElementById('myIndex').innerHTML = uiNode === null ? "" : uiNode.idx;
    document.getElementById('myBounds').innerHTML = uiNode === null ? "" : uiNode.bds;
    document.getElementById('myPackage').innerHTML = uiNode === null ? "" : uiNode.pkg;
    document.getElementById('myClass').innerHTML = uiNode === null ? "" : uiNode.cls;
    document.getElementById('myResource-id').innerHTML = uiNode === null ? "" : uiNode.res;
    document.getElementById('myContentDesc').innerHTML = uiNode === null ? "" : uiNode.dsc;
    document.getElementById('myText').innerHTML = uiNode === null ? "" : uiNode.txt;
    document.getElementById('myClickable').innerHTML = uiNode === null ? "" : uiNode.clickable;
    document.getElementById('myLongClickable').innerHTML = uiNode === null ? "" : uiNode.longClickable;
    document.getElementById('myEditable').innerHTML = uiNode === null ? "" : uiNode.editable;
    document.getElementById('myScrollable').innerHTML = uiNode === null ? "" : uiNode.scrollable;
    document.getElementById('myCheckable').innerHTML = uiNode === null ? "" : uiNode.checkable;
    document.getElementById('myChecked').innerHTML = uiNode === null ? "" : uiNode.checked;
}

// 上一个搜索对象
let PreSearch = "";
// 搜索结果
let SearchAns = [];
// 当前展示的搜索结果在SearchAns中的下标
let SearchIndex = 0;

// 根据key:value查询
function searchJ(e) {
    let evt = window.event || e;
    if (evt.keyCode === 13) {
        if (Root !== null) {
            let curSearch = document.getElementById("searchContent").value;
            let kv = curSearch.split("@");
            if (kv.length !== 2) {
                alert("Search format error");
                return;
            }
            if (curSearch !== PreSearch) {
                PreSearch = curSearch;
                SearchAns = [];
                let key = kv[0];
                let value = kv[1];
                Root.foreach(function(node) {
                    if ((key === "depth" && value === node.dp.toString()) ||
                        (key === "index" && value === node.idx.toString()) ||
                        (key === "package" && value === node.pkg) ||
                        (key === "class" && value === node.cls) ||
                        (key === "resource-id" && value === node.res) ||
                        (key === "contentDesc" && value === node.dsc) ||
                        (key === "text" && value === node.txt) ||
                        (key === "clickable" && value === node.clickable.toString()) ||
                        (key === "longClickable" && value === node.longClickable.toString()) ||
                        (key === "editable" && value === node.editable.toString()) ||
                        (key === "scrollable" && value === node.scrollable.toString()) ||
                        (key === "checkable" && value === node.checkable.toString()) ||
                        (key === "checked" && value === node.checked.toString())) {
                        SearchAns.push(node);
                    }
                });
                SearchIndex = 0;
            }
            if (SearchAns.length === 0) {
                alert("no such node!");
                return;
            }
            setSelectedNode(SearchAns[SearchIndex], true);
            SearchIndex = (SearchIndex+1)%SearchAns.length;
        } // if (Root !== null)
    } // if (evt.keyCode === 13)
}

// 获取屏幕截图/ui tree
async function dumpJ() {
    if (!acquireLock()) return;
    // dump时初始化ui tree和之前的选择搜索状态,
    // 屏幕截图的清理交给相应的绘制函数
    cleanUiTree();
    setSelectedNode(null, false);
    PreSearch = "";
    SearchAns = [];
    SearchIndex = 0;
    await httpGet("screenshot", "", function(res) {
        drawScreen(res);
    });
    await httpGet("uitree", "", function(res) {
        drawUiTree(res);
    });
    releaseLock();
}

async function devJ() {
    if (!acquireLock()) return;

    await httpGet("devices", "", function (res) {
        let devices = JSON.parse(res);
        document.getElementById("avd").innerHTML = "";
        let avd = $('#avd');
        if (devices != null) {
            for (let i = 0; i < devices.length; i++) {
                avd.append("<option value='"+devices[i]+"'>"+devices[i]+"</option>");
            }
        }
    });

    releaseLock();
}

// init, 应用崩溃时会上传错误日志, 禁止后续ui tree的获取, 这种情况需要用户手动调用init
async function InitJ() {
    if (!acquireLock()) return;

    await httpGet("init", "", function (res) {
    });

    releaseLock();
}

// 安装/启动 antrance
async function antranceJ() {
    if (!acquireLock()) return;

    await httpGet("installstart", "projectId=xmu.wrxlab.antrance", function (res) {
        alert("antrance start!");
    });

    releaseLock();
}
