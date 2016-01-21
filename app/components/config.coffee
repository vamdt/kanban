defaults =
  nc: false
  nmas: false
  ocl: false
  nvolume: false
  nmacd: false
  mas: [
    interval: 5
    interval: 10
    interval: 20
  ]

load = ->
  s = try
    JSON.parse localStorage.getItem 'settings'
  catch
    defaults
  s || defaults

save = (settings) ->
  localStorage.setItem('settings', JSON.stringify(settings))

update = (settings) ->
  o = load()
  try
    n = JSON.parse JSON.stringify settings
    for k,v of n
      o[k] = v
    save n
  catch

module.exports =
  load: load
  save: save
  update: update
