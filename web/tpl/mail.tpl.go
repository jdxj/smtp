package tpl

const MailTpl = `<html>
<head>
	<meta charset="UTF-8">
	<title>temp mail</title>
</head>
<body>
	mail addr: {{.Addr}}
	content  : {{.Content}}
</body>
</html>`
