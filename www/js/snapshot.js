async function snapshotJ() {
    if (!acquireLock()) return;
    await httpGet("snapshot", "", function(res) {
        let snapshots = JSON.parse(res);
        document.getElementById("loadSnapshotName").innerHTML = "";
        let snapshotName = $('#loadSnapshotName');
        if (snapshots != null) {
            for (let i = 0; i < snapshots.length; i++) {
                snapshotName.append("<option value='"+snapshots[i]+"'>"+snapshots[i]+"</option>");
            }
        }
    });
    releaseLock();
}

async function loadSnapshotJ() {
    if (!acquireLock()) return;
    let loadSnapshotName = document.getElementById("loadSnapshotName").value;
    if (window.confirm("确认加载快照:"+loadSnapshotName)) {
        await httpGet("loadsnapshot", "name="+loadSnapshotName, function(res) {
            alert("successful: " + res);
        });
    }
    releaseLock();
}

async function deleteSnapshotJ() {
    if (!acquireLock()) return;
    let deleteSnapshotName = document.getElementById("loadSnapshotName").value;
    if (window.confirm("确认删除快照:"+deleteSnapshotName)) {
        await httpGet("deletesnapshot", "name=" + deleteSnapshotName, function (res) {
            alert("successful: " + res);
        });
    }
    releaseLock();
}

async function saveSnapshotJ() {
    if (!acquireLock()) return;
    let saveSnapshotName = document.getElementById("saveSnapshotName").value;
    if (window.confirm("确认保存快照:"+saveSnapshotName)) {
        await httpGet("savesnapshot", "name=" + saveSnapshotName, function (res) {
            alert("successful: " + res);
        });
    }
    releaseLock();
}