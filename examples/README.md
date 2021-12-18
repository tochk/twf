# TWF example

Example TWF usage. Contains simple login form and users pages (table and edit).

## Usage

`go run main.go -addr=":8081"`

## File structure

- `example_templates` - directory with templates (quicktemplate)
  - `custom_menu.qtpl` - template file with custom menu
  - `custom_menu.qtpl.go` - generated template
  - `generate.go` - file with generate script
- `main.go` - entrypoint
- `app.go` - base app struct
- `index_page.go` - / page
- `users_page.go` - users pages (table and editing)