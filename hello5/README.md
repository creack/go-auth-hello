# Hello5 - Token based auth

In this example, we replace the BasicAuth with a token system.
First, we create a `/login` endpoint which validates user credentials
and generates the token, then in our main `hello` endpoint, we validate it.
