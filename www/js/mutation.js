async function mutationJ() {
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("testcaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    await httpPost("mutation", "projectId="+projectId+"&caseName="+caseName,
        "", function(res) {
            alert("变异成功!");
        })

    releaseLock();
}