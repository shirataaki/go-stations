package middleware

import (
	"log"
	"net/http"
)

type RecoveryHandler struct {
	Message string
}

func NewRecoveryHandler() *RecoveryHandler {
	return &RecoveryHandler{}
}

// HTTPリクエスト処理の際、panicが発生した場合にrecoverする関数
// recover: panicが発生した時、プログラムを終了せずに復帰
func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) { // HTTPリクエストを受け取る
		defer func() { // defer文で関数を呼ぶので、ここはRecovery関数を抜けるときに実行される
			if err := recover(); err != nil { // もしpanicが発生した場合は（panicの時は渡されたエラー値が入る）
				log.Printf("panicが発生したのでrecoverします: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// 以下の処理中にpanicが発生した場合もrecover関数が使われるように
		h.ServeHTTP(w, r) // HTTPリクエストを受け取り、それに対するレスポンスを生成する
	}
	return http.HandlerFunc(fn)
}

// defer: 関数を抜ける前に必ずしておきたい処理がある時に使用
// defer で呼び出している即時関数は、スコープ内で発生したpanicに対して反応
// つまり、defer文が書かれた関数内で後に発生するpanicを見つけて処理する
