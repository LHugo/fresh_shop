import {ossUrl} from "../api/api"
function send_request() {
  var xmlhttp = null;
  if (window.XMLHttpRequest) {
    xmlhttp = new XMLHttpRequest();
  } else if (window.ActiveXObject) {
    xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
  }

  if (xmlhttp != null) {
    let tokenUrl = ossUrl
    let serverUrl = tokenUrl + '/oss/v1/oss/token'

    xmlhttp.open("GET", serverUrl, false);
    xmlhttp.send(null);
    return xmlhttp.responseText
  } else {
    alert("Your browser does not support XMLHTTP.");
  }
};
export default send_request
