This folder contains a collection of Golang concurrency patterns sourced from both the internet and the book Concurrency in Go by Katherine Budday. The primary objective of this project is to serve as a point of reference for future Golang development endeavors. These patterns have been divided into three categories, which are:

- Generators - These patterns are aimed at converting discrete values into streams of values using different methods.
- Processing -These patterns involve the manipulation of channels in various ways, such as merging, dividing, enriching them, and more.
- Signals - These patterns involve signaling different events to other goroutines

All of these patterns are implemented using older technique called Done pattern, now we could also provide context(introduced in Go 1.7) to the function and therefore extend cancellation mechanisms with timeouts, deadlines and regular cancellation