<p align="center">
  <img width="200" src="/internal/static/gourd.svg">
</p>


<h1 align="center">gourd</h1>

<br><br>
Easy serving of questions for interviews, code challenges or similar. Fully rendered server-side utilzing [`templ`](https://github.com/a-h/templ) and [`HTMX`](https://github.com/bigskysoftware/htmx).<br>

GitHub repositories are used as dynamic sources for questions to render.<br>
Available repositories can be configured ahead of and during runtime through a `config.toml` given to the `serve` command.<br>

An admin token will be generated on first start and can be used to generate session tokens for specific users and repositories.<br>
User submissions will be pushed to dedicated branches in the source repository, where they can be reviewed.

