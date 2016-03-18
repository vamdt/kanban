import EventEmitter from './event';

export default class Model extends EventEmitter {
  constructor(id) {
    super();
    Object.defineProperty(this, 'id', {
      enumerable: true,
      value: id,
    });
  }
}
