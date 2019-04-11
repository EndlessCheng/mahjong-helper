# mahjong-helper

## 安装

分下面几步：

1. 前往 [release](https://github.com/EndlessCheng/mahjong-helper/releases/latest) 页面下载程序

2. 安装浏览器扩展 Header Editor，具体操作可以参考[这篇](https://tieba.baidu.com/p/5956122477)，

    安装好扩展后点进该扩展的`管理`界面，点击`导入和导出`，在下载规则中填入 `https://jianyan.me/js/mahjong-helper.json`，点击右侧的下载按钮，然后点击下方的`保存`

3. （雀魂需要）允许本地证书通过浏览器，在浏览器（仅限 Chrome 和使用了 Chrome 内核的浏览器）中输入

    ```
    chrome://flags/#allow-insecure-localhost
    ```

    然后把高亮那一项的 Disabled 改成 Enabled（不同浏览器/版本的描述可能不一样，如果是中文的话点击「启用」按钮），之后重启浏览器

（PS：第2步发生了什么见[如何获取 WebSocket 收发的消息](#如何获取WebSocket收发的消息)）

### 从源码安装程序

您也可以选择从源码安装：

`go get -u -v github.com/EndlessCheng/mahjong-helper/...`

完成后程序生成于 `$GOPATH/bin/` 目录下


## 服务端脚本说明

分为手动/自动两种模式

### 手动挡

安装[手动挡客户端](https://github.com/EndlessCheng/mahjong-helper-gui)，可用于牌谱分析

或者在命令行输入:

- 分析何切

    `mahjong-helper 34568m 5678p 23567s`
    
- 分析鸣牌

    `mahjong-helper 33567789m 46s + 6m`

### 自动挡

启动程序，选择平台即可

自动化分析会在牌效率上考虑场上已有的牌（含舍牌、副露、宝牌指示牌）；
防守方面，提示他家手切模切，并当有人立直或多副露时，将手牌按照危险度进行排序，详见下文


## CLI

### 牌效率

综合了进张数、改良、向听前进后的进张数，越好的越靠前

显示如下（无改良时不显示改良）：

```
进张数               切哪张牌 [哪些牌是进张]
改良后的进张数加权均值 [改良数] 向听前进后的进张数的加权均值
```

可能说的有点绕，还是看看下面几个例子吧

两面一向听：

按照蓝-黄-红的顺序，进张数越多颜色越红，切牌越靠近 456 颜色越红

此例切 8s

![](img/example11.png)

包含三个复合搭子的一向听：

这里展示了本程序对于进张与向听前进后进张之间的综合判断，切 6s 最佳

![](img/example22.png)

如果听牌型很差，额外提供向听倒退的选择：

不考虑场况的话，相比 8m，切 1m 虽然向听倒退但是有断幺一役，速度其实是高于 8m 的

（两向听的切牌不区分颜色）

![](img/example41.png)

### 鸣牌判断

下图是一个鸣了红中之后（鸣出的牌不显示），听坎 5s 的例子，宝牌为 6m

上家打出了 6m 宝牌之后分析如下：

这里就可以考虑用 57m 吃，打出 9m，提升打点的同时又能维持听牌，相比向听倒退更期待 37s 的两面改良

![](img/example63.png)

### 模切与安牌显示

下图展示了某局中三家的模切情况（宝牌为红中，自家手牌为 25567m 488p 335568s）：

- 白色为手切，暗灰色为模切
- 鸣牌后打出的那张牌会用灰底白字显示，供读牌分析用
- 副露玩家的手切中张牌(3-7)会有不同颜色的高亮，用来辅助判断其听牌率
- 玩家立直后会显示对该玩家的安牌（这里上家 3p 立直），第一排为现物，第二排按照铳率由低到高排序（No Chance 和 One Chance 的安牌作为补充参考显示在后面）

![](img/example_moqie_risk.png)

补充说明：

- 铳率排序是基于巡目、筋牌、No Chance 的综合考虑结果，对于早外、One Chance 和其他特殊情况并没有考虑，请玩家自行斟酌
- 但是对于某些情况下的 No Chance 安牌，本程序是会将其视作现物的（比如 3m 为壁，剩下的 2m 在牌河和自己手里时，2m 是不会放铳的）
- 上图 8s 为安牌是因为在立直时手里没有任何安牌，拆掉了暗刻 8s
- 上图对家吃了下家打出的 5s，又在其后手切了 6p 及 4m，很有可能听牌了，这里就需要玩家自己在攻守时留意下对家的动作了
- 无安时可以考虑拆掉 One Chance 的 3s 对子

## 如何获取WebSocket收发的消息

1. 打开开发者工具，找到相关 JS 文件，保存到本地
2. 搜索 `WebSocket`, `socket`，找到 `message`, `onmessage` 等函数
3. 修改代码，使用 `XMLHttpRequest` 将收发的消息发送到（在 localhost 开启的）mahjong-helper 服务器，服务器收到消息后会自动进行相关分析
4. 将修改后的 JS 代码传至个人的 github.io 项目，拿到该 JS 文件地址
5. 安装浏览器扩展 Header Editor，重定向原 JS 文件地址到上一步中拿到的地址，具体操作可以参考[这篇](https://tieba.baidu.com/p/5956122477)
6. 允许本地证书通过浏览器，在浏览器（仅限 Chrome 内核）中输入
    
    ```
    chrome://flags/#allow-insecure-localhost
    ```
    
    然后把高亮那一项的 Disabled 改成 Enabled（不同浏览器/版本的描述可能不一样，如果是中文的话点击「启用」按钮）

7. 重启浏览器

下面说明天凤和雀魂的代码注入点

### 天凤 (tenhou)

1. 搜索 `new WebSocket`，找到下方的 `message` 函数，该函数中的 `a.data` 就是 WebSocket 收到的数据
2. 在该函数末尾添加如下代码

    ```javascript
    var req = new XMLHttpRequest();
    req.open("POST", "http://localhost:12121/");
    req.send(a.data);
    ```

### 雀魂 (majsoul)

考虑到雀魂的 WebSocket 收到的是封装后的 protobuf 二进制数据，不好解析，于是另寻他路

大致思路是根据 [liqi.json](https://github.com/EndlessCheng/mahjong-helper/blob/master/liqi.json) 文件提供的对分析玩家操作有用的字段查找相关关键字，如 `ActionDealTile` `ActionDiscardTile` `ActionChiPengGang` 等，具体修改了哪些内容可以对比雀魂的 JS 代码和我修改后的 https://jianyan.me/majsoul/code-v0.1.4.js

PS: 在网页控制台输入 `GameMgr._inRelease = 0` 即可开启调试模式
