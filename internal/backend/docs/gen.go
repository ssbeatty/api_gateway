package docs

//go:generate swag init -d ../handler --parseDependency -g admin_login.go -o ./
//go:generate swag init -d ../handler --parseDependency -g endpoints.go -o ./
