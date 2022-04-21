async function testingJ() {
    let projectId = document.getElementById("projectName").value;
    let crashCase = document.getElementById("testcaseName").value;
    let tester = document.getElementById("tester").value;
    let testNum = document.getElementById("testNum").value;
    let testPrefix = document.getElementById("testPrefix").value;
    let testParam = document.getElementById("testParam").value;
    if (projectId === "" || crashCase === "" || tester === "" || testNum === "" || testPrefix === "" || testParam === "") {
        alert("check your testing configuration");
        return
    }
    if (!acquireLock()) return;
    if (window.confirm("确认开启测试: projectId="+projectId+", crashCase="+crashCase+", tester="+tester+
        ", testNum="+testNum+", testPrefix="+testPrefix+", testParam="+testParam)) {
        await httpGet("testing", "projectId="+projectId+"&crashCase="+crashCase
            +"&tester="+tester+"&testNum="+testNum+"&testPrefix="+testPrefix+"&testParam="+testParam, function(res) {
            alert("测试结束:" + res);
        });
    }
    releaseLock();
}
