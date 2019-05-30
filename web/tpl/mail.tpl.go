package tpl

const MailTpl = `<html>
<head>
	<meta charset="UTF-8">
	<title>temp mail</title>
</head>
<body>
	邮件地址: {{.Addr}}
	邮件内容: {{.Content}}
</body>
</html>`
