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

var Tmpl = `<html>
<head>
<style>
/* CSS стиль */
.container {
    margin: 20px;
    padding: 10px;
    border: 1px solid #ccc;
    background-color: #f9f9f9;
}
</style>
</head>
<body>
<div class="container">
<h1>Storage Data:</h1>
{{range $key, $value := .Gauges}}
<p>{{ $key }}: {{ $value }}</p>
{{end}}
{{range $key, $value := .Counters}}
<p>{{ $key }}: {{ $value }}</p>
{{end}}
</div>
</body>
</html>`
