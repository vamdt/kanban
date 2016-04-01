export default class EventEmitter {
  constructor() {
    this.events = {};
  }

  on(event, cb) {
    this.events[event] = this.events[event] || [];
    if (this.events[event].indexOf(cb) < 0) {
      this.events[event].push(cb);
      return true;
    }
    return false;
  }

  off(event, cb) {
    if (!this.events[event]) {
      return;
    }
    if (!cb) {
      delete this.events[event];
      return;
    }
    const i = this.events[event].indexOf(cb);
    if (i > -1) {
      this.events[event].splice(i, 1);
    }
  }

  emit(event, data) {
    if (!this.events[event]) {
      return;
    }
    this.events[event].forEach((fn) => {
      fn(data);
    });
  }
}
