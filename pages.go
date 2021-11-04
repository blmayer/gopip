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
	    font-size: x-large;
        }
	h2 {
	    font-family: arial;
    	    text-align: center;
    	}
	input {
	    border: thin solid gray;
	    border-radius: 6px;
    	    font-size: inherit;
    	}
        footer {
            margin-top: 80px;
        }
        @media (prefers-color-scheme: dark) {
            body {
                color: #eee;
                background: #151515;
            }
            body a {
                color: #809fff;
            }
	    input {
	        background-color: #343434;
		color: #eee;
	    }
        }
    </style>
    <body>`

const end = `
        <footer>
            <hr>
            <small>
                This work is provided under the BSD 3-Clause License.
		The code for this project is available here:
                <a href="https://github.com/blmayer/gopip">GitHub</a>.
            </small>
        </footer>
    </body>
</html>`

const index = start + `
        <h2>Python 3 package builder</h2>
	<p>
	    Use this web application to build and zip the python packages
	    you want. After inserting the package(s) name(s) and clicking
	    submit, you will be prompted to download the archive.
	</p>
        <p>
            You can also use the package name in the URL for simple names
	    e.g.: <kbd>https://gopip.blmayer.dev/requests</kbd>
        </p>
        <br>
        <p>
            <em>Note:</em> Separate packages by a space character.
        </p>
        <form action="/package.zip" method="POST" download>
            <label for="package">Name(s):</label>
            <input name="package" size="40" autofocus required>
            <input type="submit">
        </form>
` + end

func errorPage(err, detail string) []byte {
	return []byte(start + `
        <h2>Failure!</h2>
        <details>
            <summary>` + err + `:</summary>
            <p>` + detail + `</p>
        </details>` + end,
	)
}
