// 记录点击位置
let PX = 0, PY = 0;

// 点击图片时获取合适的ui tree元素, 更新选择状态
window.onload = function () {
    let cvr = document.getElementById("screenshotRect");
    cvr.addEventListener("click", function (event) {
        if (!acquireLock()) return;
        let rect = cvr.getBoundingClientRect();
        let x = (event.clientX - rect.left) * (cvr.width / rect.width);
        let y = (event.clientY - rect.top) * (cvr.height / rect.height);

        // 设置PX, PY, 更新position
        PX = parseInt(x)
        PY = parseInt(y)
        document.getElementById("position").value = ""+PX+" "+PY;

        if (Root != null) {
            let bestNode = null;
            Root.foreach(function(node) {
                // 1. 包含目标坐标
                // 2. 可互动
                // 3. 面积最小
                if (node.contain(x, y)) {
                    if (bestNode === null) {
                        bestNode = node;
                    } else if (bestNode.isOperable()) {
                        if (node.isOperable() && node.areaS() < bestNode.areaS()) {
                            bestNode = node;
                        }
                    } else {
                        if (node.isOperable()) {
                            bestNode = node
                        } else if (node.areaS() < bestNode.areaS()) {
                            bestNode = node;
                        }
                    }
                }
            })
            if (bestNode !== null) {
                setSelectedNode(bestNode, true);
            }
        }
        releaseLock();
    });
};

// 绘制图片
// @param pngSrc: 要显示图片的base64编码
function drawScreen(pngSrc) {
    let cv = document.getElementById("screenshot");
    let cvr = document.getElementById("screenshotRect");
    let cxt = cv.getContext("2d");
    let cxtr = cvr.getContext("2d");

    let img = new Image();
    img.src = pngSrc;
    img.onload = function() {
        // 将画布分辨率更改为图片分辨率, 这样二者的坐标也就对应上了
        cv.setAttribute("width", img.width.toString())
        cv.setAttribute("height", img.height.toString())
        cvr.setAttribute("width", img.width.toString())
        cvr.setAttribute("height", img.height.toString())
        cxt.clearRect(0, 0, cv.width, cv.height)
        cxtr.clearRect(0, 0, cvr.width, cvr.height)
        cxtr.clearRect(0, 0, img.width, img.height)
        cxt.drawImage(img, 0, 0, img.width, img.height);
    }
}
