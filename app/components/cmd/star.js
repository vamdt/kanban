import d3 from 'd3';

export default function (sid) {
  if (sid.length < 1) {
    return;
  }

  d3.xhr('/star')
    .header('Content-Type', 'application/x-www-form-urlencoded')
    .post(`s=${sid}`);
  this.$root.$broadcast('star', sid);
}

export function isStar(sid, cb) {
  if (sid.length < 1) {
    cb('sid is empty', false);
    return;
  }

  d3.json(`/star?s=${sid}`, cb);
}
