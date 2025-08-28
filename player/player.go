package player

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

var (
	currentPlayer oto.Player
	currentFile   io.Closer
	playerMutex   sync.Mutex
)

func Play(file string) {
	playerMutex.Lock()
	defer playerMutex.Unlock()

	if currentPlayer != nil {
		currentPlayer.Close()
	}
	if currentFile != nil {
		currentFile.Close()
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening audio file:", err)
		return
	}
	currentFile = f

	var reader io.Reader

	if strings.HasSuffix(file, ".mp3") {
		d, err := mp3.NewDecoder(f)
		if err != nil {
			fmt.Println("Error creating MP3 decoder:", err)
			return
		}
		reader = d
	} else {
		// We don't support other formats for now
		return
	}

	ctx, ready, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		fmt.Println("Error creating Oto context:", err)
		return
	}
	<-ready

	player := ctx.NewPlayer(reader)
	currentPlayer = player
	player.Play()
}

func Stop() {
	playerMutex.Lock()
	defer playerMutex.Unlock()

	if currentPlayer != nil {
		currentPlayer.Close()
		currentPlayer = nil
	}
	if currentFile != nil {
		currentFile.Close()
		currentFile = nil
	}
}

// BinauralPlayer implementation (can be simplified or removed if not needed)

type BinauralPlayer struct{}

func NewBinauralPlayer() *BinauralPlayer {
	return &BinauralPlayer{}
}

func (bp *BinauralPlayer) Play(file string) {
	Play(file)
}

func (bp *BinauralPlayer) Stop() {
	Stop()
}
