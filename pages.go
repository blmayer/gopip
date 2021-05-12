package main

const index = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Python package installer</title>
        <meta charset="utf-8"/>
    </head>
    <body>
        <p>
            You can use the package name in the URL for simple names:
            <kbd>https://gopip-vjz2keikqq-uc.a.run.app/requests</kbd>
        </p>
        <p>
            For complex package names like 
            <kbd>git+https://github.com/httpie/httpie.git#egg=httpie</kbd>
            use the form below.
        </p>
        <p>
            <em>Note:</em> Separate packages by a space character.
        </p>
        <form action="/index.html" method="POST">
            <label for="package">Insert package name(s):</label>
            <input name="package" size="40" autofocus>
            <input type="submit">
        </form>
    </body>
</html>`

func successPage(url string) []byte {
	return []byte(
		`<p>
            Success! Get your package <a href="` + url + `">here</a>.
        </p>`,
	)
}
