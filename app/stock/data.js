import EventEmitter from './event';

export default class IO extends EventEmitter {
  constructor() {
    super();
    this.auto_reconnect = false;

    this.on('ready', () => {
      this.readySubscribe();
    });
  }

  data(evt) {
    const data = JSON.parse(evt.data);
    if (data && data.id) {
      this.emit(`subscribe.${data.id}`, data);
    }
    this.emit('data', evt);
  }

  connect() {
    if (this.connected) {
      return;
    }
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
      this.data(evt);
    };

    ws.onerror = (evt) => {
      this.emit('error', evt);
    };

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
    if (this.ws) {
      this.ws.close();
      this.ws = undefined;
    }
  }

  readySubscribe() {
    Object.keys(this.events).forEach((e) => {
      const ee = e.split('.');
      if (ee.length !== 2) {
        return;
      }
      if (ee[0] !== 'subscribe') {
        return;
      }
      this.ws.send(JSON.stringify({ s: ee[1] }));
    });
  }

  subscribe(sid, cb) {
    this.on(`subscribe.${sid}`, cb);
    if (this.connected) {
      this.ws.send(JSON.stringify({ s: sid }));
      return;
    }

    if (!this.auto_reconnect) {
      this.connect();
    }
  }

  unsubscribe(sid) {
    this.off(`subscribe.${sid}`);
  }
}
