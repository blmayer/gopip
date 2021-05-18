package main

const start = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Python package builder</title>
        <meta charset="utf-8"/>
        <meta name="author" content="blmayer"/>
    </head>
    <style>
        body {
            margin: 60px auto;
            max-width: 800px;
        }
        footer {
            margin-top: 80px;
        }
    </style>
    <body>`

const end = `
        <footer>
            <hr>
            <small>
                The code for this project is available here:
                <a href="https://github.com/blmayer/gopip">GitHub</a>.
            </small>
        </footer>
    </body>
</html>`

const index = start + `
        <h3>Python 3 package builder</h3>
        <p>
            You can use the package name in the URL for simple names eg.
            <kbd>https://gopip-vjz2keikqq-uc.a.run.app/requests</kbd>
        </p>
        <p>
            For complex package names like 
            <kbd>git+https://github.com/httpie/httpie.git#egg=httpie</kbd>
            use the form below.
        </p>
        <br>
        <p>
            <em>Note:</em> Separate packages by a space character.
        </p>
        <form action="/package.zip" method="POST" download>
            <label for="package">Insert package name(s):</label>
            <input name="package" size="40" autofocus>
            <input type="submit">
        </form>
` + end

func errorPage(err, detail string) []byte {
	return []byte(start + `
        <h3>Failure!</h3>
        <details>
            <summary>` + err + `:</summary>
            <p>` + detail + `</p>
        </details>` + end,
	)
}
