module github.com/SinKingCloud/sinking-go/sinking-consul

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008065154-eb6948aeb02a
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/viper v1.9.0 // indirect
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/gorm v1.21.16 // indirect
)

replace (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a => ../sinking-web
	github.com/SinKingCloud/sinking-go/sinking-websocket v0.0.0-20211008065154-eb6948aeb02a => ../sinking-websocket
)
