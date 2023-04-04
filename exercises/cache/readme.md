This is a simple implementation of a cache in Golang with a few eviction policies. The main objective of this project was to work with mutex in a somewhat realistic scenario with a cache of small/medium size and some useful functionality. To provide a better reading experience, I focused on limiting the scope of this project, so it doesn't consume too much of the reader's time. Therefore, there is no advanced logging, error handling, business rules, or extensive testing. Since this project was primarily developed for learning purposes, there might be some bugs.

Future enhancements may include:
- Additional unit and benchmark tests
- Refactoring of benchmark tests to utilize Parallel
- Load testing implementation
- Development of more cache methods to make it more similar to modern ones.