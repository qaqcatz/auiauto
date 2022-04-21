// 点击
async function clickJ() {
    let node = SelectedNode;
    if (node === null) {
        alert("you should select a node first!");
        return;
    }
    if (!acquireLock()) return;

    let desc = document.getElementById("eventDesc").value;
    let event = {id: 0, type: "click", value: "", object: node.eventObjectCode(), prefix: node.eventObjectPrefix(), desc: desc}
    await httpPost("perform", "", JSON.stringify(event), function(res) {
        addEvent(event);
        releaseLock();
        dumpJ();
    });

    releaseLock();
}

// 长按
async function longClickJ() {
    let node = SelectedNode;
    if (node === null) {
        alert("you should select a node first!");
        return;
    }
    if (!acquireLock()) return;

    let desc = document.getElementById("eventDesc").value;
    let event = {id: 0, type: "longclick", value: "", object: node.eventObjectCode(), prefix: node.eventObjectPrefix(), desc: desc}
    await httpPost("perform", "", JSON.stringify(event), function(res) {
        addEvent(event);
        releaseLock();
        dumpJ();
    });

    releaseLock();
}

// 滑动, 左上右下0123
async function scrollJ(flag) {
    let node = SelectedNode;
    if (node === null) {
        alert("you should select a node first!");
        return;
    }
    if (!acquireLock()) return;

    let desc = document.getElementById("eventDesc").value;
    let event = {id: 0, type: "scroll", value: ""+flag, object: node.eventObjectCode(), prefix: node.eventObjectPrefix(), desc: desc}
    await httpPost("perform", "", JSON.stringify(event), function(res) {
        addEvent(event);
        releaseLock();
        dumpJ();
    });

    releaseLock();
}

// back
async function backJ() {
    if (!acquireLock()) return;

    let desc = document.getElementById("eventDesc").value;
    let event = {id: 0, type: "keyevent", value: "KEYCODE_BACK", object: "", prefix: [], desc: desc}
    await httpPost("perform", "", JSON.stringify(event), function(res) {
        addEvent(event);
        releaseLock();
        dumpJ();
    });

    releaseLock();
}

// perform click/longclick/edit text/editx text/scroll 0,1,2,3/keyevent keycode/swipe x1 y1 x2 y2 ms/wait ms/rotate 0,1
async function performJ() {
    let options=$("#eventType option:selected");
    let eventType = options.val();
    let eventValue = document.getElementById("eventValue").value;
    let eventObject = "";
    let eventObjectPrefix = [];

    if (eventType === "click" || eventType === "dclick" || eventType === "longclick" || eventType === "edit" || eventType === "editx" ||
    eventType === "scroll" || eventType === "check") {
        let node = SelectedNode;
        if (node === null) {
            alert("you should select a node first!");
            return;
        }
        eventObject = node.eventObjectCode();
        eventObjectPrefix =  node.eventObjectPrefix();
    }
    if (!acquireLock()) return;

    let desc = document.getElementById("eventDesc").value;
    let event = {id: 0, type: eventType, value: eventValue, object: eventObject, prefix: eventObjectPrefix, desc: desc}
    // do not perform wait
    if (eventType === "wait") {
        addEvent(event);
    } else {
        await httpPost("perform", "", JSON.stringify(event), function(res) {
            addEvent(event);
            releaseLock();
            dumpJ();
        });
    }

    releaseLock();
}

async function performsJ() {
    let projectId = document.getElementById("projectName").value;
    let caseName = document.getElementById("testcaseName").value;

    if (projectId === "" || caseName === "") {
        alert("projectId or caseName empty");
        return
    }

    if (!acquireLock()) return;

    if (window.confirm("确认performs:"+projectId+"/"+caseName)){
        await httpPost("performs", "projectId=" + projectId + "&caseName=" + caseName,
            "", function (res) {
                alert("复现成功:" + res);
            })
    }

    releaseLock();
}

async function reInstallJ() {
    if (!acquireLock()) return;

    let project = document.getElementById("projectName").value;
    await httpGet("reinstall", "projectId="+project, function(res) {
        alert("successful: " + res);
    });

    releaseLock();
}

async function reStartJ() {
    if (!acquireLock()) return;

    let project = document.getElementById("projectName").value;
    await httpGet("restart", "projectId="+project, function(res) {
        alert("successful: " + res);
    });

    releaseLock();
}