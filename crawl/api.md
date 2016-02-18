# 股票信息

## 数据库

### 数据量

            数据量
   K线类型    日   年     30年    3000股*30年
1. 1分钟k线  240  72000 2160000    64.8 亿
2. 5分钟k线   48  14400  432000    12.96亿
3. 15分钟k线  16   4800  144000     4.32亿
4. 30分钟k线   8   2400   72000     2.16亿
5. 60分钟k线   4   1200   36000     1.08亿
6. 日k线       1    300    9000     2700万
7. 周k线
8. 月k线

暂时忽略1分钟线，需21000个表

因为表数量多，所以mongodb的设置需要相应调整

[mongod config](http://docs.mongodb.org/manual/reference/program/mongod/#bin.mongod)

存储引擎为 mmapv1 即默认引擎则
--nssize 2047

或改用 wiredTiger 引擎

### 存储规则

{stock_id} : 股票id
  带SH 或者 SZ 前缀  e.g. SH600000

{kline} : k线类型
  k5 5分钟K线
  k15 15分钟K线
  k30 30分钟K线
  k60 60分钟K线
  kday 日K线
  kweek 周K线
  kmonth 月K线

{fq} : 复权
  前复权 qfq
  后复权 hfq
  不复权 空

tdata:trading data
fdata:financial data
rdata:research data
ndata:note data

#### 数据表

K线数据表名 {stock_id}.tdata.{kline}_{fq}

数据存储格式
"time": "时间",	"data": [open, close, high, low, volume], "macd": [DIF, DEA, MACD]
"time": "010102",	"data": [16.75,17.39,17.55,16.65,66081]

#### 标注

表名: {stock_id}.ndata.{kline}_{fq}

数据存储格式
"type": "类型: rect, point, line, path",
"name": "曲线名字"
"x1": "",
"y1": "",
"x2": "",
"y2":"",
"comment": "point,line 从x1y1指向x2y2"
"config": {}

## API

### K线

1. 1分钟k线
2. 5分钟k线
3. 15分钟k线
4. 30分钟k线
5. 60分钟k线
6. 日k线
7. 周k线
8. 月k线

url: /stock/k
param:
  s: sh600000      股票代码
  k: 1-8           k线类型
  t: 1438653768500 上一次请求时间 毫秒级 无此参数时输出全部数据

response:
{
  param: {},
  code: 200,
  msg: "",
  data: {
    sh600000: []
  }
}

### 标注

url: /stock/annotate
param:
  s: sh600000      股票代码
  k: 1-8           k线类型
  t: 1438653768500 上一次请求时间 毫秒级 无此参数时输出全部数据

response:
{
  param: {},
  code: 200,
  msg: "",
  data: {
    sh600000: [
      {type:"rect", x1:"", y1:"", x2:"", y2:"", comment: ""},
      {type:"point", x1:"", y1:"", x2:"", y2:"", comment: "从x1y1指向x2y2"},
      {type:"line", x1:"", y1:"", x2:"", y2:"", comment: "从x1y1指向x2y2"}
    ]
  }
}

### 自定义曲线

url: /stock/path
param:
  name: pathname   曲线名字
  s: sh600000      股票代码
  k: 1-8           k线类型
  t: 1438653768500 上一次请求时间 毫秒级 无此参数时输出全部数据

response:
{
  param: {},
  code: 200,
  msg: "",
  data: {
    pathname: {
      sh600000: [
        ["2015-07-01 14:57:41", 5.10],
        ...
        []
      ]
    }
  }
}


## Data Source

### gu.qq.com

http://data.gtimg.cn/flashdata/hushen/4day/sh/sh000001.js

http://data.gtimg.cn/flashdata/hushen/minute/sh000001.js

http://data.gtimg.cn/flashdata/hushen/monthly/sh000001.js

http://data.gtimg.cn/flashdata/hushen/weekly/sh000001.js

http://data.gtimg.cn/flashdata/hushen/daily/15/sh000001.js

http://data.gtimg.cn/flashdata/hushen/latest/daily/sh000001.js

num:100 total:5958 start:901219 90:9 91:247 92:90 93:257 94:252 95:251 96:247 97:243 98:246 99:239 00:239 01:240 02:237 03:241 04:243 05:242 06:241 07:242 08:246 09:244 10:242 11:244 12:243 13:238 14:245 15:244 16:6\n\
150813 3869.91 3954.56 3955.79 3838.16 430073303\n\
       open    close   high    low     volume
160111 3131.85 3016.70 3166.22 3016.70 271643691\n\

### sina

http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=sh000001&scale=5&ma=no&datalen=1023


### 同花顺

//timestr, open,   high,   low,    close,  volume
//20160217,2829.76,2868.70,2824.36,2867.34,21690992000,225964250000.00,
http://d.10jqka.com.cn/v2/line/hs_600000/01/last.js
http://d.10jqka.com.cn/v2/line/hs_600000/01/2015.js
