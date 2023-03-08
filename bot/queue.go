package bot

import log "github.com/sirupsen/logrus"

type SongQueue struct {
	list    []Song
	current *Song
	Running bool
}

func (queue SongQueue) Get() []Song {
	return queue.list
}

func (queue *SongQueue) Set(list []Song) {
	queue.list = list
}

func (queue *SongQueue) Push(song Song) {
	queue.list = append(queue.list, song)
}

func (queue *SongQueue) Insert(index int, song Song) {
	if len(queue.list) == index { // nil or empty slice or after last element
		queue.list = append(queue.list, song)
		return
	}
	queue.list = append(queue.list[:index+1], queue.list[index:]...)
	queue.list[index] = song
}

func (queue SongQueue) HasNext() bool {
	return len(queue.list) > 0
}

func (queue *SongQueue) Next() Song {
	song := queue.list[0]
	queue.list = queue.list[1:]
	queue.current = &song
	return song
}

func (queue *SongQueue) Pop() bool {
	if len(queue.list) != 0 {
		queue.list = queue.list[:len(queue.list)-1]
		return true
	}
	return false
}

func (queue *SongQueue) Clear() {
	queue.list = make([]Song, 0)
	queue.Running = false
	queue.current = nil
}

func (queue *SongQueue) Start(sess *Session, callback func(track string, duration string, requestedBy string)) {
	queue.Running = true
	for queue.HasNext() && queue.Running {
		song := queue.Next()
		callback(song.Title, song.Duration, "")
		err := sess.Play(song)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warning("Session play error")
		}
		if sess.RepeatMode == RepeatModeOne {
			queue.Insert(0, song)
		}
	}
	if !queue.Running {
		//callback("Stopped playing.")
	} else {
		queue.Running = false
		//callback("Finished queue.")
	}
}

func (queue *SongQueue) Current() *Song {
	return queue.current
}

func (queue *SongQueue) Pause() {
	queue.Running = false
}

func newSongQueue() *SongQueue {
	queue := new(SongQueue)
	queue.list = make([]Song, 0)
	return queue
}
