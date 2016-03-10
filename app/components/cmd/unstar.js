import d3 from 'd3';
import unwatch from './unwatch';

export default function (sid) {
  if (sid.length < 1) {
    return;
  }

  d3.xhr(`/star?s=${sid}`).send('DELETE');
  unwatch(sid);
}
