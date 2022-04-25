async function staAnalyzeAllJ() {
    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;
    let tester = document.getElementById("tester").value;

    if (!acquireLock()) return;

    await httpGet("staanalyzeall", "analyzeFile="+analyzeFile+"&factor="+factor+"&tester="+tester, function(res) {
        let jsonData = JSON.parse(res);
        console.log(jsonData);
        document.getElementById("message").value =
            "artAvgMap="+jsonData["artAvgMap"]+";"+
            "twoAvgMap="+jsonData["twoAvgMap"]+";"+
            "combineBestMapAvgMap="+jsonData["combineBestMapAvgMap"]+";"+
            "rdBestAvgMapAvgMap="+jsonData["rdBestAvgMapAvgMap"]+";\n"+
            "artMrr="+jsonData["artMrr"]+";"+
            "twoMrr="+jsonData["twoMrr"]+";"+
            "combineBestMapMrr="+jsonData["combineBestMapMrr"]+";"+
            "rdBestAvgMidMapMrr="+jsonData["rdBestAvgMidMapMrr"]+";"+
            "rdBestAvgMaxMapMrr="+jsonData["rdBestAvgMaxMapMrr"]+";"+
            "rdBestAvg34MapMrr="+jsonData["rdBestAvg34MapMrr"]+";\n"+
            "artTop="+jsonData["artTop"]+";\n"+
            "twoTop="+jsonData["twoTop"]+";\n"+
            "combineTop="+jsonData["combineTop"]+";\n"+
            "rdBestAvgMidMapTop="+jsonData["rdBestAvgMidMapTop"]+";\n"+
            "rdBestAvgMaxMapTop="+jsonData["rdBestAvgMaxMapTop"]+";\n"+
            "rdBestAvg34MapTop="+jsonData["rdBestAvg34MapTop"];

        let staEachs = jsonData["staEachs"]

        let statisticTableC = "width: 100%;font-size: 12px;font-family:Arial,Helvetica,sans-serif;text-align: center;"
        let statisticTableThC = "color: white;border: 1px solid black;background-color: #333333;"
        document.getElementById("table").innerHTML = "" +
            "<table style=\""+statisticTableC+"\">"+
            "<thead><tr>"+
            "<th style=\""+statisticTableThC+"\">project</th>"+
            "<th style=\""+statisticTableThC+"\">rdBestAvgMap</th>"+
            "<th style=\""+statisticTableThC+"\">rdPassNum</th>"+
            "<th style=\""+statisticTableThC+"\">rdCrashNum</th>"+
            "<th style=\""+statisticTableThC+"\">rdMinMap</th>"+
            "<th style=\""+statisticTableThC+"\">rdMinRank</th>"+
            "<th style=\""+statisticTableThC+"\">rd14Map</th>"+
            "<th style=\""+statisticTableThC+"\">rd14Rank</th>"+
            "<th style=\""+statisticTableThC+"\">rdMidMap</th>"+
            "<th style=\""+statisticTableThC+"\">rdMidRank</th>"+
            "<th style=\""+statisticTableThC+"\">rd34Map</th>"+
            "<th style=\""+statisticTableThC+"\">rd34Rank</th>"+
            "<th style=\""+statisticTableThC+"\">rdMaxMap</th>"+
            "<th style=\""+statisticTableThC+"\">rdMaxRank</th>"+
            "<th style=\""+statisticTableThC+"\">artMap</th>"+
            "<th style=\""+statisticTableThC+"\">artRank</th>"+
            "<th style=\""+statisticTableThC+"\">twoMap</th>"+
            "<th style=\""+statisticTableThC+"\">twoRank</th>"+
            "<th style=\""+statisticTableThC+"\">combineBestMap</th>"+
            "<th style=\""+statisticTableThC+"\">pn</th>"+
            "<th style=\""+statisticTableThC+"\">fn</th>"+
            "<th style=\""+statisticTableThC+"\">combinePassNum</th>"+
            "<th style=\""+statisticTableThC+"\">combineCrashNum</th>"+
            "<th style=\""+statisticTableThC+"\">combineRank</th>"+
            "</tr></thead>"+
            "<tbody id=\"statisticTableBody\"></tbody>"+
            "</table>";
        let statisticTableBody = $('#statisticTableBody');
        for (let i = 0; i < staEachs.length; i++) {
            statisticTableBody.append(
                "<tr>" +
                "<td>" + staEachs[i]["projectId"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvgMap"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvgPassNum"] + "</td>" +
                "<td>" + staEachs[i]["rdBestAvgCrashNum"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvgMinMap"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvgMinRank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvg14Map"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvg14Rank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvgMidMap"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvgMidRank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvg34Map"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvg34Rank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["rdBestAvgMaxMap"]) + "</td>" +
                "<td>" + staEachs[i]["rdBestAvgMaxRank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["artMap"]) + "</td>" +
                "<td>" + staEachs[i]["artRank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["twoMap"]) + "</td>" +
                "<td>" + staEachs[i]["twoRank"] + "</td>" +
                "<td>" + parseFloat(staEachs[i]["combineBestMap"]) + "</td>" +
                "<td>" + parseFloat(staEachs[i]["pn"]) + "</td>" +
                "<td>" + parseFloat(staEachs[i]["fn"]) + "</td>" +
                "<td>" + staEachs[i]["combineBestPassNum"] + "</td>" +
                "<td>" + staEachs[i]["combineBestCrashNum"] + "</td>" +
                "<td>" + staEachs[i]["combineBestRank"] + "</td>" +
                "</tr>");
        }
    });

    releaseLock();
}

async function staRdAnalyzeJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("casePrefix").value;
    let tester = document.getElementById("tester").value;

    if (projectId === "" || casePrefix === "" || tester === "") {
        alert("projectId or casePrefix or tester empty");
        return
    }

    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;

    if (!acquireLock()) return;

    await httpGet("stardanalyze", "projectId="+projectId+"&casePrefix="+casePrefix+
        "&analyzeFile="+analyzeFile+"&factor="+factor+"&tester="+tester, function(res) {
        let jsonData = JSON.parse(res);
        document.getElementById("message").value =
            "passNum="+jsonData["passNum"]+";"+
            "crashNum="+jsonData["crashNum"]+";"+
            "randNum="+jsonData["randNum"]+";\n"+
            "bestAvgMap="+jsonData["bestAvgMap"]+";"+
            "bestAvgPassNum="+jsonData["bestAvgPassNum"]+";"+
            "bestAvgCrashNum="+jsonData["bestAvgCrashNum"]+";"+
            "bestAvgMidMap="+jsonData["bestAvgMidMap"]+";"+
            "bestAvgMidRank="+jsonData["bestAvgMidRank"]+";\n"+
            "bestAvgMaxMap="+jsonData["bestAvgMaxMap"]+";"+
            "bestAvgMaxRank="+jsonData["bestAvgMaxRank"]+";"+
            "bestAvgMinMap="+jsonData["bestAvgMinMap"]+";"+
            "bestAvgMinRank="+jsonData["bestAvgMinRank"]+";\n"+
            "bestAvg14Map="+jsonData["bestAvg14Map"]+";"+
            "bestAvg14Rank="+jsonData["bestAvg14Rank"]+";"+
            "bestAvg34Map="+jsonData["bestAvg34Map"]+";"+
            "bestAvg34Rank="+jsonData["bestAvg34Rank"];
        let xs = jsonData["passCrashAvgMapX"];
        let ys = jsonData["passCrashAvgMapY"];
        let zs = jsonData["passCrashAvgMapZ"];

        document.getElementById("table").innerHTML = "";

        // 基于准备好的dom，初始化echarts实例
        let myChart = echarts.init(document.getElementById('graph'));
        myChart.clear();
        // 指定图表的配置项和数据
        let option = null
        if (jsonData["crashNum"] === 1) {
            let echartData = [];
            for (let i = 0; i < xs.length; i++) {
                echartData.push([xs[i], parseFloat(zs[i])]);
            }
            console.log(zs);
            option = {
                title: {
                    text: projectId
                },
                tooltip: {},
                toolbox: {
                    show: true,
                    feature: {
                        saveAsImage: {
                            pixelRatio: 4
                        }
                    }
                },
                xAxis: {name:"passNum"},
                yAxis: {name:"avgMap"},
                series: [{
                    type: 'scatter',
                    animation: false,
                    data: echartData
                }]
            };
        } else {
            let echartData = [];
            for (let i = 0; i < xs.length; i++) {
                echartData.push([xs[i], ys[i], parseFloat(zs[i])]);
            }
            console.log(echartData)
            option = {
                title: {
                    text: projectId
                },
                tooltip: {},
                toolbox: {
                    show: true,
                    feature: {
                        saveAsImage: {
                            pixelRatio: 4
                        }
                    }
                },
                // 需要注意的是我们不能跟 grid 一样省略 grid3D
                grid3D: {},
                // 默认情况下, x, y, z 分别是从 0 到 1 的数值轴
                xAxis3D: {name:"passNum"},
                yAxis3D: {name:"crashNum"},
                zAxis3D: {name:"avgMap"},
                series: [{
                    type: 'scatter3D',
                    data: echartData
                }]
            };
        }
        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    });

    releaseLock();
}

async function staCombineAnalyzeJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("casePrefix").value;

    if (projectId === "" || casePrefix === "") {
        alert("projectId or casePrefix empty");
        return
    }

    let analyzeFile = document.getElementById("analyzeFile").value;
    let factor = document.getElementById("factor").value;

    if (!acquireLock()) return;

    await httpGet("stacombineanalyze", "projectId="+projectId+"&casePrefix="+casePrefix+
        "&analyzeFile="+analyzeFile+"&factor="+factor, function(res) {
        let jsonData = JSON.parse(res);
        document.getElementById("message").value =
            "n="+jsonData["n"]+";"+
            "bestMap="+jsonData["bestMap"]+";"+
            "bestPassNum="+jsonData["bestPassNum"]+";"+
            "bestCrashNum="+jsonData["bestCrashNum"]+";"+
            "bestRank="+jsonData["bestRank"];
        let xData = [];
        for (let i = 1; i <= jsonData["n"]; i++) {
            xData.push(i);
        }
        let numMaps = jsonData["numMaps"];
        let echartData = [];
        for (let i = 0; i < numMaps.length; i++) {
            let numMapsI = numMaps[i];
            let temp = [];
            for (let j = 0; j < numMapsI.length; j++) {
                temp.push(parseFloat(numMapsI[j]));
            }
            echartData.push(temp);
        }
        // box
        echartData = echarts.dataTool.prepareBoxplotData(echartData);
        console.log(echartData);

        document.getElementById("table").innerHTML = "";

        // 基于准备好的dom，初始化echarts实例
        let myChart = echarts.init(document.getElementById('graph'));
        myChart.clear();
        // 指定图表的配置项和数据
        let option = {
            title: {
                text: projectId
            },
            tooltip: {},
            toolbox: {
                show: true,
                feature: {
                    saveAsImage: {
                        pixelRatio: 4
                    }
                }
            },
            xAxis: {
                name:"caseNum",
                data: xData
            },
            yAxis: {name:"map"},
            series: [
                {
                    type: 'boxplot',
                    animation: false,
                    data: echartData.boxData
                },
                {
                    // 别忘了异常值
                    name: 'outlier',
                    type: 'scatter',
                    data: echartData.outliers
                }
            ]
        };
        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    });

    releaseLock();
}

async function staFactorsAnalyzeJ() {
    let analyzeFile = document.getElementById("analyzeFile").value;
    let casePrefix = document.getElementById("casePrefix").value;

    if (!acquireLock()) return;

    await httpGet("stafactorsanalyze", "analyzeFile="+analyzeFile+"&casePrefix="+casePrefix, function(res) {
        let jsonData = JSON.parse(res);
        console.log(jsonData);
        let topsStr = "";
        let tops = jsonData["tops"];
        for (let i = 0; i < tops.length; i++) {
            let topsI = tops[i];
            for (let j = 0; j < topsI.length; j++) {
                topsStr += topsI[j]+" ";
            }
            topsStr += "\n";
        }
        document.getElementById("message").value =
            "factors="+jsonData["factors"]+";\n"+
            "mrrs="+jsonData["mrrs"]+";\n"+
            topsStr;

        let factors = jsonData["factors"];
        let staFactorsEachs = jsonData["staFactorsEachs"];

        let statisticTableC = "width: 100%;font-size: 12px;font-family:Arial,Helvetica,sans-serif;text-align: center;"
        let statisticTableThC = "color: white;border: 1px solid black;background-color: #333333;"

        let thData = "";
        for (let i = 0; i < factors.length; i++) {
            thData += "<th style=\""+statisticTableThC+"\">"+factors[i]+"-map"+"</th>";
            thData += "<th style=\""+statisticTableThC+"\">"+factors[i]+"-rank"+"</th>";
        }

        document.getElementById("table").innerHTML = "" +
            "<table style=\""+statisticTableC+"\">"+
            "<thead><tr>"+
            "<th style=\""+statisticTableThC+"\">project</th>"+
            thData+
            "</tr></thead>"+
            "<tbody id=\"statisticTableBody\"></tbody>"+
            "</table>";

        let statisticTableBody = $('#statisticTableBody');
        for (let i = 0; i < staFactorsEachs.length; i++) {
            let tbData = "";
            for (let j = 0; j < factors.length; j++) {
                tbData += "<td>" + parseFloat(staFactorsEachs[i]["maps"][j]) + "</td>";
                tbData += "<td>" + staFactorsEachs[i]["ranks"][j] + "</td>";
            }
            statisticTableBody.append(
                "<tr>"+
                "<td>" + staFactorsEachs[i]["projectId"] + "</td>" +
                tbData+
                "</tr>");
        }

    });

    releaseLock();
}


async function checkTJ() {
    let projectId = document.getElementById("projectName").value;
    let casePrefix = document.getElementById("casePrefix").value;
    let tester = document.getElementById("tester").value;

    if (projectId === "" || casePrefix === "" || tester === "") {
        alert("projectId or casePrefix or tester empty");
        return
    }

    if (!acquireLock()) return;

    await httpGet("stardtesting", "projectId="+projectId+"&casePrefix="+casePrefix+"&tester="+tester, function(res) {
        let jsonData = JSON.parse(res);

        document.getElementById("table").innerHTML = "";

        // 基于准备好的dom，初始化echarts实例
        let myChart = echarts.init(document.getElementById('graph'));
        myChart.clear();
        // 指定图表的配置项和数据
        let option = null
        let echartData = [];
        for (let i = 0; i < jsonData.length; i++) {
            echartData.push([i,jsonData[i]]);
        }
        console.log(echartData);
        option = {
            title: {
                text: projectId
            },
            tooltip: {},
            toolbox: {
                show: true,
                feature: {
                    saveAsImage: {
                        pixelRatio: 4
                    }
                }
            },
            xAxis: {name:"passNum"},
            yAxis: {name:"avgMap"},
            series: [{
                type: 'scatter',
                animation: false,
                data: echartData
            }]
        };
        // 使用刚指定的配置项和数据显示图表。
        myChart.setOption(option);
    });

    releaseLock();
}