class EventEmitter {
  constructor() {
    this.events = {};
  }

  on(event, cb) {
    this.events[event] = this.events[event] || [];
    if (this.events[event].indexOf(cb) < 0) {
      this.events[event].push(cb);
    }
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

export default class IO extends EventEmitter {
  constructor(root) {
    this.root = root;
    this.dispatch = root.dispatch;
    this.sids = new EventEmitter();
    this.auto_reconnect = false;
  }

  connect(readyCb) {
    const protocol = location.protocol.toLowerCase() === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${location.host}/socket.io/`);
    this.ws = ws;
    ws.onopen = (evt) => {
      this.connected = true;
      this.emit('ready', evt);
    };

    ws.onclose = (evt) => {
      this.connected = false;
      this.emit('close', evt);
    };

    ws.onmessage = (evt) => {
      this.emit('data', JSON.parse(evt.data));
    };

    ws.onerror = (evt) => {
      this.emit('error', evt);
    }

    if (readyCb) {
      this.on('ready', readyCb);
    }

    this.auto_reconnect = true;
    this.on('close', () => {
      if (!this.auto_reconnect) {
        return;
      }
      setTimeout(() => {
        this.connect();
      }, 2000);
    });
  }

  close() {
    this.auto_reconnect = false;
    this.ws.close();
  }

  subscribe(sid, cb) {
    this.sids.on(`data.${sid}`, cb);
    if (this.connected) {
      this.ws.send(JSON.stringify({ s: sid }));
      return;
    }
    if (!this.auto_reconnect) {
      this.on('ready', () => {});
    }
  }

  unsubscribe(sid) {
    this.sids.off(`data.${sid}`);
  }
}
