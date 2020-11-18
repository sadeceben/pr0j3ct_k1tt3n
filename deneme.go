func SayHello(done chan int) {
    for i := 0; i < 10; i++ {
        fmt.Print(i, " ")
    }
    if done != nil {
        done <- 0 // Signal that we're done
    }
}

func main() {
    SayHello(nil) // Passing nil: we don't want notification here
    done := make(chan int)
    go SayHello(done)
    <-done // Wait until done signal arrives
}
