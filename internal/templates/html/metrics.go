package html

import "html/template"

var MetricsTemplate = template.Must(template.New("metrics").Parse(`
<html> 
<head> 
    <title>metrics</title>
	<body>{{.metrics}}</title>
</body>
</html>
`))
