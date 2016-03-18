import EventEmitter from './event';

class Typings {
  constructor(name) {
    Object.defineProperty(this, 'name', {
      value: name.toLowerCase(),
    });
    const props = [
      'Data',
      'Line',
    ];
    props.forEach((name) => {
      Object.defineProperty(this, name, {
        value: [],
      });
    });
  }
}

class Level {
  constructor(name, prev) {
    Object.defineProperty(this, 'name', {
      value: name.toLowerCase(),
    });

    if (prev) {
      Object.defineProperty(this, 'prev', {
        value: prev.toLowerCase(),
      });
    }

    Object.defineProperty(this, 'data', {
      value: [],
    });

    const props = [
      'Typing',
      'Segment',
      'Hub',
    ];
    props.forEach((name) => {
      Object.defineProperty(this, name, {
        value: new Typings(name),
      });
    });
  }
}

export default class Model extends EventEmitter {
  constructor(id) {
    super();
    Object.defineProperty(this, 'id', {
      value: id,
    });

    Object.defineProperty(this, 'name', {
      value: '',
      writable: true,
    });

    const levels = [
      'm1s',
      'm5s',
      'm30s',
      'days',
      'weeks',
      'months',
    ];
    let prev = false;
    levels.forEach((lname) => {
      Object.defineProperty(this, lname, {
        value: new Level(lname, prev),
      });
      prev = lname;
    });
  }
}
