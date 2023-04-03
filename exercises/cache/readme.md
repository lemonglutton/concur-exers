This is simple implementation of cache in Golang, with few evicion policies. The main aim of this project was to work a bit with mutext on some more or less real life scenario, which has small/medium size and some cool functionality.
To make reading a more pleasant experience, I tried to narrow the scope of this project as much as possible, so it doesn't take too much of the reader's time. This is why there is no sophisticated logging, error handling, more "business" rules, more tests and so on. This project might have some bug, since it was written purely for sake of learning concurrency.

Maybe in the future:
- More unit tests
- More benchmark tests
- Refactor benchmark for using Parallel
- Add load testing
- Implement more cache methods, to make this cache more similar to modern onces
