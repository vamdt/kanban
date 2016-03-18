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
        levels.forEach((k) => {
          Object.keys(def_level).forEach((p) => {
            assert.property(m[k], p);
            if (Array.isArray(def_level[p])) {
              assert.isArray(m[k][p], `${k}.${p}`);
            } else {
              assert.typeOf(m[k][p], typeof def_level[p], `${k}.${p}`);
            }
          });
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
});
