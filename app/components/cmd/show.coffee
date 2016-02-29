module.exports = ->
  return if arguments.length < 1
  opts =
    mas: 'nmas'
    candle: 'nc'
    volume: 'nvolume'
    macd: 'nmacd'
    typing: 'ntyping'
    handcraft: 'handcraft'
  args = {}
  v = off
  args = Array.apply(null, arguments)
  if args.length > 1
    if args[args.length-1].toLowerCase() == 'false'
      v = on
      args.pop()
  for o in args
    unless opts[o]
      continue
    args[opts[o]] = v
  @$root.$broadcast('param_change', args)
