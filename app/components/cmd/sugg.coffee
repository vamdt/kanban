d3 = require 'd3'
module.exports = (sid) ->
  if sid.length < 1
    return
  for s in @stocks
    if s.sid == sid
      @show_stock(s)
      return

  d3.text '/search?s='+sid, (error, data) =>
    if error
      console.log error
      return
    info = data.split(';')
    info.forEach (v, i) ->
      v = v.split(',')
      info[i] =
        sid: v[3]
        name: v[4]
    if info.length is 1
      return @show_stock(info[0])
    @sugg = info
