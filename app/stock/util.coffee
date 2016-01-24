d3 = require 'd3'

parseDate = d3.time.format("%Y-%m-%dT%XZ").parse

extend = ->
  dest = {}
  for i in arguments when i
    for k,v of i
      dest[k] = dest[k] || v
  dest

filter = (src, range) ->
  if (src||[]).length < 1
    return []
  if (range||[]).length < 2
    return []

  for d in range
    d.date = d.date || parseDate(d.Time)
  for d in src
    d.date = d.date || parseDate(d.Time)

  start_date = range[0].date
  end_date = range[range.length-1].date

  bisect = d3.bisector((d) -> +d.date)
  istart = bisect.left(src, +start_date)
  iend = bisect.right(src, +end_date)
  istart = Math.max(istart - 1, 0)
  src = src.slice istart, iend+1

  hash = {}
  hash[+d.date] = i for d, i in range

  indexOfFun = (start, end) ->
    (date) ->
      idate = +date
      if start > idate
        return -1
      if idate > end
        return hash[end]+1
      if hash.hasOwnProperty idate
        return hash[idate]
      bisect.right(range, +date)

  indexOf = indexOfFun(+start_date, +end_date)

  for d in src
    d.i = indexOf d.date
    if d.ETime
      d.edate = d.edate || parseDate(d.ETime)
      d.ei = indexOf d.edate

  src

merge_with_key = (o, n, k) ->
  if not o
    return n
  if !Array.isArray(n[k]) or n[k].length < 1
    return o
  if !Array.isArray(o[k]) or o[k].length < 1
    o[k] = n[k]
  else
    ndate = +n[k][0].date
    odate = +o[k][o[k].length-1].date
    o0date = +o[k][0].date
    if odate < ndate
      console.log 'merge_data with concat () + ()'
      o[k] = o[k].concat n[k]
    else if o0date >= ndate
      o[k] = n[k]
    else
      bisect = d3.bisector((d) -> +d.date)
      i = bisect.left(o[k], ndate)
      o[k] = o[k].slice(0, i).concat(n[k])
  o

data_init = (n) ->
  return n unless n
  for k in ['m1s', 'm5s', 'm30s', 'days', 'weeks', 'months'] when n[k] and n[k].data
    n[k].data.forEach (d) -> d.date = d.date || parseDate(d.time)
    for name in ['Typing', 'Segment', 'Hub'] when n[k][name]
      for dn in ['Data', 'Line'] when n[k][name][dn]
        n[k][name][dn].forEach (d) -> d.date = d.date || parseDate(d.Time)
  n

merge_data = (o, n) ->
  return o if not n
  n = data_init n
  return n if not o
  o = data_init o

  for k in ['m1s', 'm5s', 'm30s', 'days', 'weeks', 'months'] when n[k] and n[k].data
    if not o[k]
      o[k] = n[k]
    else
      o[k] = merge_with_key o[k], n[k], 'data'
      for name in ['Typing', 'Segment', 'Hub'] when n[k][name]
        if not o[k][name]
          o[k][name] = n[k][name]
          continue
        for dn in ['Data', 'Line'] when n[k][name][dn]
          o[k][name] = merge_with_key o[k][name], n[k][name], dn
  o


[cup, cdown, ceq] = ["#f00", "#080", "#000"]

kColor = (d, i, data) ->
  if d.open == d.close
    if i and data
      if data[i] and data[i-1]
        return cup if data[i].open >= data[i-1].close
        return cdown if data[i].open < data[i-1].close
    return ceq
  if d.open > d.close
    return cdown
  cup

module.exports =
  kColor: kColor
  parseDate: parseDate
  extend: extend
  filter: filter
  merge_data: merge_data
  merge_with_key: merge_with_key
