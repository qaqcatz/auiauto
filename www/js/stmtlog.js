async function stmtlogJ() {
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("saveCaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    if (window.confirm("确认stmtlog:"+projectId+"/"+caseName)) {
        await httpGet("stmtlognow", "projectId=" + projectId + "&caseName=" + caseName, function (res) {
            alert("获取日志成功:" + res);
        });
    }

    releaseLock();
}