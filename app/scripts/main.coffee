KLine = require './stock'
kl = new KLine(container: '#container')

do ->
  getQuery = (key) ->
    search = location.search.slice(1)
    arr = search.split "&"
    for i in arr
      ar = i.split "="
      if ar[0] == key
        return ar[1]
    return ''

  s = getQuery 's'
  if not s.length
    return
  fq = getQuery 'fq'
  k = getQuery 'k'
  opacity = getQuery 'opacity'

  kl.param
    s: s
    k: k
    fq: fq
    opacity: opacity

  kl.init()
