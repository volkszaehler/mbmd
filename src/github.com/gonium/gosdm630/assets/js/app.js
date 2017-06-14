var statusapp = new Vue({
  el: '#status',
  delimiters: ['${', '}'],
  data: {
	message: 'Loading...'
  }
})

var meterdata = {}

var meterapp = new Vue({
  el: '#meters',
  delimiters: ['${', '}'],
  data: {
	meterdata: meterdata
  }
})


var isActive = true;

$().ready(function () {
  pollServer();
});

function pollServer() {
  if (isActive) {
	window.setTimeout(function () {
      $.ajax({
		url: window.location.href + "/firehose?timeout=45&category=all",
		type: "GET",
		success: function (result) {
		  // extract the last update
		  payload = result["events"][0]["data"]
		  timestamp = payload["ReadTimestamp"]
		  id = payload["DeviceId"]
		  iec61850 = payload["IEC61850"]
		  reading = payload["Value"]
		  // put into statusapp
		  statusapp.message = timestamp + ": " + id + "/" + iec61850 + " - " + reading
		  // update data table
		  var datadict = meterdata[id]
		  //console.log(datadict)
		  if (!datadict) {
			// this is the first time we touch this meter, create an
			// empty dict
			var datadict = {}
			meterdata[id] = datadict
		  }
		  datadict[iec61850] = reading
		  pollServer();
		},
		error: function () {
		  statusapp.message("Error retrieving updates")
		  pollServer();
		}
	  });
	}, 1);
  }
}
