package html

import "html/template"

var MetricsTemplate = template.Must(template.New("storage").Parse(`
<html> 
<head> 
    <title>Metrics</title>
	<body>
	<h1>Metrics</h1>
	<h2>Gauges</h2>
	<div>{{.gaugeBody}}</div>
	<h2>Counters</h2>
	<div>{{.counterBody}}</div>
	</body>
</html>
`))
