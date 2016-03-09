sina finance api

http://money.finance.sina.com.cn

* 5min kline

url

/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?
symbol=sh600000
scale=5
ma=no
datalen=1023

res

[{
  day:"2015-07-20 15:00:00",
  open:"16.680",
  high:"16.690",
  low:"16.670",
  close:"16.680",
  volume:"5848304"
}]

* 15min kline

url

/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?
symbol=sh600000
scale=15
ma=no
datalen=1023

res

[{
  day:"2015-07-20 15:00:00",
  open:"16.680",
  high:"16.720",
  low:"16.670",
  close:"16.680",
  volume:"15399504"
}]

1. 历史分笔数据

表名 {stock_id}.tick

stock_id 中包含交易所信息
e.g. SH600000

Time          成交时间
Price         成交价
Change        价格变动
Volume        成交量(手)
Turnover      成交额(元)
Type          性质

sina api

http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s

2. 今日分笔数据

表名 t_tick_{stock_id}

sina api

http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=11&symbol=sh600000&rn=1438653768484

referer

http://finance.sina.com.cn/realstock/company/sh600000/nc.shtml

sina hq api

http://hq.sinajs.cn/rn=1438655910468&list=sh600000

0：大秦铁路，股票名字；
1：27.55，今日开盘价；
2：27.25，昨日收盘价；
3：26.91，当前价格；
4：27.55，今日最高价；
5：26.20，今日最低价；
6：26.91，竞买价，即“买一”报价；
7：26.92，竞卖价，即“卖一”报价；
8：22114263，成交的股票数，由于股票交易以一百股为基本单位，所以在使用时，通常把该值除以一百；
9：589824680，成交金额，单位为“元”，为了一目了然，通常以“万元”为成交金额的单位，所以通常把该值除以一万；
10：4695，“买一”申请4695股，即47手；
11：26.91，“买一”报价；
12：57590，“买二”
13：26.90，“买二”
14：14700，“买三”
15：26.89，“买三”
16：14300，“买四”
17：26.88，“买四”
18：15100，“买五”
19：26.87，“买五”
20：3100，“卖一”申报3100股，即31手；
21：26.92，“卖一”报价
(22, 23), (24, 25), (26,27), (28, 29)分别为“卖二”至“卖四的情况”
30：2008-01-11，日期；
31：15:05:32，时间；

业务逻辑在
http://vip.stock.finance.sina.com.cn/quotes_service/view/js/detail_a.js

根据最新的成交额和成交量计算出ticks信息

3. 新浪日数据
http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?&rand=9000&symbol=sz002241&end_date=&begin_date=&type=plain
date open high close low volume

4. 腾讯实时行情
http://qt.gtimg.cn/r=0.8409869808238&q=s_sz000559,s_sz000913,s_sz002048,s_sz002085,s_sz002126,s_sz002284,s_sh600001,s_sh600003,s_sh600004

5. other
http://ifzq.gtimg.cn/appstock/app/fqkline/get?p=1&param=sz000819,week,,,10000,hfq
http://ifzq.gtimg.cn/appstock/app/kline/kline?p=1&param=sh600765,day,,,320     开盘、收盘、最高、最低，成交量
http://ifzq.gtimg.cn/appstock/indicators/MACD/D1?market=SH&code=600765&args=12-26-9&start=&end=&limit=320&fq=bfq

http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=sz002405&scale=5&ma=no&datalen=1023

http://data.gtimg.cn/flashdata/hushen/daily/15/sh600604.js
