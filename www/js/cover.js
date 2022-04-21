async function coverJ() {
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("testcaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    let coverEventId = document.getElementById("coverEventId").value;

    if (coverEventId !== "all") {
        await httpGet("cover", "projectId="+projectId+"&caseName="+caseName+"&eventId="+coverEventId, function(res) {
            drawSourceTree(JSON.parse(res));
            // document.getElementById("analysisListView").innerHTML = "";
        });
    } else {
        await httpGet("coverAll", "projectId="+projectId+"&caseName="+caseName, function(res) {
            let coverMsg = JSON.parse(res);
            let ans = "";
            if (coverMsg != null) {
                for (let i = 0; i < coverMsg.length; i++) {
                    ans += coverMsg[i] + " ";
                }
            }
            alert(ans);
            console.log(ans);
        });
    }



    releaseLock();
}
