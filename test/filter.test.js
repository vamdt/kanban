/* eslint-env mocha */
/* global assert */

import { filter } from '../app/stock/util';

describe('Util', () => {
  describe('filter', () => {
    it('should return [] when src is empty', () => {
      const data = filter([], []);
      assert.typeOf(data, 'array', 'data should be an array');
      assert.lengthOf(data, 0);
    });

    it('should return [] when range.length < 2', () => {
      const src = [0, 1, 2];
      const range = [0];
      const data = filter(src, range);
      assert.typeOf(data, 'array', 'data should be an array');
      assert.lengthOf(data, 0);
    });

    it('should have one item before range[0] and after range[last]', () => {
      const src = [
        {
          date: 20010820,
          ETime: '2001-08-20T01:00:00Z',
        },
        {
          date: 20010821,
          ETime: '2001-08-21T01:00:00Z',
        },
        {
          date: 20010822,
          ETime: '2001-08-22T01:00:00Z',
        },
        {
          date: 20010823,
          ETime: '2001-08-23T01:00:00Z',
        },
        {
          date: 20010824,
          ETime: '2001-08-24T01:00:00Z',
        },
        {
          date: 20010825,
          ETime: '2001-08-25T01:00:00Z',
        },
      ];
      const range = [
        { date: 20010822 },
        { date: 20010823 },
      ];
      const data = filter(src, range);
      assert.typeOf(data, 'array', 'data should be an array');
      assert.lengthOf(data, 4, JSON.stringify({
        data,
        range,
      }));
      assert.equal(data[0], src[1]);
      assert.equal(data[3], src[4]);
    });
  });
});
