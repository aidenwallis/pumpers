var connection;
var connected = false;
var firstConnect = true;
var reconnectInterval;

var totalCount = document.getElementById('total-count');
var dayCount = document.getElementById('day-count');
var hourCount = document.getElementById('hour-count');
var minuteCount = document.getElementById('minute-count');

function wsErr(error) {
  console.error(error);
  connected = false;
  reconnectInterval = setInterval(connect, 5000);
}

function connect() {
  if (connected) {
    reconnectInterval = undefined;

    return;
  }

  if (connection) {
    connection.close();
  }

  connection = new WebSocket('wss://pumpers.aidenwallis.co.uk/ws');

  connection.onopen = function() {
    connected = true;
    if (!firstConnect) {
      getCount();
    }
    firstConnect = false;
    if (reconnectInterval) {
      clearInterval(reconnectInterval);
      reconnectInterval = undefined;
    }
    console.log('WebSocket connected');
  };

  connection.onclose = wsErr;

  connection.onerror = wsErr;

  connection.onmessage = function(event) {
    try {
      var data = JSON.parse(event.data);
      return updateCounts(data.t, data.d, data.h, data.m);
    } catch (ex) {
      console.error(ex);
    }
  };
};

function getCount() {
  return fetch('/count')
    .then(function(response) { return response.json(); })
    .then(function(data) {
      return updateCounts(data.t, data.d, data.h, data.m);
    });
}

function updateCounts(total, day, hour, minute) {
  totalCount.textContent = total;
  dayCount.textContent = day;
  hourCount.textContent = hour;
  minuteCount.textContent = minute;
}

(function() {

  connect();
  getCount();

})();
