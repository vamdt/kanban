cmds = [
  'sugg'
  'unwatch'
  'watch'
  'show'
]

handlers = {}
for e in cmds
  handlers[e] = require './'+e

module.exports =
  ready: ->
    vm = @
    for e, func of handlers
      do (e, func) ->
        vm.$on e, (opt) ->
          func.apply vm, opt
