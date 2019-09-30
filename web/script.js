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

var currentStats = { t: 0, h: 0, d: 0, m: 0 };

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
      if (data.t !== undefined) {
        currentStats.t = data.t;
      }
      if (data.d !== undefined) {
        currentStats.d = data.d;
      }
      if (data.h !== undefined) {
        currentStats.h = data.h;
      }
      if (data.h !== undefined) {
        currentStats.h = data.h;
      }
      return updateCounts();
    } catch (ex) {
      console.error(ex);
    }
  };
};

function getCount() {
  return fetch('/count')
    .then(function(response) { return response.json(); })
    .then(function(data) {
      currentStats = data;
      return updateCounts();
    });
}

function updateCounts() {
  totalCount.textContent = currentStats.t;
  dayCount.textContent = currentStats.d;
  hourCount.textContent = currentStats.h;
  minuteCount.textContent = currentStats.m;
}

(function() {

  connect();
  getCount();

})();
