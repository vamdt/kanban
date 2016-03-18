import Model from '../app/stock/model';

describe('Model', () => {
  describe('property', () => {
    const def_level = {
      name: '',
      data: [],
      Typing: {
        Data: [],
        Line: [],
      },
      Segment: {
        Data: [],
        Line: [],
      },
      Hub: {
        Data: [],
        Line: [],
      },
    };
    const defaults = {
      id: '',
      name: '',
      m1s: def_level,
      m5s: def_level,
      m30s: def_level,
      days: def_level,
      weeks: def_level,
      months: def_level,
    };

    const m = new Model('007');
    const props = [
      'id',
      'name',
    ];
    it('should has expect property', () => {
      Object.keys(defaults).forEach((k) => {
        assert.property(m, k);
        assert.typeOf(m[k], typeof defaults[k], k);
      });
    });

    const levels = [
      'm1s',
      'm5s',
      'm30s',
      'days',
      'weeks',
      'months',
    ];

    describe('levels', () => {
      it('should has expect property', () => {
        let prev = false;
        levels.forEach((k) => {
          Object.keys(def_level).forEach((p) => {
            assert.property(m[k], p);
            if (Array.isArray(def_level[p])) {
              assert.isArray(m[k][p], `${k}.${p}`);
            } else {
              assert.typeOf(m[k][p], typeof def_level[p], `${k}.${p}`);
            }
          });
          if (prev) {
            assert.equal(m[k].prev, prev);
          }
          prev = k;
        });
      });

      it('should has expect prev', () => {
        let prev = false;
        levels.forEach((k) => {
          if (prev) {
            assert.equal(m[k].prev, prev);
          }
          prev = k;
        });
      });
    });
  });

  describe('id', () => {
    it('should readonly', () => {
      const id = '007';
      const m = new Model(id);
      assert.equal(m.id, id);
      assert.throws(() => m.id = '001', 'Cannot assign');
    });
  });

  describe('assign', () => {
    it('should assign .data currect', () => {
      const m = new Model('007');
      const o = {
        m1s: {
          data: [
            { time: '2001-08-20T00:00:00Z' },
          ],
        },
      };
      m.assign(o);
      assert.instanceOf(m.m1s.data[0].date, Date);
      assert.equal(m.m1s.data[0].time, o.m1s.data[0].time);
      assert.typeOf(m.m1s.data[0].time, 'string', 'data[].time should be a time string');
      assert.lengthOf(m.m1s.data[0].time, 20, 'should be a time string with format %Y-%m-%dT%XZ');
    });

    describe('merge .data', () => {
      it('should merge pop the repeat data', () => {
        const m = new Model('007');
        const o = {
          m1s: {
            data: [
              { time: '2001-08-20T00:00:00Z' },
            ],
          },
        };
        m.assign(o);
        m.assign(o);
        assert.equal(m.m1s.data[0].time, o.m1s.data[0].time);
        assert.lengthOf(m.m1s.data, 1);
      });

      it('should merge currect', () => {
        const m = new Model('007');
        const o = {
          m1s: {
            data: [
              { time: '2001-08-20T00:00:00Z' },
              { time: '2001-08-21T00:00:00Z' },
              { time: '2001-08-22T00:00:00Z' },
              { time: '2001-08-23T00:00:00Z' },
            ],
          },
        };
        const n = {
          m1s: {
            data: [
              { time: '2001-08-21T00:00:00Z' },
              { time: '2001-08-24T00:00:00Z' },
            ],
          },
        };
        m.assign(o);
        assert.equal(m.m1s.data[0].time, o.m1s.data[0].time);
        assert.lengthOf(m.m1s.data, o.m1s.data.length);
        m.assign(n);
        assert.equal(m.m1s.data[0].time, o.m1s.data[0].time);
        assert.lengthOf(m.m1s.data, 3);
        assert.equal(m.m1s.data[2].time, n.m1s.data[1].time);
      });
    });
  });
});
