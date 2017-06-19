
var meterapp = new Vue({
  el: '#meters',
  delimiters: ['${', '}'],
  data: {
	meters: {},
	time: 'n/a',
	message: 'Loading...'
  }
})


var isActive = true;

$().ready(function () {
  pollServer();
});

	
function timeConverter(UNIX_timestamp){
  var date = new Date(UNIX_timestamp);
  // Hours part from the timestamp
  var hours = date.getHours();
  // Minutes part from the timestamp
  var minutes = "0" + date.getMinutes();
  // Seconds part from the timestamp
  var seconds = "0" + date.getSeconds();
  return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}

function pollServer() {
  if (isActive) {
	window.setTimeout(function () {
      $.ajax({
		url: window.location.href + "/firehose?timeout=45&category=all",
		type: "GET",
		success: function (result) {
		  // extract the last update
		  var payload = result["events"][0]["data"]
		  var timestamp = result["events"][0]["timestamp"]
		  var time = timeConverter(timestamp)
		  var id = payload["DeviceId"]
		  var iec61850 = payload["IEC61850"]
		  var reading = payload["Value"].toFixed(2)
		  // put into statusline & update page
		  meterapp.message = time + ": " + id + "/" + iec61850 + " - " + reading
		  meterapp.time = time
		  // update data table
		  var datadict = meterapp.meters[id]
		  if (!datadict) {
			// this is the first time we touch this meter, create an
			// empty dict
			var datadict = {}
		  }
		  datadict[iec61850] = reading
		  // make update reactive, see
		  // https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
		  Vue.set(meterapp.meters, id, datadict)
		  pollServer();
		},
		error: function () {
		  meterapp.message = "Error retrieving updates"
		  pollServer();
		}
	  });
	}, 1);
  }
}
