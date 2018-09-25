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

var statusapp = new Vue({
	el: '#status',
	delimiters: ['${', '}'],
	data: {
		meterstatus: {}
	}
})

var fixed = d3.format("~.2f")
var si = d3.format("~s")

$().ready(function () {
	connectSocket();
});

function convert_date(unixtimestamp){
	var date = new Date(unixtimestamp);
	var day = "0" + date.getDate();
	var month = "0" + date.getMonth();
	var year = date.getFullYear();
	return year + '/' + month.substr(-2) + '/' + day.substr(-2);
}

function convert_time(unixtimestamp){
	var date = new Date(unixtimestamp);
	var hours = date.getHours();
	var minutes = "0" + date.getMinutes();
	var seconds = "0" + date.getSeconds();
	return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}

function statusUpdate(meter) {
	var meterid = meter["Id"]
	var metertype = meter["Type"]
	var meterstatus = meter["Status"]

	// update data table
	var datadict = statusapp.meterstatus[meterid]
	if (!datadict) {
		// this is the first time we touch this meter, create an
		// empty dict
		var datadict = {}
	}

	datadict["Id"] = meter["Id"]
	datadict["Type"] = meter["Type"]
	datadict["Status"] = meter["Status"]

	// make update reactive, see
	// https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
	Vue.set(statusapp.meterstatus, meterid, datadict)
}

function meterUpdate(data) {
	timeapp.time = convert_time(data["Timestamp"])
	timeapp.date = convert_date(data["Timestamp"])

	// extract the last update
	var id = data["DeviceId"]
	var iec61850 = data["IEC61850"]
	var reading = fixed(data["Value"])
	// put into statusline
	meterapp.message = "Received #" + id + " / " + iec61850 + ": " + si(reading)
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
}

function processMessage(data) {
	if (data.Modbus) {
		if (data.ConfiguredMeters && data.ConfiguredMeters.length) {
			statusUpdate(data.ConfiguredMeters[0]);
		}
	}
	else if (data.DeviceId) {
		meterUpdate(data);
	}
}

function connectSocket() {
	var ws, loc = window.location;
	var protocol = loc.protocol == "https:" ? "wss:" : "ws:"

	ws = new WebSocket(protocol + "//" + loc.hostname + (loc.port ? ":" + loc.port : "") + "/ws");

	ws.onerror = function(evt) {
		// console.warn("Connection error");
		ws.close();
	}
	ws.onclose = function (evt) {
		// console.warn("Connection closed");
		window.setTimeout(connectSocket, 100);
	};
	ws.onmessage = function (evt) {
		var json = JSON.parse(evt.data);
		processMessage(json);
	};
}
