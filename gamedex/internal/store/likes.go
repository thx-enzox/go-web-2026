package store

import "sync"

var (
	mu    sync.Mutex
	likes = map[int]int{}
)

func AddLike(gameID int) int {
	mu.Lock()
	defer mu.Unlock()

	likes[gameID]++
	return likes[gameID]
}

func Likes(gameID int) int {
	mu.Lock()
	defer mu.Unlock()

	return likes[gameID]
}

func AllLikes() map[int]int {
	mu.Lock()
	defer mu.Unlock()

	result := make(map[int]int)
	for id, count := range likes {
		result[id] = count
	}

	return result
}
