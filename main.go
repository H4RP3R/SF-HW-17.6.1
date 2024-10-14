// Напишите код, в котором имеются два канала сообщений из целых чисел, так,
// чтобы приём сообщений из них никогда не приводил к блокировке и чтобы
// вероятность приёма сообщения из первого канала была выше в 2 раза,
// чем из второго.
// *Если хотите, можете написать код, который бы демонстрировал это соотношение.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"math/rand"
)

const gNum int = 12

var (
	wg      sync.WaitGroup
	sigChan = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // Handle Ctrl+C
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	msgCounter := NewCounter("chan1", "chan2")

	go func(c1, c2 chan int) {
		for {
			for i := 0; i < 3; i++ {
				if i < 2 {
					select {
					case <-c1:
						msgCounter.Increment("chan1")
					default:
					}
				} else {
					select {
					case <-c2:
						msgCounter.Increment("chan2")
					default:
					}
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
	}(c1, c2)

	sendRandIntToChan := func(c chan int) {
		for {
			select {
			case <-sigChan:
				msgCounter.PrintStats()
				os.Exit(0)
			default:
				c <- rand.Intn(10)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	for i := 0; i < gNum; i++ {
		wg.Add(1)
		if i%2 != 0 {
			go sendRandIntToChan(c1)
		} else {
			go sendRandIntToChan(c2)
		}
	}

	printCounter := func() {
		for {
			fmt.Print("\033[H\033[2J") // Clear terminal. Not sure about Win
			fmt.Println("Ctrl+C to quit")
			fmt.Printf("From channel 1: %d\n", msgCounter.Read("chan1"))
			fmt.Printf("From channel 2: %d\n", msgCounter.Read("chan2"))
			time.Sleep(20 * time.Millisecond)
		}
	}

	go printCounter()
	wg.Wait()
}
