<p align="center">
  <img width="200" src="/internal/static/gourd.png">
</p>


<h1 align="center">gourd</h1>

<br><br>
Easy serving of questions for interviews, code challenges or similar. Fully rendered server-side utilzing [`templ`](https://github.com/a-h/templ) and [`HTMX`](https://github.com/bigskysoftware/htmx).<br>

## Running gourd
Make sure a Postgres docker container is running. You can use the `run_db.sh` found [here](run_db.sh) if you want.

Run `make serve` from the directory root.
> If no `--config=pathExcludingFilename` is given to `serve`, it will default to looking in the current directory and fail if there is no `config.toml` there.

## Configuration

Everything is configured through a `config.toml` that you have to pass to the `gourd serve` command.

An [example configuration](config_example.toml):

```toml
ApplicationTitle = "gourd" # Main title for the application
ApplicationSubtitle = "Online Assessment" # Subtitle for the Application
LogoPath = "../custom/logo.png" # Local path to a logo shown on the login mask
ServerPort = 8080 # Port to run the server on

# Database configuration
[DB]
Password = "pwd"
User = "local"
Name = "gourd_db"
Host = "localhost"
Port = "5432"

# Sources array
[[Sources]]
URL = "https://github.com/some/special_repository.git" # URL of the GH repository
LocalPath = "../special_repository" # Local path to the repository, including the name
DisplayName = "Gourd Base Course" # Name displayed in the dropdown selection
# Repository Login Credentials
Username = "myusername"
PAT = "mypersonalaccesstoken"
```

The application listens to changes to the `config.toml` at runtime and applies them (such as adding or removing a source).

GitHub repositories are used as dynamic sources for questions to render. Currently, those repositories are assumed to be private,
so a PAT and Username are required.<br>
As of now, the repository structure needs to be:
```
/part_XX
   code.ext
   question.md
/part_XX
   ...
```
Where XX is replaced by a zero-padded number starting from 1, and ext is one of [go | kt | java | js | ts | py | rs].


An admin token will be generated on first start of the application and can be used to generate session tokens for specific users and repositories.<br>
User submissions will be pushed to dedicated branches in the source repository, where they can be reviewed.

## Screenshots

<img src="/screenshots/gourd_landing.png">
<img src="/screenshots/gourd_question.png">
