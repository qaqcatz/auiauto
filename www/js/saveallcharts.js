let AllProjectsId = [
    "ActivityDiary-118",
    "AlarmClock_389",
    "AmazeFileManager-1796",
    "AmazeFileManager-1837",
    "and-bible-261",
    "and-bible-375",
    "and-bible-480",
    "Anki-Android-4586",
    "Anki-Android-4589",
    "Anki-Android-4977",
    "Anki-Android-5638",
    "Anki-Android-5756",
    "Anki-Android-6145",
    "AntennaPod-3138",
    "AntennaPod_4645",
    "AnyMemo_422",
    "AnyMemo_440",
    "APhotoManager_139",
    "collect-3222",
    "commons-1391",
    "commons-1581",
    "commons-2123",
    "dagger-46",
    "Easy_xkcd_134",
    "FirefoxLite-4881",
    "FirefoxLite-4942",
    "FirefoxLite-5085",
    "geohashdroid-73",
    "Images-to-PDF_585",
    "Images-to-PDF_771",
    "k-9_3255",
    "markor_194",
    "nextcloud-1918",
    "nextcloud-4026",
    "nextcloud-4792",
    "nextcloud-5173",
    "open-event-attendee-2198",
    "openlauncher-67",
    "opentasks_629",
    "osmeditor4android-729",
    "Scarlet-Notes-114",
    "screenrecorder-32",
    "Simple-Music-Player_128",
    "Simple-Music-Player_204",
    "ToGoZip_10",
    "ToGoZip_20"
]

let AnalyzeFiles = ["rootcause", "fix"]

async function saveAllChartsJ() {
    if (!acquireLock()) return;
    if (window.confirm("确认save all charts")){
        for (let i0 = 0; i0 < AllProjectsId.length; i0++) {
            for (let j0 = 0; j0 < AnalyzeFiles.length; j0++) {
                // 基于准备好的dom，初始化echarts实例
                let myChart = echarts.init(document.getElementById('graph'));
                // rd
                let projectId = AllProjectsId[i0];
                let casePrefix = "rd";
                let tester = "monkey";
                let analyzeFile = AnalyzeFiles[j0];
                let factor = "Ochiai";

                let flag = false;
                let threeD = "";
                await httpGet("stardanalyze", "projectId="+projectId+"&casePrefix="+casePrefix+
                    "&analyzeFile="+analyzeFile+"&factor="+factor+"&tester="+tester, function(res) {
                    flag = true;
                    let jsonData = JSON.parse(res);
                    let xs = jsonData["passCrashAvgMapX"];
                    let ys = jsonData["passCrashAvgMapY"];
                    let zs = jsonData["passCrashAvgMapZ"];
                    document.getElementById("table").innerHTML = "";
                    myChart.clear();
                    // 指定图表的配置项和数据
                    let option = null;
                    if (jsonData["crashNum"] === 1) {
                        let echartData = [];
                        for (let i = 0; i < xs.length; i++) {
                            echartData.push([xs[i], parseFloat(zs[i])]);
                        }
                        option = {
                            title: {
                                text: projectId
                            },
                            tooltip: {},
                            xAxis: {name:"passNum"},
                            yAxis: {name:"avgMap"},
                            series: [{
                                type: 'scatter',
                                animation: false,
                                data: echartData
                            }]
                        };
                    } else {
                        // echart导出3d图像可能会有问题, 需要自己手动拉一下
                        threeD = "3D";
                        let echartData = [];
                        for (let i = 0; i < xs.length; i++) {
                            echartData.push([xs[i], ys[i], parseFloat(zs[i])]);
                        }
                        option = {
                            title: {
                                text: projectId
                            },
                            tooltip: {},
                            // 需要注意的是我们不能跟 grid 一样省略 grid3D
                            grid3D: {},
                            // 默认情况下, x, y, z 分别是从 0 到 1 的数值轴
                            xAxis3D: {name:"passNum"},
                            yAxis3D: {name:"crashNum"},
                            zAxis3D: {name:"avgMap"},
                            animation: false,
                            series: [{
                                type: 'scatter3D',
                                data: echartData
                            }]
                        };
                    }
                    // 使用刚指定的配置项和数据显示图表。
                    myChart.setOption(option);
                });
                if (flag === false) {
                    releaseLock();
                    return;
                }
                flag = false
                await httpPost("savecharts", "casePrefix="+casePrefix+"&analyzeType=random"+
                    "&analyzeFile="+analyzeFile+"&factor="+factor+"&projectId="+threeD+projectId,
                    JSON.stringify({
                        base64: myChart.getDataURL({ pixelRatio: 4}),
                    }), function (res) {
                        flag = true;
                    });
                if (flag === false) {
                    releaseLock();
                    return;
                }
                // combine
                casePrefix = "art";
                flag = false;
                await httpGet("stacombineanalyze", "projectId="+projectId+"&casePrefix="+casePrefix+
                    "&analyzeFile="+analyzeFile+"&factor="+factor, function(res) {
                    flag = true
                    let jsonData = JSON.parse(res);
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
                    myChart.clear();
                    // 指定图表的配置项和数据
                    let option = {
                        title: {
                            text: projectId
                        },
                        tooltip: {},
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
                if (flag === false) {
                    releaseLock();
                    return;
                }
                flag = false
                await httpPost("savecharts", "casePrefix="+casePrefix+"&analyzeType=combine"+
                    "&analyzeFile="+analyzeFile+"&factor="+factor+"&projectId="+projectId,
                    JSON.stringify({
                        base64: myChart.getDataURL({ pixelRatio: 4}),
                    }), function (res) {
                        flag = true;
                    });
                if (flag === false) {
                    releaseLock();
                    return;
                }
            }
        }
    }
    alert("success!")
    releaseLock();
}