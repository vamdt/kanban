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

业务逻辑在
http://vip.stock.finance.sina.com.cn/quotes_service/view/js/detail_a.js

根据最新的成交额和成交量计算出ticks信息

