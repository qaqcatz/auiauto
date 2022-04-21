class UiNode {
    constructor(prefix, dp, idx, bds, pkg, cls, res, dsc, txt, op, sta, nodeId) {
        this.preifx = prefix;
        this.dp = dp;
        this.idx = idx;
        let sp = bds.split("@");
        this.bds = bds
        this.left = parseInt(sp[0]);
        this.top = parseInt(sp[1]);
        this.right = parseInt(sp[2]);
        this.bottom = parseInt(sp[3]);
        this.pkg = pkg;
        this.cls = cls;
        this.res = res;
        this.dsc = dsc;
        this.txt = txt;
        this.clickable = false;
        this.longClickable = false;
        this.editable = false;
        this.scrollable = false;
        this.checkable = false;
        this.op = op;
        if ((op & 1) !== 0) {
            this.clickable = true;
        }
        if ((op & 2) !== 0) {
            this.longClickable = true;
        }
        if ((op & 4) !== 0) {
            this.editable = true;
        }
        if ((op & 8) !== 0) {
            this.scrollable = true;
        }
        if ((op & 16) !== 0) {
            this.checkable = true;
        }
        this.sta = sta;
        this.checked = false;
        if ((sta & 1) !== 0) {
            this.checked = true;
        }
        this.nodeId = nodeId
        this.children = [];
        this.father = null;
        this.htmlNode = null;
    }

    addChild(child) {
        this.children.push(child);
        child.father = this;
    }

    toShortString(showPrefix, v="|", h="-") {
        let ans = "";
        if (showPrefix) {
            ans = this.preifx;
            ans = ans.replace(/\|/g, v);
            ans = ans.replace(/-/g, h);
        }
        ans += this.dp + "@" + this.idx + "["+this.left+","+this.top+","+this.right+","+this.bottom+"]";
        if (this.res === "") ans += this.cls;
        else ans += this.res;
        if (this.clickable) ans += "@c"
        if (this.longClickable) ans += "@lc"
        if (this.editable) ans += "@e"
        if (this.scrollable) ans += "@s"
        if (this.checkable) ans += "@ck"
        return ans;
    }

    contain(x, y) {
        return this.left <= x && x <= this.right && this.top <= y && y <= this.bottom;
    }

    centerX() {
        return (this.right+this.left)/2;
    }

    centerY() {
        return (this.bottom+this.top)/2;
    }

    areaS() {
        return (this.right-this.left)*(this.bottom-this.top)
    }

    isOperable() {
        return this.op !== 0;
    }

    eventObjectCode() {
        return this.pkg+"@"+this.cls+"@"+this.res+"@"+this.op+"@"+this.dp;
    }

    eventObjectPrefix() {
        let ans = [];
        let cur = this;
        while (cur != null) {
            ans.push(cur.idx);
            cur = cur.father;
        }
        // reverse
        for (let i = 0, j = ans.length-1; i < j; i++, j--) {
            let t = ans[i];
            ans[i] = ans[j];
            ans[j] = t;
        }
        return ans;
    }
}

class UiTree {
    constructor() {
        this.nodeNum = 0;
        this.children = [];
        this.htmlNode = null;
    }

    addChild(child) {
        this.children.push(child);
    }

    foreach(func) {
        for (let i = 0; i < this.children.length; i++) {
            this.foreachDFS(func, this.children[i]);
        }
    }

    foreachDFS(func, cur) {
        func(cur);
        for (let i = 0; i < cur.children.length; i++) {
            this.foreachDFS(func, cur.children[i]);
        }
    }

    debug() {
        this.foreach(function(cur) {
            console.log(cur.toShortString(true, "|", " "));
        })
    }

    BindHtml() {
        let ulObj = document.createElement("ul");
        this.htmlNode = ulObj
        ulObj.setAttribute("class", "treeUl")
        this.foreach(function(cur) {
            let liObj = document.createElement("li");
            ulObj.appendChild(liObj)
            cur.htmlNode = liObj
            liObj.setAttribute("class", "treeLi")
            liObj.innerHTML = cur.toShortString(true, "|", "&ensp;");
            liObj.onclick = function() {
                setSelectedNode(cur, false);
            }
        })
    }
}

// 当前ui tree的根节点
let Root = null;

function createUiTree(data) {
    Root = new UiTree();
    let parser = new DOMParser();
    let xmlDoc = parser.parseFromString(data,"text/xml");
    let rt = $(xmlDoc).find("rt");
    $(rt).children().each(function (i) {
        Root.addChild(createUiTreeDFS("", this));
    });
    Root.BindHtml();
}

function createUiTreeDFS(prefix, node) {
    let dp = parseInt($(node).attr("dp"));
    let idx = parseInt($(node).attr("idx"));
    let bds = $(node).attr("bds");
    let pkg = $(node).attr("pkg");
    let cls = $(node).attr("cls");
    let res = $(node).attr("res");
    let dsc = $(node).attr("dsc");
    let txt = $(node).attr("txt");
    let op = parseInt($(node).attr("op"));
    let sta = parseInt($(node).attr("sta"));
    let nodeId = Root.nodeNum++;

    let uiNode = new UiNode(prefix, dp, idx, bds, pkg, cls, res, dsc, txt, op, sta, nodeId);

    $(node).children().each(function (i) {
        uiNode.addChild(createUiTreeDFS(prefix+"|--", this));
    })

    return uiNode
}

// 绘制ui tree
function drawUiTree(data) {
    createUiTree(data);
    let treeView = document.getElementById("treeView");
    treeView.innerHTML = "";
    treeView.appendChild(Root.htmlNode);
}

// 清除ui tree, 网络请求失败时上次的ui tree还在的话会影响体验, 因此在dump发送请求前主动清空一次
function cleanUiTree() {
    Root = null;
    let treeView = document.getElementById("treeView");
    treeView.innerHTML = "";
}