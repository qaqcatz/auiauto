class SourceNode {
    constructor(prefix, name, fullName, totalNum, coverNum, nodeId) {
        this.prefix = prefix;
        this.name = name;
        this.fullName = fullName;
        this.totalNum = totalNum;
        this.coverNum = coverNum;
        this.children = [];
        this.nodeId = nodeId;
        this.htmlNode = null;
    }

    addChild(child) {
        this.children.push(child)
    }

    toShortString(v="|", h="-") {
        let ans = "";
        ans = this.prefix;
        ans = ans.replace(/\|/g, v);
        ans = ans.replace(/-/g, h);
        ans += this.name+"("+this.coverNum+"/"+this.totalNum+")"
        return ans;
    }
}

class SourceTree {
    constructor() {
        this.nodeNum = 0;
        this.root = null;
        this.htmlNode = null;
    }

    setRoot(root) {
        this.root = root;
    }

    foreach(func) {
        this.foreachDFS(func, this.root);
    }

    foreachDFS(func, cur) {
        func(cur);
        for (let i = 0; i < cur.children.length; i++) {
            this.foreachDFS(func, cur.children[i]);
        }
    }

    debug() {
        this.foreach(function(cur) {
            console.log(cur.toShortString("|", " "));
        })
    }

    bindHtml() {
        let ulObj = document.createElement("ul");
        this.htmlNode = ulObj
        ulObj.setAttribute("class", "treeUl")
        this.foreach(function(cur) {
            let liObj = document.createElement("li");
            ulObj.appendChild(liObj)
            cur.htmlNode = liObj
            liObj.setAttribute("class", "treeLi")
            liObj.innerHTML = cur.toShortString("|", "&ensp;");
            if (cur.children.length === 0 && cur.coverNum !== 0) {
                liObj.onclick = function() {
                    setSelectedSourceNode(cur, false).then(r => {});
                }
            } else if (cur.coverNum === 0) {
                liObj.setAttribute("class", "treeLiUnSelected");
            } else {
                liObj.setAttribute("class", "treeLiPackage");
            }
        })
    }

    getSourceNode(className) {
        let ans = null;
        this.foreach(function(cur) {
            if (cur.children.length === 0 && cur.fullName === className) {
                ans = cur;
            }
        })
        return ans;
    }
}

let SRCTree = null;

// 上一个搜索对象
let SourcePreSearch = "";
// 搜索结果
let SourceSearchAns = [];
// 当前展示的搜索结果在SearchAns中的下标
let SourceSearchIndex = 0;

// 根据key:value查询
function sourceSearchJ(e) {
    let evt = window.event || e;
    if (evt.keyCode === 13) {
        if (SRCTree !== null) {
            let curSearch = document.getElementById("sourceSearchContent").value;
            if (curSearch !== SourcePreSearch) {
                SourcePreSearch = curSearch;
                SourceSearchAns = [];
                SRCTree.foreach(function(srcNode){
                    if (curSearch === "cover") {
                        if (srcNode.children.length === 0 && srcNode.coverNum !== 0) {
                            SourceSearchAns.push(srcNode);
                        }
                    } else {
                        if (srcNode.fullName.endsWith(curSearch)) {
                            SourceSearchAns.push(srcNode);
                        }
                    }
                });
                SourceSearchIndex = 0;
            }
            if (SourceSearchAns.length === 0) {
                alert("no such node!");
                return;
            }
            setSelectedSourceNode(SourceSearchAns[SourceSearchIndex], true).then(r => {});
            SourceSearchIndex = (SourceSearchIndex+1)%SourceSearchAns.length;
        } // if (SRCTree !== null)
    } // if (evt.keyCode === 13)
}

function drawSourceTreeDFS(prefix, jsonNode) {
    let nodeId = SRCTree.nodeNum++;
    let sourceNode = new SourceNode(prefix, jsonNode["name"], jsonNode["fullName"],
        jsonNode["totalNum"], jsonNode["coverNum"], nodeId);
    let jsonChildren = jsonNode["children"];
    for (let i = 0; i < jsonChildren.length; i++) {
        sourceNode.addChild(drawSourceTreeDFS(prefix+"|--", jsonChildren[i]))
    }
    return sourceNode;
}

function drawSourceTree(jsonTree) {
    SRCTree = new SourceTree();
    let root = drawSourceTreeDFS("", jsonTree["root"]);
    SRCTree.setRoot(root);
    SRCTree.bindHtml();
    let sourceTreeView = document.getElementById("sourceTreeView");
    sourceTreeView.innerHTML = "";
    sourceTreeView.appendChild(SRCTree.htmlNode);
    SourcePreSearch = "";
    SourceSearchAns = [];
    SourceSearchIndex = 0;
    document.getElementById('sourceListView').innerHTML = ""
}

let SelectedSourceNode = null;

async function setSelectedSourceNode(sourceNode, jump, susLine = 0) {
    if (SelectedSourceNode !== null) {
        SelectedSourceNode.htmlNode.setAttribute("class", "treeLi"); // 清除之前的选择状态
    }
    SelectedSourceNode = sourceNode;
    if (SelectedSourceNode != null) {
        SelectedSourceNode.htmlNode.setAttribute("class", "treeLiSelected")
        if (jump) {
            document.getElementById('sourceTree').scrollTop=
                $(SRCTree.htmlNode).height() * ((SelectedSourceNode.nodeId) / SRCTree.nodeNum);
        }
        if (SelectedSourceNode.children.length !== 0 || SelectedSourceNode.coverNum === 0) {
            return;
        }
        await httpGet("ccl", "projectId="+SRCTree.root.name+"&dotClassPath="+SelectedSourceNode.fullName, function(res) {
            let temp = JSON.parse(res);
            if (temp != null) {
                let codes = temp["codes"];
                let codesType = temp["codesType"];

                let ulObj = document.createElement("ul");
                ulObj.setAttribute("class", "treeUl")
                for (let i = 0; i < codes.length; i++) {
                    let liObj = document.createElement("li");
                    ulObj.appendChild(liObj);
                    liObj.setAttribute("class", "treeLi")
                    liObj.innerHTML = "<pre style='display: inline'>"+formatInt61(i+1)+codes[i]+"</pre>";
                    if (0 <= codesType[i] && codesType[i] <= 4) {
                        liObj.setAttribute("class", "treeLiCovered"+codesType[i]);
                    }
                }

                let sourceListView = document.getElementById("sourceListView");
                sourceListView.innerHTML = "";
                sourceListView.appendChild(ulObj);

                if (codes.length !== 0 && 1 <= susLine && susLine <= codes.length) {
                    document.getElementById('sourceList').scrollTop=
                        $(ulObj).height() * ((susLine-10 <= 0 ? 0 : susLine-10) / codes.length);
                }

                document.getElementById("sourceHead").innerHTML = SelectedSourceNode.fullName
            }
        });
    }
}

function formatInt61(i) {
    let strI = ""+i;
    if (strI.length <= 6) {
        let blank = "";
        for (let i = 0; i < 6-strI.length; i++) {
            blank += " ";
        }
        return blank+strI+"|";
    }
    return strI.substring(strI.length-6, strI.length)+"|";
}