kline = require '../app/stock/kline'

describe 'KLine', ->
  describe 'merge_data', ->
    describe 'merge_with_key', ->
      it 'should return n when o is null', ->
        o = off
        n = {}
        m = kline.merge_with_key o, n, ''
        assert.equal(m, n)

      it 'should copy key to o when donot conatin key in o', ->
        o = {}
        n = {key:[1]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.equal(m.key, n.key)

      it 'should copy key to o when o[key] is not a array', ->
        o = {key: 'no'}
        n = {key: [1]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.equal(m.key, n.key)

      it 'should copy key to o when o[key] is empty', ->
        o = {key: []}
        n = {key: [1]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.equal(m.key, n.key)

      it 'should concat n[key] when o[last].date lt n', ->
        o = {key: [{date: 1}]}
        n = {key: [{date: 2}]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.lengthOf(m.key, 2)
        assert.equal(m.key[0], o.key[0])
        assert.equal(m.key[1], n.key[0])

      it 'should drop o[key] when o[0].date gt n[0].date', ->
        o = {key: [{date: 2}]}
        n = {key: [{date: 1}]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.lengthOf(m.key, 1)
        assert.equal(m.key[0], n.key[0])

      it 'should drop some data of o[key] when o[some].date gt n[0].date', ->
        o = {key: [{date: 1}, {date: 2}]}
        n = {key: [{date: 2}, {date: 3}]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.lengthOf(m.key, 3)
        assert.equal(m.key[0], o.key[0])
        assert.equal(m.key[1], n.key[0])

      it 'should drop the first data of o[key] when o[some].date eq n[0].date', ->
        o = {key: [{date: 2}, {date: 3}]}
        n = {key: [{date: 2}, {date: 3}]}
        m = kline.merge_with_key o, n, 'key'
        assert.equal(m, o)
        assert.lengthOf(m.key, 2)
        assert.equal(m.key[0], n.key[0])
        assert.equal(m.key[1], n.key[1])

    it 'should return o when n is null', ->
      o = 1
      n = off
      m = kline.merge_data o, n
      assert.equal(m, o)

    it 'should return n when o is null', ->
      o = off
      n = 1
      m = kline.merge_data o, n
      assert.equal(m, n)

    it 'should copy key from n when o[key] is null', ->
      o = {}
      n =
        m1s:
          data: []
      assert.typeOf(n.m1s.data, 'array', 'data should be an array')
      m = kline.merge_data o, n
      assert.equal(m.m1s, n.m1s)

    it 'should copy key from n when o[key].data is empty', ->
      o =
        m1s:
          data:[]
      n =
        m1s:
          data: [{time: '2001-08-20T00:00:00Z'}]
      assert.typeOf(n.m1s.data[0].time, 'string', 'data[].time should be a time string')
      assert.lengthOf(n.m1s.data[0].time, 20, 'should be a time string with format %Y-%m-%dT%XZ')
      m = kline.merge_data o, n
      assert.typeOf(m.m1s.data, 'array', 'data should be an array')
      assert.equal(m.m1s.data, n.m1s.data)

    it 'should merge .data currect', ->
      o =
        m1s:
          data:[{time: '2001-08-20T00:00:00Z'}, {time: '2001-08-21T00:00:00Z'}]
      n =
        m1s:
          data: [{time: '2001-08-21T00:00:00Z'}, {time: '2001-08-22T00:00:00Z'}]
      assert.typeOf(n.m1s.data[0].time, 'string', 'data[].time should be a time string')
      assert.lengthOf(n.m1s.data[0].time, 20, 'should be a time string with format %Y-%m-%dT%XZ')
      m = kline.merge_data o, n
      assert.typeOf(n.m1s.data[0].date, 'date', 'data[].date should be a time')
      assert.lengthOf(m.m1s.data, 3)
      assert.equal(m.m1s[0], o.m1s[0])
      assert.equal(m.m1s[1], n.m1s[0])

  describe 'filter', ->

    it 'should return [] when src is empty', ->
      src = []
      range = []
      data = kline.filter src, range
      assert.typeOf(data, 'array', 'data should be an array')
      assert.lengthOf(data, 0)

    it 'should return [] when range.length < 2', ->
      src = [0,1,2]
      range = [0]
      data = kline.filter src, range
      assert.typeOf(data, 'array', 'data should be an array')
      assert.lengthOf(data, 0)

    it 'should have one item before range[0] and after range[last]', ->
      src = [
        {Time: '2001-08-20T00:00:00Z', ETime: '2001-08-20T01:00:00Z'}
        {Time: '2001-08-21T00:00:00Z', ETime: '2001-08-21T01:00:00Z'}
        {Time: '2001-08-22T00:00:00Z', ETime: '2001-08-22T01:00:00Z'}
        {Time: '2001-08-23T00:00:00Z', ETime: '2001-08-23T01:00:00Z'}
        {Time: '2001-08-24T00:00:00Z', ETime: '2001-08-24T01:00:00Z'}
        {Time: '2001-08-25T00:00:00Z', ETime: '2001-08-25T01:00:00Z'}
      ]
      range = [
        {Time: '2001-08-22T00:00:00Z'}
        {Time: '2001-08-23T00:00:00Z'}
      ]
      data = kline.filter src, range
      assert.typeOf(data, 'array', 'data should be an array')
      assert.lengthOf(data, 4, JSON.stringify(data: data, range: range))
      assert.typeOf(data[0].date, 'date')
      assert.equal(data[0], src[1])
      assert.equal(data[3], src[4])
