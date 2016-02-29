module.exports = ->
  return if arguments.length < 1

  num = 0
  for o in arguments
    i = -1
    i = j for ss, j in @stocks when ss.sid == o
    if i > -1
      @stocks.splice(i, 1)
      num++
  if num
    localStorage.setItem('stocks', JSON.stringify(@stocks))
