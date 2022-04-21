async function projectJ() {
    if (!acquireLock()) return;
    await httpGet("projects", "", function(res) {
        let projects = JSON.parse(res);
        document.getElementById("projectName").innerHTML = "";
        let projectName = $('#projectName');
        if (projects != null) {
            for (let i = 0; i < projects.length; i++) {
                projectName.append("<option value='"+projects[i]+"'>"+projects[i]+"</option>");
            }
        }
    });
    releaseLock();
}

async function testcaseJ() {
    if (!acquireLock()) return;
    let project = document.getElementById("projectName").value;
    await httpGet("testcases", "project="+project, function(res) {
        let testcases = JSON.parse(res);
        document.getElementById("testcaseName").innerHTML = "";
        let testcaseName = $('#testcaseName');
        if (testcases != null) {
            for (let i = 0; i < testcases.length; i++) {
                testcaseName.append("<option value='"+testcases[i]+"'>"+testcases[i]+"</option>");
            }
        }
    });
    releaseLock();
}