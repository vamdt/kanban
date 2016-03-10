import d3 from 'd3';
import watch from './watch';

export default function (sid) {
  if (sid.length < 1) {
    return;
  }

  d3.xhr('/star').post(`s=${sid}`);
  watch(sid);
}
