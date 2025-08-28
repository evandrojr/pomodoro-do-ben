
package player

import (
	"log"
	"os/exec"
	"sync"
)

var (
	currentBinauralCmd *exec.Cmd
	binauralMutex      sync.Mutex
)

func Play(file string) {
	cmd := exec.Command("mpv", "--no-video", file)
	err := cmd.Run()
	if err != nil {
		log.Println("Error playing audio:", err)
	}
}

// NewBinauralPlayer retorna uma instância do player de áudios binaurais
func NewBinauralPlayer() *BinauralPlayer {
	return &BinauralPlayer{}
}

// BinauralPlayer gerencia a reprodução de áudios binaurais
type BinauralPlayer struct{}

// Play inicia a reprodução de um áudio binaural
func (bp *BinauralPlayer) Play(file string) {
	binauralMutex.Lock()
	defer binauralMutex.Unlock()
	
	// Parar áudio atual se estiver tocando
	if currentBinauralCmd != nil && currentBinauralCmd.Process != nil {
		currentBinauralCmd.Process.Kill()
	}
	
	// Iniciar novo áudio
	cmd := exec.Command("mpv", "--no-video", "--loop", file)
	currentBinauralCmd = cmd
	
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Println("Error playing binaural audio:", err)
		}
	}()
}

// Stop para a reprodução do áudio binaural atual
func (bp *BinauralPlayer) Stop() {
	binauralMutex.Lock()
	defer binauralMutex.Unlock()
	
	if currentBinauralCmd != nil && currentBinauralCmd.Process != nil {
		currentBinauralCmd.Process.Kill()
		currentBinauralCmd = nil
	}
}
