var meterapp = new Vue({
  el: '#meters',
  delimiters: ['${', '}'],
  data: {
		meters: {},
		message: 'Loading...'
  }
})

var timeapp = new Vue({
  el: '#time',
  delimiters: ['${', '}'],
  data: {
		time: 'n/a',
		date: 'n/a'
  }
})



var isActive = true;

$().ready(function () {
  pollServer();
});

function convert_date(UNIX_timestamp){
  var date = new Date(UNIX_timestamp);
  var day = "0" + date.getDate();
  var month = "0" + date.getMonth();
  var year = date.getFullYear();
  return year + '/' + month.substr(-2) + '/' + day.substr(-2);
}

function convert_time(UNIX_timestamp){
  var date = new Date(UNIX_timestamp);
  var hours = date.getHours();
  var minutes = "0" + date.getMinutes();
  var seconds = "0" + date.getSeconds();
  return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}

function pollServer(since_time) {
	if (since_time == undefined) {
		since_time = Date.now()
	}
  var loc = window.location;
  if (isActive) {
		var firehose =  loc.protocol + "//" + loc.hostname + (loc.port? ":"+loc.port : "") + "/firehose?timeout=45&category=all&since_time="+since_time;
		$.ajax({
			url: firehose,
			type: "GET",
			success: function (result) {
				var timestamp = result["events"][0]["timestamp"]
				timeapp.time = convert_time(timestamp)
				timeapp.date = convert_date(timestamp)
				// extract the last update
				var payload = result["events"][0]["data"]
				var id = payload["DeviceId"]
				var iec61850 = payload["IEC61850"]
				var reading = payload["Value"].toFixed(2)
				// put into statusline
				meterapp.message = "Received " + id + " / " + reading + " - " + iec61850
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
				pollServer(timestamp);
			},
			error: function () {
				meterapp.message = "Error retrieving updates"
				pollServer();
			}
	  });
  }
}
