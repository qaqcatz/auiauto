let Events = []

function addEvent(event) {
    let eventListTableBody = $('#eventListTableBody');
    event.id = Events.length+1
    eventListTableBody.append("<tr id='eventListTableBody"+
        eventListTableBody.children().length+ "'><td>"
        +event["id"]+"</td><td>"
        +event["type"]+"</td><td>"
        +event["value"]+"</td><td>"
        +event["desc"]+"</td></tr>");
    Events.push(event);
    let eventList = document.getElementById("eventList");
    eventList.scrollTop = eventList.scrollHeight;
}

async function saveJ() {
    if (Events.length === 0) {
        alert("empty!")
        return;
    }
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("saveCaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    let initSnapshot = document.getElementById("initSnapshot").value;
    let initTestcase = document.getElementById("initTestcase").value;
    let reInstall = document.getElementById("isReInstall").checked;
    let reStart = document.getElementById("isReStart").checked;
    if (window.confirm("确认save testcase:"+projectId+"/"+caseName)){
        await httpPost("savecase", "projectId=" + projectId + "&caseName=" + caseName,
            JSON.stringify({
                initSnapshot: initSnapshot,
                initTestcase: initTestcase,
                reInstall: reInstall,
                reStart: reStart,
                events: Events
            }), function (res) {
                alert("上传成功");
                document.getElementById("eventListTableBody").innerHTML = "";
                while (Events.length !== 0) {
                    Events.pop();
                }
            });
    }

    releaseLock();
}

async function loadJ() {
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("testcaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    await httpGet("loadcase", "projectId="+projectId+"&caseName="+caseName,  function(res) {
        let obj = JSON.parse(res);
        let events = obj["events"];
        document.getElementById("eventListTableBody").innerHTML = "";
        Events = [];
        for (let i = 0; i < events.length; i++) {
            addEvent(events[i]);
        }

        document.getElementById("initSnapshot").value = obj["initSnapshot"];
        document.getElementById("initTestcase").value = obj["initTestcase"];
        document.getElementById("isReInstall").checked = obj["reInstall"];
        document.getElementById("isReStart").checked = obj["reStart"];
        alert("load successful");
    });

    releaseLock();
}

function rmJ() {
    if (Events.length === 0) {
        alert("empty!")
        return;
    }
    Events.pop();
    $("#eventListTableBody tr:last").remove();
}