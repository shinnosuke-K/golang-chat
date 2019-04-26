package main

type room struct {
	// 転送するためのメッセージを保持
	forward chan []byte
	// チャットの参加
	join chan *client
	// チャットの退出
	leave chan *client
	// 在室しているクライアント一覧
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// すべてのクライアントにメッセージ転送
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
