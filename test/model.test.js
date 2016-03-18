import Model from '../app/stock/model';

describe('Model', () => {
  describe('id', () => {
    it('should readonly', () => {
      const id = '007';
      const m = new Model(id);
      assert.equal(m.id, id);
      assert.throws(() => m.id = '001', 'Cannot assign');
    });
  });
});
