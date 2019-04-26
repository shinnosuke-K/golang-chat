package main

type room struct {
	forward chan []byte
	// チャットの参加
	join chan *client
	// チャットの退出
	leave chan *client
	// 在室しているクライアント一覧
	clients map[*client]bool
}
