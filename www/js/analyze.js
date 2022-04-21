let SelectedAnalysisNode = null;

async function analyzeJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("analyzePrefix").value;

    if (projectId === "" || casePrefix === "") {
        alert("projectId or caseName empty");
        return
    }

    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;

    if (!acquireLock()) return;

    await httpGet("analyze", "projectId="+projectId+"&casePrefix="+casePrefix+
        "&analyzeFile="+analyzeFile+"&factor="+factor, function(res) {
        let temp = JSON.parse(res);
        drawSourceTree(temp["sourceTree"]);
        let susSlice = temp["susCodesSlice"];
        if (susSlice != null) {
            let ulObj = document.createElement("ul");
            ulObj.setAttribute("class", "treeUl")
            for (let i = 0; i < susSlice.length; i++) {
                let susCode = susSlice[i];
                let liObj = document.createElement("li");
                ulObj.appendChild(liObj);
                liObj.setAttribute("class", "treeLi")
                liObj.innerHTML = (i+1)+"["+susCode["rank"]+"]"+"("+susCode["a11"]+"/"+susCode["a10"]+"/"+
                    susCode["a01"]+"/"+susCode["a00"]+"/"+susCode["value"]+")"+
                    susCode["classShortName"]+":"+susCode["line"]+"["+susCode["className"]+"]";

                liObj.onclick = function() {
                    let tempNode = SRCTree.getSourceNode(susCode["className"]);
                    if (tempNode === null || tempNode.totalNum < susCode["line"]) {
                        alert("未找到 "+susCode["className"]+":"+susCode["line"]);
                        return;
                    }
                    if (SelectedAnalysisNode !== null) {
                        SelectedAnalysisNode.setAttribute("class", "treeLi"); // 清除之前的选择状态
                    }
                    SelectedAnalysisNode = this;
                    SelectedAnalysisNode.setAttribute("class", "treeLiSelected")
                    setSelectedSourceNode(tempNode, true, susCode["line"]);
                }
            }
            let analysisListView = document.getElementById("analysisListView");
            analysisListView.innerHTML = "";
            analysisListView.appendChild(ulObj);
        }
    });

    releaseLock();
}

async function combineAnalyzeJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("analyzePrefix").value;

    if (projectId === "" || casePrefix === "") {
        alert("projectId or caseName empty");
        return
    }

    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;

    if (!acquireLock()) return;

    if (window.confirm("是否开启组合分析?")) {
        await httpGet("combineanalyze", "projectId=" + projectId + "&casePrefix=" + casePrefix +
            "&analyzeFile=" + analyzeFile + "&factor=" + factor, function (res) {
            alert("success: " + res);
        });
    }

    releaseLock();
}

async function rdAnalyzeJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("analyzePrefix").value;
    // 借用一下testNum
    let randNum = document.getElementById("testNum").value;
    if (projectId === "" || casePrefix === "" || randNum === "") {
        alert("projectId or casePrefix or randNum empty");
        return
    }
    let tester = document.getElementById("tester").value;
    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;

    if (!acquireLock()) return;

    if (window.confirm("是否开启随机采样分析?")) {
        await httpGet("rdanalyze", "projectId=" + projectId + "&casePrefix=" + casePrefix +
            "&randNum=" + randNum + "&tester=" + tester +
            "&analyzeFile=" + analyzeFile + "&factor=" + factor, function (res) {
            alert("success: " + res);
        });
    }

    releaseLock();
}