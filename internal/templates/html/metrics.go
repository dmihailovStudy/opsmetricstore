package html

var MetricsTemplate = `<html>
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
