import { mergeWithKey, mergeData } from '../app/stock/util';

describe('Util', () => {
  describe('merge_with_key', () => {
    it('should return n when o is null', () => {
      const o = false;
      const n = {};
      const m = mergeWithKey(o, n, '');
      assert.equal(m, n);
    });

    it('should copy key to o when donot conatin key in o', () => {
      const o = {};
      const n = { key: [1] };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.equal(m.key, n.key);
    });

    it('should copy key to o when o[key] is not a array', () => {
      const o = { key: 'no' };
      const n = { key: [1] };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.equal(m.key, n.key);
    });

    it('should copy key to o when o[key] is empty', () => {
      const o = { key: [] };
      const n = { key: [1] };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.equal(m.key, n.key);
    });

    it('should concat n[key] when o[last].date lt n', () => {
      const o = {
        key: [ { date: 1 } ],
      };
      const n = {
        key: [ { date: 2 } ],
      };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.lengthOf(m.key, 2);
      assert.equal(m.key[0], o.key[0]);
      assert.equal(m.key[1], n.key[0]);
    });

    it('should drop o[key] when o[0].date gt n[0].date', () => {
      const o = {
        key: [ { date: 2 } ],
      };
      const n = {
        key: [ { date: 1 } ],
      };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.lengthOf(m.key, 1);
      assert.equal(m.key[0], n.key[0]);
    });

    it('should drop some data of o[key] when o[some].date gt n[0].date', () => {
      const o = {
        key: [
          { date: 1 },
          { date: 2 },
        ],
      };
      const n = {
        key: [
          { date: 2 },
          { date: 3 },
        ],
      };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.lengthOf(m.key, 3);
      assert.equal(m.key[0], o.key[0]);
      assert.equal(m.key[1], n.key[0]);
    });

    it('should drop the first data of o[key] when o[some].date eq n[0].date', () => {
      const o = {
        key: [
          { date: 2 },
          { date: 3 },
        ],
      };
      const n = {
        key: [
          { date: 2 },
          { date: 3 },
        ],
      };
      const m = mergeWithKey(o, n, 'key');
      assert.equal(m, o);
      assert.lengthOf(m.key, 2);
      assert.equal(m.key[0], n.key[0]);
      assert.equal(m.key[1], n.key[1]);
    });
  });

  describe('merge_data', () => {
    it('should return o when n is null', () => {
      const o = 1;
      const n = false;
      const m = mergeData(o, n);
      assert.equal(m, o);
    });

    it('should return n when o is null', () => {
      const o = false;
      const n = 1;
      const m = mergeData(o, n);
      assert.equal(m, n);
    });

    it('should copy key from n when o[key] is null', () => {
      const o = {};
      const n = {
        m1s: {
          data: [],
        },
      };
      assert.typeOf(n.m1s.data, 'array', 'data should be an array');
      const m = mergeData(o, n);
      assert.equal(m.m1s, n.m1s);
    });

    it('should copy key from n when o[key].data is empty', () => {
      const o = {
        m1s: {
          data: [],
        },
      };
      const n = {
        m1s: {
          data: [
            { Time: '2001-08-20T00:00:00Z' },
          ],
        },
      };
      assert.typeOf(n.m1s.data[0].Time, 'string', 'data[].time should be a time string');
      assert.lengthOf(n.m1s.data[0].Time, 20, 'should be a time string with format %Y-%m-%dT%XZ');
      const m = mergeData(o, n);
      assert.typeOf(m.m1s.data, 'array', 'data should be an array');
      assert.equal(m.m1s.data, n.m1s.data);
    });

    it('should merge .data currect', () => {
      const o = {
        m1s: {
          data: [
            { Time: '2001-08-20T00:00:00Z' },
            { Time: '2001-08-21T00:00:00Z' },
          ],
        },
      };
      const n = {
        m1s: {
          data: [
            { Time: '2001-08-21T00:00:00Z' },
            { Time: '2001-08-22T00:00:00Z' },
          ],
        },
      };
      assert.typeOf(n.m1s.data[0].Time, 'string', 'data[].time should be a time string');
      assert.lengthOf(n.m1s.data[0].Time, 20, 'should be a time string with format %Y-%m-%dT%XZ');
      const m = mergeData(o, n);
      assert.typeOf(n.m1s.data[0].date, 'date', 'data[].date should be a time');
      assert.lengthOf(m.m1s.data, 3);
      assert.equal(m.m1s[0], o.m1s[0]);
      assert.equal(m.m1s[1], n.m1s[0]);
    });
  });
});
