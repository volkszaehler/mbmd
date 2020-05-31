Vue.component('row', {
	template: '#measurement',
	delimiters: ["${", "}"],
	props: {
		data: Object,
		title: String,
		val: String,
		sum: Boolean,
	},
	data: function () {
		self = this;

		// p determines if the argumnt is non-null
		let p = function (i) {
			return self.data[i] !== undefined && self.data[i] !== null && self.data[i] !== "";
		}

		// val returns addable value: null, NaN and empty are converted to 0
		let v = function (i) {
			let v = parseFloat(self.data[i]);
			return isNaN(v) ? 0 : v;
		}

		// total sum, phase or string definition
		let l123 = p(this.val);
		let l1 = p(this.val+"L1") || p(this.val+"S1");
		let l2 = p(this.val+"L2") || p(this.val+"S2");
		let l3 = p(this.val+"L3") || p(this.val+"S3");

		let valsum;
		if (this.sum) {
			if (l123) {
				valsum = v(this.val);
			} else if (p(this.val+"L1") || p(this.val+"L2") || p(this.val+"L3")) {
				valsum = v(this.val+"L1") + v(this.val+"L2") + v(this.val+"L3");
			} else {
				valsum = v(this.val+"S1") + v(this.val+"S2") + v(this.val+"S3");
			}
		}

		return {
			display: l123 || l1 || l2 || l3,
			l1: l1,
			l2: l2,
			l3: l3,
			val1: (p(this.val+"L1") ? v(this.val+"L1") : v(this.val+"S1")).toFixed(2),
			val2: (p(this.val+"L2") ? v(this.val+"L2") : v(this.val+"S2")).toFixed(2),
			val3: (p(this.val+"L3") ? v(this.val+"L3") : v(this.val+"S3")).toFixed(2),
			valsum: valsum.toFixed(2),
		};
	},
});

let sort = {
	methods: {
		sorted: function (theMap) {
			var devs = Object.keys(theMap);
			devs.sort();
			var res = {};
			devs.forEach(function (key) {
				res[key] = theMap[key];
			});
			return res;
		}
	}
}

var dataapp = new Vue({
	el: '#realtime',
	delimiters: ['${', '}'],
	mixins: [sort],
	data: {
		meters: {},
		message: 'Loading...'
	},
	methods: {
		// pop returns true if it was called with any non-null argument
		pop: function () {
			for(var i=0; i<arguments.length; i++) {
				if (arguments[i] !== undefined && arguments[i] !== null && arguments[i] !== "") {
					return true;
				}
			}
			return false;
		},

		// val returns addable value: null, NaN and empty are converted to 0
		val: function (v) {
			v = parseFloat(v);
			return isNaN(v) ? 0 : v;
		}
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
	mixins: [sort],
	data: {
		meters: {}
	}
})

var fixed = d3.format(".2f")
var si = d3.format(".3~s")

$().ready(function () {
	connectSocket();
});

function convertDate(unixtimestamp){
	var date = new Date(unixtimestamp);
	var day = "0" + date.getDate();
	var month = "0" + (date.getMonth() + 1);
	var year = date.getFullYear();
	return year + '/' + month.substr(-2) + '/' + day.substr(-2);
}

function convertTime(unixtimestamp){
	var date = new Date(unixtimestamp);
	var hours = date.getHours();
	var minutes = "0" + date.getMinutes();
	var seconds = "0" + date.getSeconds();
	return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}

function updateTime(data) {
	timeapp.date = convertDate(data["Timestamp"])
	timeapp.time = convertTime(data["Timestamp"])
}

function updateStatus(status) {
	var id = status["Device"]
	status["Status"] = status["Online"] ? "online" : "offline"

	// update data table
	var dict = statusapp.meters[id] || {}
	dict = Object.assign(dict, status)

	// make update reactive, see
	// https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
	Vue.set(statusapp.meters, id, dict)
}

function updateData(data) {
	// extract the last update
	var id = data["Device"]
	var type = data["IEC61850"]
	var value = fixed(data["Value"])

	// create or update data table
	var dict = dataapp.meters[id] || {}
	dict[type] = value

	// put into statusline
	dataapp.message = "Received " + id + " / " + type + ": " + si(value)

	// make update reactive, see
	// https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
	Vue.set(dataapp.meters, id, dict)
}

function processMessage(data) {
	if (data.Meters && data.Meters.length) {
		for (var i=0; i<data.Meters.length; i++) {
			updateStatus(data.Meters[i]);
		}
	}
	else if (data.Device) {
		updateTime(data);
		updateData(data);
	}
}

function connectSocket() {
	var ws, loc = window.location;
	var protocol = loc.protocol == "https:" ? "wss:" : "ws:"

	// ws = new WebSocket(protocol + "//" + loc.hostname + (loc.port ? ":" + loc.port : "") + "/ws");
	ws = new WebSocket("ws://localhost:8081/ws");

	ws.onerror = function(evt) {
		ws.close();
	}
	ws.onclose = function (evt) {
		window.setTimeout(connectSocket, 1000);
	};
	ws.onmessage = function (evt) {
		var json = JSON.parse(evt.data);
		processMessage(json);
	};
}
