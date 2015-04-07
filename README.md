# bat
Go implement CLI, cURL-like tool for humans. Bat can be used for testing, debugging, and generally interacting with HTTP servers.

![](images/logo.png)


- [Main Features](#Main Features)
- [Installation](#Installation)
- [Usage](#Usage)
- [HTTP Method](#HTTP Method)
- [Request URL](#Request URL)
- [Request Items](#Request Items)
- [JSON](#JSON)
- [Forms](#Forms)
- [HTTP Headers](#HTTP Headers)

## Main Features

- Expressive and intuitive syntax
- Built-in JSON support
- Forms and file uploads
- HTTPS, proxies, and authentication
- Arbitrary request data
- Custom headers

## Installation

	go get -u github.com/astaxie/bat
	
make sure the `$GOPATH/bin` is added into `$PATH`

## Usage

Hello World:

	$ bat beego.me

Synopsis:

	bat [flags] [METHOD] URL [ITEM [ITEM]]
	
See also `bat --help`.	

### Examples

Custom [HTTP method](#HTTP Method), [HTTP headers](HTTP Headers) and [JSON](#JSON) data:

	$ bat PUT example.org X-API-Token:123 name=John

Submitting forms:

	$ bat -form=true POST example.org hello=World
	
See the request that is being sent using one of the output options:

	$ bat -v example.org

Use Github API to post a comment on an issue with authentication:

	$ bat -a USERNAME POST https://api.github.com/repos/astaxie/bat/issues/1/comments body='HTTPie is awesome!'

Upload a file using redirected input:

	$ bat example.org < file.json
	
Download a file and save it via redirected output:

	$ bat example.org/file > file
	
Download a file wget style:

	$ bat --download example.org/file

Set a custom Host header to work around missing DNS records:

	$ bat localhost:8000 Host:example.com
	
What follows is a detailed documentation. It covers the command syntax, advanced usage, and also features additional examples.
	
## HTTP Method
The name of the HTTP method comes right before the URL argument:

	$ bat DELETE example.org/todos/7
	
Which looks similar to the actual Request-Line that is sent:

DELETE /todos/7 HTTP/1.1

When the METHOD argument is omitted from the command, bat defaults to either GET (with no request data) or POST (with request data).

## Request URL

## Request Items

## JSON

## Forms

## HTTP Headers