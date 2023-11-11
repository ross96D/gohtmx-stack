package templates

const Templ = `package views

templ Index(name string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>%s</title>
		</head>
		<body>
			<h1>This is {name}</h1>
		</body>
	</html>
}`
