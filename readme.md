# A starter template using GO, Fiber, Gorm, TailwindCSS, DaisyUI
On the first start a config.json will be created.
On the first start if the user database is empty a admin user will be created with email admin@admin.com and password "password"

## Tailwind
Tailwind: npx tailwindcss -i ./input.css -o ./public/output.css --watch

## Run Server
if installed air  
1. cd into the folder
`air`

else 
`go run server.go`

or 
`go build server.go && /.server` in production
