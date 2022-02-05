package db

var DB *PrismaClient

func Connect() {
	DB = NewClient()
	if err := DB.Connect(); err != nil {
		panic(err)
	}
}

func Disconnect() {
	DB.Disconnect()
}
